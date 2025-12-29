package disease_activity_keywords

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"plantheon-backend/models/activity_keywords"
	"plantheon-backend/models/diseases"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupDiseaseActivityKeywordRoutes sets up the routes
func SetupDiseaseActivityKeywordRoutes(router *gin.RouterGroup) {
	// Add activity keyword to disease
	router.POST("/diseases/:disease_id/activity-keywords", AddActivityKeywordToDiseaseHandler)
	
	// Remove activity keyword from disease
	router.DELETE("/diseases/:disease_id/activity-keywords/:keyword_id", RemoveActivityKeywordFromDiseaseHandler)
	
	// Get all activity keywords for a disease
	router.GET("/diseases/:disease_id/activity-keywords", GetActivityKeywordsForDiseaseHandler)
	
	// Set all activity keywords for a disease (replace existing)
	router.PUT("/diseases/:disease_id/activity-keywords", SetActivityKeywordsForDiseaseHandler)
}

// AddActivityKeywordToDiseaseHandler adds an activity keyword to a disease
func AddActivityKeywordToDiseaseHandler(c *gin.Context) {
	diseaseID := c.Param("disease_id")
	
	var req AddActivityKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check if activity keyword exists
	_, err := activity_keywords.GetActivityKeywordByID(req.ActivityKeywordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy từ khóa hoạt động"})
		return
	}
	
	if err := AddActivityKeywordToDisease(diseaseID, req.ActivityKeywordID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Thêm từ khóa hoạt động vào bệnh thất bại"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Thêm từ khóa hoạt động thành công"})
}

// RemoveActivityKeywordFromDiseaseHandler removes an activity keyword from a disease
func RemoveActivityKeywordFromDiseaseHandler(c *gin.Context) {
	diseaseID := c.Param("disease_id")
	keywordID := c.Param("keyword_id")
	
	if err := RemoveActivityKeywordFromDisease(diseaseID, keywordID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa từ khóa hoạt động khỏi bệnh thất bại"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Xóa từ khóa hoạt động thành công"})
}

