// models/Url.go
package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums" // ตรวจสอบ path

	"github.com/google/uuid"
)

// UrlType defines the type of content the URL points to.
type UrlType string

const (
	UrlTypeLandingPages UrlType = "landing_pages"
	UrlTypeFaqPages     UrlType = "faq_pages"
	UrlTypePartnerPages UrlType = "partner_pages"
)

type Url struct {
	ID        uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Path      string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_url_path_lang_mode_alias"` // The actual URL path segment
	Type      UrlType             `gorm:"type:varchar(50);not null"`
	ContentID uuid.UUID           `gorm:"type:uuid;not null;index:idx_url_path_lang_mode_alias"`                            // ID of the Page (LandingPage, FaqPage, etc.)
	Language  *enums.PageLanguage `gorm:"type:varchar(10);index:idx_url_path_lang_mode_alias"`                              // Optional: for language-specific URLs
	Mode      enums.PageMode      `gorm:"type:varchar(50);not null;default:'Published';index:idx_url_path_lang_mode_alias"` // Published, Preview
	IsAlias   bool                `gorm:"default:false;index:idx_url_path_lang_mode_alias"`                                 // True if this is an additional alias, false for main URL
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"updated_at"`	
}
