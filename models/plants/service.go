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
	err := service.db.Find(&plants).Error
	return plants, err
}
