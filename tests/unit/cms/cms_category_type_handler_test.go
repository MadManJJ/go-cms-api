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
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCMSCategoryTypeService struct {
	mock.Mock
}

func (m *MockCMSCategoryTypeService) CreateCategoryType(req dto.CreateCategoryTypeRequest) (*dto.CategoryTypeResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryTypeResponse), args.Error(1)
}

func (m *MockCMSCategoryTypeService) GetCategoryTypeByID(idStr string) (*dto.CategoryTypeResponse, error) {
	args := m.Called(idStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryTypeResponse), args.Error(1)
}

func (m *MockCMSCategoryTypeService) GetCategoryTypeByCode(code string) (*dto.CategoryTypeResponse, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryTypeResponse), args.Error(1)
}

func (m *MockCMSCategoryTypeService) ListCategoryTypes(isActive *bool) ([]dto.CategoryTypeResponse, error) {
	args := m.Called(isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.CategoryTypeResponse), args.Error(1)
}

func (m *MockCMSCategoryTypeService) UpdateCategoryType(idStr string, req dto.UpdateCategoryTypeRequest) (*dto.CategoryTypeResponse, error) {
	args := m.Called(idStr, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryTypeResponse), args.Error(1)
}

func (m *MockCMSCategoryTypeService) DeleteCategoryType(idStr string) error {
	args := m.Called(idStr)
	return args.Error(0)
}

func (m *MockCMSCategoryTypeService) GetCategoryTypeWithDetails(categoryTypeIDStr string, languageCodeStr string) (*dto.CategoryTypeWithDetailsResponse, error) {
	args := m.Called(categoryTypeIDStr, languageCodeStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryTypeWithDetailsResponse), args.Error(1)
}

