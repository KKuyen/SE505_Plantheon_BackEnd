package plants

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plant struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;type:text"`
	Description *string   `json:"description" gorm:"type:text"`
	ImageURL    *string   `json:"image_url" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Plant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate will update the UpdatedAt field.
func (p *Plant) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
