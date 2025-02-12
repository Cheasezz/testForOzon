package pg

import "github.com/Cheasezz/testForOzon/pkg/postgres"

type Repo struct {
	*TestRepo
}

func NewRepo(db *postgres.Postgres) *Repo {
	return &Repo{
		TestRepo: NewTestRepo(db),
	}
}
