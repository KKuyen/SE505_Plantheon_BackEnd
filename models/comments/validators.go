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