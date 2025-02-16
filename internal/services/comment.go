package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/Cheasezz/testForOzon/internal/repositories/loaders"
	"github.com/google/uuid"
)

var errToLongtext = errors.New("comment is too long (max 2000 characters)")
var errCmntAreProh = errors.New("comments are prohibited")

type Comment interface {
	CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error)
	GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error)
	GetReplies(ctx context.Context, commentId uuid.UUID, limit, offset int) ([]*core.Comment, error)
	RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error)
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

func (s *CommentService) GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetRootComments service func call")

	rootComments, err := s.repo.GetRootComments(ctx, postId, limit, offset)
	if err != nil {
		return nil, err
	}

	return rootComments, nil
}

func (s *CommentService) GetReplies(ctx context.Context, commentId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("getFlatReplies service func call")

	replies, err := s.repo.GetRepliesById(ctx, commentId, limit, offset)
	if err != nil {
		return nil, err
	}

	return replies, nil
}

func (s *CommentService) RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error) {
	fmt.Println("RepliesCount service func call")

	dataloader, ok := ctx.Value(loaders.DataLoadersContextKey).(*loaders.DataLoaders)
	if ok && dataloader.RepliesCountLoaderByID != nil {
		count, err := dataloader.RepliesCountLoaderByID.L.Load(ctx, commentId)

		if err != nil {
			return 0, err
		}
		return count, nil
	}

	count, err := s.repo.RepliesCount(ctx, commentId)
	if err != nil {
		return 0, err
	}
	return count, err
}
