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

type MockPartnerService struct {
	mock.Mock
}

func (m *MockPartnerService) CreatePartnerPage(partnerPage *models.PartnerPage) (*models.PartnerPage, error) {
	args := m.Called(partnerPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerPage), args.Error(1)
}

func (m *MockPartnerService) FindPartnerPages(rawQuery string, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error) {
	args := m.Called(rawQuery, sort, page, limit, language)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.PartnerPage), args.Get(1).(int64), args.Error(2)
}

func (m *MockPartnerService) FindPartnerPageById(id uuid.UUID) (*models.PartnerPage, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerPage), args.Error(1)
}

func (m *MockPartnerService) UpdatePartnerContent(updatedPartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error) {
	args := m.Called(updatedPartnerContent, prevContentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerContent), args.Error(1)
}

func (m *MockPartnerService) DeletePartnerPage(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPartnerService) FindContentByPartnerPageId(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error) {
	args := m.Called(pageId, language, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerContent), args.Error(1)
}

func (m *MockPartnerService) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
	args := m.Called(pageId, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerContent), args.Error(1)
}

func (m *MockPartnerService) DeleteContentByPartnerPageId(pageId uuid.UUID, language, mode string) error {
	args := m.Called(pageId, language, mode)
	return args.Error(0)
}

func (m *MockPartnerService) DuplicatePartnerPage(pageId uuid.UUID) (*models.PartnerPage, error) {
	args := m.Called(pageId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerPage), args.Error(1)
}

func (m *MockPartnerService) DuplicatePartnerContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	args := m.Called(contentId, newRevision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerContent), args.Error(1)
}

func (m *MockPartnerService) RevertPartnerContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	args := m.Called(revisionId, newRevision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PartnerContent), args.Error(1)
}

func (m *MockPartnerService) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	args := m.Called(pageId, categoryTypeCode, language, mode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockPartnerService) FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error) {
	args := m.Called(pageId, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Revision), args.Error(1)
}

func (m *MockPartnerService) PreviewPartnerContent(pageId uuid.UUID, partnerContentPreview *models.PartnerContent) (string, error) {
	args := m.Called(pageId, partnerContentPreview)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}	
	return args.Get(0).(string), args.Error(1)
}

