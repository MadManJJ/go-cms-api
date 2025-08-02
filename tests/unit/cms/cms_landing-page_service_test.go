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

type MockCMSLandingPageRepo struct {
	createLandingPage                        func(landingPage *models.LandingPage) (*models.LandingPage, error)
	findAllLandingPage                       func(query dto.LandingPageQuery, sort string, page, limit int, language string) ([]models.LandingPage, int64, error)
	findLandingPageById                      func(id uuid.UUID) (*models.LandingPage, error)
	updateLandingContent                     func(updateLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error)
	deleteLandingPage                        func(id uuid.UUID) error
	findContentByLandingPageId               func(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error)
	findLatestContentByPageId                func(pageId uuid.UUID, language string) (*models.LandingContent, error)
	createContentForLandingPage              func(landingContent *models.LandingContent, lang string, mode string) (*models.LandingContent, error)
	deleteLandingContent                     func(pageId uuid.UUID, lang, mode string) error
	duplicateLandingPage                     func(pageId uuid.UUID) (*models.LandingPage, error)
	duplicateLandingContentToAnotherLanguage func(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error)
	revertLandingContent                     func(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error)
	getCategory                              func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	getRevisionByLandingPageId               func(pageId uuid.UUID, language string) ([]models.Revision, error)
	isUrlAliasDuplicate                      func(urlAlias string, pageId uuid.UUID) (bool, error)
	getPageIdByContentId                     func(contentId uuid.UUID) (uuid.UUID, error)
	createLandingContentPreview               func(landingContentPreview *models.LandingContent) (*models.LandingContent, error)
	updateLandingContentPreview               func(landingContentPreview *models.LandingContent) (*models.LandingContent, error)	
	findLandingContentPreviewById             func(pageId uuid.UUID, language string) (*models.LandingContent, error)
}

func (m *MockCMSLandingPageRepo) CreateLandingPage(landingPage *models.LandingPage) (*models.LandingPage, error) {
	return m.createLandingPage(landingPage)
}

func (m *MockCMSLandingPageRepo) FindAllLandingPage(query dto.LandingPageQuery, sort string, page, limit int, language string) ([]models.LandingPage, int64, error) {
	return m.findAllLandingPage(query, sort, page, limit, language)
}

func (m *MockCMSLandingPageRepo) FindLandingPageById(id uuid.UUID) (*models.LandingPage, error) {
	return m.findLandingPageById(id)
}

func (m *MockCMSLandingPageRepo) UpdateLandingContent(updateLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error) {
	return m.updateLandingContent(updateLandingContent, prevContentId)
}

func (m *MockCMSLandingPageRepo) DeleteLandingPage(id uuid.UUID) error {
	return m.deleteLandingPage(id)
}

func (m *MockCMSLandingPageRepo) FindContentByLandingPageId(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error) {
	return m.findContentByLandingPageId(pageId, language, mode)
}

func (m *MockCMSLandingPageRepo) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.LandingContent, error) {
	return m.findLatestContentByPageId(pageId, language)
}

func (m *MockCMSLandingPageRepo) CreateContentForLandingPage(landingContent *models.LandingContent, lang string, mode string) (*models.LandingContent, error) {
	return m.createContentForLandingPage(landingContent, lang, mode)
}

func (m *MockCMSLandingPageRepo) DeleteLandingContent(pageId uuid.UUID, lang, mode string) error {
	return m.deleteLandingContent(pageId, lang, mode)
}

func (m *MockCMSLandingPageRepo) DuplicateLandingPage(pageId uuid.UUID) (*models.LandingPage, error) {
	return m.duplicateLandingPage(pageId)
}

func (m *MockCMSLandingPageRepo) DuplicateLandingContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	return m.duplicateLandingContentToAnotherLanguage(contentId, newRevision)
}

func (m *MockCMSLandingPageRepo) RevertLandingContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	return m.revertLandingContent(revisionId, newRevision)
}

func (m *MockCMSLandingPageRepo) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	return m.getCategory(pageId, categoryTypeCode, language, mode)
}

