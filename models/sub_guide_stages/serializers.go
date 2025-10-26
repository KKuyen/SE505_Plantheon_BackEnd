package sub_guide_stages

// SubGuideStageResponse represents sub guide stage response
type SubGuideStageResponse struct {
	ID             string  `json:"id"`
	GuideStagesID  *string `json:"guide_stages_id"`
	Title          *string `json:"title"`
	StartDayOffset int     `json:"start_day_offset"`
	EndDayOffset   int     `json:"end_day_offset"`
}

// CreateSubGuideStageRequest represents sub guide stage creation request
type CreateSubGuideStageRequest struct {
	GuideStagesID  string  `json:"guide_stages_id" binding:"required"`
	Title          *string `json:"title"`
	StartDayOffset int     `json:"start_day_offset"`
	EndDayOffset   int     `json:"end_day_offset"`
}

// UpdateSubGuideStageRequest represents sub guide stage update request
type UpdateSubGuideStageRequest struct {
	Title          *string `json:"title"`
	StartDayOffset *int    `json:"start_day_offset"`
	EndDayOffset   *int    `json:"end_day_offset"`
}

// ToSubGuideStageResponse converts SubGuideStage to SubGuideStageResponse
func (sgs *SubGuideStage) ToSubGuideStageResponse() SubGuideStageResponse {
	return SubGuideStageResponse{
		ID:             sgs.ID,
		GuideStagesID:  sgs.GuideStagesID,
		Title:          sgs.Title,
		StartDayOffset: sgs.StartDayOffset,
		EndDayOffset:   sgs.EndDayOffset,
	}
}
