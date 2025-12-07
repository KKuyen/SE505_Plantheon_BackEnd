package blog_tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateBlogTagHandler creates a new blog tag
func CreateBlogTagHandler(c *gin.Context) {
	var req BlogTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := ValidateBlogTagRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &BlogTag{Name: req.Name}
	if err := CreateBlogTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog tag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Blog tag created successfully",
		"data":    tag.ToBlogTagResponse(),
	})
}

// UpdateBlogTagHandler updates an existing blog tag
func UpdateBlogTagHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req BlogTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := ValidateBlogTagRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &BlogTag{ID: id, Name: req.Name}
	if err := UpdateBlogTag(tag); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Blog tag updated successfully",
		"data":    tag.ToBlogTagResponse(),
	})
}

// DeleteBlogTagHandler deletes a blog tag by id
func DeleteBlogTagHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := DeleteBlogTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Blog tag deleted successfully",
	})
}

// GetAllBlogTagsHandler returns all blog tags (public)
func GetAllBlogTagsHandler(c *gin.Context) {
	tags, err := GetAllBlogTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blog tags"})
		return
	}

	responses := make([]BlogTagResponse, 0, len(tags))
	for _, t := range tags {
		responses = append(responses, t.ToBlogTagResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"blog_tags": responses,
			"count":     len(responses),
		},
	})
}
