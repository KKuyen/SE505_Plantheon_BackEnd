package plants

import (
	"errors"
	"strings"
)

// ValidateCreatePlantRequest validates the create plant request
func ValidateCreatePlantRequest(req *CreatePlantRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name cannot be empty")
	}

	if len(req.Name) > 255 {
		return errors.New("name cannot exceed 255 characters")
	}

	if req.Description != nil && len(*req.Description) > 5000 {
		return errors.New("description cannot exceed 5000 characters")
	}

	if req.ImageURL != nil && len(*req.ImageURL) > 1000 {
		return errors.New("image_url cannot exceed 1000 characters")
	}

	return nil
}

// ValidateUpdatePlantRequest validates the update plant request
func ValidateUpdatePlantRequest(req *UpdatePlantRequest) error {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return errors.New("name cannot be empty")
	}

	if req.Name != nil && len(*req.Name) > 255 {
		return errors.New("name cannot exceed 255 characters")
	}

	if req.Description != nil && len(*req.Description) > 5000 {
		return errors.New("description cannot exceed 5000 characters")
	}

	if req.ImageURL != nil && len(*req.ImageURL) > 1000 {
		return errors.New("image_url cannot exceed 1000 characters")
	}

	// Check if at least one field is being updated
	if req.Name == nil && req.Description == nil && req.ImageURL == nil {
		return errors.New("at least one field must be provided for update")
	}

	return nil
}