// GetActivityKeywordsForDiseaseHandler gets all activity keywords for a disease with full details
func GetActivityKeywordsForDiseaseHandler(c *gin.Context) {
	diseaseID := c.Param("disease_id")
	
	keywordIDs, err := GetActivityKeywordsByDiseaseID(diseaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách từ khóa hoạt động thất bại"})
		return
	}
	
	if len(keywordIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}
	
	// Fetch full activity keyword details
	keywords, err := activity_keywords.GetActivityKeywordsByIDs(keywordIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy chi tiết từ khóa hoạt động thất bại"})
		return
	}
	
	responses := make([]activity_keywords.ActivityKeywordResponse, len(keywords))
	for i, kw := range keywords {
		responses[i] = kw.ToActivityKeywordResponse()
	}
	
	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetDiseasesForKeywordHandler gets all diseases that have a specific activity keyword
func GetDiseasesForKeywordHandler(c *gin.Context) {
	keywordID := c.Param("keyword_id")
	
	// Check if keyword exists
	_, err := activity_keywords.GetActivityKeywordByID(keywordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy từ khóa hoạt động"})
		return
	}
	
	diseaseIDs, err := GetDiseasesByActivityKeywordID(keywordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách bệnh thất bại"})
		return
	}
	
	if len(diseaseIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}
	
	// Fetch full disease details
	diseaseList, err := diseases.GetDiseasesByIDs(diseaseIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy chi tiết bệnh thất bại"})
		return
	}
	
	responses := make([]diseases.DiseaseResponse, len(diseaseList))
	for i, d := range diseaseList {
		responses[i] = d.ToDiseaseResponse()
	}
	
	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// SetActivityKeywordsForDiseaseHandler sets all activity keywords for a disease (replaces existing)
func SetActivityKeywordsForDiseaseHandler(c *gin.Context) {
	diseaseID := c.Param("disease_id")
	
	var req SetActivityKeywordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate all activity keywords exist
	for _, keywordID := range req.ActivityKeywordIDs {
		_, err := activity_keywords.GetActivityKeywordByID(keywordID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy từ khóa hoạt động: " + keywordID})
			return
		}
	}
	
	if err := SetActivityKeywordsForDisease(diseaseID, req.ActivityKeywordIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Thiết lập từ khóa hoạt động thất bại"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Thiết lập từ khóa hoạt động thành công"})
}

// CSVImportError represents an error during CSV import
type CSVImportError struct {
	Row     int    `json:"row"`
	Error   string `json:"error"`
	Details string `json:"details"`
}

// ImportActivityKeywordsFromCSVHandler imports activity keywords and disease mappings from CSV
func ImportActivityKeywordsFromCSVHandler(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Không có tệp được tải lên",
		})
		return
	}

	// Check file extension
	filename := strings.ToLower(file.Filename)
	if !strings.HasSuffix(filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Chỉ hỗ trợ tệp .csv",
		})
		return
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Mở tệp thất bại",
		})
		return
	}
	defer src.Close()

	// Read CSV
	reader := csv.NewReader(src)
	rows, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Đọc tệp CSV thất bại: %v", err),
		})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tệp CSV phải có ít nhất 2 dòng (tiêu đề + dữ liệu)",
		})
		return
	}

	// Expected columns: ClassName, keywordName, keywordDescription, keywordType, keywordDayOffset, KeywordIsFreeTime, keywordHourTime, timeDuration
	var errors []CSVImportError
	successCount := 0
	keywordCache := make(map[string]string) // keywordName -> keywordID

	// Process each row (skip header)
	for i, row := range rows[1:] {
		rowNum := i + 2 // +2 because: +1 for skipping header, +1 for 1-based indexing

		if len(row) < 8 {
			errors = append(errors, CSVImportError{
				Row:     rowNum,
				Error:   "Dòng không hợp lệ",
				Details: "Dòng phải có ít nhất 8 cột",
			})
			continue
		}

		className := strings.TrimSpace(row[0])
		keywordName := strings.TrimSpace(row[1])
		keywordDescription := strings.TrimSpace(row[2])
		keywordType := strings.TrimSpace(row[3])
		keywordDayOffsetStr := strings.TrimSpace(row[4])
		keywordIsFreeTimeStr := strings.TrimSpace(row[5])
		keywordHourTimeStr := strings.TrimSpace(row[6])
		timeDurationStr := strings.TrimSpace(row[7])

		// Validate required fields
		if className == "" {
			errors = append(errors, CSVImportError{
				Row:     rowNum,
				Error:   "Thiếu tên lớp",
				Details: "Cần có tên lớp",
			})
			continue
		}

		if keywordName == "" {
			errors = append(errors, CSVImportError{
				Row:     rowNum,
				Error:   "Thiếu tên từ khóa",
				Details: "Cần có tên từ khóa",
			})
			continue
		}

		if keywordType == "" {
			errors = append(errors, CSVImportError{
				Row:     rowNum,
				Error:   "Thiếu loại từ khóa",
				Details: "Cần có loại từ khóa",
			})
			continue
		}

		// Parse keywordDayOffset
		keywordDayOffset := 0
		if keywordDayOffsetStr != "" {
			offset, err := strconv.Atoi(keywordDayOffsetStr)
			if err != nil {
				errors = append(errors, CSVImportError{
					Row:     rowNum,
					Error:   "Số ngày không hợp lệ",
					Details: fmt.Sprintf("Phải là số: %v", err),
				})
				continue
			}
			keywordDayOffset = offset
		}

		// Parse keywordIsFreeTime
		keywordIsFreeTime := false
		if keywordIsFreeTimeStr != "" {
			lowerStr := strings.ToLower(keywordIsFreeTimeStr)
			if lowerStr == "true" || lowerStr == "1" || lowerStr == "yes" {
				keywordIsFreeTime = true
			}
		}

		// Parse keywordHourTime
		var keywordHourTime *int
		if keywordHourTimeStr != "" {
			hourTime, err := strconv.Atoi(keywordHourTimeStr)
			if err != nil {
				errors = append(errors, CSVImportError{
					Row:     rowNum,
					Error:   "Giờ không hợp lệ",
					Details: fmt.Sprintf("Phải là số: %v", err),
				})
				continue
			}
			if hourTime < 0 || hourTime > 23 {
				errors = append(errors, CSVImportError{
					Row:     rowNum,
					Error:   "Giờ không hợp lệ",
					Details: "Phải từ 0 đến 23",
				})
				continue
			}
			keywordHourTime = &hourTime
		}

		// Parse timeDuration
		var timeDuration *int
		if timeDurationStr != "" {
			duration, err := strconv.Atoi(timeDurationStr)
			if err != nil {
				errors = append(errors, CSVImportError{
					Row:     rowNum,
					Error:   "Thời lượng không hợp lệ",
					Details: fmt.Sprintf("Phải là số: %v", err),
				})
				continue
			}
			if duration < 0 {
				errors = append(errors, CSVImportError{
					Row:     rowNum,
					Error:   "Thời lượng không hợp lệ",
					Details: "Phải là số không âm",
				})
				continue
			}
			timeDuration = &duration
		}

		// Find disease by ClassName
		disease, err := diseases.GetDiseaseByClassName(className)
		if err != nil {
			errors = append(errors, CSVImportError{
				Row:     rowNum,
				Error:   "Không tìm thấy bệnh",
				Details: fmt.Sprintf("Không có bệnh với tên lớp '%s'", className),
			})
			continue
		}

		// Check if keyword already exists in cache or database
		var keywordID string
		if cachedID, exists := keywordCache[keywordName]; exists {
			keywordID = cachedID
		} else {
			// Search for existing keyword by name
			existingKeywords, err := activity_keywords.SearchActivityKeywords(keywordName)
			if err == nil && len(existingKeywords) > 0 {
				// Find exact match
				for _, kw := range existingKeywords {
					if kw.Name == keywordName {
						keywordID = kw.ID
						keywordCache[keywordName] = keywordID
						break
					}
				}
			}

			// If not found, create new keyword
			if keywordID == "" {
				var desc *string
				if keywordDescription != "" {
					desc = &keywordDescription
				}

				newKeyword := activity_keywords.ActivityKeyword{
					Name:           keywordName,
					Description:    desc,
					Type:           keywordType,
					BaseDaysOffset: keywordDayOffset,
					IsFreeTime:     keywordIsFreeTime,
					HourTime:       keywordHourTime,
					TimeDuration:   timeDuration,
				}

				if err := activity_keywords.CreateActivityKeyword(&newKeyword); err != nil {
					errors = append(errors, CSVImportError{
						Row:     rowNum,
						Error:   "Tạo từ khóa thất bại",
						Details: fmt.Sprintf("Lỗi: %v", err),
					})
					continue
				}

				keywordID = newKeyword.ID
				keywordCache[keywordName] = keywordID
			}
		}

		// Create disease-activity keyword relationship
		if err := AddActivityKeywordToDisease(disease.ID, keywordID); err != nil {
			errors = append(errors, CSVImportError{
				Row:     rowNum,
				Error:   "Liên kết từ khóa với bệnh thất bại",
				Details: fmt.Sprintf("Lỗi: %v", err),
			})
			continue
		}

		successCount++
	}

	// Return response
	response := gin.H{
		"success": successCount,
		"total":   len(rows) - 1, // Exclude header
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["message"] = fmt.Sprintf("Nhập khẩu %d trên %d dòng với %d lỗi", successCount, len(rows)-1, len(errors))
	} else {
		response["message"] = fmt.Sprintf("Nhập khẩu thành công tất cả %d dòng", successCount)
	}

	c.JSON(http.StatusOK, response)
}
