package inmemory

import (
	"context"
	"errors"
	"time"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/google/uuid"
)

var errCmntDisabled = errors.New("comments are disabled for this post")

type CommentRepo struct {
	comments *GenericMap[string, core.Comment]
	posts    *GenericMap[string, core.Post]
}

func NewCommentRepo(pr *PostRepo) *CommentRepo {
	return &CommentRepo{
		comments: &GenericMap[string, core.Comment]{},
		posts:    pr.posts,
	}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error) {
	post, err := r.posts.Load(comment.PostId.String())
	if err != nil {
		return nil, err
	}

	if !post.CommentsAllowed {
		return nil, errCmntDisabled
	}

	comment.Id = uuid.New()
	comment.CreatedAt = time.Now()

	r.comments.Store(comment.Id.String(), comment)
	return &comment, nil
}
