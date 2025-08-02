package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockEmailCategoryService struct {
	mock.Mock
}

func (m *MockEmailCategoryService) CreateCategory(req dto.CreateEmailCategoryRequest) (*dto.EmailCategoryResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EmailCategoryResponse), args.Error(1)
}

func (m *MockEmailCategoryService) GetCategoryByID(idStr string) (*dto.EmailCategoryResponse, error) {
	args := m.Called(idStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EmailCategoryResponse), args.Error(1)
}

func (m *MockEmailCategoryService) GetCategoryByTitle(title string) (*dto.EmailCategoryResponse, error) {
	args := m.Called(title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EmailCategoryResponse), args.Error(1)
}

func (m *MockEmailCategoryService) ListCategories() ([]dto.EmailCategoryResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.EmailCategoryResponse), args.Error(1)
}

func (m *MockEmailCategoryService) UpdateCategory(idStr string, req dto.UpdateEmailCategoryRequest) (*dto.EmailCategoryResponse, error) {
	args := m.Called(idStr, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EmailCategoryResponse), args.Error(1)
}

func (m *MockEmailCategoryService) DeleteCategory(idStr string) error {
	args := m.Called(idStr)
	return args.Error(0)
}


func TestCMSEmailCategoryHandler(t *testing.T) {
	mockService := &MockEmailCategoryService{}
	handler := cmsHandler.NewEmailCategoryHandler(mockService)

	app := fiber.New()
	app.Post("/cms/email-categories", handler.HandleCreateEmailCategory)
	app.Get("/cms/email-categories", handler.HandleListEmailCategories)
	app.Get("/cms/email-categories/:id", handler.HandleGetEmailCategory)
	app.Patch("/cms/email-categories/:id", handler.HandleUpdateEmailCategory)
	app.Delete("/cms/email-categories/:id", handler.HandleDeleteEmailCategory)

	t.Run("POST /cms/email-categories HandleCreateEmailCategory", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()

		createEmailCategoryRequest := dto.CreateEmailCategoryRequest{
			Title: mockEmailCategory.Title,
		}

		mockEmailCategoryResponse := dto.EmailCategoryResponse{
			ID:        mockEmailCategory.ID.String(),
			Title:     mockEmailCategory.Title,
			CreatedAt: mockEmailCategory.CreatedAt,
			UpdatedAt: mockEmailCategory.UpdatedAt,
		}

		body, err := json.Marshal(createEmailCategoryRequest)
		require.NoError(t, err)

		t.Run("successfully create email category", func(t *testing.T) {
			mockService.On("CreateCategory", mock.AnythingOfType("dto.CreateEmailCategoryRequest")).Return(&mockEmailCategoryResponse, nil)

			req := httptest.NewRequest("POST", "/cms/email-categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to create email category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateCategory", mock.AnythingOfType("dto.CreateEmailCategoryRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("POST", "/cms/email-categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("GET /cms/email-categories HandleListEmailCategories", func(t *testing.T) {
		emailCategoryResponse := dto.EmailCategoryResponse{
			ID:        "1",
			Title:     "Email Category 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockEmailCategories := []dto.EmailCategoryResponse{
			emailCategoryResponse,
			emailCategoryResponse,
		}

		t.Run("successfully list email categories", func(t *testing.T) {
			mockService.On("ListCategories").Return(mockEmailCategories, nil)

			req := httptest.NewRequest("GET", "/cms/email-categories", nil)
	
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to list email categories", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("ListCategories").Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", "/cms/email-categories", nil)
		
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("GET /cms/email-categories/:id HandleGetEmailCategory", func(t *testing.T) {
		emailCategoryResponse := dto.EmailCategoryResponse{
			ID:        "1",
			Title:     "Email Category 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		categoryId := uuid.New()

		t.Run("successfully get email category", func(t *testing.T) {
			mockService.On("GetCategoryByID", categoryId.String()).Return(&emailCategoryResponse, nil)

			req := httptest.NewRequest("GET", "/cms/email-categories/"+categoryId.String(), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get email category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetCategoryByID", categoryId.String()).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", "/cms/email-categories/"+categoryId.String(), nil)
		
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("PATCH /cms/email-categories/:id HandleUpdateEmailCategory", func(t *testing.T) {
		updateEmailCategoryRequest := dto.UpdateEmailCategoryRequest{
			Title: "Updated Email Category",
		}

		body, err := json.Marshal(updateEmailCategoryRequest)
		require.NoError(t, err)		

		emailCategoryResponse := dto.EmailCategoryResponse{
			ID:        "1",
			Title:     "Email Category 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		emailCategoryResponse.Title = updateEmailCategoryRequest.Title

		categoryId := uuid.New()

		t.Run("successfully update email category", func(t *testing.T) {
			mockService.On("UpdateCategory", categoryId.String(), updateEmailCategoryRequest).Return(&emailCategoryResponse, nil)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/email-categories/%s", categoryId.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to update email category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("UpdateCategory", categoryId.String(), updateEmailCategoryRequest).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/email-categories/%s", categoryId.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
		
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("DELETE /cms/email-categories/:id HandleDeleteEmailCategory", func(t *testing.T) {
		categoryId := uuid.New()

		t.Run("successfully delete email category", func(t *testing.T) {
			mockService.On("DeleteCategory", categoryId.String()).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/email-categories/%s", categoryId.String()), nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to delete email category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("DeleteCategory", categoryId.String()).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/email-categories/%s", categoryId.String()), nil)
			req.Header.Set("Content-Type", "application/json")
		
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})
}
