package activities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Description     *string   `json:"description" gorm:"type:text"`
	Description2    *string   `json:"description2" gorm:"type:text"`
	Description3    *string   `json:"description3" gorm:"type:text"`
	TimeStart       *time.Time `json:"time_start" gorm:"type:timestamp"`
    TimeEnd         *time.Time `json:"time_end" gorm:"type:timestamp"`
    Day             *bool      `json:"day" gorm:"type:boolean"`
	Money           *float64  `json:"money" gorm:"type:decimal(15,2)"`
    Type            string    `json:"type" gorm:"type:varchar(255);not null"`
	Title           string    `json:"title" gorm:"not null;type:varchar(255)"`
	IsRepeat        *string   `json:"is_repeat" gorm:"type:varchar(50)"`
    Repeat          *string   `json:"repeat" gorm:"type:varchar(50)"`
	EndRepeatDay    *time.Time `json:"end_repeat_day" gorm:"type:timestamp"`
	AlertTime       *string   `json:"alert_time" gorm:"type:varchar(50)"`
	Object          *string   `json:"object" gorm:"type:varchar(255)"`
	Amount          *int      `json:"amount" gorm:"type:integer"`
	Unit            *string   `json:"unit" gorm:"type:varchar(50)"`
	Purpose         *string   `json:"purpose" gorm:"type:text"`
	TargetPerson    *string   `json:"target_person" gorm:"type:varchar(255)"`
	SourcePerson    *string   `json:"source_person" gorm:"type:varchar(255)"`
	AttachedLink    *string   `json:"attached_link" gorm:"type:text"`
	Note            *string   `json:"note" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

