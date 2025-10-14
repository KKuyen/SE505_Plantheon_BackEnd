package scan_history

import (
	"time"
)

// ScanHistoryResponse represents scan history response
type ScanHistoryResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	DiseaseID string    `json:"disease_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateScanHistoryRequest represents scan history creation request
type CreateScanHistoryRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	DiseaseID string `json:"disease_id" binding:"required"`
}

// ToScanHistoryResponse converts ScanHistory model to ScanHistoryResponse
func (s *ScanHistory) ToScanHistoryResponse() ScanHistoryResponse {
	return ScanHistoryResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		DiseaseID: s.DiseaseID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
