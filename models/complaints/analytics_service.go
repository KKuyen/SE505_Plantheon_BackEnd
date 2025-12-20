package complaints

import (
	"plantheon-backend/common"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"column:id"`
	Username string `gorm:"column:username"`
	FullName string `gorm:"column:full_name"`
	Avatar   string `gorm:"column:avatar"`
}

func (User) TableName() string {
	return "users"
}

// AnalyticsService handles analytics operations for complaints
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService creates a new analytics service instance
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		db: common.GetDB(),
	}
}

// DiseaseBasicInfo represents basic disease information for analytics
type DiseaseBasicInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ClassName string `json:"class_name"`
	ImageLink string `json:"image_link"`
	PlantName string `json:"plant_name"`
}

// UserBasicInfo represents basic user information for analytics
type UserBasicInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}

// ProblematicDisease represents a disease with many complaints
type ProblematicDisease struct {
	Disease        *DiseaseBasicInfo `json:"disease"`
	ComplaintCount int64             `json:"complaint_count"`
	VerifiedCount  int64             `json:"verified_count"`
	AvgConfidence  float64           `json:"avg_confidence"`
	ErrorRate      float64           `json:"error_rate"` // % times AI was wrong
}

// TrendData represents daily complaint trends
type TrendData struct {
	Date            string  `json:"date"`
	ComplaintCount  int64   `json:"complaint_count"`
	VerifiedCount   int64   `json:"verified_count"`
	AvgConfidence   float64 `json:"avg_confidence"`
}

// OverallStats represents overall system statistics
type OverallStats struct {
	TotalComplaints     int64   `json:"total_complaints"`
	PendingComplaints   int64   `json:"pending_complaints"`
	VerifiedComplaints  int64   `json:"verified_complaints"`
	AICorrectRate       float64 `json:"ai_correct_rate"` // % times AI was actually right
}

// ContributorStats represents user contribution statistics
type ContributorStats struct {
	User           *UserBasicInfo `json:"user"`
	ComplaintCount int64          `json:"complaint_count"`
	VerifiedCount  int64          `json:"verified_count"`
	CorrectCount   int64          `json:"correct_count"` // Times user suggestion was correct
}

// GetMostProblematicDiseases returns diseases with most complaints
func GetMostProblematicDiseases(limit int) ([]ProblematicDisease, error) {
	service := NewAnalyticsService()
	
	if limit <= 0 {
		limit = 10
	}

	// First get the stats
	type tempResult struct {
		DiseaseID      string
		ComplaintCount int64
		VerifiedCount  int64
		AvgConfidence  float64
		ErrorRate      float64
	}
	
	var tempResults []tempResult
	err := service.db.Model(&Complaint{}).
		Select(`
			predicted_disease_id as disease_id,
			COUNT(*) as complaint_count,
			SUM(CASE WHEN is_verified = true THEN 1 ELSE 0 END) as verified_count,
			AVG(confidence_score) as avg_confidence,
			(SUM(CASE WHEN is_verified = true AND predicted_disease_id != verified_disease_id THEN 1 ELSE 0 END)::float / 
			 NULLIF(SUM(CASE WHEN is_verified = true THEN 1 ELSE 0 END), 0) * 100) as error_rate
		`).
		Where("target_type = ? AND predicted_disease_id IS NOT NULL", ComplaintTypeScan).
		Group("predicted_disease_id").
		Order("complaint_count DESC").
		Limit(limit).
		Scan(&tempResults).Error
	
	if err != nil {
		return nil, err
	}

	// Now populate disease info
	results := make([]ProblematicDisease, 0, len(tempResults))
	for _, temp := range tempResults {
		var disease Disease
		if err := service.db.Where("id = ?", temp.DiseaseID).First(&disease).Error; err == nil {
			// Get first image from array
			var imageLink string
			if len(disease.ImageLink) > 0 {
				imageLink = disease.ImageLink[0]
			}
			
			results = append(results, ProblematicDisease{
				Disease: &DiseaseBasicInfo{
					ID:        disease.ID,
					Name:      disease.Name,
					ClassName: disease.ClassName,
					ImageLink: imageLink,
					PlantName: disease.PlantName,
				},
				ComplaintCount: temp.ComplaintCount,
				VerifiedCount:  temp.VerifiedCount,
				AvgConfidence:  temp.AvgConfidence,
				ErrorRate:      temp.ErrorRate,
			})
		}
	}

	return results, nil
}

// GetComplaintTrends returns daily complaint and verification trends
func GetComplaintTrends(days int) ([]TrendData, error) {
	service := NewAnalyticsService()
	var results []TrendData

	if days <= 0 {
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)

	err := service.db.Model(&Complaint{}).
		Select(`
			DATE(created_at) as date,
			COUNT(*) as complaint_count,
			SUM(CASE WHEN is_verified = true THEN 1 ELSE 0 END) as verified_count,
			AVG(confidence_score) as avg_confidence
		`).
		Where("target_type = ? AND created_at >= ?", ComplaintTypeScan, startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error

	return results, err
}

// GetOverallStats returns overall system statistics
func GetOverallStats() (*OverallStats, error) {
	service := NewAnalyticsService()
	var stats OverallStats

	// Get total complaints
	err := service.db.Model(&Complaint{}).
		Where("target_type = ?", ComplaintTypeScan).
		Count(&stats.TotalComplaints).Error
	if err != nil {
		return nil, err
	}

	// Get pending complaints (unverified)
	err = service.db.Model(&Complaint{}).
		Where("target_type = ? AND is_verified = ?", ComplaintTypeScan, false).
		Count(&stats.PendingComplaints).Error
	if err != nil {
		return nil, err
	}

	// Get verified complaints
	err = service.db.Model(&Complaint{}).
		Where("target_type = ? AND is_verified = ?", ComplaintTypeScan, true).
		Count(&stats.VerifiedComplaints).Error
	if err != nil {
		return nil, err
	}

	// Get AI correct rate
	var correctRate struct {
		Total   int64
		Correct int64
	}
	err = service.db.Model(&Complaint{}).
		Select(`
			COUNT(*) as total,
			SUM(CASE WHEN predicted_disease_id = verified_disease_id THEN 1 ELSE 0 END) as correct
		`).
		Where("target_type = ? AND is_verified = ? AND predicted_disease_id IS NOT NULL AND verified_disease_id IS NOT NULL",
			ComplaintTypeScan, true).
		Scan(&correctRate).Error
	if err != nil {
		return nil, err
	}

	if correctRate.Total > 0 {
		stats.AICorrectRate = float64(correctRate.Correct) / float64(correctRate.Total) * 100
	}

	return &stats, nil
}

// GetTopContributors returns users who contributed most verified complaints
func GetTopContributors(limit int) ([]ContributorStats, error) {
	service := NewAnalyticsService()
	
	if limit <= 0 {
		limit = 10
	}

	// First get the stats
	type tempResult struct {
		UserID         string
		ComplaintCount int64
		VerifiedCount  int64
		CorrectCount   int64
	}
	
	var tempResults []tempResult
	err := service.db.Model(&Complaint{}).
		Select(`
			user_id,
			COUNT(*) as complaint_count,
			SUM(CASE WHEN is_verified = true THEN 1 ELSE 0 END) as verified_count,
			SUM(CASE WHEN is_verified = true AND user_suggested_disease_id = verified_disease_id THEN 1 ELSE 0 END) as correct_count
		`).
		Where("target_type = ?", ComplaintTypeScan).
		Group("user_id").
		Order("verified_count DESC").
		Limit(limit).
		Scan(&tempResults).Error

	if err != nil {
		return nil, err
	}

	// Now populate user info
	results := make([]ContributorStats, 0, len(tempResults))
	for _, temp := range tempResults {
		var user User
		if err := service.db.Where("id = ?", temp.UserID).First(&user).Error; err == nil {
			results = append(results, ContributorStats{
				User: &UserBasicInfo{
					ID:       user.ID,
					Username: user.Username,
					FullName: user.FullName,
					Avatar:   user.Avatar,
				},
				ComplaintCount: temp.ComplaintCount,
				VerifiedCount:  temp.VerifiedCount,
				CorrectCount:   temp.CorrectCount,
			})
		}
	}

	return results, nil
}
