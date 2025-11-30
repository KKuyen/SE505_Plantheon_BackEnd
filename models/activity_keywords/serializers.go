package activity_keywords

// ActivityKeywordResponse represents the response structure for activity keyword
type ActivityKeywordResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	Type           string  `json:"type"`
	BaseDaysOffset int     `json:"base_days_offset"`
	IsFreeTime     bool    `json:"is_free_time"`
	HourTime       *int    `json:"hour_time"`
	EndHourTime    *int    `json:"end_hour_time"`
	TimeDuration   *int    `json:"time_duration"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// ToActivityKeywordResponse converts ActivityKeyword to ActivityKeywordResponse
func (ak *ActivityKeyword) ToActivityKeywordResponse() ActivityKeywordResponse {
	return ActivityKeywordResponse{
		ID:             ak.ID,
		Name:           ak.Name,
		Description:    ak.Description,
		Type:           ak.Type,
		BaseDaysOffset: ak.BaseDaysOffset,
		IsFreeTime:     ak.IsFreeTime,
		HourTime:       ak.HourTime,
		EndHourTime:    ak.EndHourTime,
		TimeDuration:   ak.TimeDuration,
		CreatedAt:      ak.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      ak.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ActivityKeywordCreateRequest represents the request structure for creating activity keyword
type ActivityKeywordCreateRequest struct {
	Name           string  `json:"name" binding:"required"`
	Description    *string `json:"description"`
	Type           string  `json:"type" binding:"required"`
	BaseDaysOffset int     `json:"base_days_offset"`
	IsFreeTime     bool    `json:"is_free_time"`
	HourTime       *int    `json:"hour_time"`
	EndHourTime    *int    `json:"end_hour_time"`
	TimeDuration   *int    `json:"time_duration"`
}

// ActivityKeywordUpdateRequest represents the request structure for updating activity keyword
type ActivityKeywordUpdateRequest struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Type           *string `json:"type"`
	BaseDaysOffset *int    `json:"base_days_offset"`
	IsFreeTime     *bool   `json:"is_free_time"`
	HourTime       *int    `json:"hour_time"`
	EndHourTime    *int    `json:"end_hour_time"`
	TimeDuration   *int    `json:"time_duration"`
}
