package services

import (
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/Cheasezz/testForOzon/pkg/logger"
	"github.com/Cheasezz/testForOzon/pkg/pubsub"
)

type Services struct {
	Post
	Comment
}

func New(repos *repositories.Repositories, pubsub pubsub.IPubSub, log logger.Logger) *Services {
	return &Services{
		Post:    NewPostService(repos),
		Comment: NewCommentService(repos.CommentRepo, pubsub, log),
	}
}
