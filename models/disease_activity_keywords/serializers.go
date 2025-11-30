package disease_activity_keywords

// DiseaseActivityKeywordResponse represents the response structure
type DiseaseActivityKeywordResponse struct {
	ID                string `json:"id"`
	DiseaseID         string `json:"disease_id"`
	ActivityKeywordID string `json:"activity_keyword_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// ToDiseaseActivityKeywordResponse converts model to response
func (dak *DiseaseActivityKeyword) ToDiseaseActivityKeywordResponse() DiseaseActivityKeywordResponse {
	return DiseaseActivityKeywordResponse{
		ID:                dak.ID,
		DiseaseID:         dak.DiseaseID,
		ActivityKeywordID: dak.ActivityKeywordID,
		CreatedAt:         dak.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         dak.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// AddActivityKeywordRequest represents the request to add an activity keyword to a disease
type AddActivityKeywordRequest struct {
	ActivityKeywordID string `json:"activity_keyword_id" binding:"required"`
}

// SetActivityKeywordsRequest represents the request to set all activity keywords for a disease
type SetActivityKeywordsRequest struct {
	ActivityKeywordIDs []string `json:"activity_keyword_ids" binding:"required"`
}
