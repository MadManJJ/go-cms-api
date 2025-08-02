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

type MockCMSFaqPageRepo struct {
	createFaqPage                         func(faqPage *models.FaqPage) (*models.FaqPage, error)
	findAllFaqPage                        func(query dto.FaqPageQuery, sort string, page, limit int, language string) ([]models.FaqPage, int64, error)
	findFaqPageById                       func(id uuid.UUID) (*models.FaqPage, error)
	updateFaqContent                      func(updateFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error)
	deleteFaqPage                         func(id uuid.UUID) error
	findContentByFaqPageId                func(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error)
	findLatestContentByPageId             func(pageId uuid.UUID, language string) (*models.FaqContent, error)
	createContentForFaqPage               func(faqContent *models.FaqContent, lang string, mode string) (*models.FaqContent, error) // Deprecated
	deleteFaqContent                      func(pageId uuid.UUID, lang, mode string) error
	duplicateFaqPage                      func(pageId uuid.UUID) (*models.FaqPage, error)
	duplicateFaqContentToAnotherLanguage  func(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error)
	revertFaqContent                      func(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error)
	getCategory                           func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	getRevisionByFaqPageId                func(pageId uuid.UUID, language string) ([]models.Revision, error)
	isUrlDuplicate                        func(url string, pageId uuid.UUID) (bool, error)
	isUrlAliasDuplicate                   func(urlAlias string, pageId uuid.UUID) (bool, error)
	getPageIdByContentId                  func(contentId uuid.UUID) (uuid.UUID, error)
	createFaqContentPreview               func(faqContentPreview *models.FaqContent) (*models.FaqContent, error)
	updateFaqContentPreview               func(faqContentPreview *models.FaqContent) (*models.FaqContent, error)
	findFaqContentPreviewById             func(pageId uuid.UUID, language string) (*models.FaqContent, error)

}

func (m *MockCMSFaqPageRepo) CreateFaqPage(faqPage *models.FaqPage) (*models.FaqPage, error) {
	return m.createFaqPage(faqPage)
}

func (m *MockCMSFaqPageRepo) FindAllFaqPage(query dto.FaqPageQuery, sort string, page, limit int, language string) ([]models.FaqPage, int64, error) {
	return m.findAllFaqPage(query, sort, page, limit, language)
}

func (m *MockCMSFaqPageRepo) FindFaqPageById(id uuid.UUID) (*models.FaqPage, error) {
	return m.findFaqPageById(id)
}

func (m *MockCMSFaqPageRepo) UpdateFaqContent(updateFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error) {
	return m.updateFaqContent(updateFaqContent, prevContentId)
}

func (m *MockCMSFaqPageRepo) DeleteFaqPage(id uuid.UUID) error {
	return m.deleteFaqPage(id)
}

func (m *MockCMSFaqPageRepo) FindContentByFaqPageId(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error) {
	return m.findContentByFaqPageId(pageId, language, mode)
}

func (m *MockCMSFaqPageRepo) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.FaqContent, error) {
	return m.findLatestContentByPageId(pageId, language)
}

func (m *MockCMSFaqPageRepo) CreateContentForFaqPage(faqContent *models.FaqContent, lang string, mode string) (*models.FaqContent, error) {
	return m.createContentForFaqPage(faqContent, lang, mode)
}

func (m *MockCMSFaqPageRepo) DeleteFaqContent(pageId uuid.UUID, lang, mode string) error {
	return m.deleteFaqContent(pageId, lang, mode)
}

func (m *MockCMSFaqPageRepo) DuplicateFaqPage(pageId uuid.UUID) (*models.FaqPage, error) {
	return m.duplicateFaqPage(pageId)
}

func (m *MockCMSFaqPageRepo) DuplicateFaqContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	return m.duplicateFaqContentToAnotherLanguage(contentId, newRevision)
}

