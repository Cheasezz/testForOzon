package services

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/repositories"
)

type Test interface {
	CheckTest(ctx context.Context) error
}

type TestService struct {
	repo *repositories.Repositories
}

func NewTestService(db *repositories.Repositories) *TestService {
	return &TestService{repo: db}
}

func (s *TestService) CheckTest(ctx context.Context) error {
	if err := s.repo.GetTest(ctx); err != nil {
		return err
	}
	fmt.Println("checkTest func call")
	return nil
}
