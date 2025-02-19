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

//go:generate mockgen -package repositories -source=repositories.go -destination=mocks_repositories.go

type PostRepo interface {
	CreatePost(ctx context.Context, post core.Post) (*core.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*core.Post, error)
	GetPost(ctx context.Context, postId uuid.UUID) (*core.Post, error)
}

type CommentRepo interface {
	CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error)
	CommentForPostAllowed(ctx context.Context, postId uuid.UUID) (bool, error)
	GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error)
	GetRepliesById(ctx context.Context, parentCommentId uuid.UUID, limit, offset int) ([]*core.Comment, error)
	RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error)
	GetRepliesCounts(ctx context.Context, ids []uuid.UUID) (map[string]int, error)
}

type DBases struct {
	Psql *postgres.Postgres
}

type Repositories struct {
	*DBases
	PostRepo
	CommentRepo
}

func New(cfg *config.Config) (*Repositories, error) {
	switch cfg.APP.MainStorage {

	case "postgres":
		psql, err := postgres.New(cfg.PG.URL)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to postgres: %v", err)
		}

		return &Repositories{
			DBases:      &DBases{Psql: psql},
			PostRepo:    pg.NewPostRepo(psql),
			CommentRepo: pg.NewCommentRepo(psql),
		}, nil

	case "memory":
		postRepo := inmemory.NewPostRepo()

		return &Repositories{
			PostRepo:    postRepo,
			CommentRepo: inmemory.NewCommentRepo(postRepo),
		}, nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", cfg.APP.MainStorage)
	}
}
