package scan_history

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// ValidateCreateScanHistoryRequest validates scan history creation request
func ValidateCreateScanHistoryRequest(req *CreateScanHistoryRequest) error {
	// Validate UserID
	req.UserID = strings.TrimSpace(req.UserID)
	if req.UserID == "" {
		return errors.New("user id is required")
	}

	// Validate DiseaseID
	req.DiseaseID = strings.TrimSpace(req.DiseaseID)
	if req.DiseaseID == "" {
		return errors.New("disease id is required")
	}
	
	// Validate UUID format for DiseaseID
	if _, err := uuid.Parse(req.DiseaseID); err != nil {
		return errors.New("disease id must be a valid UUID")
	}

	return nil
}