func (m *MockCMSLandingPageRepo) GetRevisionByLandingPageId(pageId uuid.UUID, language string) ([]models.Revision, error) {
	return m.getRevisionByLandingPageId(pageId, language)
}

func (m *MockCMSLandingPageRepo) IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error) {
	return m.isUrlAliasDuplicate(urlAlias, pageId)
}

func (m *MockCMSLandingPageRepo) GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error) {
	return m.getPageIdByContentId(contentId)
}

func (m *MockCMSLandingPageRepo) CreateLandingContentPreview(landingContentPreview *models.LandingContent) (*models.LandingContent, error) {
	return m.createLandingContentPreview(landingContentPreview)
}

func (m *MockCMSLandingPageRepo) UpdateLandingContentPreview(landingContentPreview *models.LandingContent) (*models.LandingContent, error) {
	return m.updateLandingContentPreview(landingContentPreview)
}

func (m *MockCMSLandingPageRepo) FindLandingContentPreviewById(pageId uuid.UUID, language string) (*models.LandingContent, error) {
	return m.findLandingContentPreviewById(pageId, language)
}

func TestCMSService_CreateLandingPage(t *testing.T) {
	t.Run("successfully create landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()

		pageId := uuid.New()

		createdLandingPage := mockLandingPage
		createdLandingPage.ID = pageId

		landingRepo := &MockCMSLandingPageRepo{
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createLandingPage: func(landingPage *models.LandingPage) (*models.LandingPage, error) {
				return createdLandingPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualLandingPage, err := service.CreateLandingPage(mockLandingPage)
		assert.NoError(t, err)
		assert.Equal(t, createdLandingPage, actualLandingPage)
	})	

	t.Run("failed to create landing page: url alias is duplicated", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()

		pageId := uuid.New()

		createdLandingPage := mockLandingPage
		createdLandingPage.ID = pageId

		landingRepo := &MockCMSLandingPageRepo{
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},			
			createLandingPage: func(landingPage *models.LandingPage) (*models.LandingPage, error) {
				return createdLandingPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualLandingPage, err := service.CreateLandingPage(mockLandingPage)
		assert.Error(t, err)
		assert.Nil(t, actualLandingPage)
	})		

	t.Run("failed to create landing page: internal server error", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()

		pageId := uuid.New()

		createdLandingPage := mockLandingPage
		createdLandingPage.ID = pageId

		landingRepo := &MockCMSLandingPageRepo{
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createLandingPage: func(landingPage *models.LandingPage) (*models.LandingPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		actualLandingPage, err := service.CreateLandingPage(mockLandingPage)
		assert.Error(t, err)
		assert.Nil(t, actualLandingPage)
	})		
}

func TestCMSService_FindLandingPages(t *testing.T) {
	// Arguments
	rawQuery := `{
		"title": "Reset Password",
		"category_keywords": "security",
		"status": "Published"
	}`
	sort := "sort test"
	page := 1
	limit := 10
	language := "en"		

	t.Run("successfully find landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockLandingPages := []models.LandingPage{*mockLandingPage, *mockLandingPage}

		landingRepo := &MockCMSLandingPageRepo{
			findAllLandingPage: func(query dto.LandingPageQuery, sort string, page, limit int, language string) ([]models.LandingPage, int64, error) {
				return mockLandingPages, 2, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualLandingPages, totalCount, err := service.FindLandingPages(rawQuery, sort, page, limit, language)
		assert.NoError(t, err)
		assert.Equal(t, mockLandingPages, actualLandingPages)
		assert.Equal(t, int64(2), totalCount)		
	})	

	t.Run("failed to find landing page", func(t *testing.T) {
		landingRepo := &MockCMSLandingPageRepo{
			findAllLandingPage: func(query dto.LandingPageQuery, sort string, page, limit int, language string) ([]models.LandingPage, int64, error) {
				return nil, 0, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualLandingPages, totalCount, err := service.FindLandingPages(rawQuery, sort, page, limit, language)
		assert.Error(t, err)
		assert.Nil(t, actualLandingPages)	
		assert.Equal(t, int64(0), totalCount)
	})		
}

func TestCMSService_FindLandingPageById(t *testing.T) {
	t.Run("successfully find landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		pageId := uuid.New()

		landingRepo := &MockCMSLandingPageRepo{
			findLandingPageById: func(id uuid.UUID) (*models.LandingPage, error) {
				return mockLandingPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualLandingPage, err := service.FindLandingPageById(pageId)
		assert.NoError(t, err)
		assert.Equal(t, mockLandingPage, actualLandingPage)	
	})

	t.Run("successfully find landing page", func(t *testing.T) {
		pageId := uuid.New()

		landingRepo := &MockCMSLandingPageRepo{
			findLandingPageById: func(id uuid.UUID) (*models.LandingPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualLandingPage, err := service.FindLandingPageById(pageId)
		assert.Error(t, err)
		assert.Nil(t, actualLandingPage)
	})		
}

func TestCMSService_UpdateLandingContent(t *testing.T) {
	t.Run("successfully update landing content", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		emailContent := helpers.InitializeMockEmailContent()
		emailCategory := emailContent.EmailCategory

		pageId := uuid.New()
		contentId := uuid.New()

		updatedContent := mockContent
		updatedContent.ID = contentId
		updatedContent.PageID = pageId

		landingRepo := &MockCMSLandingPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			updateLandingContent: func(updateLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error) {
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
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualLandingPage, err := service.UpdateLandingContent(mockContent, contentId)
		assert.NoError(t, err)
		assert.Equal(t, updatedContent, actualLandingPage)
	})	

	t.Run("failed to update landing content", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]

		pageId := uuid.New()
		contentId := uuid.New()

		updatedContent := mockContent
		updatedContent.ID = contentId
		updatedContent.PageID = pageId

		landingRepo := &MockCMSLandingPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			updateLandingContent: func(updateLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		actualLandingPage, err := service.UpdateLandingContent(mockContent, contentId)
		assert.Error(t, err)
		assert.Nil(t, actualLandingPage)
	})		
}

func TestCMSService_DeleteLandingPage(t *testing.T) {
	t.Run("successfully delete landing page", func(t *testing.T) {
		pageId := uuid.New()

		landingRepo := &MockCMSLandingPageRepo{
			deleteLandingPage: func(id uuid.UUID) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		err := service.DeleteLandingPage(pageId)
		assert.NoError(t, err)
	})	

	t.Run("failed to delete landing page", func(t *testing.T) {
		pageId := uuid.New()

		landingRepo := &MockCMSLandingPageRepo{
			deleteLandingPage: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		err := service.DeleteLandingPage(pageId)
		assert.Error(t, err)
	})		
}

func TestCMSService_FindContentByLandingPageId(t *testing.T) {
	t.Run("successfully find landing content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		landingPage := helpers.InitializeMockLandingPage()
		landingContent := landingPage.Contents[0]

		landingRepo := &MockCMSLandingPageRepo{
			findContentByLandingPageId: func(pageId uuid.UUID, language, mode string) (*models.LandingContent, error) {
				return landingContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualLandingContent, err := service.FindContentByLandingPageId(pageId, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, landingContent, actualLandingContent)
	})

	t.Run("failed to find landing content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		landingRepo := &MockCMSLandingPageRepo{
			findContentByLandingPageId: func(pageId uuid.UUID, language, mode string) (*models.LandingContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualLandingContent, err := service.FindContentByLandingPageId(pageId, language, mode)
		assert.Error(t, err)
		assert.Nil(t, actualLandingContent)
	})	
}

func TestCMSService_FindLatestLandingContentByPageId(t *testing.T) {
	t.Run("successfully find latest landing content by page id", func(t *testing.T) {
		pageId := uuid.New()
		landingPage := helpers.InitializeMockLandingPage()
		landingContent := landingPage.Contents[0]
		language := string(enums.PageLanguageEN)

		landingRepo := &MockCMSLandingPageRepo{
			findLatestContentByPageId: func(pageId uuid.UUID, language string) (*models.LandingContent, error) {
				return landingContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualLandingContent, err := service.FindLatestContentByPageId(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, landingContent, actualLandingContent)
	})

	t.Run("failed to find latest landing content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		landingRepo := &MockCMSLandingPageRepo{
			findLatestContentByPageId: func(pageId uuid.UUID, language string) (*models.LandingContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualLandingContent, err := service.FindLatestContentByPageId(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, actualLandingContent)
	})	
}

func TestCMSService_DeleteContentByLandingPageId(t *testing.T) {
	t.Run("successfully delete landing content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		landingRepo := &MockCMSLandingPageRepo{
			deleteLandingContent: func(pageId uuid.UUID, lang, mode string) error {
				return nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		err := service.DeleteContentByLandingPageId(pageId, language, mode)
		assert.NoError(t, err)
	})	

	t.Run("successfully delete landing content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		

		landingRepo := &MockCMSLandingPageRepo{
			deleteLandingContent: func(pageId uuid.UUID, lang, mode string) error {
				return errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		err := service.DeleteContentByLandingPageId(pageId, language, mode)
		assert.Error(t, err)
	})	
}

func TestCMSService_DuplicateLandingPage(t *testing.T) {
	t.Run("successfully duplicate landing page", func(t *testing.T) {
		pageId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()

		landingRepo := &MockCMSLandingPageRepo{
			duplicateLandingPage: func(pageId uuid.UUID) (*models.LandingPage, error) {
				return mockLandingPage, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		landingPage, err := service.DuplicateLandingPage(pageId)
		assert.NoError(t, err)
		assert.Equal(t, mockLandingPage, landingPage)
	})	

	t.Run("failed to duplicate landing page", func(t *testing.T) {
		pageId := uuid.New()

		landingRepo := &MockCMSLandingPageRepo{
			duplicateLandingPage: func(pageId uuid.UUID) (*models.LandingPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)

		landingPage, err := service.DuplicateLandingPage(pageId)
		assert.Error(t, err)
		assert.Nil(t, landingPage)
	})		
}

func TestCMSService_DuplicateLandingContentToAnotherLanguage(t *testing.T) {
	t.Run("successfully duplicate landing content to another language", func(t *testing.T) {
		pageId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		mockRevision := mockContent.Revision

		landingRepo := &MockCMSLandingPageRepo{
			duplicateLandingContentToAnotherLanguage: func(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
				return mockContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)			

		landingContent, err := service.DuplicateLandingContentToAnotherLanguage(pageId, mockRevision)
		assert.NoError(t, err)
		assert.Equal(t, mockContent, landingContent)
	})	

	t.Run("failed to duplicate landing content to another language", func(t *testing.T) {
		pageId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		mockRevision := mockContent.Revision

		landingRepo := &MockCMSLandingPageRepo{
			duplicateLandingContentToAnotherLanguage: func(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		landingContent, err := service.DuplicateLandingContentToAnotherLanguage(pageId, mockRevision)
		assert.Error(t, err)
		assert.Nil(t, landingContent)
	})		
}

func TestCMSService_RevertLandingContent(t *testing.T) {
	t.Run("successfully revert landing content", func(t *testing.T) {
		revisionId := uuid.New()
		landingPage := helpers.InitializeMockLandingPage()
		landingContent := landingPage.Contents[0]
		revision := landingContent.Revision

		landingRepo := &MockCMSLandingPageRepo{
			revertLandingContent: func(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
				return landingContent, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualLandingContent, err := service.RevertLandingContent(revisionId, revision)
		assert.NoError(t, err)
		assert.Equal(t, landingContent, actualLandingContent)
	})	

	t.Run("failed to revert landing content", func(t *testing.T) {
		revisionId := uuid.New()
		landingPage := helpers.InitializeMockLandingPage()
		landingContent := landingPage.Contents[0]
		revision := landingContent.Revision

		landingRepo := &MockCMSLandingPageRepo{
			revertLandingContent: func(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualLandingContent, err := service.RevertLandingContent(revisionId, revision)
		assert.Error(t, err)
		assert.Nil(t, actualLandingContent)
	})	
}

func TestCMSService_GetLandingCategory(t *testing.T) {
	t.Run("successfully get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		
		categoryTypeCode := "some code"		

		landingPage := helpers.InitializeMockLandingPage()
		landingContent := landingPage.Contents[0]
		category := landingContent.Categories[0]
		categories := []models.Category{*category, *category}

		landingRepo := &MockCMSLandingPageRepo{
			getCategory: func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
				return categories, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualCategories, err := service.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, categories, actualCategories)
	})
	
	t.Run("failed to get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)		
		categoryTypeCode := "some code"		

		landingRepo := &MockCMSLandingPageRepo{
			getCategory: func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualCategories, err := service.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.Error(t, err)
		assert.Nil(t, actualCategories)
	})	
}

func TestCMSService_FindLandingRevisions(t *testing.T) {
	t.Run("successfully get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		landingPage := helpers.InitializeMockLandingPage()
		landingContent := landingPage.Contents[0]
		revision := landingContent.Revision
		revisions := []models.Revision{*revision, *revision}

		landingRepo := &MockCMSLandingPageRepo{
			getRevisionByLandingPageId: func(pageId uuid.UUID, language string) ([]models.Revision, error) {
				return revisions, nil
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualRevisions, err := service.FindRevisions(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, revisions, actualRevisions)
	})	

	t.Run("failed to get catories", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		landingRepo := &MockCMSLandingPageRepo{
			getRevisionByLandingPageId: func(pageId uuid.UUID, language string) ([]models.Revision, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)
	
		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	
		
		actualRevisions, err := service.FindRevisions(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, actualRevisions)
	})		
}

func TestCMSService_PreviewLandingContent(t *testing.T) {
	cfg := config.New()

	urls := cfg.App.FrontendURLS
	parts := strings.Split(urls, ",")
	appUrl := parts[1]	

	t.Run("successfully preview content: create new content", func(t *testing.T) {
		pageId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]

		createdContent := mockContent
		createdContent.ID = uuid.New()

		landingRepo := &MockCMSLandingPageRepo{
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			findLandingContentPreviewById: func(pageId uuid.UUID, language string) (*models.LandingContent, error) {
				return nil, gorm.ErrRecordNotFound
			},
			createLandingContentPreview: func(landingContentPreview *models.LandingContent) (*models.LandingContent, error) {
				return createdContent, nil
			},
		}

		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()		
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		expectPreviewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(mockContent.Language),
			"landing",
			createdContent.ID,
		)

		assert.NoError(t, err)

		previewUrl, err := service.PreviewLandingContent(pageId, mockContent)
		assert.NoError(t, err)
		assert.Equal(t, expectPreviewUrl, previewUrl)
	})	

	t.Run("successfully preview content: update existing content", func(t *testing.T) {
		pageId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]

		contentId := uuid.New()

		existingContent := mockContent
		existingContent.ID = contentId

		updatedContent := existingContent
		updatedContent.Title = "Updated Title"
		updatedContent.UrlAlias = "Updated URL Alias"

		landingRepo := &MockCMSLandingPageRepo{
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			findLandingContentPreviewById: func(pageId uuid.UUID, language string) (*models.LandingContent, error) {
				return existingContent, nil
			},
			updateLandingContentPreview: func(landingContentPreview *models.LandingContent) (*models.LandingContent, error) {
				return updatedContent, nil
			},
		}

		emailContentRepo := &MockCMSEmailContentRepo{}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}
		cfg := config.New()		
		emailSendingService := services.NewEmailSendingService(cfg, emailCategoryRepo, emailContentRepo)		

		service := services.NewCMSLandingPageService(landingRepo, emailSendingService, emailContentRepo, emailCategoryRepo, cfg)	

		expectPreviewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(updatedContent.Language),
			"landing",
			updatedContent.ID,
		)

		assert.NoError(t, err)

		previewUrl, err := service.PreviewLandingContent(pageId, mockContent)
		assert.NoError(t, err)
		assert.Equal(t, expectPreviewUrl, previewUrl)
	})		
}