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
	err := service.db.Where("id = ?", id).First(&complaint).Error
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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&complaints).Error
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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&complaints).Error
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
func GetComplaintsCount(status string, targetType string) (int64, error) {
	service := NewComplaintService()
	var count int64

	query := service.db.Model(&Complaint{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	err := query.Count(&count).Error
	return count, err
}
