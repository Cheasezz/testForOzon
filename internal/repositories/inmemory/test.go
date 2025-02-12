package inmemory

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
)

type TestRepo struct {
	data map[int]core.Test
}

func NewTestRepo() *TestRepo {
	return &TestRepo{data: make(map[int]core.Test)}
}

func (r *TestRepo) GetTest(ctx context.Context) error {
	fmt.Println("GetTest func call (from inmemory)")
	return nil
}
