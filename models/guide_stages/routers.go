package guide_stages

import (
	"net/http"

	"plantheon-backend/models/sub_guide_stages"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGuideStageHandler handles guide stage creation
func CreateGuideStageHandler(c *gin.Context) {
	var req CreateGuideStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
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
		StartDayOffset: req.StartDayOffset,
		EndDayOffset:   req.EndDayOffset,
		ImageURL:       req.ImageURL,
	}

	if err := CreateGuideStage(guideStage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create guide stage",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Guide stage created successfully",
		"data":    guideStage.ToGuideStageResponse(),
	})
}

// GetGuideStageByIDHandler handles getting guide stage by ID
func GetGuideStageByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Guide stage ID is required",
		})
		return
	}

	guideStage, err := GetGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Guide stage not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get guide stage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": guideStage.ToGuideStageResponse(),
	})
}

// GetGuideStagesByPlantIDHandler handles getting all guide stages for a plant
func GetGuideStagesByPlantIDHandler(c *gin.Context) {
	plantID := c.Param("plant_id")
	if plantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Plant ID is required",
		})
		return
	}

	guideStages, err := GetGuideStagesByPlantID(plantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get guide stages",
		})
		return
	}

	// Convert to response format and populate sub guide stages
	var response []GuideStageResponse
	for _, gs := range guideStages {
		guideStageResponse := gs.ToGuideStageResponse()
		
		// Get sub guide stages for this guide stage
		subGuideStages, err := sub_guide_stages.GetSubGuideStagesByGuideStageID(gs.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get sub guide stages",
			})
			return
		}
		
		// Convert sub guide stages to response format
		subGuideStageResponses := make([]sub_guide_stages.SubGuideStageResponse, 0, len(subGuideStages))
		for _, sgs := range subGuideStages {
			subGuideStageResponses = append(subGuideStageResponses, sgs.ToSubGuideStageResponse())
		}
		
		guideStageResponse.SubGuideStages = subGuideStageResponses
		response = append(response, guideStageResponse)
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
			"error": "Guide stage ID is required",
		})
		return
	}

	// Get existing guide stage
	guideStage, err := GetGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Guide stage not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get guide stage",
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

	// Validate day offsets after update
	if guideStage.EndDayOffset < guideStage.StartDayOffset {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "end_day_offset must be greater than or equal to start_day_offset",
		})
		return
	}

	// Save updated guide stage
	if err := UpdateGuideStage(guideStage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update guide stage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Guide stage updated successfully",
		"data":    guideStage.ToGuideStageResponse(),
	})
}

// DeleteGuideStageHandler handles guide stage deletion
func DeleteGuideStageHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Guide stage ID is required",
		})
		return
	}

	// Check if guide stage exists
	_, err := GetGuideStageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Guide stage not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get guide stage",
		})
		return
	}

	// Delete guide stage
	if err := DeleteGuideStage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete guide stage",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Guide stage deleted successfully",
	})
}
