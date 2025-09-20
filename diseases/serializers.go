package diseases

import (
	"time"
)

// DiseaseResponse represents disease response
type DiseaseResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ClassName  string    `json:"class_name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Solution    string    `json:"solution"`
	ImageLink   []string  `json:"image_link"`
	PlantName   string    `json:"plant_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateDiseaseRequest represents disease creation request
type CreateDiseaseRequest struct {
	Name        string   `json:"name" binding:"required"`
	ClassName  string   `json:"class_name" binding:"required"`
	Type        string   `json:"type" binding:"required"`
	Description string   `json:"description"`
	Solution    string   `json:"solution"`
	ImageLink   []string `json:"image_link"`
	PlantName   string   `json:"plant_name"`
}

// UpdateDiseaseRequest represents disease update request
type UpdateDiseaseRequest struct {
	Name        string   `json:"name"`
	ClassName  string   `json:"class_name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Solution    string   `json:"solution"`
	ImageLink   []string `json:"image_link"`
	PlantName   string   `json:"plant_name"`
}

// ExcelDiseaseRow represents a single row from Excel file
type ExcelDiseaseRow struct {
	Number      int      `json:"number"`       // Row number for display only
	Name        string   `json:"name"`
	ClassName   string   `json:"class_name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Solution    string   `json:"solution"`
	ImageLink   []string `json:"image_link"`
	PlantName   string   `json:"plant_name"`
}

// ExcelImportResponse represents response for Excel import
type ExcelImportResponse struct {
	TotalRows     int                `json:"total_rows"`
	SuccessCount  int                `json:"success_count"`
	ErrorCount    int                `json:"error_count"`
	Errors        []ExcelImportError `json:"errors"`
	CreatedDiseases []DiseaseResponse `json:"created_diseases"`
}

// ExcelImportError represents error for a specific row
type ExcelImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// DiseasesListResponse represents paginated diseases response
type DiseasesListResponse struct {
	Diseases []DiseaseResponse `json:"diseases"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Pages    int               `json:"pages"`
}

// ToDiseaseResponse converts Disease model to DiseaseResponse
func (d *Disease) ToDiseaseResponse() DiseaseResponse {
	return DiseaseResponse{
		ID:          d.ID,
		Name:        d.Name,
		ClassName:   d.ClassName,
		Type:        d.Type,
		Description: d.Description,
		Solution:    d.Solution,
		ImageLink:   []string(d.ImageLink),
		PlantName:   d.PlantName,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// ToDiseasesListResponse converts diseases slice to paginated response
func ToDiseasesListResponse(diseases []Disease, total int64, page, limit int) DiseasesListResponse {
	diseaseResponses := make([]DiseaseResponse, len(diseases))
	for i, disease := range diseases {
		diseaseResponses[i] = disease.ToDiseaseResponse()
	}

	pages := int(total) / limit
	if int(total)%limit != 0 {
		pages++
	}

	return DiseasesListResponse{
		Diseases: diseaseResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
		Pages:    pages,
	}
}
