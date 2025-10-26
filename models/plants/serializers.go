package plants

import (
	"time"
)

// PlantResponse represents the API response for a plant
type PlantResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ToPlantResponse converts a Plant model to PlantResponse
func (p *Plant) ToPlantResponse() PlantResponse {
	return PlantResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

// CreatePlantRequest represents the request body for creating a plant
type CreatePlantRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
}

// UpdatePlantRequest represents the request body for updating a plant
type UpdatePlantRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
}
