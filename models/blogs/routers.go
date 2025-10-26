package blogs

import (
	"net/http"
	"time"

	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"
)

func GetNewsHandler(c *gin.Context) {
	news, err := GetAllNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get news",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": news,
	})
}

func GetNewsByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	news, err := GetNewsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": news,
	})
}

func CreateNewsHandler(c *gin.Context) {
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

	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := ValidateCreateNewsRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	blog := &Blog{
		Title:       req.Title,
		Content:     req.Content,
		CoverImageURL: req.CoverImageURL,
		Status:      req.Status,
		UserID:      user.ID,
	}

	if req.Status == "published" {
		now := time.Now()
		blog.PublishedAt = &now
	}

	if err := CreateNews(blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create news",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "News created successfully",
		"data": gin.H{
			"id":             blog.ID,
			"title":          blog.Title,
			"content":        blog.Content,
			"cover_image_url": blog.CoverImageURL,
			"status":         blog.Status,
			"published_at":   blog.PublishedAt,
			"user_id":        blog.UserID,
			"created_at":     blog.CreatedAt,
		},
	})
}

func UpdateNewsHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := ValidateCreateNewsRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	blog := &Blog{
		ID:          id,
		Title:       req.Title,
		Content:     req.Content,
		CoverImageURL: req.CoverImageURL,
		Status:      req.Status,
	}

	if req.Status == "published" {
		// Check if already published
		existingBlog, err := GetNewsByID(id)
		if err == nil && existingBlog.PublishedAt == nil {
			now := time.Now()
			blog.PublishedAt = &now
		}
	}

	if err := UpdateNews(blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "News updated successfully",
		"data": gin.H{
			"id":             blog.ID,
			"title":          blog.Title,
			"content":        blog.Content,
			"cover_image_url": blog.CoverImageURL,
			"status":         blog.Status,
			"updated_at":     blog.UpdatedAt,
		},
	})
}

func DeleteNewsHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := DeleteNewsByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "News deleted successfully",
	})
}
