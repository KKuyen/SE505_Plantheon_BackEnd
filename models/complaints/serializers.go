package complaints

import (
	"time"
)

// ComplaintResponse represents the response for a complaint
type ComplaintResponse struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	TargetID   string     `json:"target_id"`
	TargetType string     `json:"target_type"`
	Category   string     `json:"category"`
	Content    string     `json:"content"`
	Status     string     `json:"status"`
	AdminNotes string     `json:"admin_notes,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *string    `json:"resolved_by,omitempty"`
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

// UpdateComplaintStatusRequest represents the request to update complaint status (admin only)
type UpdateComplaintStatusRequest struct {
	Status     string `json:"status" binding:"required"` // PENDING, REVIEWED, RESOLVED, REJECTED
	AdminNotes string `json:"admin_notes"`               // Optional notes from admin
}
