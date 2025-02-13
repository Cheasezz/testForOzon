package services

import (
	"github.com/Cheasezz/testForOzon/internal/repositories"
)

type Services struct {
	Post
}

func New(repos *repositories.Repositories) *Services {
	return &Services{
		Post: NewPostService(repos),
	}
}
