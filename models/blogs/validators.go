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

	// Validate blog tag id if provided
	if req.BlogTagID != nil {
		blogTagID := strings.TrimSpace(*req.BlogTagID)
		if blogTagID == "" {
			req.BlogTagID = nil
		} else {
			if _, err := uuid.Parse(blogTagID); err != nil {
				return errors.New("blog_tag_id must be a valid UUID")
			}
			req.BlogTagID = &blogTagID
		}
	}

	// Validate description
	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if len(desc) > 1000 {
			return errors.New("news description must be less than 1000 characters")
		}
		req.Description = &desc
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

	// Validate sub_guide_stages_id if provided
	if req.SubGuideStagesID != nil && *req.SubGuideStagesID != "" {
		id := strings.TrimSpace(*req.SubGuideStagesID)
		if _, err := uuid.Parse(id); err != nil {
			return errors.New("sub_guide_stages_id must be a valid UUID")
		}
		req.SubGuideStagesID = &id
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