func TestCMSCategoryTypeHandler(t *testing.T) {
	mockService := &MockCMSCategoryTypeService{}
	handler := cmsHandler.NewCMSCategoryTypeHandler(mockService)

	app := fiber.New()	
	app.Post("/cms/category-types", handler.HandleCreateCategoryType)
	app.Get("/cms/category-types", handler.HandleListCategoryTypes)
	app.Get("/cms/category-types/:id", handler.HandleGetCategoryType)
	app.Patch("/cms/category-types/:id", handler.HandleUpdateCategoryType)
	app.Delete("/cms/category-types/:id", handler.HandleDeleteCategoryType)
	app.Get("/cms/category-types/:categoryTypeId/categories", handler.HandleListCategoriesForType)

	t.Run("POST /cms/category-types HandleCreateCategoryType", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		body, err := json.Marshal(mockCategoryType)
		require.NoError(t, err)

		t.Run("successfully create category type", func(t *testing.T) {
			mockService.On("CreateCategoryType", mock.AnythingOfType("dto.CreateCategoryTypeRequest")).Return(&dto.CategoryTypeResponse{
				ID: mockCategoryType.ID.String(),
				TypeCode: mockCategoryType.TypeCode,
				Name: &mockCategoryType.Name,
				IsActive: mockCategoryType.IsActive,
				CreatedAt: mockCategoryType.CreatedAt,
				UpdatedAt: mockCategoryType.UpdatedAt,
			}, nil)

			req := httptest.NewRequest("POST", "/cms/category-types", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		})

		t.Run("failed to create category type", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateCategoryType", mock.AnythingOfType("dto.CreateCategoryTypeRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("POST", "/cms/category-types", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})		
	})
	
	t.Run("GET /cms/category-types HandleListCategoryTypes", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		body, err := json.Marshal(mockCategoryType)
		require.NoError(t, err)

		t.Run("successfully list category types", func(t *testing.T) {
			mockService.On("ListCategoryTypes", mock.AnythingOfType("*bool")).Return([]dto.CategoryTypeResponse{
				{
					ID: mockCategoryType.ID.String(),
					TypeCode: mockCategoryType.TypeCode,
					Name: &mockCategoryType.Name,
					IsActive: mockCategoryType.IsActive,
					CreatedAt: mockCategoryType.CreatedAt,
					UpdatedAt: mockCategoryType.UpdatedAt,
				},
			}, nil)

			req := httptest.NewRequest("GET", "/cms/category-types", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to list category types", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("ListCategoryTypes", mock.AnythingOfType("*bool")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", "/cms/category-types", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})
	
	t.Run("GET /cms/category-types/:id HandleGetCategoryType", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		body, err := json.Marshal(mockCategoryType)
		require.NoError(t, err)

		t.Run("successfully get category type", func(t *testing.T) {
			mockService.On("GetCategoryTypeByID", mock.AnythingOfType("string")).Return(&dto.CategoryTypeResponse{
				ID: mockCategoryType.ID.String(),
				TypeCode: mockCategoryType.TypeCode,
				Name: &mockCategoryType.Name,
				IsActive: mockCategoryType.IsActive,
				CreatedAt: mockCategoryType.CreatedAt,
				UpdatedAt: mockCategoryType.UpdatedAt,
			}, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/category-types/%s", mockCategoryType.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to get category type", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetCategoryTypeByID", mock.AnythingOfType("string")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/category-types/%s", mockCategoryType.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})
	
	t.Run("PATCH /cms/category-types/:id HandleUpdateCategoryType", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()
		updateCategoryTypeRequest := dto.UpdateCategoryTypeRequest{
			Name: &mockCategoryType.Name,
			IsActive: &mockCategoryType.IsActive,
		}
		body, err := json.Marshal(updateCategoryTypeRequest)
		require.NoError(t, err)

		t.Run("successfully update category type", func(t *testing.T) {
			mockService.On("UpdateCategoryType", mock.AnythingOfType("string"), mock.AnythingOfType("dto.UpdateCategoryTypeRequest")).Return(&dto.CategoryTypeResponse{
				ID: mockCategoryType.ID.String(),
				TypeCode: mockCategoryType.TypeCode,
				Name: &mockCategoryType.Name,
				IsActive: mockCategoryType.IsActive,
				CreatedAt: mockCategoryType.CreatedAt,
				UpdatedAt: mockCategoryType.UpdatedAt,
			}, nil)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/category-types/%s", mockCategoryType.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to update category type", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("UpdateCategoryType", mock.AnythingOfType("string"), mock.AnythingOfType("dto.UpdateCategoryTypeRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/category-types/%s", mockCategoryType.ID.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})
	
	t.Run("DELETE /cms/category-types/:id HandleDeleteCategoryType", func(t *testing.T) {
		categoryTypeId := uuid.New()

		t.Run("successfully delete category type", func(t *testing.T) {
			mockService.On("DeleteCategoryType", mock.AnythingOfType("string")).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/category-types/%s", categoryTypeId.String()), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		})

		t.Run("failed to delete category type", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("DeleteCategoryType", mock.AnythingOfType("string")).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/category-types/%s", categoryTypeId.String()), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})
	
	t.Run("GET /cms/category-types/:id/details/:languageCode HandleListCategoriesForType", func(t *testing.T) {
		categoryTypeId := uuid.New()
		languageCode := "en"
		description := "Description 1"
		name := "name"
		
		t.Run("successfully list categories for type", func(t *testing.T) {
			mockService.On("GetCategoryTypeWithDetails", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&dto.CategoryTypeWithDetailsResponse{
				ID: categoryTypeId.String(),
				TypeCode: "typeCode",
				Name: &name,
				IsActive: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Categories: []dto.CategoryResponse{
					{
						ID: "categoryUuid",
						Name: "categoryName",
						Description: &description,
						Weight: 1,
						PublishStatus: enums.PublishStatusPublished,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
			}, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/category-types/%s/categories?lang=%s", categoryTypeId.String(), languageCode), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to list categories for type", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetCategoryTypeWithDetails", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/category-types/%s/categories?lang=%s", categoryTypeId.String(), languageCode), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})
}