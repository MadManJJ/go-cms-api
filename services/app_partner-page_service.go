package services

import (
	"strings"

	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
)

type AppPartnerPageServiceInterface interface {
	GetPartnerPage(slug string, isAlias bool, selectParam string, language string) (*models.PartnerPage, error)
	GetPartnerContentPreview(id uuid.UUID) (*models.PartnerContent, error)
}

type AppPartnerPageService struct {
	repo repositories.AppPartnerPageRepositoryInterface
}

func NewAppPartnerPageService(repo repositories.AppPartnerPageRepositoryInterface) *AppPartnerPageService {
	return &AppPartnerPageService{repo: repo}
}

func (s *AppPartnerPageService) GetPartnerPage(slug string, isAlias bool, selectParam string, language string) (*models.PartnerPage, error) {
	var preloads []string

	var preloadMap = map[string]string{
  // User input         // Actual Preload
		"revision":           "Contents.Revision",
		"categories":         "Contents.Categories",
		"components":         "Contents.Components",
		"metatag":            "Contents.MetaTag",
	}	

	// Split the selectParam by comma and trim spaces and check if they exist in the preloadMap
	// If they do, append them to the preloads slice
	for _, item := range strings.Split(selectParam, ",") {
		if preload, ok := preloadMap[strings.ToLower(strings.TrimSpace(item))]; ok {
			preloads = append(preloads, preload)
		}
	}	

	result, err := s.repo.GetPartnerPageBySlug(slug, preloads, isAlias, language)
	if err != nil {
		return nil, err
	}	

	return result, nil
}

func (s *AppPartnerPageService) GetPartnerContentPreview(id uuid.UUID) (*models.PartnerContent, error) {
	partnerContent, err := s.repo.GetPartnerContentPreview(id)
	if err != nil {
		return nil, err
	}

	return partnerContent, nil
}