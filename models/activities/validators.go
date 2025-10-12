package activities

import (
	"errors"
	"strings"
)

// ValidateCreateActivityRequest validates create activity request
func ValidateCreateActivityRequest(req *CreateActivityRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}

    if strings.TrimSpace(req.Type) == "" {
        return errors.New("type is required")
    }

	if len(req.Title) > 255 {
		return errors.New("title must be less than 255 characters")
	}

	if req.Description != nil && len(*req.Description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}

	if req.Description2 != nil && len(*req.Description2) > 1000 {
		return errors.New("description2 must be less than 1000 characters")
	}

	if req.Description3 != nil && len(*req.Description3) > 1000 {
		return errors.New("description3 must be less than 1000 characters")
	}

    if len(req.Type) > 255 {
        return errors.New("type must be less than 255 characters")
    }

	if req.IsRepeat != nil && len(*req.IsRepeat) > 50 {
		return errors.New("is_repeat must be less than 50 characters")
	}

	if req.AlertTime != nil && len(*req.AlertTime) > 50 {
		return errors.New("alert_time must be less than 50 characters")
	}

	if req.Object != nil && len(*req.Object) > 255 {
		return errors.New("object must be less than 255 characters")
	}

	if req.Unit != nil && len(*req.Unit) > 50 {
		return errors.New("unit must be less than 50 characters")
	}

	if req.Purpose != nil && len(*req.Purpose) > 1000 {
		return errors.New("purpose must be less than 1000 characters")
	}

	if req.TargetPerson != nil && len(*req.TargetPerson) > 255 {
		return errors.New("target_person must be less than 255 characters")
	}

	if req.SourcePerson != nil && len(*req.SourcePerson) > 255 {
		return errors.New("source_person must be less than 255 characters")
	}

	if req.AttachedLink != nil && len(*req.AttachedLink) > 1000 {
		return errors.New("attached_link must be less than 1000 characters")
	}

	if req.Note != nil && len(*req.Note) > 1000 {
		return errors.New("note must be less than 1000 characters")
	}

	if req.Money != nil && *req.Money < 0 {
		return errors.New("money must be non-negative")
	}

	if req.Amount != nil && *req.Amount < 0 {
		return errors.New("amount must be non-negative")
	}

	return nil
}

// ValidateUpdateActivityRequest validates update activity request
func ValidateUpdateActivityRequest(req *UpdateActivityRequest) error {
	if req.Title != nil {
		if strings.TrimSpace(*req.Title) == "" {
			return errors.New("title cannot be empty")
		}
		if len(*req.Title) > 255 {
			return errors.New("title must be less than 255 characters")
		}
	}

	if req.Description != nil && len(*req.Description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}

	if req.Description2 != nil && len(*req.Description2) > 1000 {
		return errors.New("description2 must be less than 1000 characters")
	}

	if req.Description3 != nil && len(*req.Description3) > 1000 {
		return errors.New("description3 must be less than 1000 characters")
	}

	if req.Type != nil && len(*req.Type) > 255 {
		return errors.New("type must be less than 255 characters")
	}

	if req.IsRepeat != nil && len(*req.IsRepeat) > 50 {
		return errors.New("is_repeat must be less than 50 characters")
	}

	if req.AlertTime != nil && len(*req.AlertTime) > 50 {
		return errors.New("alert_time must be less than 50 characters")
	}

	if req.Object != nil && len(*req.Object) > 255 {
		return errors.New("object must be less than 255 characters")
	}

	if req.Unit != nil && len(*req.Unit) > 50 {
		return errors.New("unit must be less than 50 characters")
	}

	if req.Purpose != nil && len(*req.Purpose) > 1000 {
		return errors.New("purpose must be less than 1000 characters")
	}

	if req.TargetPerson != nil && len(*req.TargetPerson) > 255 {
		return errors.New("target_person must be less than 255 characters")
	}

	if req.SourcePerson != nil && len(*req.SourcePerson) > 255 {
		return errors.New("source_person must be less than 255 characters")
	}

	if req.AttachedLink != nil && len(*req.AttachedLink) > 1000 {
		return errors.New("attached_link must be less than 1000 characters")
	}

	if req.Note != nil && len(*req.Note) > 1000 {
		return errors.New("note must be less than 1000 characters")
	}

	if req.Money != nil && *req.Money < 0 {
		return errors.New("money must be non-negative")
	}

	if req.Amount != nil && *req.Amount < 0 {
		return errors.New("amount must be non-negative")
	}

	return nil
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(page, limit int) (int, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit, nil
}
