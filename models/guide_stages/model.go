package guide_stages

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GuideStage struct {
	ID            string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PlantID       string     `json:"plant_id" gorm:"not null;type:uuid"`
	StageTitle    string     `json:"stage_title" gorm:"not null;type:text"`
	StartDayOffset int       `json:"start_day_offset" gorm:"not null"`
	EndDayOffset  int        `json:"end_day_offset" gorm:"not null"`
	ImageURL      *string    `json:"image_url" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (gs *GuideStage) BeforeCreate(tx *gorm.DB) error {
	if gs.ID == "" {
		gs.ID = uuid.New().String()
	}
	if gs.CreatedAt.IsZero() {
		gs.CreatedAt = time.Now()
	}
	return nil
}
