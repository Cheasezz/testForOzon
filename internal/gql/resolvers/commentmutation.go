package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
)

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input core.CommentCreateInput) (*core.Comment, error) {
	// В самом начале проверяет, можно ли оставлять комменты под постом.
	// Добавит коммент с id поста и коммента под которым оставлен.
	// Если id коммента под которым оставлен создаваемый коммент не указан (речь про input.ParentId),
	// То создаваемый коммент считает корневым, т.е относится именно к посту
	comment, err := r.env.Services.Comment.CreateComment(ctx, input)
	if err != nil {
		return nil, err
	}
	fmt.Println("CreateComment mutation resolver")
	return comment, nil
}
