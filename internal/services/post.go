package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/google/uuid"
)

type Post interface {
	CreatePost(ctx context.Context, input core.PostCreateInput) (*core.Post, error)
	GetPosts(ctx context.Context, id *uuid.UUID, limit, offset *int) ([]*core.Post, error)
}

type PostService struct {
	repo *repositories.Repositories
}

func NewPostService(db *repositories.Repositories) *PostService {
	return &PostService{repo: db}
}

func (s *PostService) CreatePost(ctx context.Context, input core.PostCreateInput) (*core.Post, error) {
	newPost := core.Post{
		Id:              uuid.New(),
		UserId:          input.UserId,
		CreatedAt:       time.Now().UTC(),
		Title:           input.Title,
		Content:         input.Content,
		CommentsAllowed: input.CommentsAllowed,
	}

	post, err := s.repo.CreatePost(ctx, newPost)
	if err != nil {
		return nil, err
	}
	fmt.Println("CreatePost post service func call")
	return post, nil
}

func (s *PostService) GetPosts(ctx context.Context, id *uuid.UUID, limit, offset *int) ([]*core.Post, error) {
	posts, err := s.repo.GetPosts(ctx, id, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
