package server

import (
	"context"
	"errors"
	"flag"
	"github.com/morticius/accuracy/internal/server/config"
	"github.com/morticius/accuracy/internal/server/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var pathToConfigFile string

func init() {
	flag.StringVar(&pathToConfigFile, "c", "", "path to config file")
	flag.Parse()
	if pathToConfigFile == "" {
		log.Fatalln("missing required flag -c with path to config file")
	}

}

func Run() {
	config.InitConfig(pathToConfigFile)

	router := routes.GetRouter("file")

	server := &http.Server{
		Addr: config.Get().Addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
