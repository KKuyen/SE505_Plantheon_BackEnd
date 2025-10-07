package activities

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateActivityHandler handles activity creation
func CreateActivityHandler(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := ValidateCreateActivityRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create activity
    activity := &Activity{
		Description:     req.Description,
		Description2:    req.Description2,
		Description3:    req.Description3,
		TimeStart:       req.TimeStart,
		TimeEnd:         req.TimeEnd,
		Day:             req.Day,
		Money:           req.Money,
        Type:            req.Type,
		Title:           req.Title,
		IsRepeat:        req.IsRepeat,
		EndRepeatDay:    req.EndRepeatDay,
		AlertTime:       req.AlertTime,
		Object:          req.Object,
		Amount:          req.Amount,
		Unit:            req.Unit,
		Purpose:         req.Purpose,
		TargetPerson:    req.TargetPerson,
		SourcePerson:    req.SourcePerson,
		AttachedLink:    req.AttachedLink,
		Note:            req.Note,
	}

	if err := CreateActivityRecord(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create activity",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Activity created successfully",
		"data":    activity.ToActivityResponse(),
	})
}

// GetActivity handles getting activity by ID
func GetActivity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Activity ID is required",
		})
		return
	}

	activity, err := GetActivityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get activity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": activity.ToActivityResponse(),
	})
}

// GetActivities handles getting all activities with pagination
func GetActivities(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	activityType := c.Query("type")
	search := c.Query("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// Validate pagination
	page, limit, _ = ValidatePaginationParams(page, limit)
	offset := (page - 1) * limit

	var activities []Activity
	var total int64

	// Handle different query types
	if search != "" {
		activities, total, err = SearchActivities(search, offset, limit)
	} else if activityType != "" {
		activities, total, err = GetActivitiesByType(activityType, offset, limit)
	} else {
		activities, total, err = GetAllActivities(offset, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get activities",
		})
		return
	}

	response := ToActivitiesListResponse(activities, total, page, limit)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetAllActivitiesHandler handles getting all activities without pagination
func GetAllActivitiesHandler(c *gin.Context) {
	// Parse query parameters for filtering
	activityType := c.Query("type")
	search := c.Query("search")

	var activities []Activity
	var total int64
	var err error

	// Handle different query types and get count
	if search != "" {
		activities, err = SearchAllActivities(search)
		if err == nil {
			total = int64(len(activities))
		}
	} else if activityType != "" {
		activities, err = GetAllActivitiesByTypeWithoutPagination(activityType)
		if err == nil {
			total = int64(len(activities))
		}
	} else {
		activities, err = GetAllActivitiesWithoutPagination()
		if err == nil {
			total = int64(len(activities))
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get activities",
		})
		return
	}

	// Convert to response format
	var response []ActivityResponse
	for _, activity := range activities {
		response = append(response, activity.ToActivityResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"activities": response,
			"total":      total,
			"count":      len(response),
		},
	})
}

// GetActivitiesCountHandler handles getting activities count only
func GetActivitiesCountHandler(c *gin.Context) {
	// Parse query parameters for filtering
	activityType := c.Query("type")
	search := c.Query("search")

	var count int64
	var err error

	// Handle different query types
	if search != "" {
		count, err = SearchActivitiesCount(search)
	} else if activityType != "" {
		count, err = GetActivitiesCountByType(activityType)
	} else {
		count, err = GetActivitiesCount()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get activities count",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"count": count,
		},
	})
}

// GetActivitiesByDayHandler returns all activities of a specific day (UTC) with minimal info
// GET /api/v1/activities/by-day?date=YYYY-MM-DD
func GetActivitiesByDayHandler(c *gin.Context) {
    dateStr := c.Query("date")
    if dateStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "date is required (YYYY-MM-DD)"})
        return
    }
    day, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected YYYY-MM-DD"})
        return
    }

    activities, err := GetActivitiesByDay(day)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities"})
        return
    }

    var response []ActivityDayItem
    for _, a := range activities {
        response = append(response, ActivityDayItem{
            ID:        a.ID,
            Title:     a.Title,
            TimeStart: a.TimeStart,
            TimeEnd:   a.TimeEnd,
            Day:       a.Day,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "date":       dateStr,
            "activities": response,
            "count":      len(response),
        },
    })
}

