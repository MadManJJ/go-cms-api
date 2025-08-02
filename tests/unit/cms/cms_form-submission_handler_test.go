package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCMSFormSubmissionService struct {
	mock.Mock
}

func (m *MockCMSFormSubmissionService) CreateFormSubmission(formId uuid.UUID, formSubmission *models.FormSubmission) (*models.FormSubmission, error) {
	args := m.Called(formId, formSubmission)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FormSubmission), args.Error(1)
}

func (m *MockCMSFormSubmissionService) GetFormSubmissions(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error) {
	args := m.Called(formId, sort, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.FormSubmission), args.Get(1).(int64), args.Error(2)
}

func (m *MockCMSFormSubmissionService) GetFormSubmission(submissionId uuid.UUID) (*models.FormSubmission, error) {
	args := m.Called(submissionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FormSubmission), args.Error(1)
}

func TestCMSFormSubmissionHandler(t *testing.T) {
	mockService := &MockCMSFormSubmissionService{}
	handler := cmsHandler.NewCMSFormSubmissionHandler(mockService)
	userId := uuid.New()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": userId.String(),
		}
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/cms/forms/:formId/submissions", handler.HandleCreateFormSubmission)
	app.Get("/cms/forms/:formId/submissions", handler.HandleGetFormSubmissions)
	app.Get("/cms/forms/submissions/:submissionId", handler.HandleGetFormSubmission)

	t.Run("POST /cms/forms/:formId/submissions HandleCreateFormSubmission", func(t *testing.T) {
		formId := uuid.New()

		mockFormSubmission := helpers.InitializeMockFormSubmission()

		createdFormSubmission := *mockFormSubmission

		body, err := json.Marshal(mockFormSubmission)
		require.NoError(t, err)				

		t.Run("successfully create form submission", func(t *testing.T) {
			mockService.On("CreateFormSubmission", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.FormSubmission")).Return(&createdFormSubmission, nil)		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/forms/%s/submissions", formId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")	
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})

		t.Run("failed to create form submission: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			body, err := json.Marshal("invalid body")
			require.NoError(t, err)					

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/forms/%s/submissions", formId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})

		t.Run("failed to create form submission: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateFormSubmission", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.FormSubmission")).Return(nil, errs.ErrInternalServerError)		

			req := httptest.NewRequest("POST", fmt.Sprintf("/cms/forms/%s/submissions", formId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")	
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})		
	})

	t.Run("GET /cms/forms/:formId/submissions HandleGetFormSubmissions", func(t *testing.T) {
		mockFormSubmission := helpers.InitializeMockFormSubmission()
		sort := "title:asc"
		page := 1
		limit := 10
		formId := uuid.New()
		
		t.Run("successfully get form submissions", func(t *testing.T) {
			mockService.On("GetFormSubmissions", formId, sort, page, limit).Return([]*models.FormSubmission{mockFormSubmission}, int64(1), nil)
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/forms/%s/submissions?sort=%s&page=%d&limit=%d", formId, sort, page, limit), nil)		

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})

		t.Run("failed to get form submissions: invalid formId", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFormSubmissions", formId, sort, page, limit).Return([]*models.FormSubmission{mockFormSubmission}, int64(0), nil)	

			formIdStr := "invalid_formId"	
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/forms/%s/submissions?sort=%s&page=%d&limit=%d", formIdStr, sort, page, limit), nil)			
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})			

		t.Run("failed to get form submissions: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFormSubmissions", formId, sort, page, limit).Return(nil, int64(0), errs.ErrInternalServerError)			
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/forms/%s/submissions?sort=%s&page=%d&limit=%d", formId, sort, page, limit), nil)			
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})			
	})

	t.Run("GET /cms/forms/submissions/:submissionId HandleGetFormSubmission", func(t *testing.T) {
		mockFormSubmission := helpers.InitializeMockFormSubmission()
		submissionId := uuid.New()
		
		t.Run("successfully get form submission", func(t *testing.T) {
			mockService.On("GetFormSubmission", submissionId).Return(mockFormSubmission, nil)
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/forms/submissions/%s", submissionId), nil)			
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})

		t.Run("failed to get form submission: invalid submissionId", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFormSubmission", submissionId).Return(mockFormSubmission, nil)	

			submissionIdStr := "invalid_submissionId"			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/forms/submissions/%s", submissionIdStr), nil)				
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})			

		t.Run("failed to get form submission: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFormSubmission", submissionId).Return(nil, errs.ErrInternalServerError)		
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/forms/submissions/%s", submissionId), nil)				
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)			
		})			
	})
}