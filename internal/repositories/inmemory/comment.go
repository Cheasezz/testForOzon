package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/google/uuid"
)

var (
	errCmntDisabled       = errors.New("comments are disabled for this post")
	errOffsetToBid        = errors.New("offset for commentaries pagination is to big")
	errParentCmntNotExist = errors.New("parent comment does not exist")
)

type CommentRepo struct {
	comments *GenericMap[string, *core.Comment]
	posts    *GenericMap[string, *core.Post]
}

func NewCommentRepo(pr *PostRepo) *CommentRepo {
	return &CommentRepo{
		comments: &GenericMap[string, *core.Comment]{},
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
	r.comments.m.Range(func(_, value interface{}) bool {
		comment := value.(*core.Comment)
		if comment.PostId.String() == postId.String() && comment.ParentId == nil {
			rootComments = append(rootComments, comment)
		}
		return true
	})

	sort.Slice(rootComments, func(i, j int) bool {
		return rootComments[i].CreatedAt.After(rootComments[j].CreatedAt)
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

	var allComments []*core.Comment
	//Собераем комментарии ответы для соответствующего parentId
	r.comments.m.Range(func(_, value interface{}) bool {
		comment := value.(*core.Comment)
		if comment.ParentId != nil && comment.ParentId.String() == parentCommentId.String() {
			allComments = append(allComments, comment)
		}
		return true
	})

	//Настройка пагинации
	start := offset
	end := offset + limit
	if start > len(allComments) {
		return []*core.Comment{}, errOffsetToBid

	}
	if end > len(allComments) {
		end = len(allComments)
	}

	return allComments[start:end], nil
}

func (r *CommentRepo) RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error) {
	count := 0
	r.comments.m.Range(func(_, value interface{}) bool {
		cmnt := value.(*core.Comment)
		if cmnt.ParentId != nil && cmnt.ParentId.String() == commentId.String() {
			count++
		}
		return true
	})
	return count, nil
}
