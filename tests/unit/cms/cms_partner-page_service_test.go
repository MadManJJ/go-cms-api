package tests

import (
	"strings"
	"testing"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockCMSPartnerPageRepo struct {
	createPartnerPage                        func(partnerPage *models.PartnerPage) (*models.PartnerPage, error)
	findAllPartnerPage                       func(query dto.PartnerPageQuery, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error)
	findPartnerPageById                      func(id uuid.UUID) (*models.PartnerPage, error)
	updatePartnerContent                     func(updatePartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error)
	deletePartnerPage                        func(id uuid.UUID) error
	findContentByPartnerPageId               func(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error)
	findLatestContentByPageId                func(pageId uuid.UUID, language string) (*models.PartnerContent, error)
	createContentForPartnerPage              func(partnerContent *models.PartnerContent, lang string, mode string) (*models.PartnerContent, error) // Deprecated
	deletePartnerContent                     func(pageId uuid.UUID, lang, mode string) error
	duplicatePartnerPage                     func(pageId uuid.UUID) (*models.PartnerPage, error)
	duplicatePartnerContentToAnotherLanguage func(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error)
	revertPartnerContent                     func(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error)
	getCategory                              func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	getRevisionByPartnerPageId               func(pageId uuid.UUID, language string) ([]models.Revision, error)
	isUrlDuplicate                           func(url string, pageId uuid.UUID) (bool, error)
	isUrlAliasDuplicate                      func(urlAlias string, pageId uuid.UUID) (bool, error)
	getPageIdByContentId                     func(contentId uuid.UUID) (uuid.UUID, error)
	createPartnerContentPreview               func(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error)
	updatePartnerContentPreview               func(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error)	
	findPartnerContentPreviewById             func(pageId uuid.UUID, language string) (*models.PartnerContent, error)	
}

func (m *MockCMSPartnerPageRepo) CreatePartnerPage(partnerPage *models.PartnerPage) (*models.PartnerPage, error) {
	return m.createPartnerPage(partnerPage)
}

func (m *MockCMSPartnerPageRepo) FindAllPartnerPage(query dto.PartnerPageQuery, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error) {
	return m.findAllPartnerPage(query, sort, page, limit, language)
}

func (m *MockCMSPartnerPageRepo) FindPartnerPageById(id uuid.UUID) (*models.PartnerPage, error) {
	return m.findPartnerPageById(id)
}

func (m *MockCMSPartnerPageRepo) UpdatePartnerContent(updatePartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error) {
	return m.updatePartnerContent(updatePartnerContent, prevContentId)
}

func (m *MockCMSPartnerPageRepo) DeletePartnerPage(id uuid.UUID) error {
	return m.deletePartnerPage(id)
}

func (m *MockCMSPartnerPageRepo) FindContentByPartnerPageId(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error) {
	return m.findContentByPartnerPageId(pageId, language, mode)
}

func (m *MockCMSPartnerPageRepo) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
	return m.findLatestContentByPageId(pageId, language)
}

func (m *MockCMSPartnerPageRepo) CreateContentForPartnerPage(partnerContent *models.PartnerContent, lang string, mode string) (*models.PartnerContent, error) {
	return m.createContentForPartnerPage(partnerContent, lang, mode)
}

func (m *MockCMSPartnerPageRepo) DeletePartnerContent(pageId uuid.UUID, lang, mode string) error {
	return m.deletePartnerContent(pageId, lang, mode)
}

func (m *MockCMSPartnerPageRepo) DuplicatePartnerPage(pageId uuid.UUID) (*models.PartnerPage, error) {
	return m.duplicatePartnerPage(pageId)
}

func (m *MockCMSPartnerPageRepo) DuplicatePartnerContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	return m.duplicatePartnerContentToAnotherLanguage(contentId, newRevision)
}

func (m *MockCMSPartnerPageRepo) RevertPartnerContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	return m.revertPartnerContent(revisionId, newRevision)
}

