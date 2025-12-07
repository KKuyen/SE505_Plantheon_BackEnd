package comments

import (
	"errors"
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
func AddComment(comment *Comments, parentID *string) error {
	service := NewCommentService()
	return service.db.Transaction(func(tx *gorm.DB) error {
		comment.ParentID = nil
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// Determine parent: if provided use it, otherwise self-reference
		resolvedParent := comment.ID
		if parentID != nil && *parentID != "" {
			resolvedParent = *parentID
		}

		if err := tx.Model(comment).Update("parent_id", resolvedParent).Error; err != nil {
			return err
		}

		comment.ParentID = &resolvedParent

		return nil
	})
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

// LikeComment adds a like by user to a comment and increments like_num.
func LikeComment(commentID string, userID string) error {
	service := NewCommentService()
	return service.db.Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Table("comment_likes").Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&exists).Error; err != nil {
			return err
		}
		if exists > 0 {
			return nil
		}

		if err := tx.Create(&CommentLike{CommentID: commentID, UserID: userID}).Error; err != nil {
			return err
		}

		if err := tx.Model(&Comments{}).Where("id = ?", commentID).UpdateColumn("like_num", gorm.Expr("like_num + 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

// UnlikeComment removes a like by user from a comment and decrements like_num.
func UnlikeComment(commentID string, userID string) error {
	service := NewCommentService()
	return service.db.Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Table("comment_likes").Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return nil
		}

		if err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&CommentLike{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&Comments{}).Where("id = ? AND like_num > 0", commentID).UpdateColumn("like_num", gorm.Expr("like_num - 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetCommentsByPostID retrieves comments flat (including parent_id, is_like, is_me) ordered by newest first.
func GetCommentsByPostID(postID string, viewerID string) ([]CommentResponse, error) {
	service := NewCommentService()
	var comments []Comments
	if err := service.db.Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}

	// Build like map for viewer
	likedMap := make(map[string]bool)
	if viewerID != "" && len(comments) > 0 {
		ids := make([]string, len(comments))
		for i, c := range comments {
			ids[i] = c.ID
		}
		type LikeRow struct {
			CommentID string
		}
		var likeRows []LikeRow
		if err := service.db.Table("comment_likes").Select("comment_id").Where("user_id = ? AND comment_id IN ?", viewerID, ids).Find(&likeRows).Error; err == nil {
			for _, lr := range likeRows {
				likedMap[lr.CommentID] = true
			}
		}
	}

	// Preload user info
	userCache := make(map[string]struct {
		FullName string
		Avatar   string
	})
	resolveUser := func(userID string) (string, string) {
		if cached, ok := userCache[userID]; ok {
			return cached.FullName, cached.Avatar
		}
		var info struct {
			FullName string
			Avatar   string
		}
		if err := service.db.Table("users").Select("full_name, avatar").Where("id = ?", userID).First(&info).Error; err != nil {
			info.FullName = "Unknown User"
			info.Avatar = ""
		}
		userCache[userID] = info
		return info.FullName, info.Avatar
	}

	var responses []CommentResponse
	for _, c := range comments {
		pid := ""
		if c.ParentID != nil {
			pid = *c.ParentID
		}
		fullName, avatar := resolveUser(c.UserID)
		responses = append(responses, CommentResponse{
			ID:        c.ID,
			PostID:    c.PostID,
			UserID:    c.UserID,
			ParentID:  pid,
			FullName:  fullName,
			Avatar:    avatar,
			Content:   c.Content,
			LikeNum:   c.LikeNum,
			IsLike:    likedMap[c.ID],
			IsMe:      viewerID != "" && c.UserID == viewerID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return responses, nil
}

// ValidateParentComment ensures a provided parent belongs to the same post.
func ValidateParentComment(tx *gorm.DB, parentID string, postID string) error {
	var parent Comments
	if err := tx.Where("id = ?", parentID).First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("parent comment not found")
		}
		return err
	}
	if parent.PostID != postID {
		return errors.New("parent comment does not belong to the same post")
	}
	return nil
}