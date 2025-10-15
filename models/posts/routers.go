package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreatePostHandler(c *gin.Context) {
	// Lấy user ID từ JWT token context (cách tốt hơn)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := ValidateCreatePostRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	post := &Post{
		Content:   req.Content,
		ImageLink: pq.StringArray(req.ImageLink),
		UserID:    userIDStr, // Sử dụng UserID từ JWT token
		Tags:      pq.StringArray(req.Tags),
	}
	
	if err := CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data": gin.H{
			"id":      post.ID,
			"user_id": post.UserID,
			"content": post.Content,
			"tags":    post.Tags,
		},
	})

}