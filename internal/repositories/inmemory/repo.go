package inmemory

type Repo struct {
	*PostRepo
	*CommentRepo
}

func NewRepo() *Repo {
	postRepo := NewPostRepo()
	commentRepo := NewCommentRepo(postRepo)

	return &Repo{
		PostRepo:    postRepo,
		CommentRepo: commentRepo,
	}
}
