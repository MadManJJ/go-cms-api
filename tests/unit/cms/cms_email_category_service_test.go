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

type MockCMSEmailCategoryRepo struct {
	create        func(category *models.EmailCategory) (*models.EmailCategory, error)
	findByID      func(id uuid.UUID) (*models.EmailCategory, error)
	findByTitle   func(title string) (*models.EmailCategory, error)
	findAll       func() ([]models.EmailCategory, error)
	update        func(category *models.EmailCategory) (*models.EmailCategory, error)
	delete        func(id uuid.UUID) error
	isTitleUnique func(title string, excludeID uuid.UUID) (bool, error)
}

func (m *MockCMSEmailCategoryRepo) Create(category *models.EmailCategory) (*models.EmailCategory, error) {
	return m.create(category)
}

func (m *MockCMSEmailCategoryRepo) FindByID(id uuid.UUID) (*models.EmailCategory, error) {
	return m.findByID(id)
}

func (m *MockCMSEmailCategoryRepo) FindByTitle(title string) (*models.EmailCategory, error) {
	return m.findByTitle(title)
}

func (m *MockCMSEmailCategoryRepo) FindAll() ([]models.EmailCategory, error) {
	return m.findAll()
}

func (m *MockCMSEmailCategoryRepo) Update(category *models.EmailCategory) (*models.EmailCategory, error) {
	return m.update(category)
}

func (m *MockCMSEmailCategoryRepo) Delete(id uuid.UUID) error {
	return m.delete(id)
}

func (m *MockCMSEmailCategoryRepo) IsTitleUnique(title string, excludeID uuid.UUID) (bool, error) {
	return m.isTitleUnique(title, excludeID)
}

func TestCMSService_CreateEmailCategory(t *testing.T) {
	mockEmailCategory := helpers.InitializeMockEmailCategory()

	createEmailCategoryRequest := dto.CreateEmailCategoryRequest{
		Title: mockEmailCategory.Title,
	}

	t.Run("successfully create email category", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			isTitleUnique: func(title string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
			create: func(category *models.EmailCategory) (*models.EmailCategory, error) {
				return mockEmailCategory, nil 
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.CreateCategory(createEmailCategoryRequest)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.Equal(t, mockEmailCategory.Title, actualEmailCategory.Title)
	})

	t.Run("failed to create email category: title is not unique", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			isTitleUnique: func(title string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.CreateCategory(createEmailCategoryRequest)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})

	t.Run("failed to create email category: internal server error", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			isTitleUnique: func(title string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
			create: func(category *models.EmailCategory) (*models.EmailCategory, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.CreateCategory(createEmailCategoryRequest)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})
}

func TestCMSService_GetEmailCategoryByID(t *testing.T) {
	mockEmailCategory := helpers.InitializeMockEmailCategory()
	mockEmailCategory.ID = uuid.New()

	emailCategoryResponse := &dto.EmailCategoryResponse{
		ID:        mockEmailCategory.ID.String(),
		Title:     mockEmailCategory.Title,
		CreatedAt: mockEmailCategory.CreatedAt,
		UpdatedAt: mockEmailCategory.UpdatedAt,
	}

	t.Run("successfully get email category by ID", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return mockEmailCategory, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.GetCategoryByID(mockEmailCategory.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.Equal(t, emailCategoryResponse, actualEmailCategory)
	})

	t.Run("failed to get email category by ID: not found", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.GetCategoryByID(mockEmailCategory.ID.String())
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})

	t.Run("failed to get email category by ID: internal server error", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.GetCategoryByID(mockEmailCategory.ID.String())
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})
}

func TestCMSService_GetEmailCategoryByTitle(t *testing.T) {
	emailCategoryId := uuid.New()

	emailCategory := helpers.InitializeMockEmailCategory()
	emailCategory.ID = emailCategoryId
	
	t.Run("successfully get email category by title", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByTitle: func(title string) (*models.EmailCategory, error) {
				return emailCategory, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.GetCategoryByTitle(emailCategory.Title)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.Equal(t, actualEmailCategory.ID, emailCategory.ID.String())
	})

	t.Run("failed to get email category by title: not found", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByTitle: func(title string) (*models.EmailCategory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.GetCategoryByTitle(emailCategory.Title)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})

	t.Run("failed to get email category by title: internal server error", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByTitle: func(title string) (*models.EmailCategory, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.GetCategoryByTitle(emailCategory.Title)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})
}

