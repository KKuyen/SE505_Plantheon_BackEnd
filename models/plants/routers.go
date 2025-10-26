package plants

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupPlantRoutes sets up all plant-related routes
func SetupPlantRoutes(router *gin.RouterGroup) {
	router.POST("", CreatePlantHandler)
	router.GET("", GetAllPlantsHandler)
	router.GET("/:id", GetPlantByIDHandler)
	router.PUT("/:id", UpdatePlantHandler)
	router.DELETE("/:id", DeletePlantHandler)
}

// CreatePlantHandler handles creating a new plant
func CreatePlantHandler(c *gin.Context) {
	var req CreatePlantRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
	if err := ValidateCreatePlantRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create plant
	plant := Plant{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	}

	if err := CreatePlantRecord(&plant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create plant",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": plant.ToPlantResponse(),
	})
}

// GetAllPlantsHandler handles getting all plants
func GetAllPlantsHandler(c *gin.Context) {
	plants, err := GetAllPlants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get plants",
		})
		return
	}

	// Convert to response format
	var response []PlantResponse
	for _, plant := range plants {
		response = append(response, plant.ToPlantResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"plants": response,
			"count":  len(response),
		},
	})
}

// GetPlantByIDHandler handles getting a plant by ID
func GetPlantByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Plant ID is required",
		})
		return
	}

	plant, err := GetPlantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Plant not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": plant.ToPlantResponse(),
	})
}

// UpdatePlantHandler handles updating a plant
func UpdatePlantHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Plant ID is required",
		})
		return
	}

	var req UpdatePlantRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
	if err := ValidateUpdatePlantRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}

	// Update plant
	plant, err := UpdatePlant(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update plant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": plant.ToPlantResponse(),
	})
}

// DeletePlantHandler handles deleting a plant
func DeletePlantHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Plant ID is required",
		})
		return
	}

	// Delete plant
	if err := DeletePlant(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Plant not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plant deleted successfully",
	})
}
