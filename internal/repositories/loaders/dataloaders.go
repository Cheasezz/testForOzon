package loaders

import "github.com/Cheasezz/testForOzon/internal/repositories"

const (
	DataLoadersContextKey          string = "DataLoadersContextKey"
	repliesCountLoaderByIDMaxBatch int    = 100
)

type DataLoaders struct {
	RepliesCountLoaderByID *RepliesCountLoader
}

func NewDataLoaders(repos *repositories.Repositories) *DataLoaders {
	return &DataLoaders{
		RepliesCountLoaderByID: NewRepliesCountLoader(repos, repliesCountLoaderByIDMaxBatch),
	}
}
