package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Content   string         `json:"content" gorm:"type:text"`
	ImageLink pq.StringArray `json:"image_link" gorm:"type:text[]"`
	UserID    string         `json:"user_id" gorm:"not null;type:uuid"`
	LikeNum   int            `json:"like_number" gorm:"default:0"`
	CommentNum int           `json:"comment_number" gorm:"default:0"`
	ShareNum   int            `json:"share_number" gorm:"default:0"`
	Tags       pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}