package diseases

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// CreateDiseaseHandler handles disease creation
func CreateDiseaseHandler(c *gin.Context) {
	var req CreateDiseaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng yêu cầu không hợp lệ",
		})
		return
	}

	// Validate request
	if err := ValidateCreateDiseaseRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create disease
	disease := &Disease{
		Name:        req.Name,
		ClassName:  req.ClassName,
		Type:        req.Type,
		Description: req.Description,
		Solution:    req.Solution,
		ImageLink:   pq.StringArray(req.ImageLink),
	}

	if err := CreateDiseaseRecord(disease); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Tạo bệnh thất bại",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo bệnh thành công",
		"data":    disease.ToDiseaseResponse(),
	})
}

// GetDisease handles getting disease by ID
func GetDisease(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID bệnh",
		})
		return
	}

	disease, err := GetDiseaseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy bệnh",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin bệnh thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": disease.ToDiseaseResponse(),
	})
}

// GetDiseases handles getting all diseases with pagination
func GetDiseases(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	diseaseType := c.Query("type")
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

	var diseases []Disease
	var total int64

	// Handle different query types
	if search != "" {
		diseases, total, err = SearchDiseases(search, offset, limit)
	} else if diseaseType != "" {
		diseases, total, err = GetDiseasesByType(diseaseType, offset, limit)
	} else {
		diseases, total, err = GetAllDiseases(offset, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy danh sách bệnh thất bại",
		})
		return
	}

	response := ToDiseasesListResponse(diseases, total, page, limit)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetAllDiseasesHandler handles getting all diseases without pagination
func GetAllDiseasesHandler(c *gin.Context) {
	// Parse query parameters for filtering
	diseaseType := c.Query("type")
	search := c.Query("search")

	var diseases []Disease
	var total int64
	var err error

	// Handle different query types and get count
	if search != "" {
		diseases, err = SearchAllDiseases(search)
		if err == nil {
			total = int64(len(diseases))
		}
	} else if diseaseType != "" {
		diseases, err = GetAllDiseasesByTypeWithoutPagination(diseaseType)
		if err == nil {
			total = int64(len(diseases))
		}
	} else {
		diseases, err = GetAllDiseasesWithoutPagination()
		if err == nil {
			total = int64(len(diseases))
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy danh sách bệnh thất bại",
		})
		return
	}

	// Convert to response format
	var response []DiseaseResponse
	for _, disease := range diseases {
		response = append(response, disease.ToDiseaseResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"diseases": response,
			"total":    total,
			"count":    len(response),
		},
	})
}

// GetDiseasesCountHandler handles getting diseases count only
func GetDiseasesCountHandler(c *gin.Context) {
	// Parse query parameters for filtering
	diseaseType := c.Query("type")
	search := c.Query("search")

	var count int64
	var err error

	// Handle different query types
	if search != "" {
		count, err = SearchDiseasesCount(search)
	} else if diseaseType != "" {
		count, err = GetDiseasesCountByType(diseaseType)
	} else {
		count, err = GetDiseasesCount()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy số lượng bệnh thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"count": count,
		},
	})
}

// UpdateDiseaseHandler handles disease update
func UpdateDiseaseHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có ID bệnh",
		})
		return
	}

	// Get existing disease
	disease, err := GetDiseaseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy bệnh",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin bệnh thất bại",
		})
		return
	}

	var req UpdateDiseaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Định dạng yêu cầu không hợp lệ",
		})
		return
	}

	// Validate request
	if err := ValidateUpdateDiseaseRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update disease fields if provided
	if req.Name != "" {
		disease.Name = req.Name
	}
	if req.ClassName != "" {
		disease.ClassName = req.ClassName
	}
	if req.Type != "" {
		disease.Type = req.Type
	}
	if req.Description != "" {
		disease.Description = req.Description
	}
	if req.Solution != "" {
		disease.Solution = req.Solution
	}
	if req.ImageLink != nil {
		disease.ImageLink = pq.StringArray(req.ImageLink)
	}

	// Save updated disease
	if err := UpdateDisease(disease); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cập nhật bệnh thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật bệnh thành công",
		"data":    disease.ToDiseaseResponse(),
	})
}

// DeleteDiseaseHandler handles disease deletion
func DeleteDiseaseHandler(c *gin.Context) {
	ClassName := c.Param("ClassName")
	if ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có tên lớp bệnh",
		})
		return
	}

	// Check if disease exists
	_, err := GetDiseaseByClassName(ClassName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy bệnh",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin bệnh thất bại",
		})
		return
	}

	// Delete disease
	if err := DeleteDisease(ClassName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Xóa bệnh thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa bệnh thành công",
	})
}

// GetDiseaseByClassNameHandler handles getting disease by class name
func GetDiseaseByClassNameHandler(c *gin.Context) {
	ClassName := c.Param("ClassName")
	if ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cần có tên lớp",
		})
		return
	}

	disease, err := GetDiseaseByClassName(ClassName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Không tìm thấy bệnh",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy thông tin bệnh thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": disease,
	})
}

