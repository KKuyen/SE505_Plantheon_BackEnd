package blog_tags

import (
	"errors"
	"strings"
)

// ValidateBlogTagRequest validates create/update requests
func ValidateBlogTagRequest(req *BlogTagRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) > 200 {
		return errors.New("name must be less than 200 characters")
	}
	req.Name = name
	return nil
}
