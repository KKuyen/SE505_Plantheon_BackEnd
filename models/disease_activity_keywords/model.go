package disease_activity_keywords

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DiseaseActivityKeyword represents the many-to-many relationship between diseases and activity_keywords
type DiseaseActivityKeyword struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DiseaseID         string    `json:"disease_id" gorm:"type:uuid;not null;index"`
	ActivityKeywordID string    `json:"activity_keyword_id" gorm:"type:uuid;not null;index"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (dak *DiseaseActivityKeyword) BeforeCreate(tx *gorm.DB) error {
	if dak.ID == "" {
		dak.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (DiseaseActivityKeyword) TableName() string {
	return "disease_activity_keywords"
}