// GetActivitiesCalendarByMonthHandler returns an array sized by days in month
// Each element contains list of activities having time_start or day matching that date
// Query: GET /api/v1/activities/calendar?year=2025&month=9
func GetActivitiesCalendarByMonthHandler(c *gin.Context) {
    yearStr := c.Query("year")
    monthStr := c.Query("month")

    if yearStr == "" || monthStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "year and month are required",
        })
        return
    }

    year, err := strconv.Atoi(yearStr)
    if err != nil || year < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
        return
    }
    month, err := strconv.Atoi(monthStr)
    if err != nil || month < 1 || month > 12 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
        return
    }

    // Get all activities in month
    acts, err := GetActivitiesByMonthYear(year, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities"})
        return
    }

    // Prepare map date -> items
    // Determine number of days in month
    firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    nextMonth := firstDay.AddDate(0, 1, 0)
    daysInMonth := int(nextMonth.Add(-time.Nanosecond).Day())

    dayMap := make(map[string][]ActivityCalendarItem)
    for d := 1; d <= daysInMonth; d++ {
        key := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
        dayMap[key] = []ActivityCalendarItem{}
    }

    for _, a := range acts {
        // Candidate dates from time_start and day (both in UTC)
        if a.TimeStart != nil {
            key := time.Date(a.TimeStart.UTC().Year(), a.TimeStart.UTC().Month(), a.TimeStart.UTC().Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
            if _, ok := dayMap[key]; ok {
                dayMap[key] = append(dayMap[key], ActivityCalendarItem{
                    Title: a.Title,
                })
            }
        }
        if a.Day != nil {
            key := time.Date(a.Day.UTC().Year(), a.Day.UTC().Month(), a.Day.UTC().Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
            if _, ok := dayMap[key]; ok {
                dayMap[key] = append(dayMap[key], ActivityCalendarItem{
                    Title: a.Title,
                })
            }
        }
    }

    // Build ordered array by days
    var days []ActivityCalendarDay
    for d := 1; d <= daysInMonth; d++ {
        key := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
        days = append(days, ActivityCalendarDay{
            Date:       key,
            Activities: dayMap[key],
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "year":  year,
            "month": month,
            "days":  days,
            "count": len(days),
        },
    })
}

// UpdateActivityHandler handles activity update
func UpdateActivityHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Activity ID is required",
		})
		return
	}

	// Get existing activity
	activity, err := GetActivityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get activity",
		})
		return
	}

	// Parse update request
	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate request
	if err := ValidateUpdateActivityRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update activity fields if provided
	if req.Description != nil {
		activity.Description = req.Description
	}
	if req.Description2 != nil {
		activity.Description2 = req.Description2
	}
	if req.Description3 != nil {
		activity.Description3 = req.Description3
	}
	if req.TimeStart != nil {
		activity.TimeStart = req.TimeStart
	}
	if req.TimeEnd != nil {
		activity.TimeEnd = req.TimeEnd
	}
	if req.Day != nil {
		activity.Day = req.Day
	}
	if req.Money != nil {
		activity.Money = req.Money
	}
    if req.Type != nil {
        activity.Type = *req.Type
    }
	if req.Title != nil {
		activity.Title = *req.Title
	}
	if req.IsRepeat != nil {
		activity.IsRepeat = req.IsRepeat
	}
	if req.EndRepeatDay != nil {
		activity.EndRepeatDay = req.EndRepeatDay
	}
	if req.AlertTime != nil {
		activity.AlertTime = req.AlertTime
	}
	if req.Object != nil {
		activity.Object = req.Object
	}
	if req.Amount != nil {
		activity.Amount = req.Amount
	}
	if req.Unit != nil {
		activity.Unit = req.Unit
	}
	if req.Purpose != nil {
		activity.Purpose = req.Purpose
	}
	if req.TargetPerson != nil {
		activity.TargetPerson = req.TargetPerson
	}
	if req.SourcePerson != nil {
		activity.SourcePerson = req.SourcePerson
	}
	if req.AttachedLink != nil {
		activity.AttachedLink = req.AttachedLink
	}
	if req.Note != nil {
		activity.Note = req.Note
	}

	// Save updated activity
	if err := UpdateActivity(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update activity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Activity updated successfully",
		"data":    activity.ToActivityResponse(),
	})
}

// DeleteActivityHandler handles activity deletion
func DeleteActivityHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Activity ID is required",
		})
		return
	}

	// Check if activity exists
	_, err := GetActivityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Activity not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get activity",
		})
		return
	}

	// Delete activity
	if err := DeleteActivity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete activity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Activity deleted successfully",
	})
}
