package services

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"

	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSPartnerPageServiceInterface interface {
	CreatePartnerPage(PartnerPage *models.PartnerPage) (*models.PartnerPage, error)
	FindPartnerPages(rawQuery string, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error)
	FindPartnerPageById(id uuid.UUID) (*models.PartnerPage, error)
	UpdatePartnerContent(updatedPartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error)
	DeletePartnerPage(id uuid.UUID) error
	FindContentByPartnerPageId(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error)
	FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.PartnerContent, error)
	DeleteContentByPartnerPageId(pageId uuid.UUID, language, mode string) error
	DuplicatePartnerPage(pageId uuid.UUID) (*models.PartnerPage, error)
	DuplicatePartnerContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error)
	RevertPartnerContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error)
	GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error)
	PreviewPartnerContent(pageId uuid.UUID, partnerContentPreview *models.PartnerContent) (string, error)
}

type CMSPartnerPageService struct {
	repo                repositories.CMSPartnerPageRepositoryInterface
	emailSendingService EmailSendingServiceInterface
	emailContentRepo    repositories.EmailContentRepositoryInterface
	emailCategoryRepo   repositories.EmailCategoryRepositoryInterface
	cfg                 *config.Config
}

func NewCMSPartnerPageService(
	repo repositories.CMSPartnerPageRepositoryInterface,
	emailSendingService EmailSendingServiceInterface,
	emailContentRepo repositories.EmailContentRepositoryInterface,
	emailCategoryRepo repositories.EmailCategoryRepositoryInterface,
	cfg *config.Config,
) *CMSPartnerPageService {
	return &CMSPartnerPageService{
		repo:                repo,
		emailSendingService: emailSendingService,
		emailContentRepo:    emailContentRepo,
		emailCategoryRepo:   emailCategoryRepo,
		cfg:                 cfg,
	}
}

