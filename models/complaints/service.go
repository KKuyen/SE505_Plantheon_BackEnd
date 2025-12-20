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

// GetTargetInfo retrieves information about the target (post or comment)
func GetTargetInfo(targetID string, targetType ComplaintType) TargetInfo {
	service := NewComplaintService()
	info := TargetInfo{ID: targetID}

	if targetType == ComplaintTypePost {
		var post struct {
			ID        string
			Content   string
			UserID    string
			IsDeleted bool
			CreatedAt time.Time
		}
		if err := service.db.Table("posts").
			Select("id, content, user_id, is_deleted, created_at").
			Where("id = ?", targetID).First(&post).Error; err == nil {
			info.Content = post.Content
			info.UserID = post.UserID
			info.IsDeleted = post.IsDeleted
			info.CreatedAt = post.CreatedAt

			// Get user info
			var user struct {
				FullName string
				Avatar   string
			}
			if err := service.db.Table("users").
				Select("full_name, avatar").
				Where("id = ?", post.UserID).First(&user).Error; err == nil {
				info.UserName = user.FullName
				info.Avatar = user.Avatar
			}
		}
	} else if targetType == ComplaintTypeComment {
		var comment struct {
			ID        string
			Content   string
			UserID    string
			IsDeleted bool
			CreatedAt time.Time
		}
		if err := service.db.Table("comments").
			Select("id, content, user_id, is_deleted, created_at").
			Where("id = ?", targetID).First(&comment).Error; err == nil {
			info.Content = comment.Content
			info.UserID = comment.UserID
			info.IsDeleted = comment.IsDeleted
			info.CreatedAt = comment.CreatedAt

			// Get user info
			var user struct {
				FullName string
				Avatar   string
			}
			if err := service.db.Table("users").
				Select("full_name, avatar").
				Where("id = ?", comment.UserID).First(&user).Error; err == nil {
				info.UserName = user.FullName
				info.Avatar = user.Avatar
			}
		}
	}

	return info
}

// GetReporterInfo retrieves information about the reporter
func GetReporterInfo(userID string) ReporterInfo {
	service := NewComplaintService()
	info := ReporterInfo{ID: userID}

	var user struct {
		FullName string
		Avatar   string
	}
	if err := service.db.Table("users").
		Select("full_name, avatar").
		Where("id = ?", userID).First(&user).Error; err == nil {
		info.UserName = user.FullName
		info.Avatar = user.Avatar
	}

	return info
}

// ToComplaintWithTargetResponse converts Complaint to ComplaintWithTargetResponse
func (c *Complaint) ToComplaintWithTargetResponse() ComplaintWithTargetResponse {
	return ComplaintWithTargetResponse{
		ID:         c.ID,
		Reporter:   GetReporterInfo(c.UserID),
		Target:     GetTargetInfo(c.TargetID, c.TargetType),
		TargetType: string(c.TargetType),
		Category:   string(c.Category),
		Content:    c.Content,
		Status:     string(c.Status),
		AdminNotes: c.AdminNotes,
		ResolvedAt: c.ResolvedAt,
		ResolvedBy: c.ResolvedBy,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
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
		// If targetType is specified, filter by it (including SCAN if specified)
		query = query.Where("target_type = ?", targetType)
	} else {
		// If no targetType specified, filter out SCAN type complaints by default
		query = query.Where("target_type != ?", ComplaintTypeScan)
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

// GetComplaintsAboutMyContent gets all complaints about posts/comments owned by the user
func GetComplaintsAboutMyContent(userID string, status string, targetType string) ([]Complaint, error) {
	service := NewComplaintService()
	var complaints []Complaint

	// Get post IDs owned by user
	var postIDs []string
	service.db.Table("posts").Select("id").Where("user_id = ?", userID).Pluck("id", &postIDs)

	// Get comment IDs owned by user
	var commentIDs []string
	service.db.Table("comments").Select("id").Where("user_id = ?", userID).Pluck("id", &commentIDs)

	// If user has no posts and no comments, return empty
	if len(postIDs) == 0 && len(commentIDs) == 0 {
		return []Complaint{}, nil
	}

	// Build query for complaints about user's content
	query := service.db.Model(&Complaint{})

	// Build OR condition for posts and comments
	if len(postIDs) > 0 && len(commentIDs) > 0 {
		query = query.Where(
			"(target_id IN ? AND target_type = ?) OR (target_id IN ? AND target_type = ?)",
			postIDs, ComplaintTypePost, commentIDs, ComplaintTypeComment,
		)
	} else if len(postIDs) > 0 {
		query = query.Where("target_id IN ? AND target_type = ?", postIDs, ComplaintTypePost)
	} else {
		query = query.Where("target_id IN ? AND target_type = ?", commentIDs, ComplaintTypeComment)
	}

	// Apply optional filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	// Get all results
	err := query.Order("created_at DESC").Find(&complaints).Error
	return complaints, err
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
