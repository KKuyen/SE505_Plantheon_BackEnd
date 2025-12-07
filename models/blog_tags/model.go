package blog_tags

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlogTag represents a tag for blogs/news
type BlogTag struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null;unique;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

// BeforeCreate ensures defaults are set
func (bt *BlogTag) BeforeCreate(tx *gorm.DB) error {
	if bt.ID == "" {
		bt.ID = uuid.New().String()
	}
	if bt.CreatedAt.IsZero() {
		bt.CreatedAt = time.Now()
	}
	if bt.UpdatedAt.IsZero() {
		bt.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate refreshes updated_at
func (bt *BlogTag) BeforeUpdate(tx *gorm.DB) error {
	bt.UpdatedAt = time.Now()
	return nil
}