// Always send only 1 content
func (s *CMSPartnerPageService) CreatePartnerPage(PartnerPage *models.PartnerPage) (*models.PartnerPage, error) {
	PartnerContents := PartnerPage.Contents
	if len(PartnerContents) > 1 {
		return nil, errs.ErrTooMuchContent
	} else if len(PartnerContents) == 0 {
		return nil, errs.ErrMissingContent
	}
	// Only one content
	PartnerContent := PartnerContents[0]

	// Check if the URL is duplicate or not
	isDuplicate, err := s.repo.IsUrlDuplicate(PartnerContent.URL, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	// Check if the URL Alias is duplicate or not
	if PartnerContent.URLAlias != "" {
		isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(PartnerContent.URLAlias, uuid.Nil)
		if err != nil {
			return nil, err
		}
		if isUrlAliasDuplicate {
			return nil, errs.ErrDuplicateUrlAlias
		}
	}

	// Always need to have a revision
	if PartnerContent.Revision == nil {
		return nil, errs.ErrNoRevisionFound
	}

	// Normalize main content
	if err := helpers.NormalizePartnerContent(PartnerContent); err != nil {
		return nil, err
	}

	// Normalize revision
	if err := helpers.NormalizeRevision(PartnerContent.Revision); err != nil {
		return nil, err
	}

	// Normalize categories if present
	// for _, category := range PartnerContent.Categories {
	// 	if err := helpers.NormalizeCategory(category); err != nil {
	// 		return nil, err
	// 	}
	// }

	return s.repo.CreatePartnerPage(PartnerPage)
}

func (s *CMSPartnerPageService) FindPartnerPages(rawQuery string, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error) {
	var query dto.PartnerPageQuery
	if rawQuery != "" {
		if err := json.Unmarshal([]byte(rawQuery), &query); err != nil {
			return nil, 0, errs.ErrInvalidQuery
		}
	}

	return s.repo.FindAllPartnerPage(query, sort, page, limit, language)
}

func (s *CMSPartnerPageService) FindPartnerPageById(id uuid.UUID) (*models.PartnerPage, error) {
	return s.repo.FindPartnerPageById(id)
}

// We might not need this service
func (s *CMSPartnerPageService) UpdatePartnerPage(updatedPartnerPage *models.PartnerPage) (*models.PartnerPage, error) {
	return nil, errs.ErrInternalServerError
}

func (s *CMSPartnerPageService) UpdatePartnerContent(updatedPartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error) {

	partnerContentId, err := s.repo.GetPageIdByContentId(prevContentId)
	if err != nil {
		return nil, err
	}

	isURLDuplicate, err := s.repo.IsUrlDuplicate(updatedPartnerContent.URL, partnerContentId)
	if err != nil {
		return nil, err
	}
	if isURLDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(updatedPartnerContent.URLAlias, partnerContentId)
	if err != nil {
		return nil, err
	}
	if isUrlAliasDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	// Always need to have a revision
	if updatedPartnerContent.Revision == nil {
		return nil, errs.ErrNoRevisionFound
	}

	log.Printf("[SERVICE-IN] Language from Request: '%s'", updatedPartnerContent.Language)

	savedContent, err := s.repo.UpdatePartnerContent(updatedPartnerContent, prevContentId)
	if err != nil {
		return nil, err
	}

	if savedContent.WorkflowStatus == enums.WorkflowWaitingDesign && len(savedContent.ApprovalEmail) > 0 {
		go s.triggerApprovalNotifications(savedContent)

	}

	log.Printf("[SERVICE-OUT] Language from Repo: '%s'", savedContent.Language)

	return savedContent, nil
}

// triggerApprovalNotifications sends approval notifications to the relevant recipients.
func (s *CMSPartnerPageService) triggerApprovalNotifications(content *models.PartnerContent) {
	// --- KEY CHANGE IS HERE ---
	const approvalCategoryTitle = "Approve"
	category, err := s.emailCategoryRepo.FindByTitle(approvalCategoryTitle)
	if err != nil {
		log.Printf("CRITICAL: Failed to query for Email Category '%s': %v. Notifications will not be sent.", approvalCategoryTitle, err)
		return
	}
	if category == nil {
		log.Printf("CRITICAL: Email Category with title '%s' not found. Please create it in the CMS. Notifications will not be sent.", approvalCategoryTitle)
		return
	}
	// --- END KEY CHANGE ---

	// 2. ค้นหา EmailContent ทั้งหมดที่อยู่ใน Category "Approve" และภาษาที่ตรงกัน
	emailCategoryIDStr := category.ID.String()
	filter := dto.EmailContentFilter{
		EmailCategoryID: &emailCategoryIDStr,
		Language:        &content.Language,
	}
	templates, err := s.emailContentRepo.ListByFilters(filter)
	if err != nil {
		log.Printf("Error fetching email templates for Category '%s' (ID: %s): %v", approvalCategoryTitle, emailCategoryIDStr, err)
		return
	}

	if len(templates) == 0 {
		log.Printf("No email templates found for Category '%s' and language '%s'. No notifications sent.", approvalCategoryTitle, content.Language)
		return
	}

	previewURL, cmsEditURL := s.buildNotificationURLs("partner", content)
	authorName, authorEmail := s.getAuthorInfo(content)

	emailData := map[string]interface{}{
		"urlPreview": previewURL,
		"urlCms":     cmsEditURL,
		"pageTitle":  content.Title,
		"author":     authorName,
	}

	for _, template := range templates {
		var recipients []string

		if template.Label == "email_to_admin" {
			if len(content.ApprovalEmail) > 0 {
				recipients = content.ApprovalEmail
			} else {
				log.Printf("Skipping 'email_to_admin' because ApprovalEmail list is empty in the content.")
				continue
			}

		} else if template.Label == "email_to_user" {
			if authorEmail != "" {
				recipients = []string{authorEmail}
			} else {
				log.Printf("Skipping 'email_to_user' because author email is not available.")
				continue
			}

		} else {
			if template.SendTo != "" {
				recipients = strings.Split(template.SendTo, ",")
			} else {
				log.Printf("Skipping template '%s' because its purpose is general and SendTo field is not set.", template.Label)
				continue
			}
		}

		emailReq := dto.SendEmailRequest{
			EmailCategoryTitleOrID: approvalCategoryTitle,
			EmailContentLabel:      template.Label,
			Language:               template.Language,
			ToRecipientEmails:      recipients,
			Data:                   emailData,
		}

		go func(req dto.SendEmailRequest) {
			if err := s.emailSendingService.SendEmail(req); err != nil {
				log.Printf("Error sending email using template '%s': %v", req.EmailContentLabel, err)
			} else {
				log.Printf("Successfully queued email using template '%s' to: %v", req.EmailContentLabel, req.ToRecipientEmails)
			}
		}(emailReq)
	}
}

func (s *CMSPartnerPageService) buildNotificationURLs(pageType string, content *models.PartnerContent) (previewUrl, cmsEditUrl string) {
	// --- Build Preview URL ---
	if s.cfg.App.WebBaseURL != "" {
		baseURL := strings.TrimSuffix(s.cfg.App.WebBaseURL, "/")

		u, err := url.Parse(fmt.Sprintf("%s/preview/%s/%s", baseURL, content.Language, pageType))
		if err != nil {
			log.Printf("Error parsing WebBaseURL for preview: %v", err)
			previewUrl = "Error: Invalid Base URL"
		} else {

			q := u.Query()
			q.Set("id", content.ID.String())
			u.RawQuery = q.Encode()
			previewUrl = u.String()
		}
	} else {
		log.Println("Warning: App.WebBaseURL is not configured. Preview URL will be a relative path.")

		previewUrl = fmt.Sprintf("/preview/%s/%s?id=%s", content.Language, pageType, content.ID)
	}

	// --- Build CMS Edit URL ---
	if s.cfg.App.CMSBaseURL != "" {
		baseURL := strings.TrimSuffix(s.cfg.App.CMSBaseURL, "/")
		cmsEditUrl = fmt.Sprintf("%s/%s-pages/%s/content/%s/edit?lang=%s", baseURL, pageType, content.PageID, content.ID, content.Language)
	} else {
		log.Println("Warning: App.CMSBaseURL is not configured. CMS URL will be a relative path.")
		cmsEditUrl = fmt.Sprintf("/%s-pages/%s/content/%s/edit?lang=%s", pageType, content.PageID, content.ID, content.Language)
	}

	return
}

// In: api-develop/services/cms_partner-page_service.go

// getAuthorInfo extracts the author's name and email from the content's revision, with fallbacks.
func (s *CMSPartnerPageService) getAuthorInfo(content *models.PartnerContent) (name, email string) {
	// --- Set Default Values ---
	name = "CMS User"
	email = "linkornn2003@gmail.com" // ใช้ค่า Default ที่เป็นอีเมลของผู้ดูแลระบบ

	// --- Attempt to get info from Revision ---
	if content.Revision != nil && content.Revision.Author != "" {

		// ในอนาคต อาจจะเก็บ AuthorID ที่เป็น UUID แล้วไป join กับตาราง users แทน
		authorInfo := content.Revision.Author

		// ลองแยกชื่อและอีเมลแบบง่ายๆ (Optional)
		if strings.Contains(authorInfo, "<") && strings.Contains(authorInfo, ">") {
			parts := strings.Split(authorInfo, "<")
			name = strings.TrimSpace(parts[0])
			email = strings.Trim(parts[1], ">")
		} else {
			// ถ้าไม่มีรูปแบบ <...> ก็ใช้ค่าที่ได้มาเป็นทั้งชื่อและอีเมล
			name = authorInfo
			email = authorInfo
		}
	} else {
		log.Printf("Warning: Revision or Author information missing for PartnerContent ID %s. Using default author info.", content.ID)
	}

	// --- Future Improvement ---
	// ในอนาคตคุณอาจจะรับ User Context ที่ได้จาก Middleware เข้ามาในฟังก์ชันนี้
	// แล้วใช้เป็นข้อมูล Author ที่ถูกต้องที่สุดแทนค่า Default
	// if email == "" && c.Locals("userEmail") != nil {
	//    email = c.Locals("userEmail").(string)
	// }

	return
}

func (s *CMSPartnerPageService) DeletePartnerPage(id uuid.UUID) error {
	return s.repo.DeletePartnerPage(id)
}

func (s *CMSPartnerPageService) FindContentByPartnerPageId(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return nil, err
	}

	return s.repo.FindContentByPartnerPageId(pageId, language, mode)
}

