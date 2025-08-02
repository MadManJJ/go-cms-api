package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockCMSCategoryRepo struct {
	createCategory                   func(category *models.Category) (*models.Category, error)
	getCategoryByID                  func(categoryID uuid.UUID) (*models.Category, error)
	updateCategory                   func(category *models.Category) (*models.Category, error)
	deleteCategory                   func(categoryID uuid.UUID) error
	countCategoriesByTypeAndLanguage func(categoryTypeID uuid.UUID) (map[string]int, error)
	listCategoriesByFilter           func(filters dto.CategoryFilter) ([]models.Category, error)
}

func (m *MockCMSCategoryRepo) CreateCategory(category *models.Category) (*models.Category, error) {
	return m.createCategory(category)
}

func (m *MockCMSCategoryRepo) GetCategoryByID(categoryID uuid.UUID) (*models.Category, error) {
	return m.getCategoryByID(categoryID)
}

func (m *MockCMSCategoryRepo) UpdateCategory(category *models.Category) (*models.Category, error) {
	return m.updateCategory(category)
}

func (m *MockCMSCategoryRepo) DeleteCategory(categoryID uuid.UUID) error {
	return m.deleteCategory(categoryID)
}

func (m *MockCMSCategoryRepo) CountCategoriesByTypeAndLanguage(categoryTypeID uuid.UUID) (map[string]int, error) {
	return m.countCategoriesByTypeAndLanguage(categoryTypeID)
}

func (m *MockCMSCategoryRepo) ListCategoriesByFilter(filters dto.CategoryFilter) ([]models.Category, error) {
	return m.listCategoriesByFilter(filters)
}

func TestCMSService_MapCategoryModelToResponse(t *testing.T) {

	t.Run("successfully map category to response: CategoryType is not nil", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)

		categoryResponse, err := service.MapCategoryModelToResponse(mockCategory)
		assert.NoError(t, err)
		assert.NotNil(t, categoryResponse)

		assert.Equal(t, mockCategory.ID.String(), categoryResponse.ID)
		assert.Equal(t, mockCategory.CategoryTypeID.String(), categoryResponse.CategoryTypeID)
		assert.Equal(t, mockCategory.LanguageCode, categoryResponse.LanguageCode)
		assert.Equal(t, mockCategory.Name, categoryResponse.Name)
		assert.Equal(t, mockCategory.Description, categoryResponse.Description)
		assert.Equal(t, mockCategory.Weight, categoryResponse.Weight)
		assert.Equal(t, mockCategory.PublishStatus, categoryResponse.PublishStatus)
		assert.Equal(t, mockCategory.CreatedAt, categoryResponse.CreatedAt)
		assert.Equal(t, mockCategory.UpdatedAt, categoryResponse.UpdatedAt)

		assert.NotNil(t, categoryResponse.CategoryType)
		assert.Equal(t, mockCategory.CategoryType.ID.String(), categoryResponse.CategoryType.ID)
		assert.Equal(t, mockCategory.CategoryType.TypeCode, categoryResponse.CategoryType.TypeCode)
	})

	t.Run("nil category returns nil response", func(t *testing.T) {
		mockCategoryRepo := &MockCMSCategoryRepo{}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)

		categoryResponse, err := service.MapCategoryModelToResponse(nil)
		assert.NoError(t, err)
		assert.Nil(t, categoryResponse)
	})
}

func TestCMSService_CreateCategory(t *testing.T) {
	t.Run("successfully create category: category type is active", func(t *testing.T) {
		mockCategoryType, mockCategory := helpers.InitializeMockCategory()

		mockCategoryType.IsActive = true

		createdMockCategory := *mockCategory
		createdMockCategory.ID = uuid.New()
		createdMockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			createCategory: func(category *models.Category) (*models.Category, error) {
				return &createdMockCategory, nil
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)
		
		categoryTypeId := uuid.New()
		req := dto.CategoryCreateRequest{
			CategoryTypeID: categoryTypeId.String(),
			LanguageCode:   mockCategory.LanguageCode,
			Name:           mockCategory.Name,
			Description:   	mockCategory.Description,
			Weight:         &mockCategory.Weight,
			PublishStatus:  mockCategory.PublishStatus,
		}
		
		categoryResponse, err := service.CreateCategory(req)
		assert.NoError(t, err)
		assert.NotNil(t, categoryResponse)
		assert.NotNil(t, categoryResponse.CategoryType)
	})

	t.Run("failed to create category: category type is not active", func(t *testing.T) {
		mockCategoryType, mockCategory := helpers.InitializeMockCategory()

		mockCategoryType.IsActive = false

		mockCategoryRepo := &MockCMSCategoryRepo{}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)
		
		categoryTypeId := uuid.New()
		req := dto.CategoryCreateRequest{
			CategoryTypeID: categoryTypeId.String(),
			LanguageCode:   mockCategory.LanguageCode,
			Name:           mockCategory.Name,
			Description:   	mockCategory.Description,
			Weight:         &mockCategory.Weight,
			PublishStatus:  mockCategory.PublishStatus,
		}
		
		categoryResponse, err := service.CreateCategory(req)
		assert.Error(t, err)
		assert.Nil(t, categoryResponse)
	})	
}

