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

// ValidateCreateScanComplaintRequest validates the create scan complaint request
func ValidateCreateScanComplaintRequest(req *CreateScanComplaintRequest) error {
	// Validate predicted_disease_id (required)
	req.PredictedDiseaseID = strings.TrimSpace(req.PredictedDiseaseID)
	if err := ValidateUUID(req.PredictedDiseaseID); err != nil {
		return errors.New("predicted_disease_id must be a valid UUID")
	}

	// Validate user_suggested_disease_id (optional, but must be valid UUID if provided)
	if req.UserSuggestedDiseaseID != nil {
		*req.UserSuggestedDiseaseID = strings.TrimSpace(*req.UserSuggestedDiseaseID)
		if *req.UserSuggestedDiseaseID != "" {
			if _, err := uuid.Parse(*req.UserSuggestedDiseaseID); err != nil {
				return errors.New("user_suggested_disease_id must be a valid UUID")
			}
		}
	}

	// Validate confidence_score (required, 0.0-1.0)
	if req.ConfidenceScore < 0.0 || req.ConfidenceScore > 1.0 {
		return errors.New("confidence_score must be between 0.0 and 1.0")
	}

	// Validate target_type
	req.TargetType = strings.ToUpper(strings.TrimSpace(req.TargetType))
	if req.TargetType != string(ComplaintTypeScan) {
		return errors.New("target_type must be SCAN")
	}

	// Validate category
	req.Category = strings.ToUpper(strings.TrimSpace(req.Category))
	validCategories := map[string]bool{
		string(ComplaintCategoryScanWrongResult):   true,
		string(ComplaintCategoryScanMisidentified): true,
		string(ComplaintCategoryScanIncorrectInfo): true,
		string(ComplaintCategoryScanPoorQuality):   true,
		string(ComplaintCategoryScanOtherIssue):    true,
	}
	if !validCategories[req.Category] {
		return errors.New("category must be one of: WRONG_RESULT, MISIDENTIFIED, INCORRECT_INFO, POOR_QUALITY, OTHER_ISSUE")
	}

	// Validate content (optional but limit length)
	req.Content = strings.TrimSpace(req.Content)
	if len(req.Content) > 1000 {
		return errors.New("content must be less than 1000 characters")
	}

	// Validate image_url (required, limit length)
	req.ImageURL = strings.TrimSpace(req.ImageURL)
	if req.ImageURL == "" {
		return errors.New("image_url is required for scan complaints")
	}
	if len(req.ImageURL) > 2000 {
		return errors.New("image_url must be less than 2000 characters")
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

// ValidateVerifyComplaintRequest validates the verify complaint request
func ValidateVerifyComplaintRequest(req *VerifyComplaintRequest) error {
	// Validate verified_disease_id (required)
	req.VerifiedDiseaseID = strings.TrimSpace(req.VerifiedDiseaseID)
	if err := ValidateUUID(req.VerifiedDiseaseID); err != nil {
		return errors.New("verified_disease_id must be a valid UUID")
	}

	// Validate admin_notes (optional but limit length)
	req.AdminNotes = strings.TrimSpace(req.AdminNotes)
	if len(req.AdminNotes) > 1000 {
		return errors.New("admin_notes must be less than 1000 characters")
	}

	return nil
}
