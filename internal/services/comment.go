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
var errCmntAreProh = errors.New("comments are prohibited")

type Comment interface {
	CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error)
	GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset, depth int) ([]*core.Comment, error)
	// GetReplies(ctx context.Context, obj *core.Comment, limit, offset *int) ([]*core.Comment, error)
}

type CommentService struct {
	repo *repositories.Repositories
}

func NewCommentService(db *repositories.Repositories) *CommentService {
	return &CommentService{repo: db}
}

func (s *CommentService) CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error) {
	fmt.Println("CreateComment service func call")
	post, err := s.repo.GetPost(ctx, input.PostId)
	if err != nil {
		
		return nil, err
	}
	if !post.CommentsAllowed {
		return nil, errCmntAreProh
	} 
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

func (s *CommentService) GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset, depth int) ([]*core.Comment, error) {
	fmt.Println("GetRootComments service func call")

	rootComments, err := s.repo.GetRootComments(ctx, postId, limit, offset)
	if err != nil {
		return nil, err
	}

	var flatComments []*core.Comment

	for _, comment := range rootComments {
		flatComments = append(flatComments, comment)
		// Рекурсивно собираем вложенные комментарии, если глубина позволяет
		if depth > 1 {
			nested, err := s.getFlatReplies(ctx, comment, limit, offset, depth-1)
			if err != nil {
				return nil, err
			}
			flatComments = append(flatComments, nested...)
		}
	}

	return flatComments, nil
}

func (s *CommentService) getFlatReplies(ctx context.Context, comment *core.Comment, limit, offset, depth int) ([]*core.Comment, error) {
	fmt.Println("getFlatReplies service func call")

	replies, err := s.repo.GetRepliesById(ctx, comment.Id, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []*core.Comment
	for _, reply := range replies {
		result = append(result, reply)
		if depth > 1 {
			nested, err := s.getFlatReplies(ctx, reply, limit, offset, depth-1)
			if err != nil {
				return nil, err
			}
			result = append(result, nested...)
		}
	}
	return result, nil
}
