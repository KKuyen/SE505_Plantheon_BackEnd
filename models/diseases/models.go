package diseases

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Disease struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	ClassName   string         `json:"class_name" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Solution    string         `json:"solution" gorm:"type:text"`
	ImageLink   pq.StringArray `json:"image_link" gorm:"type:text[]"`
	PlantName   string         `json:"plant_name"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (d *Disease) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

