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
			"error": "Định dạng yêu cầu không hợp lệ",
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
		Repeat:          req.Repeat,
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
			"error": "Tạo hoạt động thất bại",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo hoạt động thành công",
		"data":    activity.ToActivityResponse(),
	})
}

// GetActivity handles getting activity by ID
func GetActivity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID hoạt động",
		})
		return
	}

	activity, err := GetActivityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy hoạt động",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin hoạt động thất bại",
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
	dateStr := c.Query("date") // Format: YYYY-MM-DD

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

	// Handle date filter first (highest priority)
	if dateStr != "" {
		selectedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Định dạng ngày không hợp lệ. Sử dụng YYYY-MM-DD",
			})
			return
		}
		activities, total, err = GetActivitiesByDateRangeWithPagination(selectedDate, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Lấy danh sách hoạt động thất bại",
			})
			return
		}
	} else if search != "" {
		activities, total, err = SearchActivities(search, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Lấy danh sách hoạt động thất bại",
			})
			return
		}
	} else if activityType != "" {
		activities, total, err = GetActivitiesByType(activityType, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Lấy danh sách hoạt động thất bại",
			})
			return
		}
	} else {
		activities, total, err = GetAllActivities(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Lấy danh sách hoạt động thất bại",
			})
			return
		}
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
	dateStr := c.Query("date") // Format: YYYY-MM-DD

	var activities []Activity
	var total int64
	var err error

	// Handle date filter first (highest priority)
	if dateStr != "" {
		selectedDate, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Định dạng ngày không hợp lệ. Sử dụng YYYY-MM-DD",
			})
			return
		}
		activities, err = GetActivitiesByDateRange(selectedDate)
		if err == nil {
			total = int64(len(activities))
		}
	} else if search != "" {
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
			"error": "Lấy danh sách hoạt động thất bại",
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
			"error": "Lấy số lượng hoạt động thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"count": count,
		},
	})
}

// GetActivitiesByDayHandler returns all activities where time_start <= date AND time_end >= date
// GET /api/v1/activities/by-day?date=YYYY-MM-DD
func GetActivitiesByDayHandler(c *gin.Context) {
    dateStr := c.Query("date")
    if dateStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có ngày (YYYY-MM-DD)"})
        return
    }
    day, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ, mong muốn YYYY-MM-DD"})
        return
    }

    // Use the new date range filter: time_start <= date AND time_end >= date
    activities, err := GetActivitiesByDateRange(day)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách hoạt động thất bại"})
        return
    }

    var response []ActivityResponse
    for _, a := range activities {
        response = append(response, a.ToActivityResponse())
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
            "error": "Cần có năm và tháng",
        })
        return
    }

    year, err := strconv.Atoi(yearStr)
    if err != nil || year < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Năm không hợp lệ"})
        return
    }
    month, err := strconv.Atoi(monthStr)
    if err != nil || month < 1 || month > 12 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Tháng không hợp lệ"})
        return
    }

    // Get all activities in month
    acts, err := GetActivitiesByMonthYear(year, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách hoạt động thất bại"})
        return
    }

    // Get all recurring activities
    var allRecurringActivities []Activity
    err = GetAllRecurringActivities(&allRecurringActivities)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách hoạt động lặp lại thất bại"})
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

    // For each regular activity, add it to ALL days it overlaps with
    for _, a := range acts {
        if a.TimeStart != nil && a.TimeEnd != nil {
            // Get start and end dates of the activity
            activityStart := time.Date(a.TimeStart.UTC().Year(), a.TimeStart.UTC().Month(), a.TimeStart.UTC().Day(), 0, 0, 0, 0, time.UTC)
            activityEnd := time.Date(a.TimeEnd.UTC().Year(), a.TimeEnd.UTC().Month(), a.TimeEnd.UTC().Day(), 0, 0, 0, 0, time.UTC)
            
            // Add activity to each day it spans
            for d := 1; d <= daysInMonth; d++ {
                currentDay := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
                
                // Check if activity overlaps with this day
                if (activityStart.Before(currentDay) || activityStart.Equal(currentDay)) && 
                   (activityEnd.After(currentDay) || activityEnd.Equal(currentDay)) {
						key := currentDay.Format("2006-01-02")
						dayMap[key] = append(dayMap[key], ActivityCalendarItem{Title: a.Title, Type: a.Type})
                }
            }
        }
    }

    // For each recurring activity, check if it should appear on each day
    for _, a := range allRecurringActivities {
        for d := 1; d <= daysInMonth; d++ {
            currentDay := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
            
            // Check if recurring activity should appear on this day
			if CheckActivityRepeatsOnDate(&a, currentDay) {
				key := currentDay.Format("2006-01-02")

				// Check if already added (avoid duplicates)
				alreadyAdded := false
				for _, item := range dayMap[key] {
						if item.Title == a.Title && item.Type == a.Type {
						alreadyAdded = true
						break
					}
				}

				if !alreadyAdded {
						dayMap[key] = append(dayMap[key], ActivityCalendarItem{Title: a.Title, Type: a.Type})
				}
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
			"error": "Cần có ID hoạt động",
		})
		return
	}

	// Get existing activity
	activity, err := GetActivityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy hoạt động",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin hoạt động thất bại",
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
	if req.Repeat != nil {
		activity.Repeat = req.Repeat
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
			"error": "Cập nhật hoạt động thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật hoạt động thành công",
		"data":    activity.ToActivityResponse(),
	})
}

