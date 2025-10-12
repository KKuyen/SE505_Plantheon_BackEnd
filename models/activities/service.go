package activities

import (
	"plantheon-backend/common"
	"time"

	"gorm.io/gorm"
)

// ActivityService handles all database operations for activities
type ActivityService struct {
	db *gorm.DB
}

// NewActivityService creates a new activity service instance
func NewActivityService() *ActivityService {
	return &ActivityService{
		db: common.GetDB(),
	}
}

// CreateActivityRecord creates a new activity
func CreateActivityRecord(activity *Activity) error {
	service := NewActivityService()
	return service.db.Create(activity).Error
}

// GetActivityByID finds activity by ID
func GetActivityByID(id string) (*Activity, error) {
	service := NewActivityService()
	var activity Activity
	err := service.db.Where("id = ?", id).First(&activity).Error
	return &activity, err
}

// GetAllActivities gets all activities with pagination
func GetAllActivities(offset, limit int) ([]Activity, int64, error) {
	service := NewActivityService()
	var activities []Activity
	var total int64
	
	// Count total records
	if err := service.db.Model(&Activity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := service.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&activities).Error
	return activities, total, err
}

// GetActivitiesByType gets activities by type with pagination
func GetActivitiesByType(activityType string, offset, limit int) ([]Activity, int64, error) {
	service := NewActivityService()
	var activities []Activity
	var total int64
	
	query := service.db.Where("type = ?", activityType)
	
	// Count total records
	if err := query.Model(&Activity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&activities).Error
	return activities, total, err
}

// SearchActivities searches activities by title or description
func SearchActivities(keyword string, offset, limit int) ([]Activity, int64, error) {
	service := NewActivityService()
	var activities []Activity
	var total int64
	
	searchQuery := "%" + keyword + "%"
	query := service.db.Where("title ILIKE ? OR description ILIKE ? OR description2 ILIKE ? OR description3 ILIKE ?", 
		searchQuery, searchQuery, searchQuery, searchQuery)
	
	// Count total records
	if err := query.Model(&Activity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&activities).Error
	return activities, total, err
}

// GetAllActivitiesWithoutPagination gets all activities without pagination
func GetAllActivitiesWithoutPagination() ([]Activity, error) {
	service := NewActivityService()
	var activities []Activity
	err := service.db.Order("created_at DESC").Find(&activities).Error
	return activities, err
}

// GetAllActivitiesByTypeWithoutPagination gets all activities by type without pagination
func GetAllActivitiesByTypeWithoutPagination(activityType string) ([]Activity, error) {
	service := NewActivityService()
	var activities []Activity
	err := service.db.Where("type = ?", activityType).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

// SearchAllActivities searches all activities by title or description without pagination
func SearchAllActivities(keyword string) ([]Activity, error) {
	service := NewActivityService()
	var activities []Activity
	searchQuery := "%" + keyword + "%"
	err := service.db.Where("title ILIKE ? OR description ILIKE ? OR description2 ILIKE ? OR description3 ILIKE ?", 
		searchQuery, searchQuery, searchQuery, searchQuery).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

// GetActivitiesCount gets total count of activities
func GetActivitiesCount() (int64, error) {
	service := NewActivityService()
	var count int64
	err := service.db.Model(&Activity{}).Count(&count).Error
	return count, err
}

// GetActivitiesCountByType gets count of activities by type
func GetActivitiesCountByType(activityType string) (int64, error) {
	service := NewActivityService()
	var count int64
	err := service.db.Model(&Activity{}).Where("type = ?", activityType).Count(&count).Error
	return count, err
}

// SearchActivitiesCount gets count of activities matching search keyword
func SearchActivitiesCount(keyword string) (int64, error) {
	service := NewActivityService()
	var count int64
	searchQuery := "%" + keyword + "%"
	err := service.db.Model(&Activity{}).Where("title ILIKE ? OR description ILIKE ? OR description2 ILIKE ? OR description3 ILIKE ?", 
		searchQuery, searchQuery, searchQuery, searchQuery).Count(&count).Error
	return count, err
}

// UpdateActivity updates activity information
func UpdateActivity(activity *Activity) error {
	service := NewActivityService()
	return service.db.Save(activity).Error
}

// DeleteActivity deletes activity by ID
func DeleteActivity(id string) error {
	service := NewActivityService()
	return service.db.Where("id = ?", id).Delete(&Activity{}).Error
}

// GetActivitiesByMonthYear returns activities whose time_start or day fall within the given month/year (UTC)
func GetActivitiesByMonthYear(year int, month int) ([]Activity, error) {
	service := NewActivityService()
    startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

    var activities []Activity
    // Half-open interval [startOfMonth, startOfNextMonth)
    err := service.db.Where(
        "time_start IS NOT NULL AND time_start >= ? AND time_start < ?",
        startOfMonth, startOfNextMonth,
    ).Order("created_at DESC").Find(&activities).Error
    return activities, err
}

// GetActivitiesByDay returns activities that match the specific day (UTC) by either time_start's date or day field
func GetActivitiesByDay(day time.Time) ([]Activity, error) {
	service := NewActivityService()
    // Normalize to date (midnight UTC)
    start := time.Date(day.UTC().Year(), day.UTC().Month(), day.UTC().Day(), 0, 0, 0, 0, time.UTC)
    next := start.AddDate(0, 0, 1)

    var activities []Activity
    err := service.db.Where(
        "time_start IS NOT NULL AND time_start >= ? AND time_start < ?",
        start, next,
    ).Order("created_at DESC").Find(&activities).Error
    return activities, err
}