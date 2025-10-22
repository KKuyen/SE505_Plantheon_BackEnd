package post_likes

import (
	"time"

	"gorm.io/gorm"
)

type PostLike struct {
	UserID    string    `json:"user_id" gorm:"not null;type:uuid"`
	PostID    string    `json:"post_id" gorm:"not null;type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (pl *PostLike) BeforeCreate(tx *gorm.DB) error {
	if pl.CreatedAt.IsZero() {
		pl.CreatedAt = time.Now()
	}
	return nil
}

// TableName returns the table name for PostLike
func (PostLike) TableName() string {
	return "post_likes"
}
