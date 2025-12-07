package comments

import (
	"errors"
)

func ValidateCreateCommentRequest(req *CreateCommentRequest) error {
	if req.Content == "" {
		return errors.New("content is required")
	}
	if req.ParentID != nil && *req.ParentID == "" {
		return errors.New("parent_id, if provided, cannot be empty")
	}
	return nil
}

func ValidateUpdateCommentRequest(req *UpdateCommentRequest) error {
	if req.Content == "" {
		return errors.New("content is required")
	}
	return nil
}
