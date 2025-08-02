package services

import (
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/repositories"
)

// CMSService handles business logic for CMS domain
type CMSService struct {
	repo repositories.CMSRepository
}

// NewCMSService creates a new instance of CMSService
func NewCMSService(repo repositories.CMSRepository) *CMSService {
	return &CMSService{
		repo: repo,
	}
}

// GetTestData retrieves test data from repository
func (s *CMSService) GetTestData() (*dto.CMSData, error) {
	return s.repo.GetTestData()
}

// GetAdditionalData retrieves additional data from repository
func (s *CMSService) GetAdditionalData() (*dto.CMSData, error) {
	return s.repo.GetAdditionalData()
}