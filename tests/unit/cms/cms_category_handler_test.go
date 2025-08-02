package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCMSCategoryService struct {
	mock.Mock
}

func (m *MockCMSCategoryService) CreateCategory(req dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func (m *MockCMSCategoryService) GetCategoryByUUID(uuidStr string) (*dto.CategoryResponse, error) {
	args := m.Called(uuidStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func (m *MockCMSCategoryService) ListAllCategories(filter dto.CategoryFilter) ([]dto.CategoryResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.CategoryResponse), args.Error(1)
}

func (m *MockCMSCategoryService) UpdateCategoryByUUID(uuidStr string, req dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	args := m.Called(uuidStr, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func (m *MockCMSCategoryService) DeleteCategoryByUUID(uuidStr string) error {
	args := m.Called(uuidStr)
	return args.Error(0)
}

func (m *MockCMSCategoryService) MapCategoryModelToResponse(cat *models.Category) (*dto.CategoryResponse, error) {
	args := m.Called(cat)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func TestCMSCategoryHandler(t *testing.T) {
	mockService := &MockCMSCategoryService{}
	handler := cmsHandler.NewCMSCategoryHandler(mockService)

	app := fiber.New()
	app.Post("/cms/categories", handler.HandleCreateCategory)
	app.Get("/cms/categories", handler.HandleListAllCategories)
	app.Get("/cms/categories/:categoryUuid", handler.HandleGetCategoryByUUID)
	app.Patch("/cms/categories/:categoryUuid", handler.HandleUpdateCategory)
	app.Delete("/cms/categories/:categoryUuid", handler.HandleDeleteCategory)

	t.Run("POST /cms/categories HandleCreateCategory", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		categoryTypeId := uuid.New()
		categoryCreateRequest := dto.CategoryCreateRequest{
			Name: mockCategory.Name,
			Description: mockCategory.Description,
			CategoryTypeID: categoryTypeId.String(),
			LanguageCode: mockCategory.LanguageCode,
			Weight: &mockCategory.Weight,
			PublishStatus: mockCategory.PublishStatus,
		}

		body, err := json.Marshal(categoryCreateRequest)
		require.NoError(t, err)

		t.Run("successfully create category", func(t *testing.T) {
			mockService.On("CreateCategory", mock.AnythingOfType("dto.CategoryCreateRequest")).Return(&dto.CategoryResponse{
				ID: mockCategory.ID.String(),
				Name: mockCategory.Name,
				Description: mockCategory.Description,
				CategoryTypeID: mockCategory.CategoryType.ID.String(),
				LanguageCode: mockCategory.LanguageCode,
				Weight: mockCategory.Weight,
				PublishStatus: mockCategory.PublishStatus,
			}, nil)

			req := httptest.NewRequest("POST", "/cms/categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusCreated, resp.StatusCode)
		})

		t.Run("failed to create category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateCategory", mock.AnythingOfType("dto.CategoryCreateRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("POST", "/cms/categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})		
	})

	t.Run("GET /cms/categories HandleListAllCategories", func(t *testing.T) {
		mockCategoryType, mockCategory := helpers.InitializeMockCategory()
		categoryTypeId := mockCategoryType.ID.String()
		categoryFilter := dto.CategoryFilter{
			CategoryTypeID: &categoryTypeId,
			LanguageCode: &mockCategory.LanguageCode,
			Name: &mockCategory.Name,
			PublishStatus: &mockCategory.PublishStatus,
		}
		body, err := json.Marshal(categoryFilter)
		require.NoError(t, err)

		t.Run("successfully list all categories", func(t *testing.T) {
			mockService.On("ListAllCategories", mock.AnythingOfType("dto.CategoryFilter")).Return([]dto.CategoryResponse{
				{
					ID: mockCategory.ID.String(),
					Name: mockCategory.Name,
					Description: mockCategory.Description,
					CategoryTypeID: mockCategory.CategoryType.ID.String(),
					LanguageCode: mockCategory.LanguageCode,
					Weight: mockCategory.Weight,
					PublishStatus: mockCategory.PublishStatus,
				},
			}, nil)

			req := httptest.NewRequest("GET", "/cms/categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to list all categories", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("ListAllCategories", mock.AnythingOfType("dto.CategoryFilter")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", "/cms/categories", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})		
	})

	t.Run("GET /cms/categories/:categoryUuid HandleGetCategoryByUUID", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		t.Run("successfully get category by uuid", func(t *testing.T) {
			mockService.On("GetCategoryByUUID", mock.AnythingOfType("string")).Return(&dto.CategoryResponse{
				ID: mockCategory.ID.String(),
				Name: mockCategory.Name,
				Description: mockCategory.Description,
				CategoryTypeID: mockCategory.CategoryType.ID.String(),
				LanguageCode: mockCategory.LanguageCode,
				Weight: mockCategory.Weight,
				PublishStatus: mockCategory.PublishStatus,
			}, nil)

			req := httptest.NewRequest("GET", "/cms/categories/"+mockCategory.ID.String(), nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, resp.StatusCode)
		})	

		t.Run("failed to get category by uuid", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetCategoryByUUID", mock.AnythingOfType("string")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/categories/%s", mockCategory.ID.String()), nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})	
	})

	t.Run("PATCH /cms/categories/:categoryUuid HandleUpdateCategory", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		description := mockCategory.Description
		categoryUpdateRequest := dto.CategoryUpdateRequest{
			Name: &mockCategory.Name,
			Description: description,
			Weight: &mockCategory.Weight,
			PublishStatus: &mockCategory.PublishStatus,
		}
		body, err := json.Marshal(categoryUpdateRequest)
		require.NoError(t, err)

		t.Run("successfully update category", func(t *testing.T) {
			mockService.On("UpdateCategoryByUUID", mock.AnythingOfType("string"), mock.AnythingOfType("dto.CategoryUpdateRequest")).Return(&dto.CategoryResponse{
				ID: mockCategory.ID.String(),
				Name: mockCategory.Name,
				Description: mockCategory.Description,
				CategoryTypeID: mockCategory.CategoryType.ID.String(),
				LanguageCode: mockCategory.LanguageCode,
				Weight: mockCategory.Weight,
				PublishStatus: mockCategory.PublishStatus,
			}, nil)

			req := httptest.NewRequest("PATCH", "/cms/categories/"+mockCategory.ID.String(), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to update category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("UpdateCategoryByUUID", mock.AnythingOfType("string"), mock.AnythingOfType("dto.CategoryUpdateRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/categories/%s", mockCategory.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})	
	})

	t.Run("DELETE /cms/categories/:categoryUuid HandleDeleteCategory", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.ID = uuid.New()
		body, err := json.Marshal(mockCategory)
		require.NoError(t, err)

		t.Run("successfully delete category", func(t *testing.T) {
			mockService.On("DeleteCategoryByUUID", mock.AnythingOfType("string")).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/categories/%s", mockCategory.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		})

		t.Run("failed to delete category", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("DeleteCategoryByUUID", mock.AnythingOfType("string")).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/categories/%s", mockCategory.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})	
	})
}