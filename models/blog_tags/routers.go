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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yêu cầu không hợp lệ"})
		return
	}

	if err := ValidateBlogTagRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &BlogTag{Name: req.Name}
	if err := CreateBlogTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tạo thẻ blog thất bại"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo thẻ blog thành công",
		"data":    tag.ToBlogTagResponse(),
	})
}

// UpdateBlogTagHandler updates an existing blog tag
func UpdateBlogTagHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có ID"})
		return
	}

	var req BlogTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yêu cầu không hợp lệ"})
		return
	}

	if err := ValidateBlogTagRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &BlogTag{ID: id, Name: req.Name}
	if err := UpdateBlogTag(tag); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thẻ blog"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cập nhật thẻ blog thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thẻ blog thành công",
		"data":    tag.ToBlogTagResponse(),
	})
}

// DeleteBlogTagHandler deletes a blog tag by id
func DeleteBlogTagHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có ID"})
		return
	}

	if err := DeleteBlogTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa thẻ blog thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa thẻ blog thành công",
	})
}

// GetAllBlogTagsHandler returns all blog tags (public)
func GetAllBlogTagsHandler(c *gin.Context) {
	tags, err := GetAllBlogTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách thẻ blog thất bại"})
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
