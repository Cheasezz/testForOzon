package core

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	PostId    uuid.UUID `json:"postId" db:"post_id"`
	Id        uuid.UUID `json:"id" db:"id"`
	ParentId  uuid.UUID `json:"parentId" db:"parent_id"`
	CreatedAt time.Time `json:"createdAt" db:"crated_at"`
	Content   string    `json:"content" db:"content"`
}
