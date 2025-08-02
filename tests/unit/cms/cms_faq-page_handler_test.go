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

type MockCMSFaqPageService struct {
	mock.Mock
}

func (m *MockCMSFaqPageService) CreateFaqPage(faqPage *models.FaqPage) (*models.FaqPage, error) {
	args := m.Called(faqPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqPage), args.Error(1)
}

func (m *MockCMSFaqPageService) FindFaqPages(rawQuery string, sort string, page, limit int, language string) ([]models.FaqPage, int64, error) {
	args := m.Called(rawQuery, sort, page, limit, language)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}	
	return args.Get(0).([]models.FaqPage), args.Get(1).(int64), args.Error(2)
}

func (m *MockCMSFaqPageService) FindFaqPageById(id uuid.UUID) (*models.FaqPage, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqPage), args.Error(1)
}

func (m *MockCMSFaqPageService) UpdateFaqContent(updatedFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error) {
	args := m.Called(updatedFaqContent, prevContentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqContent), args.Error(1)
}

func (m *MockCMSFaqPageService) DeleteFaqPage(id uuid.UUID) error {
	args := m.Called(id)	
	return args.Error(0)
}

func (m *MockCMSFaqPageService) FindContentByFaqPageId(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error) {
	args := m.Called(pageId, language, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqContent), args.Error(1)
}

func (m *MockCMSFaqPageService) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.FaqContent, error) {
	args := m.Called(pageId, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqContent), args.Error(1)
}

func (m *MockCMSFaqPageService) DeleteContentByFaqPageId(pageId uuid.UUID, language, mode string) error {
	args := m.Called(pageId, language, mode)
	return args.Error(0)
}

func (m *MockCMSFaqPageService) DuplicateFaqPage(pageId uuid.UUID) (*models.FaqPage, error) {
	args := m.Called(pageId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqPage), args.Error(1)
}

func (m *MockCMSFaqPageService) DuplicateFaqContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	args := m.Called(contentId, newRevision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqContent), args.Error(1)
}

func (m *MockCMSFaqPageService) RevertFaqContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	args := m.Called(revisionId, newRevision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FaqContent), args.Error(1)
}

func (m *MockCMSFaqPageService) FindCategories(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	args := m.Called(pageId, categoryTypeCode, language, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}	
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCMSFaqPageService) FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error) {
	args := m.Called(pageId, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}	
	return args.Get(0).([]models.Revision), args.Error(1)
}

func (m *MockCMSFaqPageService) PreviewFaqContent(pageId uuid.UUID, faqContentPreview *models.FaqContent) (string, error) {
	args := m.Called(pageId, faqContentPreview)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}	
	return args.Get(0).(string), args.Error(1)
}

