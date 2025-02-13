package core

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id              uuid.UUID `json:"id" db:"id"`
	CreatedAt       time.Time `json:"createdAt" db:"crated_at"`
	Title           string    `json:"title" db:"title"`
	Content         string    `json:"content" db:"content"`
	CommentsAllowed bool      `json:"commentsAllowed" db:"comments_allowed"`
}
