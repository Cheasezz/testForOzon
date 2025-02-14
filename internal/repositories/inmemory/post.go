package inmemory

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
)

type PostRepo struct {
	posts *GenericMap[string, core.Post]
}

func NewPostRepo() *PostRepo {
	return &PostRepo{posts: &GenericMap[string, core.Post]{}}
}

func (r *PostRepo) CreatePost(ctx context.Context, post core.Post) (*core.Post, error) {
	fmt.Println("CreatePost inmemory repo func call")

	r.posts.Store(post.Id.String(), post)
	res, err := r.posts.Load(post.Id.String())
	if err != nil {
		return nil, err
	}

	return &res, nil
}
