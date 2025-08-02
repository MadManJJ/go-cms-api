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

type MockAppFaqPageRepo struct {
	getFaqPageBySlug func(slug string, preloads []string, isAlias bool, language string) (*models.FaqPage, error)
	getFaqContentPreview func(id uuid.UUID) (*models.FaqContent, error)
}

func (m *MockAppFaqPageRepo) GetFaqPageBySlug(slug string, preloads []string, isAlias bool, language string) (*models.FaqPage, error) {
	return m.getFaqPageBySlug(slug, preloads, isAlias, language)
}

func (m *MockAppFaqPageRepo) GetFaqContentPreview(id uuid.UUID) (*models.FaqContent, error) {
	return m.getFaqContentPreview(id)
}

func TestAppService_GetFaqPage(t *testing.T) {
	slug := "about/us"
	language := string(enums.PageLanguageEN)
	selectParam := "revision,categories,components,metatag"
	isAlias := true

	t.Run("successfully get faq page", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()

		repo := &MockAppFaqPageRepo{
			getFaqPageBySlug: func(slug string, preloads []string, isAlias bool, language string) (*models.FaqPage, error) {
				return mockFaqPage, nil
			},
		}

		service := services.NewAppFaqPageService(repo)

		actualFaqPage, err := service.GetFaqPage(slug, isAlias, selectParam, language)
		assert.NoError(t, err)
		assert.Equal(t, mockFaqPage, actualFaqPage)
	})	

	t.Run("failed to get faq page", func(t *testing.T) {
		repo := &MockAppFaqPageRepo{
			getFaqPageBySlug: func(slug string, preloads []string, isAlias bool, language string) (*models.FaqPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewAppFaqPageService(repo)

		actualFaqPage, err := service.GetFaqPage(slug, isAlias, selectParam, language)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})	
}

func TestAppService_GetFaqContentPreview(t *testing.T) {
	contentId := uuid.New()

	t.Run("successfully get faq content preview", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqContent := mockFaqPage.Contents[0]

		repo := &MockAppFaqPageRepo{
			getFaqContentPreview: func(id uuid.UUID) (*models.FaqContent, error) {
				return mockFaqContent, nil
			},
		}

		service := services.NewAppFaqPageService(repo)

		actualFaqContent, err := service.GetFaqContentPreview(contentId)
		assert.NoError(t, err)
		assert.Equal(t, mockFaqContent, actualFaqContent)
	})	

	t.Run("failed to get faq content preview", func(t *testing.T) {
		repo := &MockAppFaqPageRepo{
			getFaqContentPreview: func(id uuid.UUID) (*models.FaqContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewAppFaqPageService(repo)

		actualFaqContent, err := service.GetFaqContentPreview(contentId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqContent)
	})	
}