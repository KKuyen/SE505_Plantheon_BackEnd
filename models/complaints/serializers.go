package complaints

import (
	"time"
)

// DiseaseInfo represents basic disease information for nested responses
type DiseaseInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
	PlantName string `json:"plant_name"`
}

// ComplaintResponse represents the response for a complaint
type ComplaintResponse struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	TargetID   string     `json:"target_id"`
	TargetType string     `json:"target_type"`
	Category   string     `json:"category"`
	Content    string     `json:"content"`
	ImageURL   string     `json:"image_url,omitempty"`
	Status     string     `json:"status"`
	AdminNotes string     `json:"admin_notes,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *string    `json:"resolved_by,omitempty"`
	
	// Scan-specific fields
	PredictedDiseaseID     *string      `json:"predicted_disease_id,omitempty"`
	PredictedDisease       *DiseaseInfo `json:"predicted_disease,omitempty"`
	UserSuggestedDiseaseID *string      `json:"user_suggested_disease_id,omitempty"`
	UserSuggestedDisease   *DiseaseInfo `json:"user_suggested_disease,omitempty"`
	VerifiedDiseaseID      *string      `json:"verified_disease_id,omitempty"`
	VerifiedDisease        *DiseaseInfo `json:"verified_disease,omitempty"`
	ConfidenceScore        *float64     `json:"confidence_score,omitempty"`
	IsVerified             bool         `json:"is_verified"`
	VerifiedBy             *string      `json:"verified_by,omitempty"` // Admin who verified scan complaint
	VerifiedAt             *time.Time   `json:"verified_at,omitempty"`
	
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TargetInfo represents information about the reported post/comment
type TargetInfo struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Avatar    string    `json:"avatar"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
}

// ReporterInfo represents information about the user who reported
type ReporterInfo struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}

// ComplaintWithTargetResponse represents complaint with target details
type ComplaintWithTargetResponse struct {
	ID         string        `json:"id"`
	Reporter   ReporterInfo  `json:"reporter"`
	Target     TargetInfo    `json:"target"`
	TargetType string        `json:"target_type"`
	Category   string        `json:"category"`
	Content    string        `json:"content"`
	Status     string        `json:"status"`
	AdminNotes string        `json:"admin_notes,omitempty"`
	ResolvedAt *time.Time    `json:"resolved_at,omitempty"`
	ResolvedBy *string       `json:"resolved_by,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// ComplaintWithTargetListResponse represents a list of complaints with target info
type ComplaintWithTargetListResponse struct {
	Complaints []ComplaintWithTargetResponse `json:"complaints"`
	Total      int64                         `json:"total"`
	Page       int                           `json:"page"`
	Limit      int                           `json:"limit"`
}

// ComplaintListResponse represents a list of complaints
type ComplaintListResponse struct {
	Complaints []ComplaintResponse `json:"complaints"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}

// CreateComplaintRequest represents the request to create a complaint
type CreateComplaintRequest struct {
	TargetID   string `json:"target_id" binding:"required"`
	TargetType string `json:"target_type" binding:"required"` // POST or COMMENT
	Category   string `json:"category" binding:"required"`    // SPAM, HARASSMENT, etc.
	Content    string `json:"content"`                        // Optional additional details
}

// CreateScanComplaintRequest represents the request to create a scan complaint
type CreateScanComplaintRequest struct {
	TargetType             string   `json:"target_type" binding:"required"`              // Must be "SCAN"
	PredictedDiseaseID     string   `json:"predicted_disease_id" binding:"required"`     // Disease that AI predicted
	UserSuggestedDiseaseID *string  `json:"user_suggested_disease_id"`                   // Disease that user thinks is correct (optional)
	ConfidenceScore        float64  `json:"confidence_score" binding:"required,min=0,max=1"` // AI confidence (0.0-1.0)
	Category               string   `json:"category" binding:"required"`                 // WRONG_RESULT, MISIDENTIFIED, etc.
	ImageURL               string   `json:"image_url" binding:"required"`                // Image URL (required for scans)
	Content                string   `json:"content"`                                     // Optional additional details
}

// VerifyComplaintRequest represents the request to verify a scan complaint (admin only)
type VerifyComplaintRequest struct {
	VerifiedDiseaseID string `json:"verified_disease_id" binding:"required"` // Ground truth disease ID
	IsVerified        bool   `json:"is_verified"`                            // Verification status
	AdminNotes        string `json:"admin_notes"`                            // Optional admin notes
}

// TrainingDataExport represents exported data for ML training
type TrainingDataExport struct {
	ImageURL              string    `json:"image_url"`
	PredictedDiseaseID    string    `json:"predicted_disease_id"`
	PredictedClassName    string    `json:"predicted_class_name"`
	VerifiedDiseaseID     string    `json:"verified_disease_id"`
	VerifiedClassName     string    `json:"verified_class_name"`
	ConfidenceScore       float64   `json:"confidence_score"`
	CreatedAt             time.Time `json:"created_at"`
}

// UpdateComplaintStatusRequest represents the request to update complaint status (admin only)
type UpdateComplaintStatusRequest struct {
	Status     string `json:"status" binding:"required"` // PENDING, REVIEWED, RESOLVED, REJECTED
	AdminNotes string `json:"admin_notes"`               // Optional notes from admin
}
