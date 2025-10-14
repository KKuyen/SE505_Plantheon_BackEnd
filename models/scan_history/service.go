package scan_history

import (
	"plantheon-backend/common"

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
	return service.db.Create(scanHistory).Error
}