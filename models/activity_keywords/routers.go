package activity_keywords

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupActivityKeywordRoutes sets up the routes for activity keywords
func SetupActivityKeywordRoutes(router *gin.Engine) {
	activityKeywordGroup := router.Group("/api/activity-keywords")
	{
		activityKeywordGroup.POST("", CreateActivityKeywordHandler)
		activityKeywordGroup.GET("", GetAllActivityKeywordsHandler)
		activityKeywordGroup.GET("/:id", GetActivityKeywordByIDHandler)
		activityKeywordGroup.PUT("/:id", UpdateActivityKeywordHandler)
		activityKeywordGroup.DELETE("/:id", DeleteActivityKeywordHandler)
		activityKeywordGroup.GET("/search", SearchActivityKeywordsHandler)
	}
}

// CreateActivityKeywordHandler handles creating a new activity keyword
func CreateActivityKeywordHandler(c *gin.Context) {
	var req ActivityKeywordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateActivityKeywordCreate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityKeyword := ActivityKeyword{
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		BaseDaysOffset: req.BaseDaysOffset,
		IsFreeTime:     req.IsFreeTime,
		HourTime:       req.HourTime,
		EndHourTime:    req.EndHourTime,
		TimeDuration:   req.TimeDuration,
	}

	if err := CreateActivityKeyword(&activityKeyword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity keyword"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": activityKeyword.ToActivityKeywordResponse()})
}

// GetAllActivityKeywordsHandler handles getting all activity keywords
func GetAllActivityKeywordsHandler(c *gin.Context) {
	activityKeywords, err := GetAllActivityKeywords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity keywords"})
		return
	}

	responses := make([]ActivityKeywordResponse, len(activityKeywords))
	for i, ak := range activityKeywords {
		responses[i] = ak.ToActivityKeywordResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetActivityKeywordByIDHandler handles getting activity keyword by ID
func GetActivityKeywordByIDHandler(c *gin.Context) {
	id := c.Param("id")

	activityKeyword, err := GetActivityKeywordByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity keyword not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": activityKeyword.ToActivityKeywordResponse()})
}

// UpdateActivityKeywordHandler handles updating activity keyword
func UpdateActivityKeywordHandler(c *gin.Context) {
	id := c.Param("id")

	var req ActivityKeywordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateActivityKeywordUpdate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if activity keyword exists
	_, err := GetActivityKeywordByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity keyword not found"})
		return
	}

	// Prepare updates map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.BaseDaysOffset != nil {
		updates["base_days_offset"] = *req.BaseDaysOffset
	}
	if req.IsFreeTime != nil {
		updates["is_free_time"] = *req.IsFreeTime
	}
	if req.HourTime != nil {
		updates["hour_time"] = *req.HourTime
	}
	if req.EndHourTime != nil {
		updates["end_hour_time"] = *req.EndHourTime
	}
	if req.TimeDuration != nil {
		updates["time_duration"] = *req.TimeDuration
	}

	if err := UpdateActivityKeyword(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity keyword"})
		return
	}

	// Fetch updated activity keyword
	updatedActivityKeyword, _ := GetActivityKeywordByID(id)
	c.JSON(http.StatusOK, gin.H{"data": updatedActivityKeyword.ToActivityKeywordResponse()})
}

// DeleteActivityKeywordHandler handles deleting activity keyword
func DeleteActivityKeywordHandler(c *gin.Context) {
	id := c.Param("id")

	// Check if activity keyword exists
	_, err := GetActivityKeywordByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity keyword not found"})
		return
	}

	if err := DeleteActivityKeyword(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity keyword"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity keyword deleted successfully"})
}

// SearchActivityKeywordsHandler handles searching activity keywords
func SearchActivityKeywordsHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword parameter is required"})
		return
	}

	activityKeywords, err := SearchActivityKeywords(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search activity keywords"})
		return
	}

	responses := make([]ActivityKeywordResponse, len(activityKeywords))
	for i, ak := range activityKeywords {
		responses[i] = ak.ToActivityKeywordResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// QueryActivityKeywordsHandler handles querying activity keywords by name and/or type with pagination
// Query params: name (optional), type (optional), page (default 1), limit (default 10, max 100)
func QueryActivityKeywordsHandler(c *gin.Context) {
	// Parse query parameters
	name := c.Query("name")
	keywordType := c.Query("type")

	// Parse page parameter
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100 // Max limit
			}
		}
	}

	activityKeywords, total, err := QueryActivityKeywords(name, keywordType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query activity keywords"})
		return
	}

	responses := make([]ActivityKeywordResponse, len(activityKeywords))
	for i, ak := range activityKeywords {
		responses[i] = ak.ToActivityKeywordResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
		"filters": gin.H{
			"name": name,
			"type": keywordType,
		},
	})
}

// GetActivityKeywordsPaginatedHandler handles getting activity keywords with pagination
// Query params: page (default 1), limit (default 10, max 100)
func GetActivityKeywordsPaginatedHandler(c *gin.Context) {
	// Parse page parameter
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100 // Max limit
			}
		}
	}

	activityKeywords, total, err := GetActivityKeywordsPaginated(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity keywords"})
		return
	}

	responses := make([]ActivityKeywordResponse, len(activityKeywords))
	for i, ak := range activityKeywords {
		responses[i] = ak.ToActivityKeywordResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
