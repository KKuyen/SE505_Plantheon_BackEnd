package blogs

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func ValidateCreateNewsRequest(req *CreateNewsRequest) error {
	// Validate title
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return errors.New("news title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("news title must be less than 200 characters")
	}

	// Validate content
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return errors.New("news content is required")
	}
	if len(req.Content) > 10000 {
		return errors.New("news content must be less than 10000 characters")
	}

	// Validate status
	if req.Status != "" && req.Status != "draft" && req.Status != "published" {
		return errors.New("status must be either 'draft' or 'published'")
	}
	if req.Status == "" {
		req.Status = "draft"
	}

	// Validate cover image URL (optional)
	if req.CoverImageURL != nil {
		url := strings.TrimSpace(*req.CoverImageURL)
		if len(url) > 500 {
			return errors.New("cover image URL must be less than 500 characters")
		}
		req.CoverImageURL = &url
	}

	return nil
}

func ValidateIdParam(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id parameter is required")
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("id parameter must be a valid UUID")
	}
	return nil
}
