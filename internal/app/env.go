package app

import (
	"github.com/Cheasezz/testForOzon/config"
	"github.com/Cheasezz/testForOzon/internal/pkg/pubsub"
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/Cheasezz/testForOzon/internal/services"
	"github.com/Cheasezz/testForOzon/pkg/logger"
)

type Env struct {
	Logger       *logger.Lg
	Databases    *repositories.DBases
	Repositories *repositories.Repositories
	Services     *services.Services
	PubSub       *pubsub.PubSub
}

func NewEnv(cfg *config.Config) (*Env, error) {
	logger := logger.New(cfg.Log.Level)

	repos, err := repositories.New(cfg)
	if err != nil {
		logger.Error("app-env-repositories.New: %s", err)
	}

	if cfg.APP.MainStorage != "memory" {
		DBMigrate(cfg, logger)
	}

	pubsub := pubsub.NewPubSub()

	services := services.New(repos,pubsub)

	env := Env{
		Logger:       logger,
		Databases:    repos.DBases,
		Repositories: repos,
		Services:     services,
		PubSub:       pubsub,
	}

	return &env, nil
}

func (env *Env) Close() {
	if env.Databases != nil {
		env.Databases.Psql.Close()
	}

}
