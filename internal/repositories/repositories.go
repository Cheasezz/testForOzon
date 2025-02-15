package repositories

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/config"
	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/internal/repositories/inmemory"
	"github.com/Cheasezz/testForOzon/internal/repositories/pg"
	"github.com/Cheasezz/testForOzon/pkg/postgres"
	"github.com/google/uuid"
)

type mainDB interface {
	CreatePost(ctx context.Context, post core.Post) (*core.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*core.Post, error)
	GetPost(ctx context.Context, postId uuid.UUID) (*core.Post, error)
	CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error)
	GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error)
	GetRepliesById(ctx context.Context, parentCommentId uuid.UUID, limit, offset int) ([]*core.Comment, error)
}

type DBases struct {
	Psql *postgres.Postgres
}

type Repositories struct {
	*DBases
	mainDB
}

func New(cfg *config.Config) (*Repositories, error) {
	switch cfg.APP.MainStorage {
	case "postgres":
		fmt.Print(cfg.PG.URL)
		psql, err := postgres.New(cfg.PG.URL)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to postgres: %v", err)
		}

		return &Repositories{DBases: &DBases{Psql: psql}, mainDB: pg.NewRepo(psql)}, nil

	case "memory":
		return &Repositories{mainDB: inmemory.NewRepo()}, nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", cfg.APP.MainStorage)
	}
}
