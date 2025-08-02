package services

import (
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/repositories"
)

// AppService handles business logic for app domain
type AppService struct {
	repo repositories.AppRepository
}

// NewAppService creates a new instance of AppService
func NewAppService(repo repositories.AppRepository) *AppService {
	return &AppService{
		repo: repo,
	}
}

// GetTestData retrieves test data from repository
func (s *AppService) GetTestData() (*dto.AppData, error) {
	return s.repo.GetTestData()
}

// GetAdditionalData retrieves additional data from repository
func (s *AppService) GetAdditionalData() (*dto.AppData, error) {
	return s.repo.GetAdditionalData()
}
