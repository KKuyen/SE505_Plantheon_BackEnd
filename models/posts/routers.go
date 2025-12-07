package posts

import (
	"log"
	"net/http"

	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreatePostHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
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

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR - Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := ValidateCreatePostRequest(&req); err != nil {
		log.Printf("ERROR - Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	post := &Post{
		Content:       req.Content,
		ImageLink:     pq.StringArray(req.ImageLink),
		UserID:        user.ID, // Sử dụng UserID từ JWT token
		DiseaseLink:   req.DiseaseLink,
		ScanHistoryID: req.ScanHistoryID,
		Tags:          pq.StringArray(req.Tags),
	}
	
	log.Printf("DEBUG - Creating post: UserID=%s, DiseaseLink=%v, ScanHistoryID=%v", post.UserID, post.DiseaseLink, post.ScanHistoryID)
	
	if err := CreatePost(post); err != nil {
		log.Printf("ERROR - Failed to create post in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data": gin.H{
			"id":        post.ID,
			"user_id":   post.UserID,
			"full_name": user.FullName, // Lấy từ user object
			"avatar":    user.Avatar,   // Lấy từ user object
			"content":   post.Content,
			"tags":      post.Tags,
			"created_at": post.CreatedAt,
		},
	})

}

func UpdatePostHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Lấy thông tin user từ JWT token context
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
		ID:        id,
		Content:   req.Content,
		ImageLink: pq.StringArray(req.ImageLink),
		Tags:      pq.StringArray(req.Tags),
		UserID:  user.ID,
	}

	if err := UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"data": gin.H{
			"id":        post.ID,
			"content":   post.Content,
			"tags":      post.Tags,
			"updated_at": post.UpdatedAt,
		},
	})
}

func GetPostsHandler(c *gin.Context) {
	// Extract viewer's user ID from JWT token
	viewerID := ""
	userInterface, exists := c.Get("user")
	if exists {
		user, ok := userInterface.(*users.User)
		if ok {
			viewerID = user.ID
		}
	}
	
	posts, err := GetAllPosts(viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get posts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})	
}

func GetPostByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Extract viewer's user ID from JWT token
	viewerID := ""
	userInterface, exists := c.Get("user")
	if exists {
		user, ok := userInterface.(*users.User)
		if ok {
			viewerID = user.ID
		}
	}
	
	post, err := GetPostByID(id, viewerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func DeletePostByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Extract user ID from JWT token
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
	
	// Get post to check ownership
	post, err := GetPostByID(id, user.ID)
	if err != nil {
		log.Printf("ERROR - Failed to get post %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}
	
	// Check if user is the owner of the post
	if post.UserID != user.ID {
		log.Printf("WARNING - User %s attempted to delete post %s owned by %s", user.ID, id, post.UserID)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to delete this post",
		})
		return
	}
	
	if err := DeletePostByID(id); err != nil {
		log.Printf("ERROR - Failed to delete post %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete post",
			"details": err.Error(),
		})
		return
	}

	log.Printf("INFO - Post %s deleted successfully by user %s", id, user.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}

func LikePostHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Extract user ID from JWT token
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
	
	if err := LikePost(id, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Post liked successfully",
	})
}

func UnlikePostHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Extract user ID from JWT token
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
	
	if err := UnlikePost(id, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
	})
}

func SharePostHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := SharePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}	
	c.JSON(http.StatusOK, gin.H{
		"message": "Post shared successfully",
	})
}

func GetPostsByUserIDHandler(c *gin.Context) {
	userId := c.Param("userId")
	if err := ValidateIdParam(userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Extract viewer's user ID from JWT token
	viewerID := ""
	userInterface, exists := c.Get("user")
	if exists {
		user, ok := userInterface.(*users.User)
		if ok {
			viewerID = user.ID
		}
	}
	
	posts, err := GetPostsByUserID(userId, viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

func GetMyPostsHandler(c *gin.Context) {
	// Extract user ID from JWT token
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
	
	// Get posts by the authenticated user's ID
	posts, err := GetPostsByUserID(user.ID, user.ID)
	if err != nil {
		log.Printf("ERROR - Failed to get posts for user %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get posts",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}
