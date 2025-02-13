package pg

import "github.com/Cheasezz/testForOzon/pkg/postgres"

const (
	postTable = "posts"
)

type Repo struct {
	*PostRepo
}

func NewRepo(db *postgres.Postgres) *Repo {
	return &Repo{
		PostRepo: NewPostRepo(db),
	}
}
