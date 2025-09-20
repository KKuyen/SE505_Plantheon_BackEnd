package diseases

import (
	"time"

	"plantheon-backend/common"

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

// CreateDiseaseRecord creates a new disease
func CreateDiseaseRecord(disease *Disease) error {
	return common.GetDB().Create(disease).Error
}

// GetDiseaseByID finds disease by ID
func GetDiseaseByID(id string) (*Disease, error) {
	var disease Disease
	err := common.GetDB().Where("id = ?", id).First(&disease).Error
	return &disease, err
}

// GetAllDiseases gets all diseases with pagination
func GetAllDiseases(offset, limit int) ([]Disease, int64, error) {
	var diseases []Disease
	var total int64
	
	// Count total records
	if err := common.GetDB().Model(&Disease{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := common.GetDB().Offset(offset).Limit(limit).Find(&diseases).Error
	return diseases, total, err
}

// GetDiseasesByType gets diseases by type with pagination
func GetDiseasesByType(diseaseType string, offset, limit int) ([]Disease, int64, error) {
	var diseases []Disease
	var total int64
	
	query := common.GetDB().Where("type = ?", diseaseType)
	
	// Count total records
	if err := query.Model(&Disease{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&diseases).Error
	return diseases, total, err
}

// SearchDiseases searches diseases by name or description
func SearchDiseases(keyword string, offset, limit int) ([]Disease, int64, error) {
	var diseases []Disease
	var total int64
	
	searchQuery := "%" + keyword + "%"
	query := common.GetDB().Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery)
	
	// Count total records
	if err := query.Model(&Disease{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&diseases).Error
	return diseases, total, err
}

// GetDiseaseByClassName gets disease by class name
func GetDiseaseByClassName(className string) (*Disease, error) {
	var disease Disease
	err := common.GetDB().Where("class_name = ?", className).First(&disease).Error
	return &disease, err
}

// UpdateDisease updates disease information
func UpdateDisease(disease *Disease) error {
	return common.GetDB().Save(disease).Error
}

// DeleteDisease deletes disease by ID
func DeleteDisease(ClassName string) error {
	return common.GetDB().Delete(&Disease{}, "class_name = ?", ClassName).Error
}
