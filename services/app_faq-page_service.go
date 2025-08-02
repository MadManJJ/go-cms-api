package services

import (
	"strings"

	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
)

type AppFaqPageServiceInterface interface {
	GetFaqPage(slug string, isAlias bool, selectParam string, language string) (*models.FaqPage, error)
	GetFaqContentPreview(id uuid.UUID) (*models.FaqContent, error)
}

type AppFaqPageService struct {
	repo repositories.AppFaqPageRepositoryInterface
}

func NewAppFaqPageService(repo repositories.AppFaqPageRepositoryInterface) *AppFaqPageService {
	return &AppFaqPageService{repo: repo}
}

func (s *AppFaqPageService) GetFaqPage(slug string, isAlias bool, selectParam string, language string) (*models.FaqPage, error) {
	var preloads []string

	var preloadMap = map[string]string{
  // User input         // Actual Preload
		"revision":           "Contents.Revision",
		"categories":         "Contents.Categories",
		"components":         "Contents.Components",
		"metatag":            "Contents.MetaTag",
	}		

	for _, item := range strings.Split(selectParam, ",") {
		if preload, ok := preloadMap[strings.ToLower(strings.TrimSpace(item))]; ok {
			preloads = append(preloads, preload)
		}
	}	

	result, err := s.repo.GetFaqPageBySlug(slug, preloads, isAlias, language)
	if err != nil {
		return nil, err
	}	

	return result, nil	
}

func (s *AppFaqPageService) GetFaqContentPreview(id uuid.UUID) (*models.FaqContent, error) {
	faqContent, err := s.repo.GetFaqContentPreview(id)
	if err != nil {
		return nil, err
	}

	return faqContent, nil
}