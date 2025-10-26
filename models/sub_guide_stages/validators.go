package sub_guide_stages

import (
	"errors"
	"strings"
)

// ValidateCreateSubGuideStageRequest validates create sub guide stage request
func ValidateCreateSubGuideStageRequest(req *CreateSubGuideStageRequest) error {
	if strings.TrimSpace(req.GuideStagesID) == "" {
		return errors.New("guide_stages_id is required")
	}

	if req.Title != nil && len(*req.Title) > 500 {
		return errors.New("title must be less than 500 characters")
	}

	if req.EndDayOffset < req.StartDayOffset {
		return errors.New("end_day_offset must be greater than or equal to start_day_offset")
	}

	return nil
}

// ValidateUpdateSubGuideStageRequest validates update sub guide stage request
func ValidateUpdateSubGuideStageRequest(req *UpdateSubGuideStageRequest) error {
	if req.Title != nil && len(*req.Title) > 500 {
		return errors.New("title must be less than 500 characters")
	}

	return nil
}
