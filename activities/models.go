package activities

import (
	"time"

	"plantheon-backend/common"

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
	Day             *time.Time `json:"day" gorm:"type:timestamp"`
	Money           *float64  `json:"money" gorm:"type:decimal(15,2)"`
    Type            string    `json:"type" gorm:"type:varchar(255);not null"`
	Title           string    `json:"title" gorm:"not null;type:varchar(255)"`
	IsRepeat        *string   `json:"is_repeat" gorm:"type:varchar(50)"`
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

// CreateActivityRecord creates a new activity
func CreateActivityRecord(activity *Activity) error {
	return common.GetDB().Create(activity).Error
}

// GetActivityByID finds activity by ID
func GetActivityByID(id string) (*Activity, error) {
	var activity Activity
	err := common.GetDB().Where("id = ?", id).First(&activity).Error
	return &activity, err
}

// GetAllActivities gets all activities with pagination
func GetAllActivities(offset, limit int) ([]Activity, int64, error) {
	var activities []Activity
	var total int64
	
	// Count total records
	if err := common.GetDB().Model(&Activity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err := common.GetDB().Offset(offset).Limit(limit).Order("created_at DESC").Find(&activities).Error
	return activities, total, err
}

// GetActivitiesByType gets activities by type with pagination
func GetActivitiesByType(activityType string, offset, limit int) ([]Activity, int64, error) {
	var activities []Activity
	var total int64
	
	query := common.GetDB().Where("type = ?", activityType)
	
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
	var activities []Activity
	var total int64
	
	searchQuery := "%" + keyword + "%"
	query := common.GetDB().Where("title ILIKE ? OR description ILIKE ? OR description2 ILIKE ? OR description3 ILIKE ?", 
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
	var activities []Activity
	err := common.GetDB().Order("created_at DESC").Find(&activities).Error
	return activities, err
}

// GetAllActivitiesByTypeWithoutPagination gets all activities by type without pagination
func GetAllActivitiesByTypeWithoutPagination(activityType string) ([]Activity, error) {
	var activities []Activity
	err := common.GetDB().Where("type = ?", activityType).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

// SearchAllActivities searches all activities by title or description without pagination
func SearchAllActivities(keyword string) ([]Activity, error) {
	var activities []Activity
	searchQuery := "%" + keyword + "%"
	err := common.GetDB().Where("title ILIKE ? OR description ILIKE ? OR description2 ILIKE ? OR description3 ILIKE ?", 
		searchQuery, searchQuery, searchQuery, searchQuery).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

// GetActivitiesCount gets total count of activities
func GetActivitiesCount() (int64, error) {
	var count int64
	err := common.GetDB().Model(&Activity{}).Count(&count).Error
	return count, err
}

// GetActivitiesCountByType gets count of activities by type
func GetActivitiesCountByType(activityType string) (int64, error) {
	var count int64
	err := common.GetDB().Model(&Activity{}).Where("type = ?", activityType).Count(&count).Error
	return count, err
}

// SearchActivitiesCount gets count of activities matching search keyword
func SearchActivitiesCount(keyword string) (int64, error) {
	var count int64
	searchQuery := "%" + keyword + "%"
	err := common.GetDB().Model(&Activity{}).Where("title ILIKE ? OR description ILIKE ? OR description2 ILIKE ? OR description3 ILIKE ?", 
		searchQuery, searchQuery, searchQuery, searchQuery).Count(&count).Error
	return count, err
}

// UpdateActivity updates activity information
func UpdateActivity(activity *Activity) error {
	return common.GetDB().Save(activity).Error
}

// DeleteActivity deletes activity by ID
func DeleteActivity(id string) error {
	return common.GetDB().Where("id = ?", id).Delete(&Activity{}).Error
}

// GetActivitiesByMonthYear returns activities whose time_start or day fall within the given month/year (UTC)
func GetActivitiesByMonthYear(year int, month int) ([]Activity, error) {
    startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

    var activities []Activity
    // Half-open interval [startOfMonth, startOfNextMonth)
    err := common.GetDB().Where(
        "(time_start IS NOT NULL AND time_start >= ? AND time_start < ?) OR (day IS NOT NULL AND day >= ? AND day < ?)",
        startOfMonth, startOfNextMonth, startOfMonth, startOfNextMonth,
    ).Order("created_at DESC").Find(&activities).Error
    return activities, err
}

// GetActivitiesByDay returns activities that match the specific day (UTC) by either time_start's date or day field
func GetActivitiesByDay(day time.Time) ([]Activity, error) {
    // Normalize to date (midnight UTC)
    start := time.Date(day.UTC().Year(), day.UTC().Month(), day.UTC().Day(), 0, 0, 0, 0, time.UTC)
    next := start.AddDate(0, 0, 1)

    var activities []Activity
    err := common.GetDB().Where(
        "(time_start IS NOT NULL AND time_start >= ? AND time_start < ?) OR (day IS NOT NULL AND day >= ? AND day < ?)",
        start, next, start, next,
    ).Order("created_at DESC").Find(&activities).Error
    return activities, err
}
