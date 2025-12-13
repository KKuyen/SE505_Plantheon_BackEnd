package complaints

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ComplaintType represents the type of content being complained about
type ComplaintType string

const (
	ComplaintTypePost    ComplaintType = "POST"
	ComplaintTypeComment ComplaintType = "COMMENT"
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
)

// ComplaintStatus represents the status of the complaint
type ComplaintStatus string

const (
	ComplaintStatusPending  ComplaintStatus = "PENDING"
	ComplaintStatusReviewed ComplaintStatus = "REVIEWED"
	ComplaintStatusResolved ComplaintStatus = "RESOLVED"
	ComplaintStatusRejected ComplaintStatus = "REJECTED"
)

// Complaint represents a user complaint about a post or comment
type Complaint struct {
	ID          string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string          `json:"user_id" gorm:"not null;type:uuid;index"`            // User who made the complaint
	TargetID    string          `json:"target_id" gorm:"not null;type:uuid;index"`          // ID of the post or comment being complained about
	TargetType  ComplaintType   `json:"target_type" gorm:"type:varchar(20);not null;index"` // POST or COMMENT
	Category    ComplaintCategory `json:"category" gorm:"type:varchar(50);not null"`        // Reason category
	Content     string          `json:"content" gorm:"type:text"`                           // Additional details from user
	Status      ComplaintStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`   // Complaint status
	AdminNotes  string          `json:"admin_notes" gorm:"type:text"`                       // Notes from admin after review
	ResolvedAt  *time.Time      `json:"resolved_at"`                                        // When the complaint was resolved
	ResolvedBy  *string         `json:"resolved_by" gorm:"type:uuid"`                       // Admin who resolved it
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
	return ComplaintResponse{
		ID:         c.ID,
		UserID:     c.UserID,
		TargetID:   c.TargetID,
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
