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

type MockAppPartnerPageRepo struct {
	getPartnerPageBySlug func(slug string, preloads []string, isAlias bool, language string) (*models.PartnerPage, error)
	getPartnerContentPreview func(id uuid.UUID) (*models.PartnerContent, error)
}

func (m *MockAppPartnerPageRepo) GetPartnerPageBySlug(slug string, preloads []string, isAlias bool, language string) (*models.PartnerPage, error) {
	return m.getPartnerPageBySlug(slug, preloads, isAlias, language)
}

func (m *MockAppPartnerPageRepo) GetPartnerContentPreview(id uuid.UUID) (*models.PartnerContent, error) {
	return m.getPartnerContentPreview(id)
}

func TestAppService_GetPartnerPage(t *testing.T) {
	slug := "about/us"
	language := string(enums.PageLanguageEN)
	selectParam := "revision,categories,components,metatag"
	isAlias := true

	t.Run("successfully get partner page", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()

		repo := &MockAppPartnerPageRepo{
			getPartnerPageBySlug: func(slug string, preloads []string, isAlias bool, language string) (*models.PartnerPage, error) {
				return mockPartnerPage, nil
			},
		}

		service := services.NewAppPartnerPageService(repo)

		actualPartnerPage, err := service.GetPartnerPage(slug, isAlias, selectParam, language)
		assert.NoError(t, err)
		assert.Equal(t, mockPartnerPage, actualPartnerPage)
	})	

	t.Run("failed to get partner page", func(t *testing.T) {
		repo := &MockAppPartnerPageRepo{
			getPartnerPageBySlug: func(slug string, preloads []string, isAlias bool, language string) (*models.PartnerPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewAppPartnerPageService(repo)

		actualPartnerPage, err := service.GetPartnerPage(slug, isAlias, selectParam, language)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPage)
	})	
}