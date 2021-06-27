package file

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/morticius/accuracy/internal/models"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UserFileRepository struct {
	pathToSave string
}

func NewUserFileRepository(pathToSave string) *UserFileRepository {
	return &UserFileRepository{pathToSave: pathToSave}
}

func (r *UserFileRepository) GetByEmail(c context.Context, email string) (*models.User, error) {
	h := md5.New()
	io.WriteString(h, email)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	filePath := r.getPathToFileFromHash(hash)

	var user models.User

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &user, fmt.Errorf("user not exists")
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &user, err
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (r *UserFileRepository) Save(c context.Context, user models.User) error {
	h := md5.New()
	io.WriteString(h, user.Email)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	filePath := r.getPathToFileFromHash(hash)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		f, err := r.createFile(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("user exists")
}

func (r *UserFileRepository) getPathToFileFromHash(hash string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", r.pathToSave, hash[0:2], hash[2:4], hash[4:6], hash)
}

func (r *UserFileRepository) createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}