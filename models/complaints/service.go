package complaints

import (
	"plantheon-backend/common"
	"time"

	"gorm.io/gorm"
)

// ComplaintService handles all database operations for complaints
type ComplaintService struct {
	db *gorm.DB
}

// NewComplaintService creates a new complaint service instance
func NewComplaintService() *ComplaintService {
	return &ComplaintService{
		db: common.GetDB(),
	}
}

// CreateComplaint creates a new complaint
func CreateComplaint(complaint *Complaint) error {
	service := NewComplaintService()
	return service.db.Create(complaint).Error
}

// GetComplaintByID gets a complaint by ID
func GetComplaintByID(id string) (*Complaint, error) {
	service := NewComplaintService()
	var complaint Complaint
	err := service.db.
		Preload("PredictedDisease").
		Preload("UserSuggestedDisease").
		Preload("VerifiedDisease").
		Where("id = ?", id).First(&complaint).Error
	return &complaint, err
}

// GetComplaintsByUserID gets all complaints made by a user
func GetComplaintsByUserID(userID string, offset, limit int) ([]Complaint, int64, error) {
	service := NewComplaintService()
	var complaints []Complaint
	var total int64

	query := service.db.Where("user_id = ?", userID)

	if err := query.Model(&Complaint{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("PredictedDisease").
		Preload("UserSuggestedDisease").
		Preload("VerifiedDisease").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&complaints).Error
	return complaints, total, err
}

// GetComplaintsByTargetID gets all complaints for a specific target (post or comment)
func GetComplaintsByTargetID(targetID string, targetType ComplaintType) ([]Complaint, error) {
	service := NewComplaintService()
	var complaints []Complaint
	err := service.db.Where("target_id = ? AND target_type = ?", targetID, targetType).
		Order("created_at DESC").Find(&complaints).Error
	return complaints, err
}

// GetAllComplaints gets all complaints with pagination and optional filters (admin)
func GetAllComplaints(offset, limit int, status string, targetType string) ([]Complaint, int64, error) {
	service := NewComplaintService()
	var complaints []Complaint
	var total int64

	query := service.db.Model(&Complaint{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("PredictedDisease").
		Preload("UserSuggestedDisease").
		Preload("VerifiedDisease").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&complaints).Error
	return complaints, total, err
}

// GetPendingComplaints gets all pending complaints (admin)
func GetPendingComplaints(offset, limit int) ([]Complaint, int64, error) {
	return GetAllComplaints(offset, limit, string(ComplaintStatusPending), "")
}

// UpdateComplaintStatus updates the status of a complaint (admin)
func UpdateComplaintStatus(id string, status ComplaintStatus, adminNotes string, adminID string) error {
	service := NewComplaintService()
	
	updates := map[string]interface{}{
		"status":      status,
		"admin_notes": adminNotes,
		"updated_at":  time.Now(),
	}

	if status == ComplaintStatusResolved || status == ComplaintStatusRejected {
		now := time.Now()
		updates["resolved_at"] = &now
		updates["resolved_by"] = &adminID
	}

	return service.db.Model(&Complaint{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteComplaint deletes a complaint by ID
func DeleteComplaint(id string) error {
	service := NewComplaintService()
	return service.db.Delete(&Complaint{}, "id = ?", id).Error
}

// CheckDuplicateComplaint checks if user already complained about the same target
func CheckDuplicateComplaint(userID, targetID string, targetType ComplaintType) (bool, error) {
	service := NewComplaintService()
	var count int64
	err := service.db.Model(&Complaint{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	return count > 0, err
}

// GetComplaintsCount gets the count of complaints with optional filters
func GetComplaintsCount(status string, targetType string, isVerified *bool) (int64, error) {
	service := NewComplaintService()
	var count int64

	query := service.db.Model(&Complaint{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}
	if isVerified != nil {
		query = query.Where("is_verified = ?", *isVerified)
	}

	err := query.Count(&count).Error
	return count, err
}

// VerifyComplaint verifies a scan complaint and sets ground truth
func VerifyComplaint(id, verifiedDiseaseID string, isVerified bool, adminNotes, adminID string) error {
	service := NewComplaintService()
	
	now := time.Now()
	updates := map[string]interface{}{
		"verified_disease_id": verifiedDiseaseID,
		"is_verified":         isVerified,
		"verified_by":         &adminID,
		"verified_at":         &now,
		"admin_notes":         adminNotes,
		"updated_at":          now,
	}

	return service.db.Model(&Complaint{}).Where("id = ?", id).Updates(updates).Error
}

// GetUnverifiedScanComplaints gets scan complaints that haven't been verified
func GetUnverifiedScanComplaints(offset, limit int) ([]Complaint, int64, error) {
	service := NewComplaintService()
	var complaints []Complaint
	var total int64

	query := service.db.Model(&Complaint{}).
		Where("target_type = ? AND is_verified = ?", ComplaintTypeScan, false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("PredictedDisease").
		Preload("UserSuggestedDisease").
		Preload("VerifiedDisease").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&complaints).Error
	return complaints, total, err
}

// GetVerifiedScanComplaints gets all verified scan complaints for ML export
func GetVerifiedScanComplaints() ([]Complaint, error) {
	service := NewComplaintService()
	var complaints []Complaint
	
	err := service.db.
		Preload("PredictedDisease").
		Preload("VerifiedDisease").
		Where("target_type = ? AND is_verified = ?", ComplaintTypeScan, true).
		Order("created_at DESC").
		Find(&complaints).Error
	
	return complaints, err
}

// ExportTrainingData exports verified complaints in ML training format
func ExportTrainingData() ([]TrainingDataExport, error) {
	complaints, err := GetVerifiedScanComplaints()
	if err != nil {
		return nil, err
	}

	exports := make([]TrainingDataExport, 0, len(complaints))
	for _, complaint := range complaints {
		// Only export if we have all required fields
		if complaint.PredictedDiseaseID != nil && 
		   complaint.VerifiedDiseaseID != nil && 
		   complaint.ConfidenceScore != nil &&
		   complaint.ImageURL != "" &&
		   complaint.PredictedDisease != nil &&
		   complaint.VerifiedDisease != nil {
			exports = append(exports, TrainingDataExport{
				ImageURL:              complaint.ImageURL,
				PredictedDiseaseID:    *complaint.PredictedDiseaseID,
				PredictedClassName:    complaint.PredictedDisease.ClassName,
				VerifiedDiseaseID:     *complaint.VerifiedDiseaseID,
				VerifiedClassName:     complaint.VerifiedDisease.ClassName,
				ConfidenceScore:       *complaint.ConfidenceScore,
				CreatedAt:             complaint.CreatedAt,
			})
		}
	}

	return exports, nil
}
