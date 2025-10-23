package comments

import (
	"errors"
)

func ValidateCreateCommentRequest(req *CreateCommentRequest) error {
	if req.Content == "" {
		return errors.New("content is required")
	}
	return nil
}

func ValidateUpdateCommentRequest(req *UpdateCommentRequest) error {
	if req.Content == "" {
		return errors.New("content is required")
	}
	return nil
}