func TestCMSService_GetCategoryByUUID(t *testing.T) {
	t.Run("successfully get category by UUID", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return mockCategory, nil
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)

		categoryResponse, err := service.GetCategoryByUUID(mockCategory.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, categoryResponse)
		assert.NotNil(t, categoryResponse.CategoryType)
	})

	t.Run("failed to get category by UUID", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)

		categoryResponse, err := service.GetCategoryByUUID(mockCategory.ID.String())
		assert.Error(t, err)
		assert.Nil(t, categoryResponse)
	})	
}

func TestCMSService_ListAllCategories(t *testing.T) {
	t.Run("successfully list all categories", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			listCategoriesByFilter: func(filters dto.CategoryFilter) ([]models.Category, error) {
				return []models.Category{*mockCategory}, nil
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)

		categoryResponse, err := service.ListAllCategories(dto.CategoryFilter{})
		assert.NoError(t, err)
		assert.NotNil(t, categoryResponse)
		assert.NotNil(t, categoryResponse[0].CategoryType)
	})

	t.Run("failed to list all categories", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			listCategoriesByFilter: func(filters dto.CategoryFilter) ([]models.Category, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)

		categoryResponse, err := service.ListAllCategories(dto.CategoryFilter{})
		assert.Error(t, err)
		assert.Nil(t, categoryResponse)
	})	
}

func TestCMSService_UpdateCategoryByUUID(t *testing.T) {
	t.Run("successfully update category by id", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		updatedMockCategory := *mockCategory
		updatedMockCategory.Name = "updated name"

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return mockCategory, nil
			},
			updateCategory: func(category *models.Category) (*models.Category, error) {
				return &updatedMockCategory, nil
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)	
		
		categoryResponse, err := service.UpdateCategoryByUUID(mockCategory.ID.String(), dto.CategoryUpdateRequest{
			Name:           &updatedMockCategory.Name,
			Description:    updatedMockCategory.Description,
			Weight:         &updatedMockCategory.Weight,
			PublishStatus:  &updatedMockCategory.PublishStatus,
		})
		assert.NoError(t, err)
		assert.NotNil(t, categoryResponse)
		assert.NotNil(t, categoryResponse.CategoryType)
		assert.Equal(t, updatedMockCategory.Name, categoryResponse.Name)
		assert.Equal(t, updatedMockCategory.Description, categoryResponse.Description)
		assert.Equal(t, updatedMockCategory.Weight, categoryResponse.Weight)
		assert.Equal(t, updatedMockCategory.PublishStatus, categoryResponse.PublishStatus)
	})

	t.Run("failed to update category by id: internal server error", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		updatedMockCategory := *mockCategory
		updatedMockCategory.Name = "updated name"

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return mockCategory, nil
			},
			updateCategory: func(category *models.Category) (*models.Category, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)	
		
		categoryResponse, err := service.UpdateCategoryByUUID(mockCategory.ID.String(), dto.CategoryUpdateRequest{
			Name:           &updatedMockCategory.Name,
			Description:    updatedMockCategory.Description,
			Weight:         &updatedMockCategory.Weight,
			PublishStatus:  &updatedMockCategory.PublishStatus,
		})
		assert.Error(t, err)
		assert.Nil(t, categoryResponse)
	})

	t.Run("failed to update category by id: category not found", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		updatedMockCategory := *mockCategory
		updatedMockCategory.Name = "updated name"

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)	
		
		categoryResponse, err := service.UpdateCategoryByUUID(mockCategory.ID.String(), dto.CategoryUpdateRequest{
			Name:           &updatedMockCategory.Name,
			Description:    updatedMockCategory.Description,
			Weight:         &updatedMockCategory.Weight,
			PublishStatus:  &updatedMockCategory.PublishStatus,
		})
		assert.Error(t, err)
		assert.Nil(t, categoryResponse)
	})	
}

func TestCMSService_DeleteCategoryByUUID(t *testing.T) {
	t.Run("successfully delete category by id", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return mockCategory, nil
			},
			deleteCategory: func(categoryID uuid.UUID) error {
				return nil
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)		
		
		err := service.DeleteCategoryByUUID(mockCategory.ID.String())
		assert.NoError(t, err)
	})

	t.Run("failed to delete category by id: internal server error", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return mockCategory, nil
			},
			deleteCategory: func(categoryID uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)		
		
		err := service.DeleteCategoryByUUID(mockCategory.ID.String())
		assert.Error(t, err)
	})

	t.Run("failed to delete category by id: category not found", func(t *testing.T) {
		_, mockCategory := helpers.InitializeMockCategory()
		mockCategory.CategoryType.ID = uuid.New()

		mockCategoryRepo := &MockCMSCategoryRepo{
			getCategoryByID: func(categoryID uuid.UUID) (*models.Category, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockCategoryTypeRepo := &MockCMSCategoryTypeRepo{}
		service := services.NewCMSCategoryService(mockCategoryRepo, mockCategoryTypeRepo)		
		
		err := service.DeleteCategoryByUUID(mockCategory.ID.String())
		assert.Error(t, err)
	})	
}