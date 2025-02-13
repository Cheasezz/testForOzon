package inmemory

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
)

type PostRepo struct {
	data *GenericMap[string, core.Post]
}

func NewPostRepo() *PostRepo {
	return &PostRepo{data: &GenericMap[string, core.Post]{}}
}

func (r *PostRepo) CreatePost(ctx context.Context, post core.Post) (*core.Post, error) {

	r.data.Store(post.Id.String(), post)
	res, err := r.data.Load(post.Id.String())
	if err != nil {
		return nil, err
	}
	fmt.Println("CreatePost inmemory repo func call")
	fmt.Println(res)
	return &res, nil
}