func TestCMSService_ListEmailCategories(t *testing.T) {
	emailCategory := helpers.InitializeMockEmailCategory()
	emailCategory.ID = uuid.New()

	t.Run("successfully list email categories", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findAll: func() ([]models.EmailCategory, error) {
				return []models.EmailCategory{*emailCategory}, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategories, err := service.ListCategories()
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategories)
		assert.Equal(t, actualEmailCategories[0].ID, emailCategory.ID.String())
	})

	t.Run("failed to list email categories: internal server error", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findAll: func() ([]models.EmailCategory, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategories, err := service.ListCategories()
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategories)
	})
}

func TestCMSService_UpdateEmailCateory(t *testing.T) {
	emailCategory := helpers.InitializeMockEmailCategory()
	emailCategory.ID = uuid.New()
	
	updatedEmailCategory := *emailCategory
	updatedEmailCategory.Title = "Updated Title"

	updateEmailCategoryRequest := dto.UpdateEmailCategoryRequest{
		Title: updatedEmailCategory.Title,
	}

	t.Run("successfully update email category", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return emailCategory, nil
			},
			isTitleUnique: func(title string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
			update: func(category *models.EmailCategory) (*models.EmailCategory, error) {
				return &updatedEmailCategory, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.UpdateCategory(emailCategory.ID.String(), updateEmailCategoryRequest)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.Equal(t, actualEmailCategory.ID, emailCategory.ID.String())
	})

	t.Run("failed to update email category: not found", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return emailCategory, nil
			},
			isTitleUnique: func(title string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},			
			update: func(category *models.EmailCategory) (*models.EmailCategory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.UpdateCategory(emailCategory.ID.String(), updateEmailCategoryRequest)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})

	t.Run("failed to update email category: internal server error", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return emailCategory, nil
			},
			isTitleUnique: func(title string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},			
			update: func(category *models.EmailCategory) (*models.EmailCategory, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		actualEmailCategory, err := service.UpdateCategory(emailCategory.ID.String(), updateEmailCategoryRequest)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
	})
}

func TestCMSService_DeleteEmailCategory(t *testing.T) {
	emailCategory := helpers.InitializeMockEmailCategory()
	emailCategory.ID = uuid.New()
	
	t.Run("successfully delete email category", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			delete: func(id uuid.UUID) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{
			deleteByCategoryID: func(id uuid.UUID) error {
				return nil
			},
		}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		err := service.DeleteCategory(emailCategory.ID.String())
		assert.NoError(t, err)
	})

	t.Run("successfully delete email category with no associated email contents", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			delete: func(id uuid.UUID) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{
			deleteByCategoryID: func(id uuid.UUID) error {
				return gorm.ErrRecordNotFound
			},
		}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		err := service.DeleteCategory(emailCategory.ID.String())
		assert.NoError(t, err)
	})	

	t.Run("failed to delete email category: not found", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			delete: func(id uuid.UUID) error {
				return gorm.ErrRecordNotFound
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{
			deleteByCategoryID: func(id uuid.UUID) error {
				return nil
			},
		}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		err := service.DeleteCategory(emailCategory.ID.String())
		assert.Error(t, err)
	})

	t.Run("failed to delete email category: internal server error", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			delete: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{
			deleteByCategoryID: func(id uuid.UUID) error {
				return nil
			},
		}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		err := service.DeleteCategory(emailCategory.ID.String())
		assert.Error(t, err)
	})

	t.Run("failed to delete email category: internal server error at email content", func(t *testing.T) {
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			delete: func(id uuid.UUID) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{
			deleteByCategoryID: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
	
		service := services.NewEmailCategoryService(emailCategoryRepo, emailContentRepo)

		err := service.DeleteCategory(emailCategory.ID.String())
		assert.Error(t, err)
	})
}