func (r *CMSPartnerPageService) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	return r.repo.FindLatestContentByPageId(pageId, language)
}

func (s *CMSPartnerPageService) DeleteContentByPartnerPageId(pageId uuid.UUID, language, mode string) error {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return err
	}

	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return err
	}

	return s.repo.DeletePartnerContent(pageId, language, mode)
}

func (s *CMSPartnerPageService) DuplicatePartnerPage(pageId uuid.UUID) (*models.PartnerPage, error) {
	PartnerPage, err := s.repo.DuplicatePartnerPage(pageId)

	if err != nil {
		return nil, err
	}

	return PartnerPage, nil
}

func (s *CMSPartnerPageService) DuplicatePartnerContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	err := helpers.NormalizeRevision(newRevision)

	if err != nil {
		return nil, err
	}

	return s.repo.DuplicatePartnerContentToAnotherLanguage(contentId, newRevision)
}

func (r *CMSPartnerPageService) RevertPartnerContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	return r.repo.RevertPartnerContent(revisionId, newRevision)
}

func (s *CMSPartnerPageService) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
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

func (s *CMSPartnerPageService) FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	revisions, err := s.repo.GetRevisionByPartnerPageId(pageId, language)
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

func (s *CMSPartnerPageService) PreviewPartnerContent(pageId uuid.UUID, partnerContentPreview *models.PartnerContent) (string, error) {
	urls := s.cfg.App.FrontendURLS // "http://localhost:8000,http://localhost:3000"
	parts := strings.Split(urls, ",")
	appUrl := parts[1] // "http://localhost:3000"

	// Default value for preview content
	partnerContentPreview.Mode = enums.PageModePreview
	partnerContentPreview.PublishStatus = enums.PublishStatusNotPublished
	partnerContentPreview.WorkflowStatus = enums.WorkflowUnPublished
	partnerContentPreview.ExpiredAt = time.Now().Add(2 * time.Hour)

	if err := helpers.NormalizePartnerContent(partnerContentPreview); err != nil {
		return "", err
	}

	// Check if the URL is duplicate or not
	isUrlDuplicate, err := s.repo.IsUrlDuplicate(partnerContentPreview.URL, pageId)
	if err != nil {
		return "", err
	}
	if isUrlDuplicate {
		return "", errs.ErrDuplicateURL
	}

	// Check if the URL Alias is duplicate or not
	if partnerContentPreview.URLAlias != "" {
		isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(partnerContentPreview.URLAlias, pageId)
		if err != nil {
			return "", err
		}
		if isUrlAliasDuplicate {
			return "", errs.ErrDuplicateUrlAlias
		}
	}

	existingPreviewContent, err := s.repo.FindPartnerContentPreviewById(pageId, string(partnerContentPreview.Language))

	// Internal error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	// Record not found, so we create a new one
	if err == gorm.ErrRecordNotFound {
		partnerContentPreview.PageID = pageId
		createdPartnerContentPreview, err := s.repo.CreatePartnerContentPreview(partnerContentPreview)
		if err != nil {
			return "", err
		}

		previewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(createdPartnerContentPreview.Language),
			"partner",
			createdPartnerContentPreview.ID,
		)

		if err != nil {
			return "", err
		}

		return previewUrl, nil
	}

	// Attach the old id, so we can save the new one in its place
	partnerContentPreview.ID = existingPreviewContent.ID
	partnerContentPreview.PageID = existingPreviewContent.PageID
	updatedPartnerContentPreview, err := s.repo.UpdatePartnerContentPreview(partnerContentPreview)
	if err != nil {
		return "", err
	}

	previewUrl, err := helpers.BuildPreviewURL(
		appUrl,
		string(updatedPartnerContentPreview.Language),
		"partner",
		updatedPartnerContentPreview.ID,
	)

	if err != nil {
		return "", err
	}

	return previewUrl, nil
}
