package services

import (
	"github.com/Cheasezz/testForOzon/internal/repositories"
)

type Services struct {
	Test
}

func New(repos *repositories.Repositories) *Services {
	return &Services{
		Test: NewTestService(repos),
	}
}
