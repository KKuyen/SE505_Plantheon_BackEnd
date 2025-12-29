package guide_stages

import (
	"net/http"

	"plantheon-backend/models/blogs"
	"plantheon-backend/models/sub_guide_stages"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGuideStageHandler handles guide stage creation
func CreateGuideStageHandler(c *gin.Context) {
	var req CreateGuideStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng yêu cầu không hợp lệ",
		})
		return
	}

	// Validate request
	if err := ValidateCreateGuideStageRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create guide stage
	guideStage := &GuideStage{
		PlantID:        req.PlantID,
		StageTitle:     req.StageTitle,
		Description:    req.Description,
		StartDayOffset: req.StartDayOffset,
		EndDayOffset:   req.EndDayOffset,
		ImageURL:       req.ImageURL,
	}

	if err := CreateGuideStage(guideStage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Tạo giai đoạn hướng dẫn thất bại",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo giai đoạn hướng dẫn thành công",
		"data":    guideStage.ToGuideStageResponse(),
	})
}

// GetGuideStageByIDHandler handles getting guide stage by ID
func GetGuideStageByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID giai đoạn hướng dẫn",
		})
		return
	}

	guideStage, err := GetGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy giai đoạn hướng dẫn",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy giai đoạn hướng dẫn thất bại",
		})
		return
	}

	// Load sub guide stages for this guide stage
	subGuideStages, err := sub_guide_stages.GetSubGuideStagesByGuideStageID(guideStage.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy danh sách giai đoạn hướng dẫn phụ thất bại",
		})
		return
	}

	// Build response with sub guide stages
	response := guideStage.ToGuideStageResponse()
	subGuideStageResponses := make([]sub_guide_stages.SubGuideStageResponse, 0, len(subGuideStages))
	for _, sgs := range subGuideStages {
		sgsResp := sgs.ToSubGuideStageResponse()

		// Load published blogs for this sub guide stage
		subBlogs, err := blogs.GetBlogsBySubGuideStageID(sgs.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Lấy danh sách bài viết cho giai đoạn hướng dẫn phụ thất bại",
			})
			return
		}

		blogSummaries := make([]blogs.BlogSummaryResponse, 0, len(subBlogs))
		for _, b := range subBlogs {
			blogSummaries = append(blogSummaries, b.ToBlogSummaryResponse())
		}
		sgsResp.Blogs = blogSummaries

		subGuideStageResponses = append(subGuideStageResponses, sgsResp)
	}
	response.SubGuideStages = subGuideStageResponses

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetGuideStagesByPlantIDHandler handles getting all guide stages for a plant
func GetGuideStagesByPlantIDHandler(c *gin.Context) {
	plantID := c.Param("plant_id")
	if plantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID cây trồng",
		})
		return
	}

	guideStages, err := GetGuideStagesByPlantID(plantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy danh sách giai đoạn hướng dẫn thất bại",
		})
		return
	}

	// Convert to response format without sub guide stages
	var response []GuideStageResponse
	for _, gs := range guideStages {
		response = append(response, gs.ToGuideStageResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"guide_stages": response,
			"count":        len(response),
		},
	})
}

// UpdateGuideStageHandler handles guide stage update
func UpdateGuideStageHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID giai đoạn hướng dẫn",
		})
		return
	}

	// Get existing guide stage
	guideStage, err := GetGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy giai đoạn hướng dẫn",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy giai đoạn hướng dẫn thất bại",
		})
		return
	}

	// Parse update request
	var req UpdateGuideStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate request
	if err := ValidateUpdateGuideStageRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update guide stage fields if provided
	if req.StageTitle != nil {
		guideStage.StageTitle = *req.StageTitle
	}
	if req.StartDayOffset != nil {
		guideStage.StartDayOffset = *req.StartDayOffset
	}
	if req.EndDayOffset != nil {
		guideStage.EndDayOffset = *req.EndDayOffset
	}
	if req.ImageURL != nil {
		guideStage.ImageURL = req.ImageURL
	}
	if req.Description != nil {
		guideStage.Description = req.Description
	}

	// Validate day offsets after update
	if guideStage.EndDayOffset < guideStage.StartDayOffset {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ngày kết thúc phải lớn hơn hoặc bằng ngày bắt đầu",
		})
		return
	}

	// Save updated guide stage
	if err := UpdateGuideStage(guideStage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cập nhật giai đoạn hướng dẫn thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật giai đoạn hướng dẫn thành công",
		"data":    guideStage.ToGuideStageResponse(),
	})
}

// DeleteGuideStageHandler handles guide stage deletion
func DeleteGuideStageHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID giai đoạn hướng dẫn",
		})
		return
	}

	// Check if guide stage exists
	_, err := GetGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy giai đoạn hướng dẫn",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy giai đoạn hướng dẫn thất bại",
		})
		return
	}

	// Delete guide stage
	if err := DeleteGuideStage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Xóa giai đoạn hướng dẫn thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa giai đoạn hướng dẫn thành công",
	})
}
