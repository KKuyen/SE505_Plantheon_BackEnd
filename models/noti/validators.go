package noti

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

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

func ValidateMarkAsSeenRequest(req *MarkNotificationAsSeenRequest) error {
	// IsRead can be true or false, always valid
	return nil
}

func ValidateCreateNotificationRequest(req *CreateNotificationRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return errors.New("content is required")
	}
	if len(req.Content) > 500 {
		return errors.New("content must be less than 500 characters")
	}

	// Validate post_id if provided
	if req.PostID != nil && *req.PostID != "" {
		if err := ValidateIdParam(*req.PostID); err != nil {
			return err
		}
	}

	return nil
}

