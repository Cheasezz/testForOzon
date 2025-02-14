package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/pkg/postgres"
	"github.com/google/uuid"
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
	fmt.Println("CreatePost pg repo func call")

	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, created_at, title, content, comments_allowed) 
												values ($1, $2, $3, $4, $5, $6) 
												RETURNING *`, postsTable)
	var createdPost core.Post
	err := r.db.Scany.Get(ctx, r.db.Pool, &createdPost, query,
		post.Id, post.UserId, post.CreatedAt, post.Title, post.Content, post.CommentsAllowed)

	if err != nil {
		return nil, err
	}

	return &createdPost, nil
}

func (r *PostRepo) GetPosts(ctx context.Context, id *uuid.UUID, limit, offset *int) ([]*core.Post, error) {
	fmt.Println("GetPosts pg repo func call")

	var posts []*core.Post
	if id != nil {
		var post core.Post
		query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postsTable)

		err := r.db.Scany.Get(ctx, r.db.Pool, &post, query, id)

		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	} else {
		query := fmt.Sprintf(`SELECT * FROM %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, postsTable)

		err := r.db.Scany.Select(ctx, r.db.Pool, &posts, query, limit, offset)

		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}
