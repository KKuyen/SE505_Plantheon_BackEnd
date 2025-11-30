package disease_activity_keywords

import (
	"plantheon-backend/common"
	"gorm.io/gorm"
)

// DiseaseActivityKeywordService handles all database operations for disease-activity keyword relationships
type DiseaseActivityKeywordService struct {
	db *gorm.DB
}

// NewDiseaseActivityKeywordService creates a new service instance
func NewDiseaseActivityKeywordService() *DiseaseActivityKeywordService {
	return &DiseaseActivityKeywordService{
		db: common.GetDB(),
	}
}

// AddActivityKeywordToDisease adds an activity keyword to a disease
func AddActivityKeywordToDisease(diseaseID, activityKeywordID string) error {
	service := NewDiseaseActivityKeywordService()
	
	// Check if relationship already exists
	var existing DiseaseActivityKeyword
	err := service.db.Where("disease_id = ? AND activity_keyword_id = ?", diseaseID, activityKeywordID).First(&existing).Error
	if err == nil {
		return nil // Already exists, no error
	}
	
	relation := DiseaseActivityKeyword{
		DiseaseID:         diseaseID,
		ActivityKeywordID: activityKeywordID,
	}
	return service.db.Create(&relation).Error
}

// RemoveActivityKeywordFromDisease removes an activity keyword from a disease
func RemoveActivityKeywordFromDisease(diseaseID, activityKeywordID string) error {
	service := NewDiseaseActivityKeywordService()
	return service.db.Where("disease_id = ? AND activity_keyword_id = ?", diseaseID, activityKeywordID).
		Delete(&DiseaseActivityKeyword{}).Error
}

// GetActivityKeywordsByDiseaseID gets all activity keywords for a specific disease
func GetActivityKeywordsByDiseaseID(diseaseID string) ([]string, error) {
	service := NewDiseaseActivityKeywordService()
	var relations []DiseaseActivityKeyword
	err := service.db.Where("disease_id = ?", diseaseID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	
	keywordIDs := make([]string, len(relations))
	for i, rel := range relations {
		keywordIDs[i] = rel.ActivityKeywordID
	}
	return keywordIDs, nil
}

// GetDiseasesByActivityKeywordID gets all diseases that have a specific activity keyword
func GetDiseasesByActivityKeywordID(activityKeywordID string) ([]string, error) {
	service := NewDiseaseActivityKeywordService()
	var relations []DiseaseActivityKeyword
	err := service.db.Where("activity_keyword_id = ?", activityKeywordID).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	
	diseaseIDs := make([]string, len(relations))
	for i, rel := range relations {
		diseaseIDs[i] = rel.DiseaseID
	}
	return diseaseIDs, nil
}

// SetActivityKeywordsForDisease sets all activity keywords for a disease (replaces existing)
func SetActivityKeywordsForDisease(diseaseID string, activityKeywordIDs []string) error {
	service := NewDiseaseActivityKeywordService()
	
	// Start transaction
	return service.db.Transaction(func(tx *gorm.DB) error {
		// Remove all existing relationships
		if err := tx.Where("disease_id = ?", diseaseID).Delete(&DiseaseActivityKeyword{}).Error; err != nil {
			return err
		}
		
		// Add new relationships
		for _, keywordID := range activityKeywordIDs {
			relation := DiseaseActivityKeyword{
				DiseaseID:         diseaseID,
				ActivityKeywordID: keywordID,
			}
			if err := tx.Create(&relation).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// DeleteAllByDiseaseID deletes all relationships for a disease
func DeleteAllByDiseaseID(diseaseID string) error {
	service := NewDiseaseActivityKeywordService()
	return service.db.Where("disease_id = ?", diseaseID).Delete(&DiseaseActivityKeyword{}).Error
}

// DeleteAllByActivityKeywordID deletes all relationships for an activity keyword
func DeleteAllByActivityKeywordID(activityKeywordID string) error {
	service := NewDiseaseActivityKeywordService()
	return service.db.Where("activity_keyword_id = ?", activityKeywordID).Delete(&DiseaseActivityKeyword{}).Error
}
