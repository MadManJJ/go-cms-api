package services

import (
	"fmt"
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
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSLandingPageServiceInterface interface {
	CreateLandingPage(LandingPage *models.LandingPage) (*models.LandingPage, error)
	FindLandingPages(rawQuery string, sort string, page, limit int, language string) ([]models.LandingPage, int64, error)
	FindLandingPageById(id uuid.UUID) (*models.LandingPage, error)
	UpdateLandingContent(updatedLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error)
	DeleteLandingPage(id uuid.UUID) error
	FindContentByLandingPageId(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error)
	FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.LandingContent, error)
	DeleteContentByLandingPageId(pageId uuid.UUID, language, mode string) error
	DuplicateLandingPage(pageId uuid.UUID) (*models.LandingPage, error)
	DuplicateLandingContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error)
	RevertLandingContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error)
	GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error)
	PreviewLandingContent(pageId uuid.UUID, landingContentPreview *models.LandingContent) (string, error)
}

type CMSLandingPageService struct {
	repo                repositories.CMSLandingPageRepositoryInterface
	emailSendingService EmailSendingServiceInterface
	emailContentRepo    repositories.EmailContentRepositoryInterface
	emailCategoryRepo   repositories.EmailCategoryRepositoryInterface
	cfg                 *config.Config
}

func NewCMSLandingPageService(
	repo repositories.CMSLandingPageRepositoryInterface,
	emailSendingService EmailSendingServiceInterface,
	emailContentRepo repositories.EmailContentRepositoryInterface,
	emailCategoryRepo repositories.EmailCategoryRepositoryInterface,
	cfg *config.Config,
) *CMSLandingPageService {
	return &CMSLandingPageService{
		repo:                repo,
		emailSendingService: emailSendingService,
		emailContentRepo:    emailContentRepo,
		emailCategoryRepo:   emailCategoryRepo,
		cfg:                 cfg,
	}
}

// Always send only 1 content
func (s *CMSLandingPageService) CreateLandingPage(LandingPage *models.LandingPage) (*models.LandingPage, error) {
	LandingContents := LandingPage.Contents
	if len(LandingContents) > 1 {
		return nil, errs.ErrTooMuchContent
	} else if len(LandingContents) == 0 {
		return nil, errs.ErrMissingContent
	}
	// Only one content
	LandingContent := LandingContents[0]

	// Always need to have a revision
	if LandingContent.Revision == nil {
		return nil, errs.ErrNoRevisionFound
	}

	// Check if the URL Alias is duplicate or not
	isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(LandingContent.UrlAlias, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if isUrlAliasDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	// Normalize main content
	if err := helpers.NormalizeLandingContent(LandingContent); err != nil {
		return nil, err
	}

	// Normalize revision
	if err := helpers.NormalizeRevision(LandingContent.Revision); err != nil {
		return nil, err
	}

	// Normalize categories if present
	// for _, category := range LandingContent.Categories {
	// 	if err := helpers.NormalizeCategory(category); err != nil {
	// 		return nil, err
	// 	}
	// }

	return s.repo.CreateLandingPage(LandingPage)
}

func (s *CMSLandingPageService) FindLandingPages(rawQuery string, sort string, page, limit int, language string) ([]models.LandingPage, int64, error) {
	var query dto.LandingPageQuery
	if rawQuery != "" {
		if err := json.Unmarshal([]byte(rawQuery), &query); err != nil {
			return nil, 0, errs.ErrInvalidQuery
		}
	}

	return s.repo.FindAllLandingPage(query, sort, page, limit, language)
}

func (s *CMSLandingPageService) FindLandingPageById(id uuid.UUID) (*models.LandingPage, error) {
	return s.repo.FindLandingPageById(id)
}

func (s *CMSLandingPageService) UpdateLandingContent(updatedLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error) {

	landingContentId, err := s.repo.GetPageIdByContentId(prevContentId)
	if err != nil {
		return nil, err
	}

	isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(updatedLandingContent.UrlAlias, landingContentId)
	if err != nil {
		return nil, err
	}
	if isUrlAliasDuplicate {
		return nil, errs.ErrDuplicateURL
	}

	// Always need to have a revision
	if updatedLandingContent.Revision == nil {
		return nil, errs.ErrNoRevisionFound
	}

	log.Printf("[SERVICE-IN] Language from Request: '%s'", updatedLandingContent.Language)

	savedContent, err := s.repo.UpdateLandingContent(updatedLandingContent, prevContentId)
	if err != nil {
		return nil, err
	}

	// Email sending logic
	if savedContent.WorkflowStatus == enums.WorkflowWaitingDesign && len(savedContent.ApprovalEmail) > 0 {

		go s.triggerApprovalNotifications(savedContent)
	}

	log.Printf("[SERVICE-OUT] Language from Repo: '%s'", savedContent.Language)

	return savedContent, nil
}

// triggerApprovalNotifications sends approval notifications to the relevant recipients.
func (s *CMSLandingPageService) triggerApprovalNotifications(content *models.LandingContent) {
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

	previewURL, cmsEditURL := s.buildNotificationURLs("landing", content)
	authorName, authorEmail := s.getAuthorInfo(content)

	emailData := map[string]interface{}{
		"urlPreview": previewURL,
		"urlCms":     cmsEditURL,
		"pageTitle":  content.Title,
		"author":     authorName,
	}
	// --- REVISED AND SIMPLIFIED LOOP LOGIC ---
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

func (s *CMSLandingPageService) buildNotificationURLs(pageType string, content *models.LandingContent) (previewUrl, cmsEditUrl string) {
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
		// สร้าง Fallback URL แบบ Relative
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

// getAuthorInfo extracts the author's name and email from the content's revision, with fallbacks.
func (s *CMSLandingPageService) getAuthorInfo(content *models.LandingContent) (name, email string) {
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
		log.Printf("Warning: Revision or Author information missing for LandingContent ID %s. Using default author info.", content.ID)
	}

	// --- Future Improvement ---
	// ในอนาคตคุณอาจจะรับ User Context ที่ได้จาก Middleware เข้ามาในฟังก์ชันนี้
	// แล้วใช้เป็นข้อมูล Author ที่ถูกต้องที่สุดแทนค่า Default
	// if email == "" && c.Locals("userEmail") != nil {
	//    email = c.Locals("userEmail").(string)
	// }

	return
}

func (s *CMSLandingPageService) DeleteLandingPage(id uuid.UUID) error {
	return s.repo.DeleteLandingPage(id)
}

func (s *CMSLandingPageService) FindContentByLandingPageId(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return nil, err
	}

	return s.repo.FindContentByLandingPageId(pageId, language, mode)
}

func (r *CMSLandingPageService) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.LandingContent, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	return r.repo.FindLatestContentByPageId(pageId, language)
}

