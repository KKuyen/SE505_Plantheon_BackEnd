package complaints

import (
	"plantheon-backend/models/diseases"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Disease type alias for cleaner code
type Disease = diseases.Disease

// ComplaintType represents the type of content being complained about
type ComplaintType string

const (
	ComplaintTypePost    ComplaintType = "POST"
	ComplaintTypeComment ComplaintType = "COMMENT"
	ComplaintTypeScan    ComplaintType = "SCAN"
)

// ComplaintCategory represents the category/reason of complaint
type ComplaintCategory string

const (
	ComplaintCategorySpam           ComplaintCategory = "SPAM"
	ComplaintCategoryHarassment     ComplaintCategory = "HARASSMENT"
	ComplaintCategoryHateSpeech     ComplaintCategory = "HATE_SPEECH"
	ComplaintCategoryViolence       ComplaintCategory = "VIOLENCE"
	ComplaintCategoryMisinformation ComplaintCategory = "MISINFORMATION"
	ComplaintCategoryInappropriate  ComplaintCategory = "INAPPROPRIATE"
	ComplaintCategoryOther          ComplaintCategory = "OTHER"
	
	// Scan-specific categories
	ComplaintCategoryScanWrongResult    ComplaintCategory = "WRONG_RESULT"
	ComplaintCategoryScanMisidentified  ComplaintCategory = "MISIDENTIFIED"
	ComplaintCategoryScanIncorrectInfo  ComplaintCategory = "INCORRECT_INFO"
	ComplaintCategoryScanPoorQuality    ComplaintCategory = "POOR_QUALITY"
	ComplaintCategoryScanOtherIssue     ComplaintCategory = "OTHER_ISSUE"
)

// ComplaintStatus represents the status of the complaint
type ComplaintStatus string

const (
	ComplaintStatusPending  ComplaintStatus = "PENDING"
	ComplaintStatusReviewed ComplaintStatus = "REVIEWED"
	ComplaintStatusResolved ComplaintStatus = "RESOLVED"
	ComplaintStatusRejected ComplaintStatus = "REJECTED"
)

// Complaint represents a user complaint about a post, comment, or scan result
type Complaint struct {
	ID          string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string          `json:"user_id" gorm:"not null;type:uuid;index"`            // User who made the complaint
	TargetID    string          `json:"target_id" gorm:"not null;type:uuid;index"`          // ID of the post, comment, or disease being complained about
	TargetType  ComplaintType   `json:"target_type" gorm:"type:varchar(20);not null;index"` // POST, COMMENT, or SCAN
	Category    ComplaintCategory `json:"category" gorm:"type:varchar(50);not null"`        // Reason category
	Content     string          `json:"content" gorm:"type:text"`                           // Additional details from user
	ImageURL    string          `json:"image_url" gorm:"type:text"`                         // Image URL for scan complaints
	Status      ComplaintStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`   // Complaint status
	AdminNotes  string          `json:"admin_notes" gorm:"type:text"`                       // Notes from admin after review
	ResolvedAt  *time.Time      `json:"resolved_at"`                                        // When the complaint was resolved
	ResolvedBy  *string         `json:"resolved_by" gorm:"type:uuid"`                       // Admin who resolved it
	
	// Scan-specific fields for ML training
	PredictedDiseaseID     *string    `json:"predicted_disease_id" gorm:"type:uuid;index"`     // Disease that AI predicted
	UserSuggestedDiseaseID *string    `json:"user_suggested_disease_id" gorm:"type:uuid;index"` // Disease that user thinks is correct (optional)
	VerifiedDiseaseID      *string    `json:"verified_disease_id" gorm:"type:uuid;index"`      // Ground truth from admin verification
	ConfidenceScore        *float64   `json:"confidence_score"`                                // AI confidence score (0.0-1.0)
	IsVerified             bool       `json:"is_verified" gorm:"default:false;index"`          // Has been verified by admin
	VerifiedBy             *string    `json:"verified_by" gorm:"type:uuid"`                    // Admin who verified (for scan complaints)
	VerifiedAt             *time.Time `json:"verified_at"`                                     // When verification happened
	
	// Disease relationships (for preloading)
	PredictedDisease       *Disease `json:"-" gorm:"foreignKey:PredictedDiseaseID;references:ID"`
	UserSuggestedDisease   *Disease `json:"-" gorm:"foreignKey:UserSuggestedDiseaseID;references:ID"`
	VerifiedDisease        *Disease `json:"-" gorm:"foreignKey:VerifiedDiseaseID;references:ID"`
	
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (c *Complaint) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Status == "" {
		c.Status = ComplaintStatusPending
	}
	return nil
}

// TableName returns the table name for Complaint
func (Complaint) TableName() string {
	return "complaints"
}

// ToComplaintResponse converts Complaint to ComplaintResponse
func (c *Complaint) ToComplaintResponse() ComplaintResponse {
	response := ComplaintResponse{
		ID:                     c.ID,
		UserID:                 c.UserID,
		TargetID:               c.TargetID,
		TargetType:             string(c.TargetType),
		Category:               string(c.Category),
		Content:                c.Content,
		ImageURL:               c.ImageURL,
		Status:                 string(c.Status),
		AdminNotes:             c.AdminNotes,
		ResolvedAt:             c.ResolvedAt,
		ResolvedBy:             c.ResolvedBy,
		PredictedDiseaseID:     c.PredictedDiseaseID,
		UserSuggestedDiseaseID: c.UserSuggestedDiseaseID,
		VerifiedDiseaseID:      c.VerifiedDiseaseID,
		ConfidenceScore:        c.ConfidenceScore,
		IsVerified:             c.IsVerified,
		VerifiedBy:             c.VerifiedBy,
		VerifiedAt:             c.VerifiedAt,
		CreatedAt:              c.CreatedAt,
		UpdatedAt:              c.UpdatedAt,
	}

	// Populate nested disease information if preloaded
	if c.PredictedDisease != nil {
		response.PredictedDisease = &DiseaseInfo{
			ID:        c.PredictedDisease.ID,
			Name:      c.PredictedDisease.Name,
			ClassName: c.PredictedDisease.ClassName,
			Type:      c.PredictedDisease.Type,
			PlantName: c.PredictedDisease.PlantName,
		}
	}

	if c.UserSuggestedDisease != nil {
		response.UserSuggestedDisease = &DiseaseInfo{
			ID:        c.UserSuggestedDisease.ID,
			Name:      c.UserSuggestedDisease.Name,
			ClassName: c.UserSuggestedDisease.ClassName,
			Type:      c.UserSuggestedDisease.Type,
			PlantName: c.UserSuggestedDisease.PlantName,
		}
	}

	if c.VerifiedDisease != nil {
		response.VerifiedDisease = &DiseaseInfo{
			ID:        c.VerifiedDisease.ID,
			Name:      c.VerifiedDisease.Name,
			ClassName: c.VerifiedDisease.ClassName,
			Type:      c.VerifiedDisease.Type,
			PlantName: c.VerifiedDisease.PlantName,
		}
	}

	return response
}

