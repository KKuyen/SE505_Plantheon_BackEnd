package guide_stages

import "time"

// Import sub guide stages for nested response
import sub_guide_stages "plantheon-backend/models/sub_guide_stages"

// GuideStageResponse represents guide stage response
type GuideStageResponse struct {
	ID             string                                   `json:"id"`
	PlantID        string                                   `json:"plant_id"`
	StageTitle     string                                   `json:"stage_title"`
	Description    *string                                  `json:"description,omitempty"`
	StartDayOffset int                                      `json:"start_day_offset"`
	EndDayOffset   int                                      `json:"end_day_offset"`
	ImageURL       *string                                  `json:"image_url"`
	CreatedAt      time.Time                                `json:"created_at"`
	SubGuideStages []sub_guide_stages.SubGuideStageResponse `json:"sub_guide_stages,omitempty"`
}

// CreateGuideStageRequest represents guide stage creation request
type CreateGuideStageRequest struct {
	PlantID        string  `json:"plant_id" binding:"required"`
	StageTitle     string  `json:"stage_title" binding:"required"`
	Description    *string `json:"description"`
	StartDayOffset int     `json:"start_day_offset"`
	EndDayOffset   int     `json:"end_day_offset"`
	ImageURL       *string `json:"image_url"`
}

// UpdateGuideStageRequest represents guide stage update request
type UpdateGuideStageRequest struct {
	StageTitle     *string `json:"stage_title"`
	Description    *string `json:"description"`
	StartDayOffset *int    `json:"start_day_offset"`
	EndDayOffset   *int    `json:"end_day_offset"`
	ImageURL       *string `json:"image_url"`
}

// ToGuideStageResponse converts GuideStage to GuideStageResponse
func (gs *GuideStage) ToGuideStageResponse() GuideStageResponse {
	return GuideStageResponse{
		ID:             gs.ID,
		PlantID:        gs.PlantID,
		StageTitle:     gs.StageTitle,
		Description:    gs.Description,
		StartDayOffset: gs.StartDayOffset,
		EndDayOffset:   gs.EndDayOffset,
		ImageURL:       gs.ImageURL,
		CreatedAt:      gs.CreatedAt,
	}
}
