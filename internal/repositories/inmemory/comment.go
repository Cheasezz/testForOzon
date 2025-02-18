package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/pkg/gSyncMap"
	"github.com/google/uuid"
)

var (
	errCmntDisabled       = errors.New("comments are disabled for this post")
	errOffsetToBid        = errors.New("offset for commentaries pagination is to big")
	errParentCmntNotExist = errors.New("parent comment does not exist")
)

type CommentRepo struct {
	comments *gSyncMap.GSyncMap[string, *core.Comment]
	posts    *gSyncMap.GSyncMap[string, *core.Post]
}

func NewCommentRepo(pr *PostRepo) *CommentRepo {
	return &CommentRepo{
		comments: gSyncMap.NewGenericSyncMap[string, *core.Comment](),
		posts:    pr.posts,
	}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error) {
	fmt.Println("Comments inmemory repo func call")

	// Проверяем, если comment.ParentId != nil, то родительский комментарий должен существовать
	if comment.ParentId != nil {
		_, err := r.comments.Load(comment.ParentId.String())
		if err != nil {
			return nil, errParentCmntNotExist
		}
	}

	comment.Id = uuid.New()
	comment.CreatedAt = time.Now()

	r.comments.Store(comment.Id.String(), &comment)
	return &comment, nil
}

func (r *CommentRepo) GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetRootComments inmemory repo func call")

	var rootComments []*core.Comment
	//Собераем комментарии для соответствующего postId
	r.comments.Range(func(key string, comment *core.Comment) bool {
		if comment.PostId.String() == postId.String() && comment.ParentId == nil {
			rootComments = append(rootComments, comment)
		}
		return true
	})

	// Сортировка сначала старые
	sort.Slice(rootComments, func(i, j int) bool {
		return rootComments[j].CreatedAt.After(rootComments[i].CreatedAt)
	})

	//Настройка пагинации
	start := offset
	end := offset + limit
	if start > len(rootComments) {
		return []*core.Comment{}, errOffsetToBid

	}
	if end > len(rootComments) {
		end = len(rootComments)
	}

	return rootComments[start:end], nil
}

func (r *CommentRepo) GetRepliesById(ctx context.Context, parentCommentId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetReplies inmemory repo func call")

	var repliesComments []*core.Comment
	//Собераем комментарии ответы для соответствующего parentId
	r.comments.Range(func(key string, comment *core.Comment) bool {
		if comment.ParentId != nil && comment.ParentId.String() == parentCommentId.String() {
			repliesComments = append(repliesComments, comment)
		}
		return true
	})

	// Сортировка сначала старые
	sort.Slice(repliesComments, func(i, j int) bool {
		return repliesComments[j].CreatedAt.After(repliesComments[i].CreatedAt)
	})

	//Настройка пагинации
	start := offset
	end := offset + limit
	if start > len(repliesComments) {
		return []*core.Comment{}, errOffsetToBid

	}
	if end > len(repliesComments) {
		end = len(repliesComments)
	}

	return repliesComments[start:end], nil
}

func (r *CommentRepo) RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error) {
	count := 0
	r.comments.Range(func(key string, comment *core.Comment) bool {
		if comment.ParentId != nil && comment.ParentId.String() == commentId.String() {
			count++
		}
		return true
	})
	return count, nil
}

func (r *CommentRepo) GetRepliesCounts(ctx context.Context, ids []uuid.UUID) (map[string]int, error) {
	countMap := make(map[string]int)
	// Инициализируем для всех ключей нулевым значением.
	for _, id := range ids {
		countMap[id.String()] = 0
	}

	r.comments.Range(func(key string, comment *core.Comment) bool {
		if comment.ParentId != nil {
			key := comment.ParentId.String()
			if _, exists := countMap[key]; exists {
				countMap[key]++
			}
		}
		return true
	})
	return countMap, nil
}
