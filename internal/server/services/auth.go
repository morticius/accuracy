package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/morticius/accuracy/internal/models"
	"github.com/morticius/accuracy/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type AuthService struct{
	userRepository repositories.UserRepository
	jwtService JWTService
}

func NewAuthService(ur repositories.UserRepository, js JWTService) *AuthService {
	return &AuthService{
		userRepository: ur,
		jwtService: js,
	}
}

type AuthCredentials struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (s *AuthService) Create(c *gin.Context) {
	var credential AuthCredentials
	err := c.Bind(&credential)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(credential.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}

	user := models.User{
		Email:    credential.Email,
		Password: string(hash),
	}

	err = s.userRepository.Save(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *AuthService) SignIn(c *gin.Context) {
	var credential AuthCredentials
	err := c.Bind(&credential)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userRepository.GetByEmail(c, credential.Email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	token := s.jwtService.GenerateToken(user.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"expiresAt": time.Now().Unix() + 3600 * 48,
	})


	fmt.Println(token)


}