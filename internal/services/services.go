package services

import (
	"github.com/Cheasezz/testForOzon/pkg/pubsub"
	"github.com/Cheasezz/testForOzon/internal/repositories"
)

type Services struct {
	Post
	Comment
}

func New(repos *repositories.Repositories, pubsub *pubsub.PubSub) *Services {
	return &Services{
		Post:    NewPostService(repos),
		Comment: NewCommentService(repos, pubsub),
	}
}
