package scan_history

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Create scan history
	scanHistory := &ScanHistory{
		UserID:    req.UserID,
		DiseaseID: req.DiseaseID,
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

// GetScanHistoriesHandler handles getting all scan histories
func GetScanHistoriesHandler(c *gin.Context) {
	scanHistories, err := GetAllScanHistories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get scan histories",
		})
		return
	}

	response := ToScanHistoriesListResponse(scanHistories)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
