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

// CheckActivityRepeatsOnDate checks if a repeating activity should appear on the given date
func CheckActivityRepeatsOnDate(activity *Activity, targetDate time.Time) bool {
	// If repeat is null or empty, it's not a repeating activity
	if activity.Repeat == nil || *activity.Repeat == "" {
		return false
	}

	// If time_start is null, we can't determine the pattern
	if activity.TimeStart == nil {
		return false
	}

	activityStart := *activity.TimeStart
	repeatType := *activity.Repeat

	// Check if target date is before the activity start date
	activityStartDay := time.Date(activityStart.Year(), activityStart.Month(), activityStart.Day(), 0, 0, 0, 0, time.UTC)
	targetDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	
	if targetDay.Before(activityStartDay) {
		return false
	}

	// If is_repeat is not null, check end_repeat_day
	if activity.IsRepeat != nil && *activity.IsRepeat != "" {
		if activity.EndRepeatDay != nil {
			endRepeatDay := time.Date(activity.EndRepeatDay.Year(), activity.EndRepeatDay.Month(), activity.EndRepeatDay.Day(), 0, 0, 0, 0, time.UTC)
			if targetDay.After(endRepeatDay) {
				return false
			}
		}
	}

	// Check repeat pattern
	switch repeatType {
	case "Hàng ngày":
		return true

	case "Hàng tuần":
		// Check if same day of week
		return activityStart.Weekday() == targetDate.Weekday()

	case "Hàng tháng":
		// Check if same day of month
		return activityStart.Day() == targetDate.Day()

	case "Hàng năm":
		// Check if same day and month
		return activityStart.Day() == targetDate.Day() && activityStart.Month() == targetDate.Month()

	default:
		return false
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

// DeleteActivities deletes multiple activities by IDs and returns number of rows affected
func DeleteActivities(ids []string) (int64, error) {
	service := NewActivityService()
	tx := service.db.Where("id IN ?", ids).Delete(&Activity{})
	return tx.RowsAffected, tx.Error
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

// GetAllRecurringActivities gets all activities where repeat is not null
func GetAllRecurringActivities(activities *[]Activity) error {
	service := NewActivityService()
	err := service.db.Where(
		"repeat IS NOT NULL AND repeat != ''",
	).Order("time_start ASC").Find(activities).Error
	return err
}

// GetMonthlyFinancialSummary calculates total income and expense for a specific month
// Only includes activities where money is not null
func GetMonthlyFinancialSummary(year int, month int) (*MonthlyFinancialSummary, error) {
	service := NewActivityService()
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	var activities []Activity
	err := service.db.Where(
		"money IS NOT NULL AND time_start IS NOT NULL AND time_start >= ? AND time_start < ?",
		startOfMonth, startOfNextMonth,
	).Order("time_start ASC").Find(&activities).Error

	if err != nil {
		return nil, err
	}

	summary := &MonthlyFinancialSummary{
		Year:         year,
		Month:        month,
		TotalIncome:  0,
		TotalExpense: 0,
		Activities:   make([]ActivityResponse, 0),
	}

	for _, activity := range activities {
		// Convert money to negative if type is EXPENSE
		activityResponse := activity.ToActivityResponse()
		if activity.Money != nil && activity.Type == "EXPENSE" {
			negativeMoney := -*activity.Money
			if negativeMoney > 0 {
				negativeMoney = -negativeMoney
			}
			activityResponse.Money = &negativeMoney
		}
		
		// Add to activities list with corrected money
		summary.Activities = append(summary.Activities, activityResponse)
		
		// Calculate totals
		if activity.Money != nil {
			moneyValue := *activity.Money
			// Convert to negative if EXPENSE
			if activity.Type == "EXPENSE" && moneyValue > 0 {
				moneyValue = -moneyValue
			}
			
			if moneyValue > 0 {
				summary.TotalIncome += moneyValue
			} else if moneyValue < 0 {
				summary.TotalExpense += moneyValue
			}
		}
	}

	summary.NetAmount = summary.TotalIncome + summary.TotalExpense // Expense đã âm rồi nên cộng

	return summary, nil
}

// GetAnnualFinancialSummary calculates financial summary for all 12 months of a year
// Optimized: Query once for the entire year instead of 12 separate queries
func GetAnnualFinancialSummary(year int) (*AnnualFinancialSummary, error) {
	service := NewActivityService()
	startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	startOfNextYear := startOfYear.AddDate(1, 0, 0)

	// Query all activities for the entire year at once
	var activities []Activity
	err := service.db.Where(
		"money IS NOT NULL AND time_start IS NOT NULL AND time_start >= ? AND time_start < ?",
		startOfYear, startOfNextYear,
	).Find(&activities).Error

	if err != nil {
		return nil, err
	}

	annualSummary := &AnnualFinancialSummary{
		Year:             year,
		MonthlySummaries: make([]MonthlyFinancialSummaryWithoutActivities, 12),
		TotalIncome:      0,
		TotalExpense:     0,
	}

	// Initialize all 12 months
	for month := 1; month <= 12; month++ {
		annualSummary.MonthlySummaries[month-1] = MonthlyFinancialSummaryWithoutActivities{
			Year:         year,
			Month:        month,
			TotalIncome:  0,
			TotalExpense: 0,
			NetAmount:    0,
		}
	}

	// Group activities by month and calculate totals
	for _, activity := range activities {
		if activity.Money != nil && activity.TimeStart != nil {
			month := int(activity.TimeStart.Month())
			monthIndex := month - 1

			moneyValue := *activity.Money
			// Convert to negative if EXPENSE
			if activity.Type == "EXPENSE" && moneyValue > 0 {
				moneyValue = -moneyValue
			}

			if moneyValue > 0 {
				annualSummary.MonthlySummaries[monthIndex].TotalIncome += moneyValue
				annualSummary.TotalIncome += moneyValue
			} else if moneyValue < 0 {
				annualSummary.MonthlySummaries[monthIndex].TotalExpense += moneyValue
				annualSummary.TotalExpense += moneyValue
			}

			annualSummary.MonthlySummaries[monthIndex].NetAmount += moneyValue
		}
	}

	annualSummary.NetAmount = annualSummary.TotalIncome + annualSummary.TotalExpense

	return annualSummary, nil
}

// GetMultiYearFinancialSummary calculates financial summary for multiple years (range)
// Optimized: Query once for all years instead of separate queries per year
func GetMultiYearFinancialSummary(startYear, endYear int) (*MultiYearFinancialSummary, error) {
	if startYear > endYear {
		startYear, endYear = endYear, startYear // Swap if start > end
	}

	service := NewActivityService()
	startOfPeriod := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)
	endOfPeriod := time.Date(endYear+1, 1, 1, 0, 0, 0, 0, time.UTC)

	// Query all activities for the entire period at once
	var activities []Activity
	err := service.db.Where(
		"money IS NOT NULL AND time_start IS NOT NULL AND time_start >= ? AND time_start < ?",
		startOfPeriod, endOfPeriod,
	).Find(&activities).Error

	if err != nil {
		return nil, err
	}

	multiYearSummary := &MultiYearFinancialSummary{
		StartYear:       startYear,
		EndYear:         endYear,
		YearlySummaries: make([]YearlyFinancialSummary, 0),
		TotalIncome:     0,
		TotalExpense:    0,
		NetAmount:       0,
	}

	// Create a map for each year
	yearMap := make(map[int]*YearlyFinancialSummary)
	for year := startYear; year <= endYear; year++ {
		yearMap[year] = &YearlyFinancialSummary{
			Year:         year,
			TotalIncome:  0,
			TotalExpense: 0,
			NetAmount:    0,
		}
	}

	// Group activities by year and calculate totals
	for _, activity := range activities {
		if activity.Money != nil && activity.TimeStart != nil {
			year := activity.TimeStart.Year()
			if yearSummary, exists := yearMap[year]; exists {
				moneyValue := *activity.Money
				// Convert to negative if EXPENSE
				if activity.Type == "EXPENSE" && moneyValue > 0 {
					moneyValue = -moneyValue
				}

				if moneyValue > 0 {
					yearSummary.TotalIncome += moneyValue
				} else if moneyValue < 0 {
					yearSummary.TotalExpense += moneyValue
				}

				yearSummary.NetAmount += moneyValue
			}
		}
	}

	// Convert map to sorted slice and calculate grand totals
	for year := startYear; year <= endYear; year++ {
		yearSummary := yearMap[year]
		multiYearSummary.YearlySummaries = append(multiYearSummary.YearlySummaries, *yearSummary)
		multiYearSummary.TotalIncome += yearSummary.TotalIncome
		multiYearSummary.TotalExpense += yearSummary.TotalExpense
	}

	multiYearSummary.NetAmount = multiYearSummary.TotalIncome + multiYearSummary.TotalExpense

	return multiYearSummary, nil
}

// GetActivitiesByDateRange gets activities where time_start <= selectedDate and time_end >= selectedDate
// Also includes recurring activities that repeat on the selected date
func GetActivitiesByDateRange(selectedDate time.Time) ([]Activity, error) {
	service := NewActivityService()
	var regularActivities []Activity
	var allRecurringActivities []Activity
	
	// Normalize selected date to start and end of day (UTC)
	startOfDay := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 23, 59, 59, 999999999, time.UTC)
	
	// Query 1: Get regular activities that overlap with the selected day
	err := service.db.Where(
		"time_start IS NOT NULL AND time_end IS NOT NULL AND time_start <= ? AND time_end >= ?",
		endOfDay, startOfDay,
	).Order("time_start ASC").Find(&regularActivities).Error
	
	if err != nil {
		return nil, err
	}
	
	// Query 2: Get all recurring activities (where repeat is not null)
	err = service.db.Where(
		"repeat IS NOT NULL AND repeat != ''",
	).Order("time_start ASC").Find(&allRecurringActivities).Error
	
	if err != nil {
		return nil, err
	}
	
	// Filter recurring activities that should appear on the selected date
	var matchingRecurringActivities []Activity
	for _, activity := range allRecurringActivities {
		if CheckActivityRepeatsOnDate(&activity, selectedDate) {
			matchingRecurringActivities = append(matchingRecurringActivities, activity)
		}
	}
	
	// Combine results (avoid duplicates by using a map)
	activityMap := make(map[string]Activity)
	
	for _, activity := range regularActivities {
		activityMap[activity.ID] = activity
	}
	
	for _, activity := range matchingRecurringActivities {
		activityMap[activity.ID] = activity
	}
	
	// Convert map back to slice
	var finalActivities []Activity
	for _, activity := range activityMap {
		finalActivities = append(finalActivities, activity)
	}
	
	return finalActivities, nil
}

// GetActivitiesByDateRangeWithPagination gets activities where time_start <= selectedDate and time_end >= selectedDate with pagination
// Also includes recurring activities that repeat on the selected date
func GetActivitiesByDateRangeWithPagination(selectedDate time.Time, offset, limit int) ([]Activity, int64, error) {
	// Get all matching activities (regular + recurring)
	allActivities, err := GetActivitiesByDateRange(selectedDate)
	if err != nil {
		return nil, 0, err
	}
	
	total := int64(len(allActivities))
	
	// Apply pagination manually
	start := offset
	end := offset + limit
	
	if start > len(allActivities) {
		return []Activity{}, total, nil
	}
	
	if end > len(allActivities) {
		end = len(allActivities)
	}
	
	paginatedActivities := allActivities[start:end]
	
	return paginatedActivities, total, nil
}