func (m *MockCMSPartnerPageRepo) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	return m.getCategory(pageId, categoryTypeCode, language, mode)
}

func (m *MockCMSPartnerPageRepo) GetRevisionByPartnerPageId(pageId uuid.UUID, language string) ([]models.Revision, error) {
	return m.getRevisionByPartnerPageId(pageId, language)
}

func (m *MockCMSPartnerPageRepo) IsUrlDuplicate(url string, pageId uuid.UUID) (bool, error) {
	return m.isUrlDuplicate(url, pageId)
}

func (m *MockCMSPartnerPageRepo) IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error) {
	return m.isUrlAliasDuplicate(urlAlias, pageId)
}

func (m *MockCMSPartnerPageRepo) GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error) {
	return m.getPageIdByContentId(contentId)
}

func (m *MockCMSPartnerPageRepo) CreatePartnerContentPreview(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error) {
	return m.createPartnerContentPreview(partnerContentPreview)
}

func (m *MockCMSPartnerPageRepo) UpdatePartnerContentPreview(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error) {
	return m.updatePartnerContentPreview(partnerContentPreview)
}

func (m *MockCMSPartnerPageRepo) FindPartnerContentPreviewById(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
	return m.findPartnerContentPreviewById(pageId, language)
}


func TestCMSService_CreatePartnerPage(t *testing.T) {
	t.Run("successfully create partner page", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()

		pageId := uuid.New()

		createdPartnerPage := mockPartnerPage
		createdPartnerPage.ID = pageId

		partnerRepo := &MockCMSPartnerPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createPartnerPage: func(partnerPage *models.PartnerPage) (*models.PartnerPage, error) {
				return createdPartnerPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualPartnerPage, err := service.CreatePartnerPage(mockPartnerPage)
		assert.NoError(t, err)
		assert.Equal(t, createdPartnerPage, actualPartnerPage)
	})

	t.Run("failed to create partner page: url is duplicated", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()

		pageId := uuid.New()

		createdPartnerPage := mockPartnerPage
		createdPartnerPage.ID = pageId

		partnerRepo := &MockCMSPartnerPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createPartnerPage: func(partnerPage *models.PartnerPage) (*models.PartnerPage, error) {
				return createdPartnerPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualPartnerPage, err := service.CreatePartnerPage(mockPartnerPage)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPage)
	})	

	t.Run("failed to create partner page: url alias is duplicated", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()

		pageId := uuid.New()

		createdPartnerPage := mockPartnerPage
		createdPartnerPage.ID = pageId

		partnerRepo := &MockCMSPartnerPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},			
			createPartnerPage: func(partnerPage *models.PartnerPage) (*models.PartnerPage, error) {
				return createdPartnerPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualPartnerPage, err := service.CreatePartnerPage(mockPartnerPage)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPage)
	})		

	t.Run("failed to create partner page: internal server error", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()

		pageId := uuid.New()

		createdPartnerPage := mockPartnerPage
		createdPartnerPage.ID = pageId

		partnerRepo := &MockCMSPartnerPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createPartnerPage: func(partnerPage *models.PartnerPage) (*models.PartnerPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualPartnerPage, err := service.CreatePartnerPage(mockPartnerPage)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPage)
	})	
}

func TestCMSService_FindPartnerPages(t *testing.T) {
	// Arguments
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
	sort := "sort test"
	page := 1
	limit := 10
	language := "en"	

	t.Run("successfully find partner page", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockPartnerPages := []models.PartnerPage{*mockPartnerPage, *mockPartnerPage}

		partnerRepo := &MockCMSPartnerPageRepo{
			findAllPartnerPage: func(query dto.PartnerPageQuery, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error) {
				return mockPartnerPages, 2, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualPartnerPages, totalCount, err := service.FindPartnerPages(rawQuery, sort, page, limit, language)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), totalCount)
		assert.Equal(t, mockPartnerPages, actualPartnerPages)
		assert.Equal(t, 2, len(actualPartnerPages))		
	})

	t.Run("failed to find partner page", func(t *testing.T) {
		partnerRepo := &MockCMSPartnerPageRepo{
			findAllPartnerPage: func(query dto.PartnerPageQuery, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error) {
				return nil, 0, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualPartnerPages, totalCount, err := service.FindPartnerPages(rawQuery, sort, page, limit, language)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPages)	
		assert.Equal(t, int64(0), totalCount)
	})	
}

