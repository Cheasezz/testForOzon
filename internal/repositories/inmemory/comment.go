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

var errCmntDisabled = errors.New("comments are disabled for this post")
var errOffsetToBid = errors.New("offset for commentaries pagination is to big")

type CommentRepo struct {
	comments *GenericMap[string, core.Comment]
	posts    *GenericMap[string, *core.Post]
}

func NewCommentRepo(pr *PostRepo) *CommentRepo {
	return &CommentRepo{
		comments: &GenericMap[string, core.Comment]{},
		posts:    pr.posts,
	}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error) {
	fmt.Println("Comments PostResolver func call")

	post, err := r.posts.Load(comment.PostId.String())
	if err != nil {
		return nil, err
	}

	if !post.CommentsAllowed {
		return nil, errCmntDisabled
	}

	comment.Id = uuid.New()
	comment.CreatedAt = time.Now()

	r.comments.Store(comment.Id.String(), comment)
	return &comment, nil
}

func (r *CommentRepo) GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetRootComments inmemory repo func call")

	var allComments []*core.Comment
	//Собераем комментарии для соответствующего postId
	r.comments.m.Range(func(_, value interface{}) bool {
		comment := value.(core.Comment)
		if comment.PostId.String() == postId.String() && comment.ParentId == nil {
			allComments = append(allComments, &comment)
		}
		return true
	})

	sort.Slice(allComments, func(i, j int) bool {
		return allComments[i].CreatedAt.After(allComments[j].CreatedAt)
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

func (r *CommentRepo) GetRepliesById(ctx context.Context, parentCommentId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetReplies inmemory repo func call")

	var allComments []*core.Comment
	//Собераем комментарии ответы для соответствующего parentId
	r.comments.m.Range(func(_, value interface{}) bool {
		comment := value.(core.Comment)
		if comment.ParentId != nil && comment.ParentId.String() == parentCommentId.String() {
			allComments = append(allComments, &comment)
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
