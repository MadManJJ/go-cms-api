package repositories

import (
	"github.com/MadManJJ/cms-api/dto"

	"gorm.io/gorm"
)

// MockCMSRepository is a mock implementation of CMSRepository
type MockCMSRepository struct{
	db *gorm.DB
}

// NewMockCMSRepository creates a new instance of MockCMSRepository
func NewMockCMSRepository(db *gorm.DB) *MockCMSRepository {
	return &MockCMSRepository{db: db}
}

// GetTestData implements the CMSRepository interface
func (r *MockCMSRepository) GetTestData() (*dto.CMSData, error) {
	// Return mock data
	return &dto.CMSData{
		ID:      "cms-123",
		Message: "Test data from CMS repository",
	}, nil
}

// GetAdditionalData implements the CMSRepository interface
func (r *MockCMSRepository) GetAdditionalData() (*dto.CMSData, error) {
	// Return mock data
	return &dto.CMSData{
		ID:      "cms-456",
		Message: "Additional data from CMS repository",
	}, nil
}
