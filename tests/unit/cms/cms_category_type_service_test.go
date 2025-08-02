package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockCMSCategoryTypeRepo struct {
	create           func(categoryType *models.CategoryType) (*models.CategoryType, error)
	findByID         func(id uuid.UUID) (*models.CategoryType, error)
	findByCode       func(code string) (*models.CategoryType, error)
	findAll          func(isActive *bool) ([]models.CategoryType, error)
	update           func(id uuid.UUID, updates map[string]interface{}) (*models.CategoryType, error)
	delete           func(id uuid.UUID) error
	isTypeCodeUnique func(typeCode string, excludeID uuid.UUID) (bool, error)
}

func (m *MockCMSCategoryTypeRepo) Create(categoryType *models.CategoryType) (*models.CategoryType, error) {
	return m.create(categoryType)
}

func (m *MockCMSCategoryTypeRepo) FindByID(id uuid.UUID) (*models.CategoryType, error) {
	return m.findByID(id)
}

func (m *MockCMSCategoryTypeRepo) FindByCode(code string) (*models.CategoryType, error) {
	return m.findByCode(code)
}

func (m *MockCMSCategoryTypeRepo) FindAll(isActive *bool) ([]models.CategoryType, error) {
	return m.findAll(isActive)
}

func (m *MockCMSCategoryTypeRepo) Update(id uuid.UUID, updates map[string]interface{}) (*models.CategoryType, error) {
	return m.update(id, updates)
}

func (m *MockCMSCategoryTypeRepo) Delete(id uuid.UUID) error {
	return m.delete(id)
}

func (m *MockCMSCategoryTypeRepo) IsTypeCodeUnique(typeCode string, excludeID uuid.UUID) (bool, error) {
	return m.isTypeCodeUnique(typeCode, excludeID)
}

func TestCMSService_CreateCategoryType(t *testing.T) {
	t.Run("successfully create category type", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			isTypeCodeUnique: func(typeCode string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
			create: func(categoryType *models.CategoryType) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.CreateCategoryType(dto.CreateCategoryTypeRequest{
			TypeCode: mockCategoryType.TypeCode,
			Name: &mockCategoryType.Name,
			IsActive: &mockCategoryType.IsActive,
		})
		assert.NoError(t, err)
		assert.NotNil(t, categoryTypeResponse.ID)
		assert.Equal(t, mockCategoryType.TypeCode, categoryTypeResponse.TypeCode)
		assert.Equal(t, mockCategoryType.Name, *categoryTypeResponse.Name)
		assert.Equal(t, mockCategoryType.IsActive, categoryTypeResponse.IsActive)
	})

	t.Run("failed to create category type: internal server error", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			isTypeCodeUnique: func(typeCode string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
			create: func(categoryType *models.CategoryType) (*models.CategoryType, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.CreateCategoryType(dto.CreateCategoryTypeRequest{
			TypeCode: mockCategoryType.TypeCode,
			Name: &mockCategoryType.Name,
			IsActive: &mockCategoryType.IsActive,
		})
		assert.Error(t, err)
		assert.Nil(t, categoryTypeResponse)
	})

	t.Run("failed to create category type: type code already exists", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			isTypeCodeUnique: func(typeCode string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.CreateCategoryType(dto.CreateCategoryTypeRequest{
			TypeCode: mockCategoryType.TypeCode,
			Name: &mockCategoryType.Name,
			IsActive: &mockCategoryType.IsActive,
		})
		assert.Error(t, err)
		assert.Nil(t, categoryTypeResponse)
	})	
}

func TestCMSService_GetCategoryTypeByID(t *testing.T) {
	t.Run("successfully get category type by id", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{
			countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
				return map[string]int{
					string(enums.PageLanguageTH): 2,
					string(enums.PageLanguageEN): 1,
				}, nil
			},
		}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.GetCategoryTypeByID(mockCategoryType.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, categoryTypeResponse.ID)
		assert.Equal(t, mockCategoryType.TypeCode, categoryTypeResponse.TypeCode)
		assert.Equal(t, mockCategoryType.Name, *categoryTypeResponse.Name)
		assert.Equal(t, mockCategoryType.IsActive, categoryTypeResponse.IsActive)
	})

	t.Run("failed to get category type by id: internal server error", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.GetCategoryTypeByID(mockCategoryType.ID.String())
		assert.Error(t, err)
		assert.Nil(t, categoryTypeResponse)
	})
}

func TestCMSService_GetCategoryTypeByCode(t *testing.T) {
	t.Run("successfully get category type by code", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByCode: func(code string) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{
			countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
				return map[string]int{
					string(enums.PageLanguageTH): 2,
					string(enums.PageLanguageEN): 1,
				}, nil
			},
		}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.GetCategoryTypeByCode(mockCategoryType.TypeCode)
		assert.NoError(t, err)
		assert.NotNil(t, categoryTypeResponse.ID)
		assert.Equal(t, mockCategoryType.TypeCode, categoryTypeResponse.TypeCode)
		assert.Equal(t, mockCategoryType.Name, *categoryTypeResponse.Name)
		assert.Equal(t, mockCategoryType.IsActive, categoryTypeResponse.IsActive)
	})

	t.Run("failed to get category type by code: internal server error", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByCode: func(code string) (*models.CategoryType, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.GetCategoryTypeByCode(mockCategoryType.TypeCode)
		assert.Error(t, err)
		assert.Nil(t, categoryTypeResponse)
	})
}