func (s *CMSLandingPageService) DeleteContentByLandingPageId(pageId uuid.UUID, language, mode string) error {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return err
	}

	mode, err = helpers.NormalizeMode(mode)
	if err != nil {
		return err
	}

	return s.repo.DeleteLandingContent(pageId, language, mode)
}

func (s *CMSLandingPageService) DuplicateLandingPage(pageId uuid.UUID) (*models.LandingPage, error) {
	landingPage, err := s.repo.DuplicateLandingPage(pageId)

	if err != nil {
		return nil, err
	}

	return landingPage, nil
}

func (s *CMSLandingPageService) DuplicateLandingContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	err := helpers.NormalizeRevision(newRevision)

	if err != nil {
		return nil, err
	}

	return s.repo.DuplicateLandingContentToAnotherLanguage(contentId, newRevision)
}

func (s *CMSLandingPageService) RevertLandingContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	return s.repo.RevertLandingContent(revisionId, newRevision)
}

func (s *CMSLandingPageService) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
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

func (s *CMSLandingPageService) FindRevisions(pageId uuid.UUID, language string) ([]models.Revision, error) {
	language, err := helpers.NormalizeLanguage(language)
	if err != nil {
		return nil, err
	}

	revisions, err := s.repo.GetRevisionByLandingPageId(pageId, language)
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

func (s *CMSLandingPageService) PreviewLandingContent(pageId uuid.UUID, landingContentPreview *models.LandingContent) (string, error) {
	urls := s.cfg.App.FrontendURLS // "http://localhost:8000,http://localhost:3000"
	parts := strings.Split(urls, ",")
	appUrl := parts[1] // "http://localhost:3000"

	// Default value for preview content
	landingContentPreview.Mode = enums.PageModePreview
	landingContentPreview.PublishStatus = enums.PublishStatusNotPublished
	landingContentPreview.WorkflowStatus = enums.WorkflowUnPublished
	landingContentPreview.ExpiredAt = time.Now().Add(2 * time.Hour)

	if err := helpers.NormalizeLandingContent(landingContentPreview); err != nil {
		return "", err
	}

	// Check if the URL Alias is duplicate or not
	isUrlAliasDuplicate, err := s.repo.IsUrlAliasDuplicate(landingContentPreview.UrlAlias, pageId)
	if err != nil {
		return "", err
	}
	if isUrlAliasDuplicate {
		return "", errs.ErrDuplicateUrlAlias
	}

	existingPreviewContent, err := s.repo.FindLandingContentPreviewById(pageId, string(landingContentPreview.Language))

	// Internal error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	// Record not found, so we create a new one
	if err == gorm.ErrRecordNotFound {
		landingContentPreview.PageID = pageId
		createdLandingContentPreview, err := s.repo.CreateLandingContentPreview(landingContentPreview)
		if err != nil {
			return "", err
		}

		previewUrl, err := helpers.BuildPreviewURL(
			appUrl,
			string(createdLandingContentPreview.Language),
			"landing",
			createdLandingContentPreview.ID,
		)

		if err != nil {
			return "", err
		}

		return previewUrl, nil
	}

	// Attach the old id, so we can save the new one in its place
	landingContentPreview.ID = existingPreviewContent.ID
	landingContentPreview.PageID = existingPreviewContent.PageID
	updatedLandingContentPreview, err := s.repo.UpdateLandingContentPreview(landingContentPreview)
	if err != nil {
		return "", err
	}

	previewUrl, err := helpers.BuildPreviewURL(
		appUrl,
		string(updatedLandingContentPreview.Language),
		"landing",
		updatedLandingContentPreview.ID,
	)

	if err != nil {
		return "", err
	}

	return previewUrl, nil
}
