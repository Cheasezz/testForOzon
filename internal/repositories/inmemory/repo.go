package inmemory

type Repo struct {
	*TestRepo
}

func NewRepo() *Repo {
	return &Repo{TestRepo: NewTestRepo()}
}
