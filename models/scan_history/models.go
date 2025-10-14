package scan_history

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScanHistory struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"not null"`
	DiseaseID string    `json:"disease_id" gorm:"not null;type:uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// BeforeCreate will set a UUID rather than numeric ID.
func (s *ScanHistory) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

