package blogs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	ID                string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title             string     `json:"title" gorm:"not null;type:text"`
	Content           string     `json:"content" gorm:"not null;type:text"`
	CoverImageURL     *string    `json:"cover_image_url" gorm:"type:text"`
	SubGuideStagesID  *string    `json:"sub_guide_stages_id" gorm:"type:uuid"`
	UserID            string     `json:"user_id" gorm:"not null;type:uuid"`
	Status            string     `json:"status" gorm:"type:varchar(50);default:'draft'"`
	PublishedAt       *time.Time `json:"published_at" gorm:"type:timestamp with time zone"`
	CreatedAt         time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Blog) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate will update the UpdatedAt field.
func (b *Blog) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
