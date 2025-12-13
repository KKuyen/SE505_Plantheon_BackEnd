package complaints

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// ValidateUUID validates a UUID string
func ValidateUUID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("id must be a valid UUID")
	}
	return nil
}

// ValidateCreateComplaintRequest validates the create complaint request
func ValidateCreateComplaintRequest(req *CreateComplaintRequest) error {
	// Validate target_id
	req.TargetID = strings.TrimSpace(req.TargetID)
	if err := ValidateUUID(req.TargetID); err != nil {
		return errors.New("target_id must be a valid UUID")
	}

	// Validate target_type
	req.TargetType = strings.ToUpper(strings.TrimSpace(req.TargetType))
	if req.TargetType != string(ComplaintTypePost) && req.TargetType != string(ComplaintTypeComment) {
		return errors.New("target_type must be POST or COMMENT")
	}

	// Validate category
	req.Category = strings.ToUpper(strings.TrimSpace(req.Category))
	validCategories := map[string]bool{
		string(ComplaintCategorySpam):           true,
		string(ComplaintCategoryHarassment):     true,
		string(ComplaintCategoryHateSpeech):     true,
		string(ComplaintCategoryViolence):       true,
		string(ComplaintCategoryMisinformation): true,
		string(ComplaintCategoryInappropriate):  true,
		string(ComplaintCategoryOther):          true,
	}
	if !validCategories[req.Category] {
		return errors.New("category must be one of: SPAM, HARASSMENT, HATE_SPEECH, VIOLENCE, MISINFORMATION, INAPPROPRIATE, OTHER")
	}

	// Validate content (optional but limit length)
	req.Content = strings.TrimSpace(req.Content)
	if len(req.Content) > 1000 {
		return errors.New("content must be less than 1000 characters")
	}

	return nil
}

// ValidateUpdateComplaintStatusRequest validates the update status request
func ValidateUpdateComplaintStatusRequest(req *UpdateComplaintStatusRequest) error {
	// Validate status
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	validStatuses := map[string]bool{
		string(ComplaintStatusPending):  true,
		string(ComplaintStatusReviewed): true,
		string(ComplaintStatusResolved): true,
		string(ComplaintStatusRejected): true,
	}
	if !validStatuses[req.Status] {
		return errors.New("status must be one of: PENDING, REVIEWED, RESOLVED, REJECTED")
	}

	// Validate admin_notes (optional but limit length)
	req.AdminNotes = strings.TrimSpace(req.AdminNotes)
	if len(req.AdminNotes) > 1000 {
		return errors.New("admin_notes must be less than 1000 characters")
	}

	return nil
}
