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

type MockAppPartnerPageService struct {
	mock.Mock
}

func (m *MockAppPartnerPageService) GetPartnerPage(slug string, isAlias bool, selectParam string, language string) (*models.PartnerPage, error) {
	args := m.Called(slug, isAlias, selectParam, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerPage), args.Error(1)	
}	

func (m *MockAppPartnerPageService) GetPartnerContentPreview(id uuid.UUID) (*models.PartnerContent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerContent), args.Error(1)	
}

func TestAppHandlerPartnerHandler(t *testing.T) {
	mockService := &MockAppPartnerPageService{}
	handler := appHandler.NewAppPartnerPageHandler(mockService)

	app := fiber.New()
	app.Get("/app/partnerpages/:languageCode/by-alias", handler.HandleGetPartnerPageByAlias)
	app.Get("/app/partnerpages/:languageCode/by-url", handler.HandleGetPartnerPageByUrl)

	t.Run("GET /app/partnerpages/:languageCode/by-alias HandleGetPartnerPageByAlias", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		slug := "about/us"
		selectParam := "revision,categories,components,metatag"
		language := string(enums.PageLanguageEN)

		t.Run("successfully get partner page (url_alias)", func(t *testing.T) {
			isAlias := true
			mockService.On("GetPartnerPage", slug, isAlias, selectParam, language).Return(mockPartnerPage, nil)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/partnerpages/%s/by-alias?url_alias=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to get partner page: internal server error (url_alias)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := true
			mockService.On("GetPartnerPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/partnerpages/%s/by-alias?url_alias=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to get partner page: not found (url_alias)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := true
			mockService.On("GetPartnerPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrNotFound)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/partnerpages/%s/by-alias?url_alias=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
			mockService.AssertExpectations(t)
		})					
	})

	t.Run("GET /app/partnerpages/:languageCode/by-urk HandleGetPartnerPageByUrl", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		slug := "about/us"
		selectParam := "revision,categories,components,metatag"
		language := string(enums.PageLanguageEN)
				
		t.Run("successfully get partner page (url)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := false
			mockService.On("GetPartnerPage", slug, isAlias, selectParam, language).Return(mockPartnerPage, nil)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/partnerpages/%s/by-url?url=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed get partner page: internal server error (url)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := false
			mockService.On("GetPartnerPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/partnerpages/%s/by-url?url=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed get partner page: not found (url)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := false
			mockService.On("GetPartnerPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrNotFound)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/partnerpages/%s/by-url?url=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
	})
}