package services

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSFaqPageServiceInterface interface {
	CreateFaqPage(faqPage *models.FaqPage) (*models.FaqPage, error)
	FindFaqPages(rawQuery string, sort string, page, limit int, language string) ([]models.FaqPage, int64, error)
	FindFaqPageById(id uuid.UUID) (*models.FaqPage, error)
	UpdateFaqContent(updatedFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error)
	DeleteFaqPage(id uuid.UUID) error
	FindContentByFaqPageId(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error)
	FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.FaqContent, error)
	DeleteContentByFaqPageId(pageId uuid.UUID, language, mode string) error
	DuplicateFaqPage(pageId uuid.UUID) (*models.FaqPage, error)
	DuplicateFaqContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error)
	RevertFaqContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error)
	FindCategories(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error)
	PreviewFaqContent(pageId uuid.UUID, faqContentPreview *models.FaqContent) (string, error)
}

type CMSFaqPageService struct {
	repo repositories.CMSFaqPageRepositoryInterface
	cfg *config.Config
}

func NewCMSFaqPageService(repo repositories.CMSFaqPageRepositoryInterface, cfg *config.Config) *CMSFaqPageService {
	return &CMSFaqPageService{
		repo: repo,
		cfg: cfg,
	}
}

// Always send only 1 content
func (s *CMSFaqPageService) CreateFaqPage(faqPage *models.FaqPage) (*models.FaqPage, error) {
	faqContents := faqPage.Contents
	if len(faqContents) > 1 {
		return nil, errs.ErrTooMuchContent
	} else if len(faqContents) == 0 {
		return nil, errs.ErrMissingContent
	}
	// Only one content
	faqContent := faqContents[0]

	// Check if the URL is duplicate or not
	isUrlDuplicate, err := s.repo.IsUrlDuplicate(faqContent.URL, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if isUrlDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	// Check if the URL Alias is duplicate or not
	if faqContent.URLAlias != "" {
		isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(faqContent.URLAlias, uuid.Nil)
		if err != nil {
			return nil, err
		}
		if isUrlAliasDuplicate {
			return nil, errs.ErrDuplicateUrlAlias
		}
	}

	// Always need to have a revision
	if faqContent.Revision == nil {
		return nil, errs.ErrNoRevisionFound
	}

	// Normalize main contents
	if err := helpers.NormalizeFaqContent(faqContent); err != nil {
		return nil, err
	}

	return s.repo.CreateFaqPage(faqPage)
}

func (s *CMSFaqPageService) FindFaqPages(rawQuery string, sort string, page, limit int, language string) ([]models.FaqPage, int64, error) {
	var query dto.FaqPageQuery
	if rawQuery != "" {
		if err := json.Unmarshal([]byte(rawQuery), &query); err != nil {
			return nil, 0, errs.ErrInvalidQuery
		}
	}

	return s.repo.FindAllFaqPage(query, sort, page, limit, language)
}

func (s *CMSFaqPageService) FindFaqPageById(id uuid.UUID) (*models.FaqPage, error) {
	return s.repo.FindFaqPageById(id)
}

func (s *CMSFaqPageService) UpdateFaqContent(updatedFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error) {
	// Check if the URL is duplicate or not
	faqPageId, err := s.repo.GetPageIdByContentId(prevContentId)
	if err != nil {
		return nil, err
	}
	isDuplicate, err := s.repo.IsUrlDuplicate(updatedFaqContent.URL, faqPageId)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	// Check if the URL Alias is duplicate or not
	if updatedFaqContent.URLAlias != "" {
		isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(updatedFaqContent.URLAlias, faqPageId)
		if err != nil {
			return nil, err
		}
		if isUrlAliasDuplicate {
			return nil, errs.ErrDuplicateUrlAlias
		}
	}

	err = helpers.NormalizeFaqContent(updatedFaqContent)
	if err != nil {
		return nil, err
	}

	// Always need to have a revision
	if updatedFaqContent.Revision == nil {
		return nil, errs.ErrNoRevisionFound
	}
	return s.repo.UpdateFaqContent(updatedFaqContent, prevContentId)
}

func (s *CMSFaqPageService) DeleteFaqPage(id uuid.UUID) error {
	return s.repo.DeleteFaqPage(id)
}

func (s *CMSFaqPageService) FindContentByFaqPageId(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return nil, err
	}

	return s.repo.FindContentByFaqPageId(pageId, language, mode)
}

