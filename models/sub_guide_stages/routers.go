package sub_guide_stages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSubGuideStageHandler handles sub guide stage creation
func CreateSubGuideStageHandler(c *gin.Context) {
	var req CreateSubGuideStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := ValidateCreateSubGuideStageRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create sub guide stage
	subGuideStage := &SubGuideStage{
		GuideStagesID:  &req.GuideStagesID,
		Title:          req.Title,
		StartDayOffset: req.StartDayOffset,
		EndDayOffset:   req.EndDayOffset,
	}

	if err := CreateSubGuideStage(subGuideStage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create sub guide stage",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sub guide stage created successfully",
		"data":    subGuideStage.ToSubGuideStageResponse(),
	})
}

// GetSubGuideStageByIDHandler handles getting sub guide stage by ID
func GetSubGuideStageByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sub guide stage ID is required",
		})
		return
	}

	subGuideStage, err := GetSubGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sub guide stage not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sub guide stage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": subGuideStage.ToSubGuideStageResponse(),
	})
}

// GetSubGuideStagesByGuideStageIDHandler handles getting all sub guide stages for a guide stage
func GetSubGuideStagesByGuideStageIDHandler(c *gin.Context) {
	guideStageID := c.Param("guide_stage_id")
	if guideStageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Guide stage ID is required",
		})
		return
	}

	subGuideStages, err := GetSubGuideStagesByGuideStageID(guideStageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sub guide stages",
		})
		return
	}

	// Convert to response format
	var response []SubGuideStageResponse
	for _, sgs := range subGuideStages {
		response = append(response, sgs.ToSubGuideStageResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"sub_guide_stages": response,
			"count":            len(response),
		},
	})
}

// UpdateSubGuideStageHandler handles sub guide stage update
func UpdateSubGuideStageHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sub guide stage ID is required",
		})
		return
	}

	// Get existing sub guide stage
	subGuideStage, err := GetSubGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sub guide stage not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sub guide stage",
		})
		return
	}

	// Parse update request
	var req UpdateSubGuideStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate request
	if err := ValidateUpdateSubGuideStageRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update sub guide stage fields if provided
	if req.Title != nil {
		subGuideStage.Title = req.Title
	}
	if req.StartDayOffset != nil {
		subGuideStage.StartDayOffset = *req.StartDayOffset
	}
	if req.EndDayOffset != nil {
		subGuideStage.EndDayOffset = *req.EndDayOffset
	}

	// Validate day offsets after update
	if subGuideStage.EndDayOffset < subGuideStage.StartDayOffset {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "end_day_offset must be greater than or equal to start_day_offset",
		})
		return
	}

	// Save updated sub guide stage
	if err := UpdateSubGuideStage(subGuideStage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update sub guide stage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sub guide stage updated successfully",
		"data":    subGuideStage.ToSubGuideStageResponse(),
	})
}

// DeleteSubGuideStageHandler handles sub guide stage deletion
func DeleteSubGuideStageHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sub guide stage ID is required",
		})
		return
	}

	// Check if sub guide stage exists
	_, err := GetSubGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sub guide stage not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sub guide stage",
		})
		return
	}

	// Delete sub guide stage
	if err := DeleteSubGuideStage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete sub guide stage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sub guide stage deleted successfully",
	})
}
