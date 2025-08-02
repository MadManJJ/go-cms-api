package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLandingService struct {
	mock.Mock
}

func (m *MockLandingService) CreateLandingPage(landingPage *models.LandingPage) (*models.LandingPage, error) {
	args := m.Called(landingPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingPage), args.Error(1)
}

func (m *MockLandingService) FindLandingPages(rawQuery string, sort string, page, limit int, language string) ([]models.LandingPage, int64, error) {
	args := m.Called(rawQuery, sort, page, limit, language)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.LandingPage), args.Get(1).(int64), args.Error(2)
}

func (m *MockLandingService) FindLandingPageById(id uuid.UUID) (*models.LandingPage, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingPage), args.Error(1)
}

func (m *MockLandingService) UpdateLandingPage(updatedLandingPage *models.LandingPage) (*models.LandingPage, error) {
	args := m.Called(updatedLandingPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingPage), args.Error(1)
}

func (m *MockLandingService) UpdateLandingContent(updatedLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error) {
	args := m.Called(updatedLandingContent, prevContentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingContent), args.Error(1)
}

func (m *MockLandingService) DeleteLandingPage(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLandingService) FindContentByLandingPageId(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error) {
	args := m.Called(pageId, language, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingContent), args.Error(1)
}

func (m *MockLandingService) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.LandingContent, error) {
	args := m.Called(pageId, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingContent), args.Error(1)
}

func (m *MockLandingService) DeleteContentByLandingPageId(pageId uuid.UUID, language, mode string) error {
	args := m.Called(pageId, language, mode)
	return args.Error(0)
}

func (m *MockLandingService) DuplicateLandingPage(pageId uuid.UUID) (*models.LandingPage, error) {
	args := m.Called(pageId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingPage), args.Error(1)
}

func (m *MockLandingService) DuplicateLandingContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	args := m.Called(contentId, newRevision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingContent), args.Error(1)
}

func (m *MockLandingService) RevertLandingContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	args := m.Called(revisionId, newRevision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LandingContent), args.Error(1)
}

func (m *MockLandingService) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	args := m.Called(pageId, categoryTypeCode, language, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockLandingService) FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error) {
	args := m.Called(pageId, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Revision), args.Error(1)
}

func (m *MockLandingService) PreviewLandingContent(pageId uuid.UUID, landingContentPreview *models.LandingContent) (string, error) {
	args := m.Called(pageId, landingContentPreview)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}	
	return args.Get(0).(string), args.Error(1)
}