func TestCMSService_FindPartnerPageById(t *testing.T) {
	t.Run("successfully find partner page", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		pageId := uuid.New()

		partnerRepo := &MockCMSPartnerPageRepo{
			findPartnerPageById: func(id uuid.UUID) (*models.PartnerPage, error) {
				return mockPartnerPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualPartnerPage, err := service.FindPartnerPageById(pageId)
		assert.NoError(t, err)
		assert.Equal(t, mockPartnerPage, actualPartnerPage)	
	})

	t.Run("successfully find partner page", func(t *testing.T) {
		pageId := uuid.New()

		partnerRepo := &MockCMSPartnerPageRepo{
			findPartnerPageById: func(id uuid.UUID) (*models.PartnerPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualPartnerPage, err := service.FindPartnerPageById(pageId)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPage)
	})	
}

func TestCMSService_UpdatePartnerContent(t *testing.T) {
	t.Run("successfully update partner content", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		emailContent := helpers.InitializeMockEmailContent()
		emailCategory := emailContent.EmailCategory

		pageId := uuid.New()
		contentId := uuid.New()

		updatedContent := mockContent
		updatedContent.ID = contentId
		updatedContent.PageID = pageId

		partnerRepo := &MockCMSPartnerPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			updatePartnerContent: func(updatePartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error) {
				return updatedContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return emailContent, nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return emailCategory, nil
			},
			findByTitle: func(title string) (*models.EmailCategory, error) {
				return emailCategory, nil
			},
		}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualPartnerPage, err := service.UpdatePartnerContent(mockContent, contentId)
		assert.NoError(t, err)
		assert.Equal(t, updatedContent, actualPartnerPage)
	})	

	t.Run("failed to update partner content", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]

		pageId := uuid.New()
		contentId := uuid.New()

		updatedContent := mockContent
		updatedContent.ID = contentId
		updatedContent.PageID = pageId

		partnerRepo := &MockCMSPartnerPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			updatePartnerContent: func(updatePartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualPartnerPage, err := service.UpdatePartnerContent(mockContent, contentId)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerPage)
	})		
}

