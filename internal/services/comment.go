package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/google/uuid"
)

var errToLongtext = errors.New("comment is too long (max 2000 characters)")

type Comment interface {
	CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error)
	GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset *int) ([]*core.Comment, error)
	GetReplies(ctx context.Context, obj *core.Comment, limit, offset *int) ([]*core.Comment, error)
}

type CommentService struct {
	repo *repositories.Repositories
}

func NewCommentService(db *repositories.Repositories) *CommentService {
	return &CommentService{repo: db}
}

func (s *CommentService) CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error) {
	fmt.Println("CreateComment service func call")

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

	comment, err := s.repo.CreateComment(ctx, newComment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset *int) ([]*core.Comment, error) {
	fmt.Println("GetRootComments service func call")

	comments, err := s.repo.GetRootComments(ctx, postId, limit, offset)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) GetReplies(ctx context.Context, obj *core.Comment, limit, offset *int) ([]*core.Comment, error) {
	fmt.Println("GetReplies service func call")

	comments, err := s.repo.GetReplies(ctx, obj, limit, offset)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