func TestCMSLandingHandler(t *testing.T) {
	mockService := &MockLandingService{}
	handler := cmsHandler.NewCMSLandingPageHandler(mockService)

	app := fiber.New()
	app.Post("/cms/landingpages", handler.HandleCreateLandingPage)
	app.Get("/cms/landingpages", handler.HandleGetLandingPages)
	app.Get("/cms/landingpages/:pageId", handler.HandleGetLandingPageById)
	app.Delete("/cms/landingpages/:pageId", handler.HandleDeleteLandingPage)
	app.Get("/cms/landingpages/:pageId/contents/:languageCode", handler.HandleGetContentByLandingPageId)
	app.Get("/cms/landingpages/:pageId/:languageCode/contents", handler.HandleGetLatestContentByLandingPageId)
	app.Delete("/cms/landingpages/:pageId/contents/:languageCode", handler.HandleDeleteLandingContentByPageId)
	app.Post("/cms/landingpages/duplicate/:pageId", handler.HandleDuplicateLandingPage)
	app.Post("/cms/landingpages/duplicate/:contentId/contents", handler.HandleDuplicateLandingContentToAnotherLanguage)	
	app.Post("/cms/landingpages/:revisionId/revision", handler.HandleRevertLandingContent)
	app.Put("/cms/landingpages/:contentId/contents", handler.HandleUpdateLandingContent)
	app.Get("/cms/landingpages/category/:categoryTypeCode/:pageId/:languageCode", handler.HandleGetCategory)
	app.Get("/cms/landingpages/revisions/:languageCode/:pageId", handler.HandleGetRevisions)
	
	t.Run("POST /cms/landingpages HandleCreateLandingPage", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		body, err := json.Marshal(mockLandingPage)
		require.NoError(t, err)			

		t.Run("successfully create landing page", func(t *testing.T) {
			mockService.On("CreateLandingPage", mock.AnythingOfType("*models.LandingPage")).Return(mockLandingPage, nil)		

			req := httptest.NewRequest("POST", "/cms/landingpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to create landing page: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateLandingPage", mock.AnythingOfType("*models.LandingPage")).Return(mockLandingPage, nil)
			
			body, err := json.Marshal("invalid body")
			require.NoError(t, err)			

			req := httptest.NewRequest("POST", "/cms/landingpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to create landing page: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateLandingPage", mock.AnythingOfType("*models.LandingPage")).Return(nil, errs.ErrInternalServerError)	

			req := httptest.NewRequest("POST", "/cms/landingpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})			
	})

	t.Run("GET /cms/landingpages HandleGetLandingPages", func(t *testing.T) {
		rawQuery := `{
			"title": "Reset Password",
			"category_keywords": "security",
			"status": "Published"
		}`
		escapedQuery := url.QueryEscape(rawQuery)
		sort := "title"
		page := 2
		limit := 5
		language := "en"
		
		expectedLandingPages := []models.LandingPage{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}	
		
		t.Run("successfully get landing page", func(t *testing.T) {
			mockService.On("FindLandingPages", rawQuery, sort, page, limit, language).Return(expectedLandingPages, int64(2), nil)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/landingpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed to get landing page: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindLandingPages", rawQuery, sort, page, limit, language).Return(nil, int64(0), errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/landingpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
		
		t.Run("failed to get landing page: invalid query", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindLandingPages", rawQuery, sort, page, limit, language).Return(nil, int64(0), errs.ErrInvalidQuery)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/landingpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})	

	t.Run("GET /cms/landingpages/:pageId HandleGetLandingPageById", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		pageId := uuid.New()		

		t.Run("successfully get landing page by id", func(t *testing.T) {
			mockService.On("FindLandingPageById", pageId).Return(mockLandingPage, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})		

		t.Run("failed get landing page by id: invalid pageId", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindLandingPageById", pageId).Return(mockLandingPage, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s", invalidPageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})	
		
		t.Run("failed get landing page by id: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindLandingPageById", pageId).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})			
	})

	t.Run("DELETE /cms/landingpages/:pageId HandleDeleteLandingPage", func(t *testing.T)	{
		pageId := uuid.New()

		t.Run("successfully delete landing page", func(t *testing.T)	{
			mockService.On("DeleteLandingPage", pageId).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/landingpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})
		
		t.Run("failed to delete landing page: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteLandingPage", pageId).Return(nil)
			invalidPageId := "1"

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/landingpages/%s", invalidPageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})	
		
		t.Run("failed to delete landing page: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteLandingPage", pageId).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/landingpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/landingpages/:pageId/contents/:languageCode HandleGetContentByLandingPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockLandingContent := mockLandingPage.Contents[0]

		t.Run("successfully get content by landing pageId", func(t *testing.T)	{
			mockService.On("FindContentByLandingPageId", pageId, language, mode).Return(mockLandingContent, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get content by landing pageId: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindContentByLandingPageId", pageId, language, mode).Return(mockLandingContent, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s/contents/%s?mode=%s", invalidPageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get content by landing pageId: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindContentByLandingPageId", pageId, language, mode).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})		

	t.Run("GET /cms/landingpages/:pageId/contents HandleGetLatestContentByLandingPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockLandingContent := mockLandingPage.Contents[0]

		t.Run("successfully get latest content", func(t *testing.T)	{
			mockService.On("FindLatestContentByPageId", pageId, language).Return(mockLandingContent, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s/%s/contents", pageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get latest content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindLatestContentByPageId", pageId, language).Return(mockLandingContent, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s/%s/contents", invalidPageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindLatestContentByPageId", pageId, language).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/%s/%s/contents", pageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	
	
	t.Run("DELETE /cms/landingpages/:pageId/contents/:languageCode HandleDeleteLandingContentByPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		t.Run("successfully revert landing content", func(t *testing.T)	{
			mockService.On("DeleteContentByLandingPageId", pageId, language, mode).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/landingpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get latest content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContentByLandingPageId", pageId, language, mode).Return(nil)
			invalidPageId := "1"

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/landingpages/%s/contents/%s?mode=%s", invalidPageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContentByLandingPageId", pageId, language, mode).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/landingpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})

	t.Run("POST /cms/landingpages/duplicate/:pageId HandleDuplicateLandingPage", func(t *testing.T)	{
		pageId := uuid.New()	
		mockLandingPage := helpers.InitializeMockLandingPage()

		t.Run("successfully duplicate landing page", func(t *testing.T)	{
			mockService.On("DuplicateLandingPage", pageId).Return(mockLandingPage, nil)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s", pageId), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to duplicate landing page: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateLandingPage", pageId).Return(mockLandingPage, nil)
			invalidPageId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s", invalidPageId), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to duplicate landing page: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateLandingPage", pageId).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s", pageId), nil)
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("POST /cms/landingpages/duplicate/:contentId/contents HandleDuplicateLandingContentToAnotherLanguage", func(t *testing.T)	{
		contentId := uuid.New()	
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully duplicate landing content to another language", func(t *testing.T)	{
			mockService.On("DuplicateLandingContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to duplicate landing content to another language: invalid body", func(t *testing.T)	{
			mockService.On("DuplicateLandingContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to duplicate landing content to another language: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateLandingContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)
			invalidContentId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s/contents", invalidContentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateLandingContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})			

	t.Run("POST /cms/landingpages/:revisionId/revision HandleRevertLandingContent", func(t *testing.T)	{
		revisionId := uuid.New()	
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully revert landing content", func(t *testing.T)	{
			mockService.On("RevertLandingContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to revert landing content: invalid body", func(t *testing.T)	{
			mockService.On("RevertLandingContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to revert landing content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("RevertLandingContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)
			invalidRevisionId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/%s/revision", invalidRevisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to revert landing content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("RevertLandingContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/landingpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})
	
	t.Run("PUT /cms/landingpages/:contentId/contents HandleUpdateLandingContent", func(t *testing.T)	{
		contentId := uuid.New()	
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]

		body, err := json.Marshal(mockContent)
		require.NoError(t, err)			

		t.Run("successfully update landing content", func(t *testing.T)	{
			mockService.On("UpdateLandingContent", mock.AnythingOfType("*models.LandingContent"), contentId).Return(mockContent, nil)				

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/landingpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to update landing content: invalid body", func(t *testing.T)	{
			mockService.On("UpdateLandingContent", mock.AnythingOfType("*models.LandingContent"), contentId).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/landingpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to update landing content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("UpdateLandingContent", mock.AnythingOfType("*models.LandingContent"), contentId).Return(mockContent, nil)
			invalidcontentId := "1"		

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/landingpages/%s/contents", invalidcontentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to update landing content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("UpdateLandingContent", mock.AnythingOfType("*models.LandingContent"), contentId).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/landingpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/landingpages/category/:categoryTypeCode/:pageId/:languageCode HandleGetCategory", func(t *testing.T)	{
		pageId := uuid.New()	
		categoryTypeCode := "TYPE_CODE"
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		category := mockContent.Categories[0]	
		
		expectedCategories := []models.Category{*category, *category}

		t.Run("successfully get categories", func(t *testing.T)	{
			mockService.On("GetCategory", pageId, categoryTypeCode, language, mode).Return(expectedCategories, nil)				

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/category/%s/%s/%s?mode=%s", categoryTypeCode, pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to get categories: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("GetCategory", pageId, categoryTypeCode, language, mode).Return(expectedCategories, nil)
			invalidpageId := "1"		

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/category/%s/%s/%s?mode=%s", categoryTypeCode, invalidpageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get categories: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("GetCategory", pageId, categoryTypeCode, language, mode).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/category/%s/%s/%s?mode=%s", categoryTypeCode,  pageId, language, mode), nil)		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/landingpages/revisions/:languageCode/:pageId HandleGetRevisions", func(t *testing.T)	{
		pageId := uuid.New()	
		language := string(enums.PageLanguageEN)
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		revision := mockContent.Revision
		expectedRevisions := []models.Revision{*revision, *revision}
		
		t.Run("successfully get revisions", func(t *testing.T)	{
			mockService.On("FindRevisions", pageId, language).Return(expectedRevisions, nil)				

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/revisions/%s/%s", language, pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to get revisions: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindRevisions", pageId, language).Return(expectedRevisions, nil)
			invalidpageId := "1"		

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/revisions/%s/%s", language, invalidpageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get revisions: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindRevisions", pageId, language).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/landingpages/revisions/%s/%s", language, pageId), nil)		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	
}