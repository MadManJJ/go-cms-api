package repositories

import "github.com/MadManJJ/cms-api/dto"

// MockAppRepository is a mock implementation of AppRepository
type MockAppRepository struct{}

// NewMockAppRepository creates a new instance of MockAppRepository
func NewMockAppRepository() *MockAppRepository {
	return &MockAppRepository{}
}

// GetTestData implements the AppRepository interface
func (r *MockAppRepository) GetTestData() (*dto.AppData, error) {
	// Return mock data
	return &dto.AppData{
		ID:      "app-123",
		Message: "Test data from App repository",
	}, nil
}

// GetAdditionalData implements the AppRepository interface
func (r *MockAppRepository) GetAdditionalData() (*dto.AppData, error) {
	// Return mock data
	return &dto.AppData{
		ID:      "app-456",
		Message: "Additional data from App repository",
	}, nil
}
