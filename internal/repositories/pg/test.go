package pg

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/pkg/postgres"
)

type TestRepo struct {
	db *postgres.Postgres
}

func NewTestRepo(db *postgres.Postgres) *TestRepo {
	return &TestRepo{db: db}
}

func (r *TestRepo) GetTest(ctx context.Context) error {
	fmt.Println("GetTest func call")
	return nil
}
