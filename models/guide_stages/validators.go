package guide_stages

import (
	"errors"
	"strings"
)

// ValidateCreateGuideStageRequest validates create guide stage request
func ValidateCreateGuideStageRequest(req *CreateGuideStageRequest) error {
	if strings.TrimSpace(req.PlantID) == "" {
		return errors.New("plant_id is required")
	}

	if strings.TrimSpace(req.StageTitle) == "" {
		return errors.New("stage_title is required")
	}

	if len(req.StageTitle) > 500 {
		return errors.New("stage_title must be less than 500 characters")
	}

	if req.EndDayOffset < req.StartDayOffset {
		return errors.New("end_day_offset must be greater than or equal to start_day_offset")
	}

	if req.ImageURL != nil && len(*req.ImageURL) > 2000 {
		return errors.New("image_url must be less than 2000 characters")
	}

	return nil
}

// ValidateUpdateGuideStageRequest validates update guide stage request
func ValidateUpdateGuideStageRequest(req *UpdateGuideStageRequest) error {
	if req.StageTitle != nil {
		if strings.TrimSpace(*req.StageTitle) == "" {
			return errors.New("stage_title cannot be empty")
		}
		if len(*req.StageTitle) > 500 {
			return errors.New("stage_title must be less than 500 characters")
		}
	}

	if req.ImageURL != nil && len(*req.ImageURL) > 2000 {
		return errors.New("image_url must be less than 2000 characters")
	}

	return nil
}
