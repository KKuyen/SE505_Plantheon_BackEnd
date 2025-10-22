package noti

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title     string     `json:"title" gorm:"type:varchar"`
	Content   string     `json:"content" gorm:"type:varchar"`
	IsRead    bool       `json:"is_read" gorm:"type:boolean"`
	PostID    *string    `json:"post_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	return nil
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "noti"
}
