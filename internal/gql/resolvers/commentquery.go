package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/google/uuid"
)

// CommentReplies is the resolver for the commentReplies field.
func (r *queryResolver) CommentReplies(ctx context.Context, commentID uuid.UUID, limit *int, offset *int) ([]*core.Comment, error) {
	// Не может быть загружено больше 50 постов за раз
	if *limit > 50 {
		*limit = 50
	}

	// Вернет реплаи к указанному коменту в соответствие с лимитом и офсетом.
	// Глубина всегда 1.
	// Реплаи будут отсортированный от старого к новому
	replies, err := r.env.Services.Comment.GetReplies(ctx, commentID, *limit, *offset)
	if err != nil {
		return nil, err
	}
	return replies, err
}
