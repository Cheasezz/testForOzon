package inmemory

import (
	"context"
	"fmt"
	"sort"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/google/uuid"
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

func (r *PostRepo) GetPosts(ctx context.Context, id *uuid.UUID, limit, offset *int) ([]*core.Post, error) {
	fmt.Println("GetPosts inmemory repo func call")

	var posts []*core.Post

	if id != nil {
		post, err := r.posts.Load(id.String())
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)

		return posts, nil
	} else {
		// Сортируем посты по `createdAt` (свежие сверху)
		r.posts.m.Range(func(_, value interface{}) bool {
			post := value.(core.Post)
			posts = append(posts, &post)
			return true
		})

		// Сортируем посты по `createdAt` (свежие сверху)
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].CreatedAt.After(posts[j].CreatedAt)
		})

		//Настройка пагинации
		start := *offset
		end := *offset + *limit
		if start > len(posts) {
			return []*core.Post{}, errOffsetToBid
		}
		if end > len(posts) {
			end = len(posts)
		}

		return posts[start:end], nil
	}
}
