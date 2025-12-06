package plants

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

// PlantService handles all database operations for plants
type PlantService struct {
	db *gorm.DB
}

// NewPlantService creates a new plant service instance
func NewPlantService() *PlantService {
	return &PlantService{
		db: common.GetDB(),
	}
}

// GetPlantByName finds plant by name
func GetPlantByName(name string) (*Plant, error) {
	service := NewPlantService()
	var plant Plant
	err := service.db.Where("name = ?", name).First(&plant).Error
	return &plant, err
}

// GetPlantByID finds plant by ID
func GetPlantByID(id string) (*Plant, error) {
	service := NewPlantService()
	var plant Plant
	err := service.db.Where("id = ?", id).First(&plant).Error
	return &plant, err
}

// CreatePlantRecord creates a new plant
func CreatePlantRecord(plant *Plant) error {
	service := NewPlantService()
	return service.db.Create(plant).Error
}

// GetAllPlants gets all plants
func GetAllPlants() ([]Plant, error) {
	service := NewPlantService()
	var plants []Plant
	err := service.db.Order("created_at ASC").Find(&plants).Error
	return plants, err
}

// UpdatePlant updates an existing plant
func UpdatePlant(id string, updates map[string]interface{}) (*Plant, error) {
	service := NewPlantService()
	var plant Plant
	
	// Check if plant exists
	if err := service.db.Where("id = ?", id).First(&plant).Error; err != nil {
		return nil, err
	}
	
	// Update the plant
	if err := service.db.Model(&plant).Updates(updates).Error; err != nil {
		return nil, err
	}
	
	// Reload the plant to get updated values
	if err := service.db.Where("id = ?", id).First(&plant).Error; err != nil {
		return nil, err
	}
	
	return &plant, nil
}

// DeletePlant deletes a plant by ID
func DeletePlant(id string) error {
	service := NewPlantService()
	
	// Check if plant exists
	var plant Plant
	if err := service.db.Where("id = ?", id).First(&plant).Error; err != nil {
		return err
	}
	
	// Delete the plant
	return service.db.Delete(&plant).Error
}
