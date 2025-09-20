package diseases

import (
	"errors"
	"strings"
)

// ValidateCreateDiseaseRequest validates disease creation request
func ValidateCreateDiseaseRequest(req *CreateDiseaseRequest) error {
	// Validate name
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return errors.New("disease name is required")
	}
	if len(req.Name) > 255 {
		return errors.New("disease name must be less than 255 characters")
	}

	// Validate class name
	req.ClassName = strings.TrimSpace(req.ClassName)
	if req.ClassName == "" {
		return errors.New("class name is required")
	}
	if len(req.ClassName) > 255 {
		return errors.New("class name must be less than 255 characters")
	}

	// Validate type
	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		return errors.New("disease type is required")
	}
	if len(req.Type) > 100 {
		return errors.New("disease type must be less than 100 characters")
	}

	// Validate description (optional)
	req.Description = strings.TrimSpace(req.Description)
	if len(req.Description) > 5000 {
		return errors.New("description must be less than 5000 characters")
	}

	// Validate solution (optional)
	req.Solution = strings.TrimSpace(req.Solution)
	if len(req.Solution) > 5000 {
		return errors.New("solution must be less than 5000 characters")
	}

	// Validate image link array (optional)
	if req.ImageLink != nil {
		if len(req.ImageLink) > 20 {
			return errors.New("image link array cannot have more than 20 items")
		}
		for i, link := range req.ImageLink {
			req.ImageLink[i] = strings.TrimSpace(link)
			if len(req.ImageLink[i]) > 500 {
				return errors.New("each image link must be less than 500 characters")
			}
		}
	}

	return nil
}

// ValidateUpdateDiseaseRequest validates disease update request
func ValidateUpdateDiseaseRequest(req *UpdateDiseaseRequest) error {
	// Validate name (optional)
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)
		if len(req.Name) > 255 {
			return errors.New("disease name must be less than 255 characters")
		}
	}

	// Validate class name (optional)
	if req.ClassName != "" {
		req.ClassName = strings.TrimSpace(req.ClassName)
		if len(req.ClassName) > 255 {
			return errors.New("class name must be less than 255 characters")
		}
	}

	// Validate type (optional)
	if req.Type != "" {
		req.Type = strings.TrimSpace(req.Type)
		if len(req.Type) > 100 {
			return errors.New("disease type must be less than 100 characters")
		}
	}

	// Validate description (optional)
	if req.Description != "" {
		req.Description = strings.TrimSpace(req.Description)
		if len(req.Description) > 5000 {
			return errors.New("description must be less than 5000 characters")
		}
	}

	// Validate solution (optional)
	if req.Solution != "" {
		req.Solution = strings.TrimSpace(req.Solution)
		if len(req.Solution) > 5000 {
			return errors.New("solution must be less than 5000 characters")
		}
	}

	// Validate image link array (optional)
	if req.ImageLink != nil {
		if len(req.ImageLink) > 20 {
			return errors.New("image link array cannot have more than 20 items")
		}
		for i, link := range req.ImageLink {
			req.ImageLink[i] = strings.TrimSpace(link)
			if len(req.ImageLink[i]) > 500 {
				return errors.New("each image link must be less than 500 characters")
			}
		}
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
	} else if limit > 100 {
		limit = 100
	}
	
	return page, limit, nil
}
