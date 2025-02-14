package services

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/google/uuid"
)

var errToLongtext = errors.New("comment is too long (max 2000 characters)")

type Comment interface {
	CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error)
}

type CommentService struct {
	repo *repositories.Repositories
}

func NewCommentService(db *repositories.Repositories) *CommentService {
	return &CommentService{repo: db}
}

func (s *CommentService) CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error) {

	if utf8.RuneCountInString(input.Content) > 2000 {
		return nil, errToLongtext
	}

	newComment := core.Comment{
		PostId:    input.PostId,
		Id:        uuid.New(),
		ParentId:  input.ParentId,
		UserId:    input.UserId,
		CreatedAt: time.Now(),
		Content:   input.Content,
	}
	return s.repo.CreateComment(ctx, newComment)
}
