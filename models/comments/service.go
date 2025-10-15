package comments

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

// CommentService handles all database operations for comments
type CommentService struct {
	db *gorm.DB
}

// NewCommentService creates a new comment service instance
func NewCommentService() *CommentService {
	return &CommentService{
		db: common.GetDB(),
	}
}

// AddComment adds a new comment to a post
func AddComment(comment *Comments) error {
	service := NewCommentService()
	return service.db.Create(comment).Error
}