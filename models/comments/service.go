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

func GetCommentByID(id string) (*Comments, error) {
	service := NewCommentService()
	var comment Comments
	if err := service.db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func UpdateComment(comment *Comments) error {
	service := NewCommentService()
	return service.db.Save(comment).Error
}

// Tăng số comment cho post
func IncreasePostCommentCount(postID string) error {
	db := common.GetDB()
	return db.Exec("UPDATE posts SET comment_num = comment_num + 1 WHERE id = ?", postID).Error
}

func DeleteCommentByID(id string) error {
	service := NewCommentService()
	return service.db.Delete(&Comments{}, "id = ?", id).Error
}