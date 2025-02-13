package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/pkg/postgres"
)

var (
	errCreatePost = errors.New("qwe")
)

type PostRepo struct {
	db *postgres.Postgres
}

func NewPostRepo(db *postgres.Postgres) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(ctx context.Context, post core.Post) (*core.Post, error) {
	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, created_at, title, content, comments_allowed) 
												values ($1, $2, $3, $4, $5, $6) 
												RETURNING *`, postTable)
	var createdPost core.Post
	err := r.db.Scany.Get(ctx, r.db.Pool, &createdPost, query, post.Id, post.UserId, post.CreatedAt, post.Title, post.Content, post.CommentsAllowed)

	if err != nil {
		return nil, err
	}
	fmt.Println("CreatePost pg repo func call")
	return &createdPost, nil
}
