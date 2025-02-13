package inmemory

type Repo struct {
	*PostRepo
}

func NewRepo() *Repo {
	return &Repo{
		PostRepo: NewPostRepo(),
	}
}