func TestCMSPartnerHandler(t *testing.T) {
	mockService := &MockPartnerService{}
	handler := cmsHandler.NewCMSPartnerPageHandler(mockService)

	app := fiber.New()
	app.Post("/cms/partnerpages", handler.HandleCreatePartnerPage)
	app.Get("/cms/partnerpages", handler.HandleGetPartnerPages)
	app.Get("/cms/partnerpages/:pageId", handler.HandleGetPartnerPageById)
	app.Delete("/cms/partnerpages/:pageId", handler.HandleDeletePartnerPage)
	app.Get("/cms/partnerpages/:pageId/contents/:languageCode", handler.HandleGetContentByPartnerPageId)
	app.Get("/cms/partnerpages/:pageId/:languageCode/contents", handler.HandleGetLatestContentByPartnerPageId)
	app.Delete("/cms/partnerpages/:pageId/contents/:languageCode", handler.HandleDeletePartnerContentByPageId)
	app.Post("/cms/partnerpages/duplicate/:pageId", handler.HandleDuplicatePartnerPage)
	app.Post("/cms/partnerpages/duplicate/:contentId/contents", handler.HandleDuplicatePartnerContentToAnotherLanguage)	
	app.Post("/cms/partnerpages/:revisionId/revision", handler.HandleRevertPartnerContent)
	app.Put("/cms/partnerpages/:contentId/contents", handler.HandleUpdatePartnerContent)
	app.Get("/cms/partnerpages/category/:categoryTypeCode/:pageId/:languageCode", handler.HandleGetCategory)
	app.Get("/cms/partnerpages/revisions/:languageCode/:pageId", handler.HandleGetRevisions)

	t.Run("POST /cms/partnerpages HandleCreatePartnerPage", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		body, err := json.Marshal(mockPartnerPage)
		require.NoError(t, err)			

		t.Run("successfully create partner page", func(t *testing.T) {
			mockService.On("CreatePartnerPage", mock.AnythingOfType("*models.PartnerPage")).Return(mockPartnerPage, nil)		

			req := httptest.NewRequest("POST", "/cms/partnerpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to create partner page: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreatePartnerPage", mock.AnythingOfType("*models.PartnerPage")).Return(mockPartnerPage, nil)
			
			body, err := json.Marshal("invalid body")
			require.NoError(t, err)			

			req := httptest.NewRequest("POST", "/cms/partnerpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		

		t.Run("failed to create partner page: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreatePartnerPage", mock.AnythingOfType("*models.PartnerPage")).Return(nil, errs.ErrInternalServerError)	

			req := httptest.NewRequest("POST", "/cms/partnerpages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})			
	})

	t.Run("GET /cms/partnerpages HandleGetPartnerPages", func(t *testing.T) {
		rawQuery := `{
			"title": "Cloud Migration Success",
			"category_partner": "Tech Partner",
			"category_keywords": "cloud, migration, success",
			"category_scale": "Enterprise",
			"category_industry": "Information Technology",
			"category_goal": "Digital Transformation",
			"category_functions": "IT Operations",
			"status": "Published"
		}`
		escapedQuery := url.QueryEscape(rawQuery)
		sort := "title"
		page := 2
		limit := 5
		language := "en"
		
		expectedPartnerPages := []models.PartnerPage{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}	
		
		t.Run("successfully get partner page", func(t *testing.T) {
			mockService.On("FindPartnerPages", rawQuery, sort, page, limit, language).Return(expectedPartnerPages, int64(2), nil)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/partnerpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed to get partner page: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindPartnerPages", rawQuery, sort, page, limit, language).Return(nil, int64(0), errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/partnerpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
		
		t.Run("failed to get partner page: invalid query", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindPartnerPages", rawQuery, sort, page, limit, language).Return(nil, int64(0), errs.ErrInvalidQuery)

			req := httptest.NewRequest("GET", fmt.Sprintf(
				"/cms/partnerpages?query=%s&sort=%s&page=%d&limit=%d&language=%s",
				escapedQuery, sort, page, limit, language,
			), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})

	t.Run("GET /cms/partnerpages/:pageId HandleGetPartnerPageById", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		pageId := uuid.New()		

		t.Run("successfully get partner page by id", func(t *testing.T) {
			mockService.On("FindPartnerPageById", pageId).Return(mockPartnerPage, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})		

		t.Run("failed get partner page by id: invalid pageId", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindPartnerPageById", pageId).Return(mockPartnerPage, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s", invalidPageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})	
		
		t.Run("failed get partner page by id: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("FindPartnerPageById", pageId).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)		
		})			
	})

	t.Run("DELETE /cms/partnerpages/:pageId HandleDeletePartnerPage", func(t *testing.T)	{
		pageId := uuid.New()

		t.Run("successfully delete partner page", func(t *testing.T)	{
			mockService.On("DeletePartnerPage", pageId).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/partnerpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})
		
		t.Run("failed to delete partner page: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeletePartnerPage", pageId).Return(nil)
			invalidPageId := "1"

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/partnerpages/%s", invalidPageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})	
		
		t.Run("failed to delete partner page: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeletePartnerPage", pageId).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/partnerpages/%s", pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})

	t.Run("GET /cms/partnerpages/:pageId/contents/:languageCode HandleGetContentByPartnerPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockPartnerContent := mockPartnerPage.Contents[0]

		t.Run("successfully get content by partner pageId", func(t *testing.T)	{
			mockService.On("FindContentByPartnerPageId", pageId, language, mode).Return(mockPartnerContent, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get content by partner pageId: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindContentByPartnerPageId", pageId, language, mode).Return(mockPartnerContent, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s/contents/%s?mode=%s", invalidPageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get content by partner pageId: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindContentByPartnerPageId", pageId, language, mode).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})		

	t.Run("GET /cms/partnerpages/:pageId/:languageCode/contents HandleGetLatestContentByPartnerPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockPartnerContent := mockPartnerPage.Contents[0]

		t.Run("successfully get latest content", func(t *testing.T)	{
			mockService.On("FindLatestContentByPageId", pageId, language).Return(mockPartnerContent, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s/%s/contents", pageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get latest content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindLatestContentByPageId", pageId, language).Return(mockPartnerContent, nil)
			invalidPageId := "1"

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s/%s/contents", invalidPageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindLatestContentByPageId", pageId, language).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/%s/%s/contents", pageId, language), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("DELETE /cms/partnerpages/:pageId/contents/:languageCode HandleDeletePartnerContentByPageId", func(t *testing.T)	{
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		t.Run("successfully revert partner content", func(t *testing.T)	{
			mockService.On("DeleteContentByPartnerPageId", pageId, language, mode).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/partnerpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			

		t.Run("failed to get latest content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContentByPartnerPageId", pageId, language, mode).Return(nil)
			invalidPageId := "1"

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/partnerpages/%s/contents/%s?mode=%s", invalidPageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContentByPartnerPageId", pageId, language, mode).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/partnerpages/%s/contents/%s?mode=%s", pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})

	t.Run("POST /cms/partnerpages/duplicate/:pageId HandleDuplicatePartnerPage", func(t *testing.T)	{
		pageId := uuid.New()	
		mockPartnerPage := helpers.InitializeMockPartnerPage()

		t.Run("successfully duplicate partner page", func(t *testing.T)	{
			mockService.On("DuplicatePartnerPage", pageId).Return(mockPartnerPage, nil)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s", pageId), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to duplicate partner page: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicatePartnerPage", pageId).Return(mockPartnerPage, nil)
			invalidPageId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s", invalidPageId), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to duplicate partner page: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicatePartnerPage", pageId).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s", pageId), nil)
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})
	
	t.Run("POST /cms/partnerpages/duplicate/:contentId/contents HandleDuplicatePartnerContentToAnotherLanguage", func(t *testing.T)	{
		contentId := uuid.New()	
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully duplicate partner content to another language", func(t *testing.T)	{
			mockService.On("DuplicatePartnerContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to duplicate partner content to another language: invalid body", func(t *testing.T)	{
			mockService.On("DuplicatePartnerContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to duplicate partner content to another language: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicatePartnerContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)
			invalidContentId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s/contents", invalidContentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get latest content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("DuplicatePartnerContentToAnotherLanguage", contentId, mock.AnythingOfType("*models.Revision")).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/duplicate/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})			

	t.Run("POST /cms/partnerpages/:revisionId/revision HandleRevertPartnerContent", func(t *testing.T)	{
		revisionId := uuid.New()	
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		mockRevision := mockContent.Revision

		body, err := json.Marshal(mockRevision)
		require.NoError(t, err)			

		t.Run("successfully revert partner content", func(t *testing.T)	{
			mockService.On("RevertPartnerContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to revert partner content: invalid body", func(t *testing.T)	{
			mockService.On("RevertPartnerContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to revert partner content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("RevertPartnerContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(mockContent, nil)
			invalidRevisionId := "1"		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/%s/revision", invalidRevisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to revert partner content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("RevertPartnerContent", revisionId, mock.AnythingOfType("*models.Revision")).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/partnerpages/%s/revision", revisionId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})
	
	t.Run("PUT /cms/partnerpages/:contentId/contents HandleUpdatePartnerContent", func(t *testing.T)	{
		contentId := uuid.New()	
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]

		body, err := json.Marshal(mockContent)
		require.NoError(t, err)			

		t.Run("successfully update partner content", func(t *testing.T)	{
			mockService.On("UpdatePartnerContent", mock.AnythingOfType("*models.PartnerContent"), contentId).Return(mockContent, nil)				

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/partnerpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to update partner content: invalid body", func(t *testing.T)	{
			mockService.On("UpdatePartnerContent", mock.AnythingOfType("*models.PartnerContent"), contentId).Return(mockContent, nil)				

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/partnerpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		

		t.Run("failed to update partner content: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("UpdatePartnerContent", mock.AnythingOfType("*models.PartnerContent"), contentId).Return(mockContent, nil)
			invalidcontentId := "1"		

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/partnerpages/%s/contents", invalidcontentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to update partner content: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("UpdatePartnerContent", mock.AnythingOfType("*models.PartnerContent"), contentId).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("PUT", fmt.Sprintf("/cms/partnerpages/%s/contents", contentId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/partnerpages/category/:categoryTypeCode/:pageId/:languageCode HandleGetCategory", func(t *testing.T)	{
		pageId := uuid.New()	
		categoryTypeCode := "TYPE_CODE"
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		category := mockContent.Categories[0]	
		
		expectedCategories := []models.Category{*category, *category}

		t.Run("successfully get categories", func(t *testing.T)	{
			mockService.On("GetCategory", pageId, categoryTypeCode, language, mode).Return(expectedCategories, nil)				

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/category/%s/%s/%s?mode=%s", categoryTypeCode, pageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to get categories: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("GetCategory", pageId, categoryTypeCode, language, mode).Return(expectedCategories, nil)
			invalidpageId := "1"		

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/category/%s/%s/%s?mode=%s", categoryTypeCode, invalidpageId, language, mode), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get categories: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("GetCategory", pageId, categoryTypeCode, language, mode).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/category/%s/%s/%s?mode=%s", categoryTypeCode,  pageId, language, mode), nil)		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})	

	t.Run("GET /cms/partnerpages/revisions/:languageCode/:pageId HandleGetRevisions", func(t *testing.T)	{
		pageId := uuid.New()	
		language := string(enums.PageLanguageEN)
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		revision := mockContent.Revision
		expectedRevisions := []models.Revision{*revision, *revision}
		
		t.Run("successfully get revisions", func(t *testing.T)	{
			mockService.On("FindRevisions", pageId, language).Return(expectedRevisions, nil)				

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/revisions/%s/%s", language, pageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})				

		t.Run("failed to get revisions: invalid pageId", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindRevisions", pageId, language).Return(expectedRevisions, nil)
			invalidpageId := "1"		

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/revisions/%s/%s", language, invalidpageId), nil)
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})		
		
		t.Run("failed to get revisions: internal server error", func(t *testing.T)	{
			mockService.ExpectedCalls = nil
			mockService.On("FindRevisions", pageId, language).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/partnerpages/revisions/%s/%s", language, pageId), nil)		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)				
		})			
	})		
}