package activity_keywords

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityKeyword struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `json:"name" gorm:"not null"`
	Description    *string   `json:"description" gorm:"type:text"`
	Type           string    `json:"type" gorm:"not null;default:'GENERAL'"`
	BaseDaysOffset int       `json:"base_days_offset" gorm:"not null;default:0"`
	IsFreeTime     bool      `json:"is_free_time" gorm:"not null;default:false"`
	HourTime       *int      `json:"hour_time"`
	EndHourTime    *int      `json:"end_hour_time"`
	TimeDuration   *int      `json:"time_duration"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (ak *ActivityKeyword) BeforeCreate(tx *gorm.DB) error {
	if ak.ID == "" {
		ak.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (ActivityKeyword) TableName() string {
	return "activity_keywords"
}