func (m *MockCMSFaqPageRepo) RevertFaqContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	return m.revertFaqContent(revisionId, newRevision)
}

func (m *MockCMSFaqPageRepo) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	return m.getCategory(pageId, categoryTypeCode, language, mode)
}

func (m *MockCMSFaqPageRepo) GetRevisionByFaqPageId(pageId uuid.UUID, language string) ([]models.Revision, error) {
	return m.getRevisionByFaqPageId(pageId, language)
}

func (m *MockCMSFaqPageRepo) IsUrlDuplicate(url string, pageId uuid.UUID) (bool, error) {
	return m.isUrlDuplicate(url, pageId)
}

func (m *MockCMSFaqPageRepo) IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error) {
	return m.isUrlAliasDuplicate(urlAlias, pageId)
}

func (m *MockCMSFaqPageRepo) GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error) {
	return m.getPageIdByContentId(contentId)
}

func (m *MockCMSFaqPageRepo) CreateFaqContentPreview(faqContentPreview *models.FaqContent) (*models.FaqContent, error) {
	return m.createFaqContentPreview(faqContentPreview)
}

func (m *MockCMSFaqPageRepo) UpdateFaqContentPreview(faqContentPreview *models.FaqContent) (*models.FaqContent, error) {
	return m.updateFaqContentPreview(faqContentPreview)
}

func (m *MockCMSFaqPageRepo) FindFaqContentPreviewById(pageId uuid.UUID, language string) (*models.FaqContent, error) {
	return m.findFaqContentPreviewById(pageId, language)
}

func TestCMSService_CreateFaqPage(t *testing.T) {
	cfg := config.New()

	t.Run("successfully create faq page", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()

		pageId := uuid.New()

		createdFaqPage := mockFaqPage
		createdFaqPage.ID = pageId

		repo := &MockCMSFaqPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createFaqPage: func(faqPage *models.FaqPage) (*models.FaqPage, error) {
				return createdFaqPage, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.CreateFaqPage(mockFaqPage)
		assert.NoError(t, err)
		assert.Equal(t, createdFaqPage, actualFaqPage)
	})

	t.Run("failed to create faq page: internal server error", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		repo := &MockCMSFaqPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},			
			createFaqPage: func(faqPage *models.FaqPage) (*models.FaqPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.CreateFaqPage(mockFaqPage)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})	

	t.Run("failed to create faq page: url is duplicated", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		repo := &MockCMSFaqPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.CreateFaqPage(mockFaqPage)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})	

	t.Run("failed to create faq page: url alias is duplicated", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		repo := &MockCMSFaqPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},			
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.CreateFaqPage(mockFaqPage)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})		
}

