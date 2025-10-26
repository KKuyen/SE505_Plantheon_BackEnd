package guide_stages

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

// GuideStageService handles all database operations for guide stages
type GuideStageService struct {
	db *gorm.DB
}

// NewGuideStageService creates a new guide stage service instance
func NewGuideStageService() *GuideStageService {
	return &GuideStageService{
		db: common.GetDB(),
	}
}

// CreateGuideStage creates a new guide stage
func CreateGuideStage(guideStage *GuideStage) error {
	service := NewGuideStageService()
	return service.db.Create(guideStage).Error
}

// GetGuideStageByID finds guide stage by ID
func GetGuideStageByID(id string) (*GuideStage, error) {
	service := NewGuideStageService()
	var guideStage GuideStage
	err := service.db.Where("id = ?", id).First(&guideStage).Error
	return &guideStage, err
}

// GetGuideStagesByPlantID gets all guide stages for a specific plant
func GetGuideStagesByPlantID(plantID string) ([]GuideStage, error) {
	service := NewGuideStageService()
	var guideStages []GuideStage
	err := service.db.Where("plant_id = ?", plantID).Order("start_day_offset ASC").Find(&guideStages).Error
	return guideStages, err
}

// UpdateGuideStage updates guide stage information
func UpdateGuideStage(guideStage *GuideStage) error {
	service := NewGuideStageService()
	return service.db.Save(guideStage).Error
}

// DeleteGuideStage deletes guide stage by ID
func DeleteGuideStage(id string) error {
	service := NewGuideStageService()
	return service.db.Where("id = ?", id).Delete(&GuideStage{}).Error
}
