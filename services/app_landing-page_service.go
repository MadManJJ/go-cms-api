package services

import (
	"strings"

	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
)

type AppLandingPageServiceInterface interface {
	GetLandingPageByUrlAlias(urlAlias string, selectParam string, language string) (*models.LandingPage, error)
	GetLandingContentPreview(id uuid.UUID) (*models.LandingContent, error)
}

type AppLandingPageService struct {
	repo repositories.AppLandingPageRepositoryInterface
}

func NewAppLandingPageService(repo repositories.AppLandingPageRepositoryInterface) *AppLandingPageService {
	return &AppLandingPageService{repo: repo}
}

func (s *AppLandingPageService) GetLandingPageByUrlAlias(urlAlias string, selectParam string, language string) (*models.LandingPage, error) {
	var preloads []string

	// Define a map of valid preloads
	var preloadMap = map[string]string{
  // User input         // Actual Preload
		"files":         			"Contents.Files",
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
	result, err := s.repo.GetLandingPageByUrlAlias(urlAlias, preloads, language)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AppLandingPageService) GetLandingContentPreview(id uuid.UUID) (*models.LandingContent, error) {
	landingContent, err := s.repo.GetLandingContentPreview(id)
	if err != nil {
		return nil, err
	}

	return landingContent, nil
}