func TestCMSService_ListCategoryTypes(t *testing.T) {
	t.Run("successfully list category types", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findAll: func(isActive *bool) ([]models.CategoryType, error) {
				return []models.CategoryType{*mockCategoryType}, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{
			countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
				return map[string]int{
					string(enums.PageLanguageTH): 2,
					string(enums.PageLanguageEN): 1,
				}, nil
			},
		}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.ListCategoryTypes(nil)
		assert.NoError(t, err)
		assert.NotNil(t, categoryTypeResponse)
		assert.Equal(t, 1, len(categoryTypeResponse))
		assert.Equal(t, mockCategoryType.TypeCode, categoryTypeResponse[0].TypeCode)
		assert.Equal(t, mockCategoryType.Name, *categoryTypeResponse[0].Name)
		assert.Equal(t, mockCategoryType.IsActive, categoryTypeResponse[0].IsActive)
	})

	t.Run("failed to list category types: internal server error", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findAll: func(isActive *bool) ([]models.CategoryType, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
		
		categoryTypeResponse,err := categoryTypeService.ListCategoryTypes(nil)
		assert.Error(t, err)
		assert.Nil(t, categoryTypeResponse)
	})
}

func TestCMSService_UpdateCategoryType(t *testing.T) {
	t.Run("successfully update category type", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		updatedCategoryType := *mockCategoryType
		updatedCategoryType.Name = "updated name"
		updatedCategoryType.IsActive = false

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
			update: func(id uuid.UUID, updates map[string]interface{}) (*models.CategoryType, error) {
				return &updatedCategoryType, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{
			countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
				return map[string]int{
					string(enums.PageLanguageTH): 2,
					string(enums.PageLanguageEN): 1,
				}, nil
			},
		}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)	
		
		categoryTypeResponse,err := categoryTypeService.UpdateCategoryType(mockCategoryType.ID.String(), dto.UpdateCategoryTypeRequest{
			Name: &updatedCategoryType.Name,
			IsActive: &updatedCategoryType.IsActive,
		})
		assert.NoError(t, err)
		assert.NotNil(t, categoryTypeResponse)
		assert.Equal(t, updatedCategoryType.TypeCode, categoryTypeResponse.TypeCode)
		assert.Equal(t, updatedCategoryType.Name, *categoryTypeResponse.Name)
		assert.Equal(t, updatedCategoryType.IsActive, categoryTypeResponse.IsActive)
	})

	t.Run("failed to update category type: internal server error", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		updatedCategoryType := *mockCategoryType
		updatedCategoryType.Name = "updated name"
		updatedCategoryType.IsActive = false

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
			update: func(id uuid.UUID, updates map[string]interface{}) (*models.CategoryType, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)		
		
		categoryTypeResponse,err := categoryTypeService.UpdateCategoryType(mockCategoryType.ID.String(), dto.UpdateCategoryTypeRequest{
			Name: &updatedCategoryType.Name,
			IsActive: &updatedCategoryType.IsActive,
		})
		assert.Error(t, err)
		assert.Nil(t, categoryTypeResponse)
	})
}