func TestCMSService_DeletePartnerPage(t *testing.T) {
	t.Run("successfully delete partner page", func(t *testing.T) {
		pageId := uuid.New()

		partnerRepo := &MockCMSPartnerPageRepo{
			deletePartnerPage: func(id uuid.UUID) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		err := service.DeletePartnerPage(pageId)
		assert.NoError(t, err)
	})	

	t.Run("failed to delete partner page", func(t *testing.T) {
		pageId := uuid.New()

		partnerRepo := &MockCMSPartnerPageRepo{
			deletePartnerPage: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		err := service.DeletePartnerPage(pageId)
		assert.Error(t, err)
	})		
}

func TestCMSService_FindContentByPartnerPageId(t *testing.T) {
	t.Run("successfully find partner content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		partnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := partnerPage.Contents[0]

		partnerRepo := &MockCMSPartnerPageRepo{
			findContentByPartnerPageId: func(pageId uuid.UUID, language, mode string) (*models.PartnerContent, error) {
				return partnerContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualPartnerContent, err := service.FindContentByPartnerPageId(pageId, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, partnerContent, actualPartnerContent)
	})

	t.Run("failed to find partner content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		partnerRepo := &MockCMSPartnerPageRepo{
			findContentByPartnerPageId: func(pageId uuid.UUID, language, mode string) (*models.PartnerContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualPartnerContent, err := service.FindContentByPartnerPageId(pageId, language, mode)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerContent)
	})	
}

func TestCMSService_FindLatestPartnerContentByPageId(t *testing.T) {
	t.Run("successfully find latest partner content by page id", func(t *testing.T) {
		pageId := uuid.New()
		partnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := partnerPage.Contents[0]
		language := string(enums.PageLanguageEN)

		partnerRepo := &MockCMSPartnerPageRepo{
			findLatestContentByPageId: func(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
				return partnerContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualPartnerContent, err := service.FindLatestContentByPageId(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerContent, actualPartnerContent)
	})

	t.Run("failed to find latest partner content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		partnerRepo := &MockCMSPartnerPageRepo{
			findLatestContentByPageId: func(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualPartnerContent, err := service.FindLatestContentByPageId(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerContent)
	})	
}

func TestCMSService_DeleteContentByPartnerPageId(t *testing.T) {
	t.Run("successfully delete partner content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		partnerRepo := &MockCMSPartnerPageRepo{
			deletePartnerContent: func(pageId uuid.UUID, lang, mode string) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		err := service.DeleteContentByPartnerPageId(pageId, language, mode)
		assert.NoError(t, err)
	})	

	t.Run("successfully delete partner content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		partnerRepo := &MockCMSPartnerPageRepo{
			deletePartnerContent: func(pageId uuid.UUID, lang, mode string) error {
				return errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		err := service.DeleteContentByPartnerPageId(pageId, language, mode)
		assert.Error(t, err)
	})		
}

func TestCMSService_DuplicatePartnerPage(t *testing.T) {
	t.Run("successfully duplicate partner page", func(t *testing.T) {
		pageId := uuid.New()

		mockPartnerPage := helpers.InitializeMockPartnerPage()

		partnerRepo := &MockCMSPartnerPageRepo{
			duplicatePartnerPage: func(pageId uuid.UUID) (*models.PartnerPage, error) {
				return mockPartnerPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		partnerPage, err := service.DuplicatePartnerPage(pageId)
		assert.NoError(t, err)
		assert.Equal(t, mockPartnerPage, partnerPage)
	})	

	t.Run("failed to duplicate partner page", func(t *testing.T) {
		pageId := uuid.New()

		partnerRepo := &MockCMSPartnerPageRepo{
			duplicatePartnerPage: func(pageId uuid.UUID) (*models.PartnerPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		partnerPage, err := service.DuplicatePartnerPage(pageId)
		assert.Error(t, err)
		assert.Nil(t, partnerPage)
	})		
}

func TestCMSService_DuplicatePartnerContentToAnotherLanguage(t *testing.T) {
	t.Run("successfully duplicate partner content to another language", func(t *testing.T) {
		pageId := uuid.New()

		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		mockRevision := mockContent.Revision

		partnerRepo := &MockCMSPartnerPageRepo{
			duplicatePartnerContentToAnotherLanguage: func(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
				return mockContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)			

		partnerContent, err := service.DuplicatePartnerContentToAnotherLanguage(pageId, mockRevision)
		assert.NoError(t, err)
		assert.Equal(t, mockContent, partnerContent)
	})	

	t.Run("failed to duplicate partner content to another language", func(t *testing.T) {
		pageId := uuid.New()

		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]
		mockRevision := mockContent.Revision

		partnerRepo := &MockCMSPartnerPageRepo{
			duplicatePartnerContentToAnotherLanguage: func(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		partnerContent, err := service.DuplicatePartnerContentToAnotherLanguage(pageId, mockRevision)
		assert.Error(t, err)
		assert.Nil(t, partnerContent)
	})		
}

func TestCMSService_RevertPartnerContent(t *testing.T) {
	t.Run("successfully revert partner content", func(t *testing.T) {
		revisionId := uuid.New()
		partnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := partnerPage.Contents[0]
		revision := partnerContent.Revision

		partnerRepo := &MockCMSPartnerPageRepo{
			revertPartnerContent: func(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
				return partnerContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualPartnerContent, err := service.RevertPartnerContent(revisionId, revision)
		assert.NoError(t, err)
		assert.Equal(t, partnerContent, actualPartnerContent)
	})	

	t.Run("failed to revert partner content", func(t *testing.T) {
		revisionId := uuid.New()
		partnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := partnerPage.Contents[0]
		revision := partnerContent.Revision

		partnerRepo := &MockCMSPartnerPageRepo{
			revertPartnerContent: func(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualPartnerContent, err := service.RevertPartnerContent(revisionId, revision)
		assert.Error(t, err)
		assert.Nil(t, actualPartnerContent)
	})	
}

func TestCMSService_GetPartnerCategory(t *testing.T) {
	t.Run("successfully get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		
		categoryTypeCode := "some code"		

		partnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := partnerPage.Contents[0]
		category := partnerContent.Categories[0]
		categories := []models.Category{*category, *category}

		partnerRepo := &MockCMSPartnerPageRepo{
			getCategory: func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
				return categories, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualCategories, err := service.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, categories, actualCategories)
	})
	
	t.Run("failed to get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		
		categoryTypeCode := "some code"		

		partnerRepo := &MockCMSPartnerPageRepo{
			getCategory: func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualCategories, err := service.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.Error(t, err)
		assert.Nil(t, actualCategories)
	})	
}

func TestCMSService_FindPartnerRevisions(t *testing.T) {
	t.Run("successfully get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		partnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := partnerPage.Contents[0]
		revision := partnerContent.Revision
		revisions := []models.Revision{*revision, *revision}

		partnerRepo := &MockCMSPartnerPageRepo{
			getRevisionByPartnerPageId: func(pageId uuid.UUID, language string) ([]models.Revision, error) {
				return revisions, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualRevisions, err := service.FindRevisions(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, revisions, actualRevisions)
	})	

	t.Run("failed to get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		partnerRepo := &MockCMSPartnerPageRepo{
			getRevisionByPartnerPageId: func(pageId uuid.UUID, language string) ([]models.Revision, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualRevisions, err := service.FindRevisions(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, actualRevisions)
	})		
}

func TestCMSService_PreviewPartnerContent(t *testing.T) {
	cfg := config.New()

	urls := cfg.App.FrontendURLS
	parts := strings.Split(urls, ",")
	appUrl := parts[1]	

	t.Run("successfully preview content: create new content", func(t *testing.T) {
		pageId := uuid.New()

		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]

		createdContent := mockContent
		createdContent.ID = uuid.New()

		partnerRepo := &MockCMSPartnerPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			findPartnerContentPreviewById: func(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
				return nil, gorm.ErrRecordNotFound
			},
			createPartnerContentPreview: func(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error) {
				return createdContent, nil
			},
		}

		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()		
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		expectPreviewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(mockContent.Language),
			"partner",
			createdContent.ID,
		)

		assert.NoError(t, err)

		previewUrl, err := service.PreviewPartnerContent(pageId, mockContent)
		assert.NoError(t, err)
		assert.Equal(t, expectPreviewUrl, previewUrl)
	})	

	t.Run("successfully preview content: update existing content", func(t *testing.T) {
		pageId := uuid.New()

		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockContent := mockPartnerPage.Contents[0]

		contentId := uuid.New()

		existingContent := mockContent
		existingContent.ID = contentId

		updatedContent := existingContent
		updatedContent.Title = "Updated Title"
		updatedContent.URL = "Updated URL"
		updatedContent.URLAlias = "Updated URL Alias"

		partnerRepo := &MockCMSPartnerPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			findPartnerContentPreviewById: func(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
				return existingContent, nil
			},
			updatePartnerContentPreview: func(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error) {
				return updatedContent, nil
			},
		}

		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()		
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSPartnerPageService(partnerRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		expectPreviewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(updatedContent.Language),
			"partner",
			updatedContent.ID,
		)

		assert.NoError(t, err)

		previewUrl, err := service.PreviewPartnerContent(pageId, mockContent)
		assert.NoError(t, err)
		assert.Equal(t, expectPreviewUrl, previewUrl)
	})		
}