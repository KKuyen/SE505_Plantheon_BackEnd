package blogs

import (
	"net/http"
	"strconv"
	"time"

	"plantheon-backend/models/blog_tags"
	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"
)

func GetNewsHandler(c *gin.Context) {
	var size *int
	if sizeParam := c.Query("size"); sizeParam != "" {
		parsed, err := strconv.Atoi(sizeParam)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "size phải là số nguyên dương",
			})
			return
		}
		size = &parsed
	}

	// Check if user is admin to determine if we should include draft posts
	userInterface, exists := c.Get("user")
	includeAllStatuses := false
	if exists {
		if user, ok := userInterface.(*users.User); ok && user.IsAdmin() {
			includeAllStatuses = true
		}
	}

	news, err := GetAllNews(size, includeAllStatuses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy danh sách tin tức thất bại",
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
			"error": "Không tìm thấy người dùng",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Định dạng người dùng không hợp lệ",
		})
		return
	}

	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Yêu cầu không hợp lệ",
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
		Title:            req.Title,
		Description:      req.Description,
		Content:          req.Content,
		CoverImageURL:    req.CoverImageURL,
		SubGuideStagesID: req.SubGuideStagesID,
		BlogTagID:        req.BlogTagID,
		Status:           req.Status,
		UserID:           user.ID,
	}

	// Validate provided blog tag id exists
	if blog.BlogTagID != nil {
		if _, err := blog_tags.GetBlogTagByID(*blog.BlogTagID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID thẻ blog không hợp lệ",
			})
			return
		}
	}

	if req.Status == "published" {
		now := time.Now()
		blog.PublishedAt = &now
	}

	if err := CreateNews(blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Tạo tin tức thất bại",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo tin tức thành công",
		"data": gin.H{
			"id":              blog.ID,
			"title":           blog.Title,
			"description":     blog.Description,
			"content":         blog.Content,
			"cover_image_url": blog.CoverImageURL,
			"blog_tag_id":     blog.BlogTagID,
			"status":          blog.Status,
			"published_at":    blog.PublishedAt,
			"user_id":         blog.UserID,
			"created_at":      blog.CreatedAt,
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
			"error": "Yêu cầu không hợp lệ",
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
		ID:               id,
		Title:            req.Title,
		Description:      req.Description,
		Content:          req.Content,
		CoverImageURL:    req.CoverImageURL,
		SubGuideStagesID: req.SubGuideStagesID,
		BlogTagID:        req.BlogTagID,
		Status:           req.Status,
	}

	// Validate provided blog tag id exists
	if blog.BlogTagID != nil {
		if _, err := blog_tags.GetBlogTagByID(*blog.BlogTagID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID thẻ blog không hợp lệ",
			})
			return
		}
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
		"message": "Cập nhật tin tức thành công",
		"data": gin.H{
			"id":              blog.ID,
			"title":           blog.Title,
			"description":     blog.Description,
			"content":         blog.Content,
			"cover_image_url": blog.CoverImageURL,
			"blog_tag_id":     blog.BlogTagID,
			"status":          blog.Status,
			"updated_at":      blog.UpdatedAt,
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
		"message": "Xóa tin tức thành công",
	})
}