// DeleteActivityHandler handles activity deletion
func DeleteActivityHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID hoạt động",
		})
		return
	}

	// Check if activity exists
	_, err := GetActivityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy hoạt động",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin hoạt động thất bại",
		})
		return
	}

	// Delete activity
	if err := DeleteActivity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Xóa hoạt động thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa hoạt động thành công",
	})
}

// DeleteActivitiesHandler handles bulk deletion of activities by IDs
func DeleteActivitiesHandler(c *gin.Context) {
	var req DeleteActivitiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng yêu cầu không hợp lệ"})
		return
	}

	if err := ValidateDeleteActivitiesRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deleted, err := DeleteActivities(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa hoạt động thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa hoạt động thành công",
		"deleted": deleted,
	})
}

// GetMonthlyFinancialSummaryHandler handles getting financial summary for a specific month
// GET /api/v1/activities/financial/monthly?year=2025&month=10
func GetMonthlyFinancialSummaryHandler(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có năm và tháng",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Năm không hợp lệ"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tháng không hợp lệ, phải từ 1-12"})
		return
	}

	summary, err := GetMonthlyFinancialSummary(year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy báo cáo tài chính thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}

// GetAnnualFinancialSummaryHandler handles getting financial summary for a whole year (12 months)
// GET /api/v1/activities/financial/annual?year=2025
func GetAnnualFinancialSummaryHandler(c *gin.Context) {
	yearStr := c.Query("year")

	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có năm",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Năm không hợp lệ"})
		return
	}

	summary, err := GetAnnualFinancialSummary(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy báo cáo tài chính năm thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}

// GetMultiYearFinancialSummaryHandler handles getting financial summary for multiple years
// GET /api/v1/activities/financial/multi-year?start_year=2020&end_year=2025
func GetMultiYearFinancialSummaryHandler(c *gin.Context) {
	startYearStr := c.Query("start_year")
	endYearStr := c.Query("end_year")

	if startYearStr == "" || endYearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có năm bắt đầu và năm kết thúc",
		})
		return
	}

	startYear, err := strconv.Atoi(startYearStr)
	if err != nil || startYear < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Năm bắt đầu không hợp lệ"})
		return
	}

	endYear, err := strconv.Atoi(endYearStr)
	if err != nil || endYear < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Năm kết thúc không hợp lệ"})
		return
	}

	summary, err := GetMultiYearFinancialSummary(startYear, endYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy báo cáo tài chính nhiều năm thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}
