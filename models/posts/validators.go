package posts

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func ValidateCreatePostRequest(req *CreatePostRequest) error {
	// Validate content
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return errors.New("post content is required")
	}
	if len(req.Content) > 5000 {
		return errors.New("post content must be less than 5000 characters")
	}
	// Validate image link array (optional)
	if req.ImageLink != nil {
		if len(req.ImageLink) > 5 {
			return errors.New("image link array cannot have more than 5 items")
		}
	}
	for i, link := range req.ImageLink {
		req.ImageLink[i] = strings.TrimSpace(link)
		if len(req.ImageLink[i]) > 500 {
			return errors.New("each image link must be less than 500 characters")
		}
	}

	// Validate tags
	if req.Tags != nil {
		if len(req.Tags) > 5 {
			return errors.New("tags array cannot have more than 5 items")
		}
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