package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
)

type EmailContent struct {
	ID              uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EmailCategoryID uuid.UUID          `gorm:"type:uuid;not null;index:idx_email_content_category_lang_label" json:"email_category_id"`
	Language        enums.PageLanguage `gorm:"type:varchar(10);not null;index:idx_email_content_category_lang_label" json:"language"`
	Label           string             `gorm:"type:varchar(255);not null;index:idx_email_content_category_lang_label" json:"label"`

	SendTo          string             `gorm:"type:varchar(255)" json:"send_to"`
	CcEmail         string             `gorm:"type:varchar(255)" json:"cc_email"`
	BccEmail        string             `gorm:"type:varchar(255)" json:"bcc_email"`
	SendFromEmail   string             `gorm:"type:varchar(100);not null" json:"send_from_email"`
	SendFromName    string             `gorm:"type:varchar(100)" json:"send_from_name"`
            
	Subject         string             `gorm:"type:varchar(255);not null" json:"subject"`
	TopImgLink      string             `gorm:"type:varchar(255)" json:"top_img_link"`
	Header          string             `gorm:"type:text" json:"header"`
	Paragraph       string             `gorm:"type:text" json:"paragraph"`
	Footer          string             `gorm:"type:text" json:"footer"`
	FooterImageLink string             `gorm:"type:varchar(255)" json:"footer_image_link"`

	CreatedAt       time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time          `gorm:"autoUpdateTime" json:"updated_at"`

	EmailCategory   *EmailCategory     `gorm:"foreignKey:EmailCategoryID;references:ID" json:"email_category,omitempty"`
}