func (r *CMSFaqPageService) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.FaqContent, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	return r.repo.FindLatestContentByPageId(pageId, language)
}

func (s *CMSFaqPageService) DeleteContentByFaqPageId(pageId uuid.UUID, language, mode string) error {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return err
	}

	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return err
	}

	return s.repo.DeleteFaqContent(pageId, language, mode)
}

func (s *CMSFaqPageService) DuplicateFaqPage(pageId uuid.UUID) (*models.FaqPage, error) {
	faqPage, err := s.repo.DuplicateFaqPage(pageId)

	if err != nil {
		return nil, err
	}

	return faqPage, nil
}

func (s *CMSFaqPageService) DuplicateFaqContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	err := helpers.NormalizeRevision(newRevision)

	if err != nil {
		return nil, err
	}	

	return s.repo.DuplicateFaqContentToAnotherLanguage(contentId, newRevision)
}

func (s *CMSFaqPageService) RevertFaqContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	err := helpers.NormalizeRevision(newRevision)

	if err != nil {
		return nil, err
	}	
		
	return s.repo.RevertFaqContent(revisionId, newRevision)
}

func (s *CMSFaqPageService) FindCategories(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}
	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategory(pageId, categoryTypeCode, language, mode)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CMSFaqPageService) FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	revisions, err := s.repo.GetRevisionByFaqPageId(pageId, language)
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

func (s *CMSFaqPageService) PreviewFaqContent(pageId uuid.UUID, faqContentPreview *models.FaqContent) (string, error) {
	urls := s.cfg.App.FrontendURLS // "http://localhost:8000,http://localhost:3000"
	parts := strings.Split(urls, ",")
	appUrl := parts[1] // "http://localhost:3000"

	// Default value for preview content
	faqContentPreview.Mode = enums.PageModePreview
	faqContentPreview.PublishStatus = enums.PublishStatusNotPublished
	faqContentPreview.WorkflowStatus = enums.WorkflowUnPublished	
	faqContentPreview.ExpiredAt = time.Now().Add(2 * time.Hour)	

	if err := helpers.NormalizeFaqContent(faqContentPreview); err != nil {
		return "", err
	}
	
	// Check if the URL is duplicate or not
	isUrlDuplicate, err := s.repo.IsUrlDuplicate(faqContentPreview.URL, pageId)
	if err != nil {
		return "", err
	}
	if isUrlDuplicate {
		return "", errs.ErrDuplicateURL
	}

	// Check if the URL Alias is duplicate or not
	if faqContentPreview.URLAlias != "" {
		isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(faqContentPreview.URLAlias, pageId)
		if err != nil {
			return "", err
		}
		if isUrlAliasDuplicate {
			return "", errs.ErrDuplicateUrlAlias
		}
	}	

	existingPreviewContent, err := s.repo.FindFaqContentPreviewById(pageId, string(faqContentPreview.Language))

	// Internal error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	// Record not found, so we create a new one
	if err == gorm.ErrRecordNotFound {
		faqContentPreview.PageID = pageId
		createdFaqContentPreview, err := s.repo.CreateFaqContentPreview(faqContentPreview)
		if err != nil {
			return "", err
		}

		previewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(createdFaqContentPreview.Language),
			"faq",
			createdFaqContentPreview.ID,
		)	

		if err != nil {
			return "", err
		}		

		return previewUrl, nil
	}

	// Attach the old id, so we can save the new one in its place
	faqContentPreview.ID = existingPreviewContent.ID
	faqContentPreview.PageID = existingPreviewContent.PageID
	updatedFaqContentPreview, err := s.repo.UpdateFaqContentPreview(faqContentPreview)
	if err != nil {
		return "", err
	}

	previewUrl, err := helpers.BuildPreviewURL(
		appUrl,
		string(updatedFaqContentPreview.Language),
		"faq",
		updatedFaqContentPreview.ID,
	)	

	if err != nil {
		return "", err
	}		

	return previewUrl, nil
}
