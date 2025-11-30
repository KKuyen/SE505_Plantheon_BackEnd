package activity_keywords

import (
	"errors"
)

// ValidateActivityKeywordCreate validates the create request
func ValidateActivityKeywordCreate(req *ActivityKeywordCreateRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	
	if req.Type == "" {
		return errors.New("type is required")
	}
	
	if req.BaseDaysOffset < 0 {
		return errors.New("base_days_offset must be non-negative")
	}
	
	if req.HourTime != nil && (*req.HourTime < 0 || *req.HourTime > 23) {
		return errors.New("hour_time must be between 0 and 23")
	}
	
	if req.EndHourTime != nil && (*req.EndHourTime < 0 || *req.EndHourTime > 23) {
		return errors.New("end_hour_time must be between 0 and 23")
	}
	
	if req.TimeDuration != nil && *req.TimeDuration < 0 {
		return errors.New("time_duration must be non-negative")
	}
	
	return nil
}

// ValidateActivityKeywordUpdate validates the update request
func ValidateActivityKeywordUpdate(req *ActivityKeywordUpdateRequest) error {
	if req.Name != nil && *req.Name == "" {
		return errors.New("name cannot be empty")
	}
	
	if req.Type != nil && *req.Type == "" {
		return errors.New("type cannot be empty")
	}
	
	if req.BaseDaysOffset != nil && *req.BaseDaysOffset < 0 {
		return errors.New("base_days_offset must be non-negative")
	}
	
	if req.HourTime != nil && (*req.HourTime < 0 || *req.HourTime > 23) {
		return errors.New("hour_time must be between 0 and 23")
	}
	
	if req.EndHourTime != nil && (*req.EndHourTime < 0 || *req.EndHourTime > 23) {
		return errors.New("end_hour_time must be between 0 and 23")
	}
	
	if req.TimeDuration != nil && *req.TimeDuration < 0 {
		return errors.New("time_duration must be non-negative")
	}
	
	return nil
}
