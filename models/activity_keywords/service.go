package activity_keywords

import (
	"plantheon-backend/common"
	"gorm.io/gorm"
)

// ActivityKeywordService handles all database operations for activity keywords
type ActivityKeywordService struct {
	db *gorm.DB
}

// NewActivityKeywordService creates a new activity keyword service instance
func NewActivityKeywordService() *ActivityKeywordService {
	return &ActivityKeywordService{
		db: common.GetDB(),
	}
}

// CreateActivityKeyword creates a new activity keyword
func CreateActivityKeyword(activityKeyword *ActivityKeyword) error {
	service := NewActivityKeywordService()
	return service.db.Create(activityKeyword).Error
}

// GetActivityKeywordByID finds activity keyword by ID
func GetActivityKeywordByID(id string) (*ActivityKeyword, error) {
	service := NewActivityKeywordService()
	var activityKeyword ActivityKeyword
	err := service.db.Where("id = ?", id).First(&activityKeyword).Error
	return &activityKeyword, err
}

// GetAllActivityKeywords gets all activity keywords
func GetAllActivityKeywords() ([]ActivityKeyword, error) {
	service := NewActivityKeywordService()
	var activityKeywords []ActivityKeyword
	err := service.db.Find(&activityKeywords).Error
	return activityKeywords, err
}

// GetActivityKeywordsByIDs gets multiple activity keywords by their IDs
func GetActivityKeywordsByIDs(ids []string) ([]ActivityKeyword, error) {
	service := NewActivityKeywordService()
	var activityKeywords []ActivityKeyword
	err := service.db.Where("id IN ?", ids).Find(&activityKeywords).Error
	return activityKeywords, err
}

// UpdateActivityKeyword updates activity keyword information
func UpdateActivityKeyword(id string, updates map[string]interface{}) error {
	service := NewActivityKeywordService()
	return service.db.Model(&ActivityKeyword{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteActivityKeyword deletes activity keyword by ID
func DeleteActivityKeyword(id string) error {
	service := NewActivityKeywordService()
	return service.db.Where("id = ?", id).Delete(&ActivityKeyword{}).Error
}

// SearchActivityKeywords searches activity keywords by name
func SearchActivityKeywords(keyword string) ([]ActivityKeyword, error) {
	service := NewActivityKeywordService()
	var activityKeywords []ActivityKeyword
	err := service.db.Where("name ILIKE ?", "%"+keyword+"%").Find(&activityKeywords).Error
	return activityKeywords, err
}

// GetActivityKeywordsPaginated gets activity keywords with pagination
func GetActivityKeywordsPaginated(page int, limit int) ([]ActivityKeyword, int64, error) {
	service := NewActivityKeywordService()
	var activityKeywords []ActivityKeyword
	var total int64

	// Count total records
	if err := service.db.Model(&ActivityKeyword{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated records
	err := service.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&activityKeywords).Error
	if err != nil {
		return nil, 0, err
	}

	return activityKeywords, total, nil
}

// QueryActivityKeywords queries activity keywords by name and/or type with pagination
func QueryActivityKeywords(name string, keywordType string, page int, limit int) ([]ActivityKeyword, int64, error) {
	service := NewActivityKeywordService()
	var activityKeywords []ActivityKeyword
	var total int64

	// Build query
	query := service.db.Model(&ActivityKeyword{})

	// Filter by name (ILIKE for case-insensitive search)
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Filter by type (exact match)
	if keywordType != "" {
		query = query.Where("type = ?", keywordType)
	}

	// Count total records matching filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated records
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&activityKeywords).Error
	if err != nil {
		return nil, 0, err
	}

	return activityKeywords, total, nil
}