func TestCMSService_FindFaqPages(t *testing.T) {
	cfg := config.New()

	// Arguments
	rawQuery := `{
		"title": "How to use",
		"category_faq": "general",
		"category_keywords": "usage,getting started",
		"status": "published"
	}`
	sort := "sort test"
	page := 1
	limit := 10
	language := "en"	

	t.Run("successfully find faq pages", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqPages := []models.FaqPage{*mockFaqPage, *mockFaqPage}
		repo := &MockCMSFaqPageRepo{
			findAllFaqPage: func(query dto.FaqPageQuery, sort string, page, limit int, language string) ([]models.FaqPage, int64, error) {
				return mockFaqPages, 2, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPages, totalCount, err := service.FindFaqPages(rawQuery, sort, page, limit, language)
		assert.NoError(t, err)
		assert.Equal(t, mockFaqPages, actualFaqPages)
		assert.Equal(t, 2, len(actualFaqPages))
		assert.Equal(t, int64(2), totalCount)
	})	

	t.Run("failed to find find faq pages", func(t *testing.T) {
		repo := &MockCMSFaqPageRepo{
			findAllFaqPage: func(query dto.FaqPageQuery, sort string, page, limit int, language string) ([]models.FaqPage, int64, error) {
				return nil, 0, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPages, totalCount, err := service.FindFaqPages(rawQuery, sort, page, limit, language)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPages)
		assert.Equal(t, int64(0), totalCount)
	})	
}

func TestCMSService_FindFaqPageById(t *testing.T) {
	cfg := config.New()

	t.Run("successfully find faq page", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		repo := &MockCMSFaqPageRepo{
			findFaqPageById: func(id uuid.UUID) (*models.FaqPage, error) {
				return mockFaqPage, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		pageId := uuid.New()

		actualFaqPage, err := service.FindFaqPageById(pageId)
		assert.NoError(t, err)
		assert.Equal(t, mockFaqPage, actualFaqPage)
	})	
	
	t.Run("successfully find faq page", func(t *testing.T) {
		repo := &MockCMSFaqPageRepo{
			findFaqPageById: func(id uuid.UUID) (*models.FaqPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		pageId := uuid.New()

		actualFaqPage, err := service.FindFaqPageById(pageId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})		
}

func TestCMSService_UpdateFaqContent(t *testing.T) {
	cfg := config.New()

	t.Run("successfully update faq content", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		pageId := uuid.New()
		contentId := uuid.New()

		updatedContent := mockContent
		updatedContent.ID = contentId
		updatedContent.PageID = pageId

		repo := &MockCMSFaqPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			updateFaqContent: func(updateFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error) {
				return updatedContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.UpdateFaqContent(mockContent, contentId)
		assert.NoError(t, err)
		assert.Equal(t, updatedContent, actualFaqPage)
	})	

	t.Run("failed to update faq content: can't find page id", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		contentId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return uuid.Nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.UpdateFaqContent(mockContent, contentId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})	
	
	t.Run("failed to update faq content: url is duplicated", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		pageId := uuid.New()
		contentId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.UpdateFaqContent(mockContent, contentId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})	

	t.Run("failed to update faq content: url alias is duplicated", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		pageId := uuid.New()
		contentId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return true, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.UpdateFaqContent(mockContent, contentId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})		

	t.Run("failed to update faq content: error at UpdateFaqContent", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		pageId := uuid.New()
		contentId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			getPageIdByContentId: func(contentId uuid.UUID) (uuid.UUID, error) {
				return pageId, nil
			},
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			updateFaqContent: func(updateFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		actualFaqPage, err := service.UpdateFaqContent(mockContent, contentId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqPage)
	})	
}

func TestCMSService_DeleteFaqPage(t *testing.T) {
	cfg := config.New()

	t.Run("successfully delete faq page", func(t *testing.T) {
		pageId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			deleteFaqPage: func(id uuid.UUID) error {
				return nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		err := service.DeleteFaqPage(pageId)
		assert.NoError(t, err)
	})
	

	t.Run("successfully delete faq page", func(t *testing.T) {
		pageId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			deleteFaqPage: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		err := service.DeleteFaqPage(pageId)
		assert.Error(t, err)
	})		
}

func TestCMSService_FindContentByFaqPageId(t *testing.T) {
	cfg := config.New()

	t.Run("successfully find faq content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		repo := &MockCMSFaqPageRepo{
			findContentByFaqPageId: func(pageId uuid.UUID, language, mode string) (*models.FaqContent, error) {
				return mockContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.FindContentByFaqPageId(pageId, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, mockContent, faqContent)
	})

	t.Run("failed to find faq content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		repo := &MockCMSFaqPageRepo{
			findContentByFaqPageId: func(pageId uuid.UUID, language, mode string) (*models.FaqContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.FindContentByFaqPageId(pageId, language, mode)
		assert.Error(t, err)
		assert.Nil(t, faqContent)
	})	

}

func TestCMSService_FindLatestFaqContentByPageId(t *testing.T) {
	cfg := config.New()

	t.Run("successfully find faq latest content by page id", func(t *testing.T) {
		pageId := uuid.New()
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		language := string(enums.PageLanguageEN)

		repo := &MockCMSFaqPageRepo{
			findLatestContentByPageId: func(pageId uuid.UUID, language string) (*models.FaqContent, error) {
				return mockContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.FindLatestContentByPageId(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, mockContent, faqContent)
	})	

	t.Run("failed to find faq latest content by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		repo := &MockCMSFaqPageRepo{
			findLatestContentByPageId: func(pageId uuid.UUID, language string) (*models.FaqContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.FindLatestContentByPageId(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, faqContent)
	})		

}

func TestCMSService_DeleteContentByFaqPageId(t *testing.T) {
	cfg := config.New()

	t.Run("successfully delete content by faq page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		repo := &MockCMSFaqPageRepo{
			deleteFaqContent: func(pageId uuid.UUID, lang, mode string) error {
				return nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		err := service.DeleteContentByFaqPageId(pageId, language, mode)
		assert.NoError(t, err)
	})	

	t.Run("failed to delete content by faq page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		repo := &MockCMSFaqPageRepo{
			deleteFaqContent: func(pageId uuid.UUID, lang, mode string) error {
				return errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		err := service.DeleteContentByFaqPageId(pageId, language, mode)
		assert.Error(t, err)
	})		

}

func TestCMSService_DuplicateFaqPage(t *testing.T) {
	cfg := config.New()

	t.Run("successfully duplicate faq page", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()

		repo := &MockCMSFaqPageRepo{
			duplicateFaqPage: func(pageId uuid.UUID) (*models.FaqPage, error) {
				return mockFaqPage, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqPage, err := service.DuplicateFaqPage(pageId)
		assert.NoError(t, err)
		assert.Equal(t, mockFaqPage, faqPage)
	})	

	t.Run("failed to duplicate faq page", func(t *testing.T) {
		pageId := uuid.New()

		repo := &MockCMSFaqPageRepo{
			duplicateFaqPage: func(pageId uuid.UUID) (*models.FaqPage, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqPage, err := service.DuplicateFaqPage(pageId)
		assert.Error(t, err)
		assert.Nil(t, faqPage)
	})		

}

func TestCMSService_DuplicateFaqContentToAnotherLanguage(t *testing.T) {
	cfg := config.New()

	t.Run("successfully duplicate faq content to another language", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		repo := &MockCMSFaqPageRepo{
			duplicateFaqContentToAnotherLanguage: func(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
				return mockContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.DuplicateFaqContentToAnotherLanguage(pageId, mockRevision)
		assert.NoError(t, err)
		assert.Equal(t, mockContent, faqContent)
	})	

	t.Run("failed to duplicate faq content to another language", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		repo := &MockCMSFaqPageRepo{
			duplicateFaqContentToAnotherLanguage: func(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.DuplicateFaqContentToAnotherLanguage(pageId, mockRevision)
		assert.Error(t, err)
		assert.Nil(t, faqContent)
	})		

}

func TestCMSService_RevertFaqContent(t *testing.T) {
	cfg := config.New()

	t.Run("successfully revert faq content", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		repo := &MockCMSFaqPageRepo{
			revertFaqContent: func(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
				return mockContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.RevertFaqContent(pageId, mockRevision)
		assert.NoError(t, err)
		assert.Equal(t, mockContent, faqContent)
	})	

	t.Run("failed to revert faq content", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		repo := &MockCMSFaqPageRepo{
			revertFaqContent: func(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		faqContent, err := service.RevertFaqContent(pageId, mockRevision)
		assert.Error(t, err)
		assert.Nil(t, faqContent)
	})		

}

func TestCMSService_FindFaqCategories(t *testing.T) {
	cfg := config.New()

	t.Run("successfully get categories", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockCategory := mockContent.Categories[0]
		mockCategoryType := mockCategory.CategoryType

		typeCode := mockCategoryType.TypeCode
		language := string(mockContent.Language)
		mode := string(mockContent.Mode)

		mockCategories := []models.Category{*mockCategory, *mockCategory}
		repo := &MockCMSFaqPageRepo{
			getCategory: func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
				return mockCategories, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		categories, err := service.FindCategories(pageId, typeCode, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, mockCategories, categories)
	})	

	t.Run("failed to get categories", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockCategory := mockContent.Categories[0]
		mockCategoryType := mockCategory.CategoryType

		typeCode := mockCategoryType.TypeCode
		language := string(mockContent.Language)
		mode := string(mockContent.Mode)

		repo := &MockCMSFaqPageRepo{
			getCategory: func(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		categories, err := service.FindCategories(pageId, typeCode, language, mode)
		assert.Error(t, err)
		assert.Nil(t, categories)
	})		

}

func TestCMSService_FindFaqRevisions(t *testing.T) {
	cfg := config.New()

	t.Run("successfully get revisions", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		mockRevision := mockContent.Revision

		language := string(mockContent.Language)		

		mockRevisions := []models.Revision{*mockRevision, *mockRevision}
		repo := &MockCMSFaqPageRepo{
			getRevisionByFaqPageId: func(pageId uuid.UUID, language string) ([]models.Revision, error) {
				return mockRevisions, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		revisions, err := service.FindRevisions(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, mockRevisions, revisions)
	})	

	t.Run("failed to get revisions", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		language := string(mockContent.Language)		

		repo := &MockCMSFaqPageRepo{
			getRevisionByFaqPageId: func(pageId uuid.UUID, language string) ([]models.Revision, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)

		revisions, err := service.FindRevisions(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, revisions)
	})		

}

func TestCMSService_PreviewFaqContent(t *testing.T) {
	cfg := config.New()

	urls := cfg.App.FrontendURLS
	parts := strings.Split(urls, ",")
	appUrl := parts[1]	

	t.Run("successfully preview content: create new content", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		createdContent := mockContent
		createdContent.ID = uuid.New()

		repo := &MockCMSFaqPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			findFaqContentPreviewById: func(pageId uuid.UUID, language string) (*models.FaqContent, error) {
				return nil, gorm.ErrRecordNotFound
			},
			createFaqContentPreview: func(faqContentPreview *models.FaqContent) (*models.FaqContent, error) {
				return createdContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)	

		expectPreviewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(mockContent.Language),
			"faq",
			createdContent.ID,
		)

		assert.NoError(t, err)

		previewUrl, err := service.PreviewFaqContent(pageId, mockContent)
		assert.NoError(t, err)
		assert.Equal(t, expectPreviewUrl, previewUrl)
	})	

	t.Run("successfully preview content: update existing content", func(t *testing.T) {
		pageId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]

		contentId := uuid.New()

		existingContent := mockContent
		existingContent.ID = contentId

		updatedContent := existingContent
		updatedContent.Title = "Updated Title"
		updatedContent.URL = "Updated URL"
		updatedContent.URLAlias = "Updated URL Alias"

		repo := &MockCMSFaqPageRepo{
			isUrlDuplicate: func(url string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			isUrlAliasDuplicate: func(urlAlias string, pageId uuid.UUID) (bool, error) {
				return false, nil
			},
			findFaqContentPreviewById: func(pageId uuid.UUID, language string) (*models.FaqContent, error) {
				return existingContent, nil
			},
			updateFaqContentPreview: func(faqContentPreview *models.FaqContent) (*models.FaqContent, error) {
				return updatedContent, nil
			},
		}

		service := services.NewCMSFaqPageService(repo, cfg)	

		expectPreviewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(updatedContent.Language),
			"faq",
			updatedContent.ID,
		)

		assert.NoError(t, err)

		previewUrl, err := service.PreviewFaqContent(pageId, mockContent)
		assert.NoError(t, err)
		assert.Equal(t, expectPreviewUrl, previewUrl)
	})		
}