package diseases

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

// DiseaseService handles all database operations for diseases
type DiseaseService struct {
	db *gorm.DB
}

// NewDiseaseService creates a new disease service instance
func NewDiseaseService() *DiseaseService {
	return &DiseaseService{
		db: common.GetDB(),
	}
}

// CreateDiseaseRecord creates a new disease
func CreateDiseaseRecord(disease *Disease) error {
	service := NewDiseaseService()
	return service.db.Create(disease).Error
}

// GetDiseaseByID finds disease by ID
func GetDiseaseByID(id string) (*Disease, error) {
	service := NewDiseaseService()
	var disease Disease
	err := service.db.Where("id = ?", id).First(&disease).Error
	return &disease, err
}

// GetAllDiseases gets all diseases with pagination
func GetAllDiseases(offset, limit int) ([]Disease, int64, error) {
	service := NewDiseaseService()
	var diseases []Disease
	var total int64
	
	// Count total records
	if err := service.db.Model(&Disease{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := service.db.Offset(offset).Limit(limit).Find(&diseases).Error
	return diseases, total, err
}

// GetDiseasesByType gets diseases by type with pagination
func GetDiseasesByType(diseaseType string, offset, limit int) ([]Disease, int64, error) {
	service := NewDiseaseService()
	var diseases []Disease
	var total int64
	
	query := service.db.Where("type = ?", diseaseType)
	
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
	service := NewDiseaseService()
	var diseases []Disease
	var total int64
	
	searchQuery := "%" + keyword + "%"
	query := service.db.Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery)
	
	// Count total records
	if err := query.Model(&Disease{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&diseases).Error
	return diseases, total, err
}

// GetAllDiseasesWithoutPagination gets all diseases without pagination
func GetAllDiseasesWithoutPagination() ([]Disease, error) {
	service := NewDiseaseService()
	var diseases []Disease
	err := service.db.Find(&diseases).Error
	return diseases, err
}

// GetAllDiseasesByTypeWithoutPagination gets all diseases by type without pagination
func GetAllDiseasesByTypeWithoutPagination(diseaseType string) ([]Disease, error) {
	service := NewDiseaseService()
	var diseases []Disease
	err := service.db.Where("type = ?", diseaseType).Find(&diseases).Error
	return diseases, err
}

// SearchAllDiseases searches all diseases by name or description without pagination
func SearchAllDiseases(keyword string) ([]Disease, error) {
	service := NewDiseaseService()
	var diseases []Disease
	searchQuery := "%" + keyword + "%"
	err := service.db.Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).Find(&diseases).Error
	return diseases, err
}

// GetDiseasesCount gets total count of diseases
func GetDiseasesCount() (int64, error) {
	service := NewDiseaseService()
	var count int64
	err := service.db.Model(&Disease{}).Count(&count).Error
	return count, err
}

// GetDiseasesCountByType gets count of diseases by type
func GetDiseasesCountByType(diseaseType string) (int64, error) {
	service := NewDiseaseService()
	var count int64
	err := service.db.Model(&Disease{}).Where("type = ?", diseaseType).Count(&count).Error
	return count, err
}

// SearchDiseasesCount gets count of diseases matching search keyword
func SearchDiseasesCount(keyword string) (int64, error) {
	service := NewDiseaseService()
	var count int64
	searchQuery := "%" + keyword + "%"
	err := service.db.Model(&Disease{}).Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).Count(&count).Error
	return count, err
}

// GetDiseaseByClassName gets disease by class name
func GetDiseaseByClassName(className string) (*Disease, error) {
	service := NewDiseaseService()
	var disease Disease
	err := service.db.Where("class_name = ?", className).First(&disease).Error
	return &disease, err
}

// UpdateDisease updates disease information
func UpdateDisease(disease *Disease) error {
	service := NewDiseaseService()
	return service.db.Save(disease).Error
}

// DeleteDisease deletes disease by ID
func DeleteDisease(ClassName string) error {
	service := NewDiseaseService()
	return service.db.Delete(&Disease{}, "class_name = ?", ClassName).Error
}