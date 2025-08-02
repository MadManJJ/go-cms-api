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

type MockAppFaqPageService struct {
	mock.Mock
}

func (m *MockAppFaqPageService) GetFaqPage(slug string, isAlias bool, selectParam string, language string) (*models.FaqPage, error) {
	args := m.Called(slug, isAlias, selectParam, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqPage), args.Error(1)	
}

func (m *MockAppFaqPageService) GetFaqContentPreview(id uuid.UUID) (*models.FaqContent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqContent), args.Error(1)	
}

func TestAppFaqHandler(t *testing.T) {
	mockService := &MockAppFaqPageService{}
	handler := appHandler.NewAppFaqPageHandler(mockService)

	app := fiber.New()
	app.Get("/app/faqpages/:languageCode/by-alias", handler.HandleGetFaqPageByAlias)
	app.Get("/app/faqpages/:languageCode/by-url", handler.HandleGetFaqPageByUrl)
	app.Get("/app/faqpages/previews/:id", handler.HandleGetFaqContentPreview)

	t.Run("GET /app/faqpages/:languageCode/by-alias HandleGetFaqPageByAlias", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		slug := "about/us"
		selectParam := "revision,categories,components,metatag"
		language := string(enums.PageLanguageEN)

		t.Run("successfully get faq page (url_alias)", func(t *testing.T) {
			isAlias := true
			mockService.On("GetFaqPage", slug, isAlias, selectParam, language).Return(mockFaqPage, nil)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/%s/by-alias?url_alias=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to get faq page: internal server error (url_alias)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := true
			mockService.On("GetFaqPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/%s/by-alias?url_alias=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to get faq page: not found (url_alias)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := true
			mockService.On("GetFaqPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrNotFound)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/%s/by-alias?url_alias=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
			mockService.AssertExpectations(t)
		})					
	})

	t.Run("GET /app/faqpages/:languageCode/by-urk HandleGetFaqPageByUrl", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		slug := "about/us"
		selectParam := "revision,categories,components,metatag"
		language := string(enums.PageLanguageEN)
				
		t.Run("successfully get faq page (url)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := false
			mockService.On("GetFaqPage", slug, isAlias, selectParam, language).Return(mockFaqPage, nil)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/%s/by-url?url=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed get faq page: internal server error (url)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := false
			mockService.On("GetFaqPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/%s/by-url?url=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed get faq page: not found (url)", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			isAlias := false
			mockService.On("GetFaqPage", slug, isAlias, selectParam, language).Return(nil, errs.ErrNotFound)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/%s/by-url?url=%s&select=%s", language, slug, selectParam), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
	})

	t.Run("GET /app/faqpages/previews/:id HandleGetFaqContentPreview", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqContent := mockFaqPage.Contents[0]
		contentId := mockFaqContent.ID

		t.Run("successfully get faq content preview", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFaqContentPreview", contentId).Return(mockFaqContent, nil)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/previews/%s", contentId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed get faq content preview: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFaqContentPreview", contentId).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("GET", fmt.Sprintf("/app/faqpages/previews/%s", contentId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})
}