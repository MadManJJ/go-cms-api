package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockAppLandingPageRepo struct {
	getLandingPageByUrlAlias func(urlAlias string, preloads []string, language string) (*models.LandingPage, error)
	getLandingContentPreview func(id uuid.UUID) (*models.LandingContent, error)
}

func (m *MockAppLandingPageRepo) GetLandingPageByUrlAlias(urlAlias string, preloads []string, language string) (*models.LandingPage, error) {
	return m.getLandingPageByUrlAlias(urlAlias, preloads, language)
}

func (m *MockAppLandingPageRepo) GetLandingContentPreview(id uuid.UUID) (*models.LandingContent, error) {
	return m.getLandingContentPreview(id)
}

func TestAppService_GetLandingPage(t *testing.T) {
	urlAlias := "about/us"
	language := string(enums.PageLanguageEN)
	selectParam := "revision,categories,components,metatag"
	
	t.Run("successfully get landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()

		repo := &MockAppLandingPageRepo{
			getLandingPageByUrlAlias: func(urlAlias string, preloads []string, language string) (*models.LandingPage, error) {
				return mockLandingPage, nil
			},
		}

		service := services.NewAppLandingPageService(repo)

		actualLandingPage, err := service.GetLandingPageByUrlAlias(urlAlias, selectParam, language)
		assert.NoError(t, err)
		assert.Equal(t, mockLandingPage, actualLandingPage)
	})	

	t.Run("failed to get landing page", func(t *testing.T) {
		repo := &MockAppLandingPageRepo{
			getLandingPageByUrlAlias: func(urlAlias string, preloads []string, language string) (*models.LandingPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewAppLandingPageService(repo)

		actualLandingPage, err := service.GetLandingPageByUrlAlias(urlAlias, selectParam, language)
		assert.Error(t, err)
		assert.Nil(t, actualLandingPage)
	})	
}