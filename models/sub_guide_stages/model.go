package sub_guide_stages

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubGuideStage struct {
	ID              string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GuideStagesID   *string `json:"guide_stages_id" gorm:"type:uuid"`
	Title           *string `json:"title" gorm:"type:text"`
	StartDayOffset  int     `json:"start_day_offset" gorm:"not null"`
	EndDayOffset    int     `json:"end_day_offset" gorm:"not null"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (sgs *SubGuideStage) BeforeCreate(tx *gorm.DB) error {
	if sgs.ID == "" {
		sgs.ID = uuid.New().String()
	}
	return nil
}
