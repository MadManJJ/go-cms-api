package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	appHandler "github.com/MadManJJ/cms-api/handlers/app"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAppLandingPageService struct {
	mock.Mock
}

func (m *MockAppLandingPageService) GetLandingPageByUrlAlias(urlAlias string, selectParam string, language string) (*models.LandingPage, error) {
	args := m.Called(urlAlias, selectParam, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingPage), args.Error(1)	
}

func (m *MockAppLandingPageService) GetLandingContentPreview(id uuid.UUID) (*models.LandingContent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingContent), args.Error(1)	
}

func TestAppLandingHandler(t *testing.T) {
	mockService := &MockAppLandingPageService{}
	handler := appHandler.NewAppLandingPageHandler(mockService)

	app := fiber.New()
	app.Get("/app/landingpages/:languageCode/by-alias", handler.HandleGetLandingPageByUrlAlias)

	t.Run("GET /app/landingpages/:languageCode/by-alias HandleGetLandingPageByAlias", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		urlAlias := "about/us"
		selectParam := "revision,categories,components,metatag"
		language := string(enums.PageLanguageEN)

		t.Run("successfully get landing page (url_alias)", func(t *testing.T) {
			mockService.On("GetLandingPageByUrlAlias", urlAlias, selectParam, language).Return(mockLandingPage, nil)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/landingpages/%s/by-alias?url_alias=%s&select=%s", language, urlAlias, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to get landing page: internal server error (url_alias)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetLandingPageByUrlAlias", urlAlias, selectParam, language).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/landingpages/%s/by-alias?url_alias=%s&select=%s", language, urlAlias, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to get landing page: not found (url_alias)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetLandingPageByUrlAlias", urlAlias, selectParam, language).Return(nil, errs.ErrNotFound)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/landingpages/%s/by-alias?url_alias=%s&select=%s", language, urlAlias, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
			mockService.AssertExpectations(t)
		})					
	})
}