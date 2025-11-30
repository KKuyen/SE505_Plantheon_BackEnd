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

// GetCommentsByPostID retrieves all comments for a specific post with user info
func GetCommentsByPostID(postID string, viewerID string) ([]CommentResponse, error) {
	service := NewCommentService()
	var comments []Comments
	if err := service.db.Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	
	var commentResponses []CommentResponse
	for _, comment := range comments {
		// Query user info for each comment
		var userInfo struct {
			FullName string
			Avatar   string
		}
		if err := service.db.Table("users").Select("full_name, avatar").Where("id = ?", comment.UserID).First(&userInfo).Error; err != nil {
			// If user not found, use default values
			userInfo.FullName = "Unknown User"
			userInfo.Avatar = ""
		}
		
		commentResponses = append(commentResponses, CommentResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			FullName:  userInfo.FullName,
			Avatar:    userInfo.Avatar,
			Content:   comment.Content,
			LikeNum:   comment.LikeNum,
			IsMe:      viewerID != "" && comment.UserID == viewerID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}
	
	return commentResponses, nil
}