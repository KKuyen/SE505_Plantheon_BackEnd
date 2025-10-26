package sub_guide_stages

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

// SubGuideStageService handles all database operations for sub guide stages
type SubGuideStageService struct {
	db *gorm.DB
}

// NewSubGuideStageService creates a new sub guide stage service instance
func NewSubGuideStageService() *SubGuideStageService {
	return &SubGuideStageService{
		db: common.GetDB(),
	}
}

// CreateSubGuideStage creates a new sub guide stage
func CreateSubGuideStage(subGuideStage *SubGuideStage) error {
	service := NewSubGuideStageService()
	return service.db.Create(subGuideStage).Error
}

// GetSubGuideStageByID finds sub guide stage by ID
func GetSubGuideStageByID(id string) (*SubGuideStage, error) {
	service := NewSubGuideStageService()
	var subGuideStage SubGuideStage
	err := service.db.Where("id = ?", id).First(&subGuideStage).Error
	return &subGuideStage, err
}

// GetSubGuideStagesByGuideStageID gets all sub guide stages for a specific guide stage
func GetSubGuideStagesByGuideStageID(guideStageID string) ([]SubGuideStage, error) {
	service := NewSubGuideStageService()
	var subGuideStages []SubGuideStage
	err := service.db.Where("guide_stages_id = ?", guideStageID).Order("start_day_offset ASC").Find(&subGuideStages).Error
	return subGuideStages, err
}

// UpdateSubGuideStage updates sub guide stage information
func UpdateSubGuideStage(subGuideStage *SubGuideStage) error {
	service := NewSubGuideStageService()
	return service.db.Save(subGuideStage).Error
}

// DeleteSubGuideStage deletes sub guide stage by ID
func DeleteSubGuideStage(id string) error {
	service := NewSubGuideStageService()
	return service.db.Where("id = ?", id).Delete(&SubGuideStage{}).Error
}
