package core

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	PostId    uuid.UUID  `json:"postId" db:"post_id"`
	Id        uuid.UUID  `json:"id" db:"id"`
	ParentId  *uuid.UUID `json:"parentId" db:"parent_id"`
	UserId    string     `json:"userId" db:"user_id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	Content   string     `json:"content" db:"content"`
}

type CommentCreateInput struct {
	UserId   string    `json:"userId" `
	PostId   uuid.UUID `json:"postId" `
	ParentId *uuid.UUID `json:"parentId" `
	Content  string    `json:"content" `
}
