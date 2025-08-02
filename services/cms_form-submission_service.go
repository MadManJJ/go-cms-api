package services

import (
	"log"
	"strings"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
)

type CMSFormSubmissionServiceInterface interface {
	CreateFormSubmission(formId uuid.UUID, formSubmission *models.FormSubmission) (*models.FormSubmission, error)
	GetFormSubmissions(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error)
	GetFormSubmission(submissionId uuid.UUID) (*models.FormSubmission, error)
}

type CMSFormSubmissionService struct {
	repo                repositories.FormSubmissionRepositoryInterface
	emailSendingService EmailSendingServiceInterface
}

func NewCMSFormSubmissionService(repo repositories.FormSubmissionRepositoryInterface, emailSendingService EmailSendingServiceInterface) CMSFormSubmissionServiceInterface {
	return &CMSFormSubmissionService{repo: repo, emailSendingService: emailSendingService}
}

func (s *CMSFormSubmissionService) CreateFormSubmission(formId uuid.UUID, formSubmission *models.FormSubmission) (*models.FormSubmission, error) {

	formSubmission.FormID = formId
	formSubmission.SubmittedAt = time.Now()
	createdFormSubmission, err := s.repo.CreateFormSubmission(formSubmission)
	if err != nil {
		log.Printf("ERROR: Failed to create submission in repo: %v", err)
		return nil, err
	}
	log.Printf("INFO: Submission %s created successfully.", createdFormSubmission.ID)

	if createdFormSubmission.Form == nil || createdFormSubmission.Form.EmailCategoryID == nil {
		log.Printf("INFO: No email will be sent for submission %s because the form is not linked to an EmailCategory.", createdFormSubmission.ID)
		return createdFormSubmission, nil
	}
	log.Printf("DEBUG: Form %s is linked to EmailCategory %s. Proceeding to send emails.", createdFormSubmission.Form.ID, createdFormSubmission.Form.EmailCategoryID)

	emailContents, err := s.repo.GetEmailContentsFormFormId(formId)
	if err != nil {
		log.Printf("WARNING: Could not get email contents for form %s. Emails will not be sent. Error: %v", formId, err)
		return createdFormSubmission, nil
	}

	for _, emailContent := range emailContents {
		if emailContent == nil {
			continue
		}

		var recipients []string
		isUserEmail := false

		if strings.Contains(strings.ToLower(emailContent.Label), "user") {

			if createdFormSubmission.SubmittedEmail != nil && *createdFormSubmission.SubmittedEmail != "" {
				recipients = append(recipients, *createdFormSubmission.SubmittedEmail)
				isUserEmail = true
			} else {
				log.Printf("INFO: Skipping email content with label '%s' because submitted_email is empty.", emailContent.Label)
				continue
			}
		} else {

			if emailContent.SendTo != "" {
				recipients = append(recipients, emailContent.SendTo)
			}
			if emailContent.CcEmail != "" {
				recipients = append(recipients, emailContent.CcEmail)
			}
			if emailContent.BccEmail != "" {
				recipients = append(recipients, emailContent.BccEmail)
			}
		}

		if len(recipients) == 0 {
			log.Printf("INFO: Skipping email content with label '%s' because there are no recipients.", emailContent.Label)
			continue
		}

		var categoryTitle string
		if emailContent.EmailCategory != nil {
			categoryTitle = emailContent.EmailCategory.Title
		}

		data := map[string]interface{}{
			"submittedData":   createdFormSubmission.SubmittedData,
			"sendFromEmail":   emailContent.SendFromEmail,
			"sendFromName":    emailContent.SendFromName,
			"subject":         emailContent.Subject,
			"topImgLink":      emailContent.TopImgLink,
			"header":          emailContent.Header,
			"paragraph":       emailContent.Paragraph,
			"footer":          emailContent.Footer,
			"footerImageLink": emailContent.FooterImageLink,
		}

		req := dto.SendEmailRequest{
			EmailCategoryTitleOrID: categoryTitle,
			EmailContentLabel:      emailContent.Label,
			Language:               emailContent.Language,
			ToRecipientEmails:      recipients,
			Data:                   data,
		}

		log.Printf("INFO: Queuing email for label '%s' to recipients: %v", req.EmailContentLabel, req.ToRecipientEmails)
		go func(req dto.SendEmailRequest, isUserMail bool) {
			if err := s.emailSendingService.SendEmail(req); err != nil {
				if isUserMail {
					log.Printf("ERROR: Failed to send user confirmation email. Label: '%s'. Error: %v", req.EmailContentLabel, err)
				} else {
					log.Printf("ERROR: Failed to send admin notification email. Label: '%s'. Error: %v", req.EmailContentLabel, err)
				}
			}
		}(req, isUserEmail)
	}

	return createdFormSubmission, nil
}
func (s *CMSFormSubmissionService) GetFormSubmissions(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error) {
	return s.repo.GetFormSubmissions(formId, sort, page, limit)
}

func (s *CMSFormSubmissionService) GetFormSubmission(submissionId uuid.UUID) (*models.FormSubmission, error) {
	return s.repo.GetFormSubmission(submissionId)
}
