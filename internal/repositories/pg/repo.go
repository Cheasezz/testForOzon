package pg

import "github.com/Cheasezz/testForOzon/pkg/postgres"

const (
	postsTable    = "posts"
	commentsTable = "posts_comments"
)

type Repo struct {
	*PostRepo
	*CommentRepo
}

func NewRepo(db *postgres.Postgres) *Repo {
	return &Repo{
		PostRepo:    NewPostRepo(db),
		CommentRepo: NewCommentRepo(db),
	}
}
