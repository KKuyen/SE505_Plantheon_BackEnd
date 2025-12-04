package scan_history

import (
	"plantheon-backend/models/diseases"
	"time"
)

// ScanHistoryResponse represents scan history response
type ScanHistoryResponse struct {
	ID        string                   `json:"id"`
	UserID    string                   `json:"user_id"`
	DiseaseID string                   `json:"disease_id"`
	Disease   diseases.DiseaseResponse `json:"disease"`
	ScanImage string                   `json:"scan_image"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// CreateScanHistoryRequest represents scan history creation request
type CreateScanHistoryRequest struct {
	DiseaseID string `json:"disease_id" binding:"required"`
	ScanImage string `json:"scan_image" binding:"required"`
}

// ScanHistoriesListResponse represents list of scan histories response
type ScanHistoriesListResponse struct {
	ScanHistories []ScanHistoryResponse `json:"scan_histories"`
	Total         int                   `json:"total"`
}

// ToScanHistoryResponse converts ScanHistory model to ScanHistoryResponse
func (s *ScanHistory) ToScanHistoryResponse() ScanHistoryResponse {
	return ScanHistoryResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		DiseaseID: s.DiseaseID,
		Disease:   s.Disease.ToDiseaseResponse(),
		ScanImage: s.ScanImage,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToScanHistoriesListResponse converts scan histories slice to list response
func ToScanHistoriesListResponse(scanHistories []ScanHistory) ScanHistoriesListResponse {
	scanHistoryResponses := make([]ScanHistoryResponse, len(scanHistories))
	for i, scanHistory := range scanHistories {
		scanHistoryResponses[i] = scanHistory.ToScanHistoryResponse()
	}

	return ScanHistoriesListResponse{
		ScanHistories: scanHistoryResponses,
		Total:         len(scanHistoryResponses),
	}
}
