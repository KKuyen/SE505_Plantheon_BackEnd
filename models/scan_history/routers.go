package scan_history

import (
	"net/http"

	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"

	"strconv"
)

// CreateScanHistoryHandler handles scan history creation
func CreateScanHistoryHandler(c *gin.Context) {
	var req CreateScanHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := ValidateCreateScanHistoryRequest(&req); err != nil {
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

	scanHistory := &ScanHistory{
		UserID:    user.ID,
		DiseaseID: req.DiseaseID,
		ScanImage: req.ScanImage,
	}

	if err := CreateScanHistoryRecord(scanHistory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create scan history",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Scan history created successfully",
		"data":    scanHistory.ToScanHistoryResponse(),
	})
}

// GetScanHistoriesHandler handles getting all scan histories for the authenticated user
func GetScanHistoriesHandler(c *gin.Context) {
	// Get user from JWT token context
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

	sizeParam := c.DefaultQuery("size", "0")
	size := 0
	if sizeParam != "0" {
		var errParse error
		size, errParse = strconv.Atoi(sizeParam)
		if errParse != nil || size < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
			return
		}
	}

	// Get scan histories filtered by user ID
	scanHistories, err := GetScanHistoriesByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get scan histories"})
		return
	}

	// Nếu size > 0 thì chỉ lấy size record mới nhất
	if size > 0 && size < len(scanHistories) {
		scanHistories = scanHistories[:size]
	}

	response := ToScanHistoriesListResponse(scanHistories)
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetScanHistoryByIDHandler handles getting a scan history by ID
func GetScanHistoryByIDHandler(c *gin.Context) {
	// Get user from JWT token context
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

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}

	scanHistory, err := GetScanHistoryByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Scan history not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get scan history",
		})
		return
	}

	// Verify that the scan history belongs to the authenticated user
	if scanHistory.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to access this scan history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": scanHistory.ToScanHistoryResponse(),
	})
}

// DeleteScanHistoryByIDHandler handles deleting a scan history by ID
func DeleteScanHistoryByIDHandler(c *gin.Context) {
	// Get user from JWT token context
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

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}

	// Get scan history to verify ownership
	scanHistory, err := GetScanHistoryByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Scan history not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get scan history",
		})
		return
	}

	// Verify that the scan history belongs to the authenticated user
	if scanHistory.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to delete this scan history",
		})
		return
	}

	// Delete the scan history
	err = DeleteScanHistoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete scan history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Scan history deleted successfully",
	})
}

// DeleteAllScanHistoriesHandler handles deleting all scan histories for the authenticated user
func DeleteAllScanHistoriesHandler(c *gin.Context) {
	// Get user from JWT token context
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

	// Delete scan histories filtered by user ID
	err := DeleteScanHistoriesByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete all scan histories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All scan histories deleted successfully",
	})
}