func TestCMSService_DeleteCategoryType(t *testing.T) {
	t.Run("successfully delete category type", func(t *testing.T) {
			mockCategoryType, _ := helpers.InitializeMockCategory()
			mockCategoryType.ID = uuid.New()

			categoryTypeRepo := &MockCMSCategoryTypeRepo{
					delete: func(id uuid.UUID) error {
							return nil
					},
			}
			categoryRepo := &MockCMSCategoryRepo{
					countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
							return map[string]int{
									string(enums.PageLanguageTH): 0,
									string(enums.PageLanguageEN): 0,
							}, nil
					},
			}
			categoryService := &MockCMSCategoryService{}
			categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
			
			err := categoryTypeService.DeleteCategoryType(mockCategoryType.ID.String())
			assert.NoError(t, err)
	})

	t.Run("failed to delete category type: invalid ID format", func(t *testing.T) {
			categoryTypeRepo := &MockCMSCategoryTypeRepo{}
			categoryRepo := &MockCMSCategoryRepo{}
			categoryService := &MockCMSCategoryService{}
			categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
			
			err := categoryTypeService.DeleteCategoryType("invalid-id")
			assert.Error(t, err)
			assert.Equal(t, "invalid ID format for delete", err.Error())
	})

	t.Run("failed to delete category type: category type not found", func(t *testing.T) {
			mockCategoryType, _ := helpers.InitializeMockCategory()
			mockCategoryType.ID = uuid.New()

			categoryTypeRepo := &MockCMSCategoryTypeRepo{
					delete: func(id uuid.UUID) error {
							return gorm.ErrRecordNotFound
					},
			}
			categoryRepo := &MockCMSCategoryRepo{
					countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
							return map[string]int{
									string(enums.PageLanguageTH): 0,
									string(enums.PageLanguageEN): 0,
							}, nil
					},
			}
			categoryService := &MockCMSCategoryService{}
			categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
			
			err := categoryTypeService.DeleteCategoryType(mockCategoryType.ID.String())
			assert.Error(t, err)
			assert.Equal(t, "category type not found for delete", err.Error())
	})

	t.Run("failed to delete category type: category type is in use", func(t *testing.T) {
			mockCategoryType, _ := helpers.InitializeMockCategory()
			mockCategoryType.ID = uuid.New()

			categoryTypeRepo := &MockCMSCategoryTypeRepo{
					delete: func(id uuid.UUID) error {
							return nil
					},
			}
			categoryRepo := &MockCMSCategoryRepo{
					countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
							return map[string]int{
									string(enums.PageLanguageTH): 2,
									string(enums.PageLanguageEN): 1,
							}, nil
					},
			}
			categoryService := &MockCMSCategoryService{}
			categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
			
			err := categoryTypeService.DeleteCategoryType(mockCategoryType.ID.String())
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot delete category type: it is still in use by 3 categories")
	})

	t.Run("failed to delete category type: internal server error", func(t *testing.T) {
			mockCategoryType, _ := helpers.InitializeMockCategory()
			mockCategoryType.ID = uuid.New()

			categoryTypeRepo := &MockCMSCategoryTypeRepo{
					delete: func(id uuid.UUID) error {
							return errs.ErrInternalServerError
					},
			}
			categoryRepo := &MockCMSCategoryRepo{
					countCategoriesByTypeAndLanguage: func(categoryTypeID uuid.UUID) (map[string]int, error) {
							return map[string]int{
									string(enums.PageLanguageTH): 0,
									string(enums.PageLanguageEN): 0,
							}, nil
					},
			}
			categoryService := &MockCMSCategoryService{}
			categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
			
			err := categoryTypeService.DeleteCategoryType(mockCategoryType.ID.String())
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to delete category type")
	})
}

func TestCMSService_GetCategoryTypeWithDetails(t *testing.T) {
	t.Run("successfully get category type with details", func(t *testing.T) {
		mockCategoryType, mockCategory := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeResponse := dto.CategoryResponse{
			ID:         mockCategory.ID.String(),
			Name:       mockCategory.Name,
			CreatedAt:  mockCategory.CreatedAt,
			UpdatedAt:  mockCategory.UpdatedAt,
		}

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)				

		categoryService.On("ListAllCategories", mock.AnythingOfType("dto.CategoryFilter")).Return([]dto.CategoryResponse{categoryTypeResponse}, nil)
		
		categoryTypeWithDetailResponse, err := categoryTypeService.GetCategoryTypeWithDetails(mockCategoryType.ID.String(), "th")
		assert.NoError(t, err)
		assert.NotNil(t, categoryTypeWithDetailResponse)
		assert.Equal(t, mockCategoryType.ID.String(), categoryTypeWithDetailResponse.ID)
		assert.Equal(t, mockCategoryType.Name, *categoryTypeWithDetailResponse.Name)
		assert.Equal(t, mockCategoryType.IsActive, categoryTypeWithDetailResponse.IsActive)
		assert.Equal(t, mockCategoryType.CreatedAt, categoryTypeWithDetailResponse.CreatedAt)
		assert.Equal(t, mockCategoryType.UpdatedAt, categoryTypeWithDetailResponse.UpdatedAt)
		assert.Equal(t, []dto.CategoryResponse{categoryTypeResponse}, categoryTypeWithDetailResponse.Categories)
	})

	t.Run("failed get category type with details: internal server error", func(t *testing.T) {
		mockCategoryType, _ := helpers.InitializeMockCategory()
		mockCategoryType.ID = uuid.New()

		categoryTypeRepo := &MockCMSCategoryTypeRepo{
			findByID: func(id uuid.UUID) (*models.CategoryType, error) {
				return mockCategoryType, nil
			},
		}
		categoryRepo := &MockCMSCategoryRepo{}
		categoryService := &MockCMSCategoryService{}
		categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)				

		categoryService.ExpectedCalls = nil
		categoryService.On("ListAllCategories", mock.AnythingOfType("dto.CategoryFilter")).Return(nil, errs.ErrInternalServerError)
		
		categoryTypeWithDetailResponse, err := categoryTypeService.GetCategoryTypeWithDetails(mockCategoryType.ID.String(), "th")
		assert.Error(t, err)
		assert.Nil(t, categoryTypeWithDetailResponse)
	})	
}