func TestCMSFaqHandler(t *testing.T) {
	mockService := &MockCMSFaqPageService{}
	handler := cmsHandler.NewCMSFaqPageHandler(mockService)

	app := fiber.New()
	app.Post("/cms/faqpages", handler.HandleCreateFaqPage)
	app.Get("/cms/faqpages", handler.HandleGetFaqPages)
	app.Get("/cms/faqpages/:pageId", handler.HandleGetFaqPageById)
	app.Delete("/cms/faqpages/:pageId", handler.HandleDeleteFaqPage)
	app.Get("/cms/faqpages/:pageId/contents/:languageCode", handler.HandleGetContentByFaqPageId)
	app.Get("/cms/faqpages/:pageId/:languageCode/contents", handler.HandleGetLatestContentByFaqPageId)
	app.Delete("/cms/faqpages/:pageId/contents/:languageCode", handler.HandleDeleteFaqContentByPageId)
	app.Post("/cms/faqpages/duplicate/:pageId/pages", handler.HandleDuplicateFaqPage)
	app.Post("/cms/faqpages/duplicate/:contentId/contents", handler.HandleDuplicateFaqContentToAnotherLanguage)
	app.Post("/cms/faqpages/:revisionId/revision", handler.HandleRevertFaqContent)
	app.Put("/cms/faqpages/:contentId/contents", handler.HandleUpdateFaqContent)
	app.Get("/cms/faqpages/category/:categoryTypeCode/:pageId/:languageCode", handler.HandleGetCategory)
	app.Get("/cms/faqpages/revisions/:languageCode/:pageId", handler.HandleGetRevisions)
	app.Post("/cms/faqpages/previews/:pageId", handler.HandlePreviewFaqContent)

	t.Run("POST /cms/faqpages HandleCreateFaqPage", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		body, err := json.Marshal(mockFaqPage)
		require.NoError(t, err)			

		t.Run("successfully create faq page", func(t *testing.T) {
			mockService.On("CreateFaqPage", mock.AnythingOfType("*models.FaqPage")).Return(mockFaqPage, nil)		

			req := httptest.NewRequest("POST", "/cms/faqpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to create faq page: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateFaqPage", mock.AnythingOfType("*models.FaqPage")).Return(mockFaqPage, nil)
			
			body, err := json.Marshal("invalid body")
			require.NoError(t, err)			

			req := httptest.NewRequest("POST", "/cms/faqpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to create faq page: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateFaqPage", mock.AnythingOfType("*models.FaqPage")).Return(nil, errs.ErrInternalServerError)	

			req := httptest.NewRequest("POST", "/cms/faqpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})			
	})

	t.Run("GET /cms/faqpages HandleGetFaqPages", func(t *testing.T) {
		rawQuery := `{"category_faq":"superfaq","category_keywords":"KeyYY222"}`
		escapedQuery := url.QueryEscape(rawQuery)
		sort := "title"
		page := 2
		limit := 5
		language := "en"

		expectedFaqPages := []models.FaqPage{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}		

		t.Run("successfully get faq page", func(t *testing.T) {
			mockService.On("FindFaqPages", rawQuery, sort, page, limit, language).Return(expectedFaqPages, int64(7), nil)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/faqpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed to get faq page: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindFaqPages", rawQuery, sort, page, limit, language).Return(nil, int64(0), errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/faqpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
		
		t.Run("failed to get faq page: invalid query", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindFaqPages", rawQuery, sort, page, limit, language).Return(nil, int64(0), errs.ErrInvalidQuery)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/faqpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})

	t.Run("GET /cms/faqpages/:pageId HandleGetFaqPageById", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		pageId := uuid.New()

		t.Run("successfully get faq page by id", func(t *testing.T) {
			mockService.On("FindFaqPageById", pageId).Return(mockFaqPage, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})

		t.Run("failed get faq page by id: invalid pageId", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindFaqPageById", pageId).Return(mockFaqPage, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s", invalidPageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})	
		
		t.Run("failed get faq page by id: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindFaqPageById", pageId).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})			
	})

	t.Run("DELETE /cms/faqpages/:pageId HandleDeleteFaqPage", func(t *testing.T)	{
		pageId := uuid.New()

		t.Run("successfully delete faq page", func(t *testing.T)	{
			mockService.On("DeleteFaqPage", pageId).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/faqpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})
		
		t.Run("failed to delete faq page: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteFaqPage", pageId).Return(nil)
			invalidPageId := "1"

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/faqpages/%s", invalidPageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})	
		
		t.Run("failed to delete faq page: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteFaqPage", pageId).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/faqpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})

	t.Run("GET /cms/faqpages/:pageId/contents/:languageCode HandleGetContentByFaqPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqContent := mockFaqPage.Contents[0]

		t.Run("successfully get content by faq pageId", func(t *testing.T)	{
			mockService.On("FindContentByFaqPageId", pageId, language, mode).Return(mockFaqContent, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get content by faq pageId: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindContentByFaqPageId", pageId, language, mode).Return(mockFaqContent, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s/contents/%s?mode=%s", invalidPageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get content by faq pageId: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindContentByFaqPageId", pageId, language, mode).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/faqpages/:pageId/:languageCode/contents HandleGetLatestContentByFaqPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqContent := mockFaqPage.Contents[0]

		t.Run("successfully get latest content", func(t *testing.T)	{
			mockService.On("FindLatestContentByPageId", pageId, language).Return(mockFaqContent, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s/%s/contents", pageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get latest content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindLatestContentByPageId", pageId, language).Return(mockFaqContent, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s/%s/contents", invalidPageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindLatestContentByPageId", pageId, language).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/%s/%s/contents", pageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("DELETE /cms/faqpages/:pageId/contents/:languageCode HandleDeleteFaqContentByPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		t.Run("successfully revert faq content", func(t *testing.T)	{
			mockService.On("DeleteContentByFaqPageId", pageId, language, mode).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/faqpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get latest content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContentByFaqPageId", pageId, language, mode).Return(nil)
			invalidPageId := "1"

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/faqpages/%s/contents/%s?mode=%s", invalidPageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContentByFaqPageId", pageId, language, mode).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/faqpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})		

	t.Run("POST /cms/faqpages/duplicate/:pageId/pages HandleDuplicateFaqPage", func(t *testing.T)	{
		pageId := uuid.New()	
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully duplicate faq page", func(t *testing.T)	{
			mockService.On("DuplicateFaqPage", pageId).Return(mockFaqPage, nil)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/pages", pageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to duplicate faq page: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateFaqPage", pageId).Return(mockFaqPage, nil)
			invalidPageId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/pages", invalidPageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to duplicate faq page: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateFaqPage", pageId).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/pages", pageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	
	
	t.Run("POST /cms/faqpages/duplicate/:contentId/contents HandleDuplicateFaqContentToAnotherLanguage", func(t *testing.T)	{
		contentId := uuid.New()	
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully duplicate faq content to another language", func(t *testing.T)	{
			mockService.On("DuplicateFaqContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to duplicate faq content to another language: invalid body", func(t *testing.T)	{
			mockService.On("DuplicateFaqContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to duplicate faq content to another language: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateFaqContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)
			invalidContentId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/contents", invalidContentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicateFaqContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})		

	t.Run("POST /cms/faqpages/:revisionId/revision HandleRevertFaqContent", func(t *testing.T)	{
		revisionId := uuid.New()	
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully revert faq content", func(t *testing.T)	{
			mockService.On("RevertFaqContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to revert faq content: invalid body", func(t *testing.T)	{
			mockService.On("RevertFaqContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to revert faq content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("RevertFaqContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)
			invalidRevisionId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/%s/revision", invalidRevisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to revert faq content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("RevertFaqContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("PUT /cms/faqpages/:contentId/contents HandleUpdateFaqContent", func(t *testing.T)	{
		contentId := uuid.New()	
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		body, err := json.Marshal(mockContent)
		require.NoError(t, err)			

		t.Run("successfully update faq content", func(t *testing.T)	{
			mockService.On("UpdateFaqContent", mock.AnythingOfType("*models.FaqContent"), contentId).Return(mockContent, nil)				

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/faqpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to update faq content: invalid body", func(t *testing.T)	{
			mockService.On("UpdateFaqContent", mock.AnythingOfType("*models.FaqContent"), contentId).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/faqpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to update faq content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("UpdateFaqContent", mock.AnythingOfType("*models.FaqContent"), contentId).Return(mockContent, nil)
			invalidcontentId := "1"		

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/faqpages/%s/contents", invalidcontentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to update faq content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("UpdateFaqContent", mock.AnythingOfType("*models.FaqContent"), contentId).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/faqpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})
	
	t.Run("GET /cms/faqpages/category/:categoryTypeCode/:pageId/:languageCode HandleGetCategory", func(t *testing.T)	{
		pageId := uuid.New()	
		categoryTypeCode := "TYPE_CODE"
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		category := mockContent.Categories[0]	
		
		expectedCategories := []models.Category{*category, *category}

		t.Run("successfully get categories", func(t *testing.T)	{
			mockService.On("FindCategories", pageId, categoryTypeCode, language, mode).Return(expectedCategories, nil)				

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/category/%s/%s/%s?mode=%s", categoryTypeCode, pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to get categories: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindCategories", pageId, categoryTypeCode, language, mode).Return(expectedCategories, nil)
			invalidpageId := "1"		

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/category/%s/%s/%s?mode=%s", categoryTypeCode, invalidpageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get categories: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindCategories", pageId, categoryTypeCode, language, mode).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/category/%s/%s/%s?mode=%s", categoryTypeCode,  pageId, language, mode), nil)		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/faqpages/revisions/:languageCode/:pageId HandleGetRevisions", func(t *testing.T)	{
		pageId := uuid.New()	
		language := string(enums.PageLanguageEN)
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		revision := mockContent.Revision
		expectedRevisions := []models.Revision{*revision, *revision}
		
		t.Run("successfully get revisions", func(t *testing.T)	{
			mockService.On("FindRevisions", pageId, language).Return(expectedRevisions, nil)				

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/revisions/%s/%s", language, pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to get revisions: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindRevisions", pageId, language).Return(expectedRevisions, nil)
			invalidpageId := "1"		

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/revisions/%s/%s", language, invalidpageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get revisions: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindRevisions", pageId, language).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/faqpages/revisions/%s/%s", language, pageId), nil)		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})		

	t.Run("POST /cms/faqpages/previews/:pageId HandlePreviewFaqContent", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		pageId := uuid.New()
		url := "fake url"

		body, err := json.Marshal(mockContent)
		require.NoError(t, err)				
		
		t.Run("successfully preview faq content", func(t *testing.T)	{
			mockService.On("PreviewFaqContent", pageId, mock.AnythingOfType("*models.FaqContent")).Return(url, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/previews/%s", pageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})	
		
		t.Run("failed to preview faq content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("PreviewFaqContent", pageId, mock.AnythingOfType("*models.FaqContent")).Return("", errs.ErrInternalServerError)		
			
			invalidPageId := "1"	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/previews/%s", invalidPageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to preview faq content: invalid body", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("PreviewFaqContent", pageId, mock.AnythingOfType("*models.FaqContent")).Return("", errs.ErrInternalServerError)		
			
			body, err := json.Marshal("invalid body")
			require.NoError(t, err)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/previews/%s", pageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to preview faq content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("PreviewFaqContent", pageId, mock.AnythingOfType("*models.FaqContent")).Return("", errs.ErrInternalServerError)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/faqpages/previews/%s", pageId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
	})
}