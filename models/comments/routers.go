package comments

import (
	"net/http"

	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"
)

func AddCommentHandler(c *gin.Context) {
	// Implementation of adding a comment to a post
	postID := c.Param("id")
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ValidateCreateCommentRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}
	
	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
		})
		return
	}

	comment := &Comments{
		PostID:  postID,
		Content: req.Content,
		UserID:  user.ID, // Sử dụng user.ID từ JWT token
	}

	if err := AddComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment added successfully",
		"data": gin.H{
			"id":        comment.ID,
			"post_id":   comment.PostID,
			"user_id":   comment.UserID,
			"full_name": user.FullName, // Thêm thông tin user
			"avatar":    user.Avatar,   // Thêm thông tin user
			"content":   comment.Content,
			"created_at": comment.CreatedAt,
		},
	})
}