package scan_history

import (
	"plantheon-backend/common"
	"plantheon-backend/models/diseases"

	"gorm.io/gorm"
)

// ScanHistoryService handles all database operations for scan histories
type ScanHistoryService struct {
	db *gorm.DB
}

// NewScanHistoryService creates a new scan history service instance
func NewScanHistoryService() *ScanHistoryService {
	return &ScanHistoryService{
		db: common.GetDB(),
	}
}

// CreateScanHistoryRecord creates a new scan history
func CreateScanHistoryRecord(scanHistory *ScanHistory) error {
	service := NewScanHistoryService()
	if err := service.db.Create(scanHistory).Error; err != nil {
		return err
	}
	
	// Load disease details manually
	var disease diseases.Disease
	if err := service.db.Where("id = ?", scanHistory.DiseaseID).First(&disease).Error; err != nil {
		return err
	}
	scanHistory.Disease = disease
	
	return nil
}

// GetScanHistoryByID gets a scan history by ID with disease details
func GetScanHistoryByID(id string) (*ScanHistory, error) {
	service := NewScanHistoryService()
	var scanHistory ScanHistory
	
	// Get scan history by ID
	if err := service.db.Where("id = ?", id).First(&scanHistory).Error; err != nil {
		return nil, err
	}
	
	// Load disease details
	var disease diseases.Disease
	if err := service.db.Where("id = ?", scanHistory.DiseaseID).First(&disease).Error; err != nil {
		return nil, err
	}
	scanHistory.Disease = disease
	
	return &scanHistory, nil
}

// GetAllScanHistories gets all scan histories with disease details
func GetAllScanHistories() ([]ScanHistory, error) {
	service := NewScanHistoryService()
	var scanHistories []ScanHistory
	
	// Get scan histories first
	if err := service.db.Order("created_at DESC").Find(&scanHistories).Error; err != nil {
		return nil, err
	}
	
	// Load disease details for each scan history
	for i := range scanHistories {
		var disease diseases.Disease
		if err := service.db.Where("id = ?", scanHistories[i].DiseaseID).First(&disease).Error; err != nil {
			return nil, err
		}
		scanHistories[i].Disease = disease
	}
	
	return scanHistories, nil
}

// DeleteScanHistoryByID deletes a scan history by ID
func DeleteScanHistoryByID(id string) error {
	service := NewScanHistoryService()
	
	// Check if scan history exists before deleting
	var scanHistory ScanHistory
	if err := service.db.Where("id = ?", id).First(&scanHistory).Error; err != nil {
		return err
	}
	
	// Delete the scan history
	if err := service.db.Delete(&scanHistory).Error; err != nil {
		return err
	}
	
	return nil
}

// DeleteAllScanHistories deletes all scan histories
func DeleteAllScanHistories() error {
	service := NewScanHistoryService()
	
	// Delete all scan histories
	if err := service.db.Where("1 = 1").Delete(&ScanHistory{}).Error; err != nil {
		return err
	}
	
	return nil
}
