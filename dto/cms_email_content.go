package dto

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"
)

type EmailContentDetailBase struct {
	SendTo          string `json:"send_to" validate:"omitempty"` // ไม่มี custom validator ก็ใช้ omitempty ไปก่อน
	CcEmail         string `json:"cc_email" validate:"omitempty"`
	BccEmail        string `json:"bcc_email" validate:"omitempty"`
	SendFromEmail   string `json:"send_from_email" validate:"required,email"`
	SendFromName    string `json:"send_from_name" validate:"omitempty,max=100"`
	Subject         string `json:"subject" validate:"required,max=255"`
	TopImgLink      string `json:"top_img_link" validate:"omitempty,url,max=255"`
	Header          string `json:"header" validate:"omitempty"`
	Paragraph       string `json:"paragraph" validate:"omitempty"`
	Footer          string `json:"footer" validate:"omitempty"`
	FooterImageLink string `json:"footer_image_link" validate:"omitempty,url,max=255"`
}

type CreateEmailContentRequest struct {
	EmailCategoryID string             `json:"email_category_id" validate:"required,uuid"`
	Language        enums.PageLanguage `json:"language" validate:"required,oneof=th en"`
	Label           string             `json:"label" validate:"required,min=3,max=100"` // ใช้ rule ทั่วไปก่อน
	EmailContentDetailBase
}

type UpdateEmailContentRequest struct {
	Language *enums.PageLanguage `json:"language,omitempty" validate:"omitempty,oneof=th en"`
	Label    *string             `json:"label,omitempty" validate:"omitempty,min=3,max=100"`
	SendTo   *string             `json:"send_to,omitempty" validate:"omitempty"`
	CcEmail  *string             `json:"cc_email,omitempty" validate:"omitempty"`
	BccEmail *string             `json:"bcc_email,omitempty" validate:"omitempty"`

	SendFromEmail   *string `json:"send_from_email,omitempty" validate:"omitempty,email"`
	SendFromName    *string `json:"send_from_name,omitempty" validate:"omitempty,max=100"`
	Subject         *string `json:"subject,omitempty" validate:"omitempty,max=255"`
	TopImgLink      *string `json:"top_img_link,omitempty" validate:"omitempty,url,max=255"`
	Header          *string `json:"header,omitempty"`
	Paragraph       *string `json:"paragraph,omitempty"`
	Footer          *string `json:"footer,omitempty"`
	FooterImageLink *string `json:"footer_image_link,omitempty" validate:"omitempty,url,max=255"`
}

type EmailContentFilter struct {
	EmailCategoryID *string             `query:"email_category_id" validate:"omitempty,uuid"`
	Language        *enums.PageLanguage `query:"language" validate:"omitempty,oneof=th en"`
	Label           *string             `query:"label" validate:"omitempty,min=3,max=100"`
}
type EmailContentResponse struct {
	ID              string                 `json:"id"`
	EmailCategoryID string                 `json:"email_category_id"`
	EmailCategory   *EmailCategoryResponse `json:"email_category"` // Optional: for richer response
	Language        enums.PageLanguage     `json:"language"`
	Label           string                 `json:"label"`
	SendTo          string                 `json:"send_to" validate:"omitempty"` // ไม่มี custom validator ก็ใช้ omitempty ไปก่อน
	CcEmail         string                 `json:"cc_email" validate:"omitempty"`
	BccEmail        string                 `json:"bcc_email" validate:"omitempty"`
	SendFromEmail   string                 `json:"send_from_email" validate:"required,email"`
	SendFromName    string                 `json:"send_from_name" validate:"omitempty,max=100"`
	Subject         string                 `json:"subject" validate:"required,max=255"`
	TopImgLink      string                 `json:"top_img_link" validate:"omitempty,url,max=255"`
	Header          string                 `json:"header" validate:"omitempty"`
	Paragraph       string                 `json:"paragraph" validate:"omitempty"`
	Footer          string                 `json:"footer" validate:"omitempty"`
	FooterImageLink string                 `json:"footer_image_link" validate:"omitempty,url,max=255"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
