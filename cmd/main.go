package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cheasezz/testForOzon/config"
	"github.com/Cheasezz/testForOzon/internal/app"
	httpHandlers "github.com/Cheasezz/testForOzon/internal/transport/http"
	"github.com/Cheasezz/testForOzon/pkg/httpserver"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	env, err := app.NewEnv(cfg)
	if err != nil {
		log.Fatalf("cannot create environment: %s", err)
	}
	defer env.Close()

	handlers := httpHandlers.New(env, cfg.HTTP.Port)

	srv := httpserver.New(handlers, httpserver.Port(cfg.HTTP.Port))
	env.Logger.Info("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-quit:
		env.Logger.Info("cmd - main - signal: %s", s.String())
	case err = <-srv.Notify():
		env.Logger.Error("cmd - main - httpServer.Notify: %s", err)
	}

	if err := srv.Shutdown(); err != nil {
		env.Logger.Error("error occured on server shutting down: %s", err.Error())
	}

	env.Logger.Info("App Shutting Down")
}
