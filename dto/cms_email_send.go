package dto

import "github.com/MadManJJ/cms-api/models/enums"

type SendEmailRequest struct {
	EmailCategoryTitleOrID string                 `json:"email_category" validate:"required"`       // Can be ID (UUID) or unique Title of the category
	EmailContentLabel      string                 `json:"email_content_label" validate:"required"`  // Label of the specific email content within the category
	Language               enums.PageLanguage     `json:"language" validate:"required,oneof=th en"` // Language of the email content
	ToRecipientEmails      []string               `json:"to_recipient_emails" validate:"omitempty,dive,email"`
	Data                   map[string]interface{} `json:"data" validate:"required"` // Placeholders and their values
}
