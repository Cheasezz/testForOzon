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

// type mainDB interface {
// 	CreatePost(ctx context.Context, post core.Post) (*core.Post, error)
// 	GetPosts(ctx context.Context, limit, offset int) ([]*core.Post, error)
// 	GetPost(ctx context.Context, postId uuid.UUID) (*core.Post, error)
// 	CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error)
// 	GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error)
// 	GetRepliesById(ctx context.Context, parentCommentId uuid.UUID, limit, offset int) ([]*core.Comment, error)
// 	RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error)
// 	GetRepliesCounts(ctx context.Context, ids []uuid.UUID) (map[string]int, error)
// }

type postRepo interface {
	CreatePost(ctx context.Context, post core.Post) (*core.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*core.Post, error)
	GetPost(ctx context.Context, postId uuid.UUID) (*core.Post, error)
}

type commentRepo interface {
	CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error)
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
	postRepo
	commentRepo
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
			postRepo:    pg.NewPostRepo(psql),
			commentRepo: pg.NewCommentRepo(psql),
		}, nil

	case "memory":
		postRepo := inmemory.NewPostRepo()

		return &Repositories{
			postRepo:    postRepo,
			commentRepo: inmemory.NewCommentRepo(postRepo),
		}, nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", cfg.APP.MainStorage)
	}
}