// ImportDiseasesFromExcelHandler handles importing diseases from Excel file
func ImportDiseasesFromExcelHandler(c *gin.Context) {
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
	if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Chỉ hỗ trợ tệp .xlsx và .csv",
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

	var rows [][]string
	var fileType string

	// Determine file type and read accordingly
	if strings.HasSuffix(filename, ".csv") {
		fileType = "CSV"
		rows, err = readCSVFile(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Đọc tệp CSV thất bại: %v", err),
			})
			return
		}
	} else {
		fileType = "Excel"
		rows, err = readExcelFile(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Đọc tệp Excel thất bại: %v", err),
			})
			return
		}
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Tệp %s phải có ít nhất 2 dòng (tiêu đề + dữ liệu)", fileType),
		})
		return
	}

	// Process rows
	var errors []ExcelImportError
	var createdDiseases []DiseaseResponse
	var updatedDiseases []DiseaseResponse

	// Skip header row (index 0), start from index 1
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1 // Row number (1-based)

		// Check if row has enough columns
		if len(row) < 7 {
			errors = append(errors, ExcelImportError{
				Row:   rowNumber,
				Error: "Dòng phải có ít nhất 8 cột",
			})
			continue
		}

		// Parse row data
		excelRow := ExcelDiseaseRow{
			Name:        strings.TrimSpace(row[0]),
			ClassName:   strings.TrimSpace(row[1]),
			Type:        strings.TrimSpace(row[2]),
			Description: strings.TrimSpace(row[3]),
			Solution:    strings.TrimSpace(row[4]),
			ImageLink:   parseStringArray(row[5]),
			PlantName:   strings.TrimSpace(row[6]),
		}

		// Validate required fields
		if excelRow.Name == "" {
			errors = append(errors, ExcelImportError{
				Row:   rowNumber,
				Error: "Cần có tên",
			})
			continue
		}

		if excelRow.ClassName == "" {
			errors = append(errors, ExcelImportError{
				Row:   rowNumber,
				Error: "Cần có tên lớp",
			})
			continue
		}

		if excelRow.Type == "" {
			errors = append(errors, ExcelImportError{
				Row:   rowNumber,
				Error: "Cần có loại",
			})
			continue
		}

		// Check if disease already exists by className
		existingDisease, err := GetDiseaseByClassName(excelRow.ClassName)
		if err == nil && existingDisease != nil {
			// Update existing disease with new data
			existingDisease.Name = excelRow.Name
			existingDisease.Type = excelRow.Type
			existingDisease.Description = excelRow.Description
			existingDisease.Solution = excelRow.Solution
			existingDisease.ImageLink = pq.StringArray(excelRow.ImageLink)
			existingDisease.PlantName = excelRow.PlantName

			if err := UpdateDisease(existingDisease); err != nil {
				errors = append(errors, ExcelImportError{
					Row:   rowNumber,
					Error: fmt.Sprintf("Cập nhật bệnh thất bại: %v", err),
				})
				continue
			}

			// Add to updated list
			updatedDiseases = append(updatedDiseases, existingDisease.ToDiseaseResponse())
			continue
		}

		// Create new disease
		disease := &Disease{
			Name:        excelRow.Name,
			ClassName:   excelRow.ClassName,
			Type:        excelRow.Type,
			Description: excelRow.Description,
			Solution:    excelRow.Solution,
			ImageLink:   pq.StringArray(excelRow.ImageLink),
			PlantName:   excelRow.PlantName,
		}

		if err := CreateDiseaseRecord(disease); err != nil {
			errors = append(errors, ExcelImportError{
				Row:   rowNumber,
				Error: fmt.Sprintf("Tạo bệnh thất bại: %v", err),
			})
			continue
		}

		// Add to created list
		createdDiseases = append(createdDiseases, disease.ToDiseaseResponse())
	}

	// Prepare response
	response := ExcelImportResponse{
		TotalRows:       len(rows) - 1, // Exclude header
		CreatedCount:    len(createdDiseases),
		UpdatedCount:    len(updatedDiseases),
		ErrorCount:      len(errors),
		Errors:          errors,
		CreatedDiseases: createdDiseases,
		UpdatedDiseases: updatedDiseases,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Nhập khẩu tệp %s hoàn tất", fileType),
		"data":    response,
	})
}

// readCSVFile reads a CSV file and returns rows
func readCSVFile(src io.Reader) ([][]string, error) {
	reader := csv.NewReader(src)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// readExcelFile reads an Excel file and returns rows
func readExcelFile(src io.Reader) ([][]string, error) {
	xlFile, err := excelize.OpenReader(src)
	if err != nil {
		return nil, err
	}

	// Get first sheet
	sheetName := xlFile.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Read all rows
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// parseStringArray parses comma-separated string into array
func parseStringArray(str string) []string {
	if str == "" {
		return []string{}
	}
	
	// Split by comma and trim spaces
	parts := strings.Split(str, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

