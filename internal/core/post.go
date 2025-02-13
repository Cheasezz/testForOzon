package core

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id              uuid.UUID `json:"id" db:"id"`
	UserId          string`json:"userId" db:"user_id"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	Title           string    `json:"title" db:"title"`
	Content         string    `json:"content" db:"content"`
	CommentsAllowed bool      `json:"commentsAllowed" db:"comments_allowed"`
}

type PostCreateInput struct {
	UserId          string`json:"userId"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	CommentsAllowed bool      `json:"commentsAllowed"`
}
