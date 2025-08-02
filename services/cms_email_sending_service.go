package services

import (
	"errors"
	"fmt"
	"html"
	"log" // For simple logging, replace with your preferred logger
	"regexp"
	"strings"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSendingServiceInterface interface {
	SendEmail(req dto.SendEmailRequest) error
}

type emailSendingService struct {
	cfg          *config.Config
	categoryRepo repositories.EmailCategoryRepositoryInterface
	contentRepo  repositories.EmailContentRepositoryInterface
}

func NewEmailSendingService(
	cfg *config.Config,
	categoryRepo repositories.EmailCategoryRepositoryInterface,
	contentRepo repositories.EmailContentRepositoryInterface,
) EmailSendingServiceInterface {
	return &emailSendingService{
		cfg:          cfg,
		categoryRepo: categoryRepo,
		contentRepo:  contentRepo,
	}
}

func substitutePlaceholders(templateString string, data map[string]interface{}) string {
	re := regexp.MustCompile(`(?i)\[([\w.-]+)\]|\{\{([\w.-]+)\}\}`)
	return re.ReplaceAllStringFunc(templateString, func(match string) string {
		keyWithBrackets := strings.Trim(match, "[]{}")
		if val, ok := data[keyWithBrackets]; ok {
			return fmt.Sprintf("%v", val)
		}
		lowerKey := strings.ToLower(keyWithBrackets)
		for k, v := range data {
			if strings.ToLower(k) == lowerKey {
				return fmt.Sprintf("%v", v)
			}
		}
		return match
	})
}

func formatEmailBody(content *models.EmailContent, data map[string]interface{}) string {
	var bodyBuilder strings.Builder

	if content.TopImgLink != "" {
		bodyBuilder.WriteString(fmt.Sprintf("<div><img src=\"%s\" alt=\"Email Top Image\" style=\"max-width: 100%%; height: auto;\"/></div>", html.EscapeString(content.TopImgLink)))
	}

	if content.Header != "" {
		// สมมติว่า Header เป็น plain text ที่จะถูกหุ้มด้วย <h1> และอาจมี placeholders
		// ถ้า Header เองเป็น HTML จาก DB ส่วนนี้อาจจะต้องปรับ
		substitutedHeaderText := substitutePlaceholders(content.Header, data)
		bodyBuilder.WriteString(fmt.Sprintf("<div><h1>%s</h1></div>", html.EscapeString(substitutedHeaderText))) // Escape text ก่อนใส่ใน h1
	}

	if content.Paragraph != "" {
		htmlParagraphWithData := substitutePlaceholders(content.Paragraph, data)
		bodyBuilder.WriteString(fmt.Sprintf("<div>%s</div>", htmlParagraphWithData)) // ใช้ HTML ที่มีข้อมูลแทนที่แล้วโดยตรง
	}

	if content.Footer != "" {
		htmlFooterWithData := substitutePlaceholders(content.Footer, data)
		bodyBuilder.WriteString(fmt.Sprintf("<div>%s</div>", htmlFooterWithData))
	}

	if content.FooterImageLink != "" {
		bodyBuilder.WriteString(fmt.Sprintf("<div><img src=\"%s\" alt=\"Email Footer Image\" style=\"max-width: 100%%; height: auto;\"/></div>", html.EscapeString(content.FooterImageLink)))
	}
	return bodyBuilder.String()
}

func (s *emailSendingService) SendEmail(req dto.SendEmailRequest) error {
	var category *models.EmailCategory
	var err error

	// 1. Find Email Category
	catID, parseErr := uuid.Parse(req.EmailCategoryTitleOrID)
	if parseErr == nil {
		category, err = s.categoryRepo.FindByID(catID)
	} else {
		category, err = s.categoryRepo.FindByTitle(req.EmailCategoryTitleOrID)
	}

	if err != nil {
		return fmt.Errorf("failed to find email category '%s': %w", req.EmailCategoryTitleOrID, err)
	}
	if category == nil {
		return fmt.Errorf("email category '%s' not found", req.EmailCategoryTitleOrID)
	}

	// 2. Find Email Content based on Category ID, Language, and Label
	content, err := s.contentRepo.FindByCategoryIDAndLanguageAndLabel(category.ID, req.Language, req.EmailContentLabel)
	if err != nil {
		return fmt.Errorf("failed to find email content for category ID '%s', lang '%s', label '%s': %w", category.ID, req.Language, req.EmailContentLabel, err)
	}
	if content == nil {
		return fmt.Errorf("email content template not found for category ID '%s', lang '%s', label '%s'", category.ID, req.Language, req.EmailContentLabel)
	}

	// 3. Determine Recipients
	var finalToEmails []*mail.Email
	if len(req.ToRecipientEmails) > 0 {
		for _, emailStr := range req.ToRecipientEmails {
			if emailStr != "" {
				finalToEmails = append(finalToEmails, mail.NewEmail("", emailStr))
			}
		}
	} else if content.SendTo != "" {
		templateTos := strings.Split(content.SendTo, ",")
		for _, emailStr := range templateTos {
			trimmedEmail := strings.TrimSpace(emailStr)
			if trimmedEmail != "" {
				finalToEmails = append(finalToEmails, mail.NewEmail("", trimmedEmail))
			}
		}
	}

	if len(finalToEmails) == 0 {
		log.Printf("No recipient email addresses for category '%s', lang '%s', label '%s'. req.ToRecipientEmails: %v, content.SendTo: %s",
			req.EmailCategoryTitleOrID, req.Language, req.EmailContentLabel, req.ToRecipientEmails, content.SendTo)
		return errors.New("no recipient email addresses provided or configured in the template")
	}

	// 4. Prepare Email Details
	fromName := content.SendFromName
	if fromName == "" {
		fromName = s.cfg.App.AppName // Default App Name from config
	}
	from := mail.NewEmail(fromName, content.SendFromEmail)
	subject := substitutePlaceholders(content.Subject, req.Data)
	htmlContent := formatEmailBody(content, req.Data)

	// 5. Construct SendGrid Message
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject

	p := mail.NewPersonalization()
	for _, toEmail := range finalToEmails {
		p.AddTos(toEmail)
	}

	if content.CcEmail != "" {
		ccEmails := strings.Split(content.CcEmail, ",")
		for _, cc := range ccEmails {
			trimmedCC := strings.TrimSpace(cc)
			if trimmedCC != "" {
				p.AddCCs(mail.NewEmail("", trimmedCC))
			}
		}
	}

	if content.BccEmail != "" {
		bccEmails := strings.Split(content.BccEmail, ",")
		for _, bcc := range bccEmails {
			trimmedBCC := strings.TrimSpace(bcc)
			if trimmedBCC != "" {
				p.AddBCCs(mail.NewEmail("", trimmedBCC))
			}
		}
	}
	message.AddPersonalizations(p)
	message.AddContent(mail.NewContent("text/html", htmlContent))

	// 6. Send Email
	apiKey := s.cfg.SendGrid.APIKey
	if apiKey == "" {
		log.Println("WARNING: SENDGRID_API_KEY is not set. Email will not be sent.")
		return errors.New("email sending is not configured (missing API key)")
	}

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email via SendGrid to %v: %v", finalToEmails, err)
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	if response.StatusCode >= 300 {
		log.Printf("SendGrid returned non-success status: %d, Body: %s, To: %v", response.StatusCode, response.Body, finalToEmails)
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	var recipientLog []string
	for _, email := range finalToEmails {
		recipientLog = append(recipientLog, email.Address)
	}
	log.Printf("Email sent successfully. Category: %s, Label: %s, Lang: %s, To: %s, Subject: %s",
		req.EmailCategoryTitleOrID, req.EmailContentLabel, req.Language, strings.Join(recipientLog, ", "), subject)
	return nil
}
