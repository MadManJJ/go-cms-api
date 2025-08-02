package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PartnerContent struct {
	ID               uuid.UUID            `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PageID           uuid.UUID            `gorm:"not null" json:"page_id"`
	Title            string               `json:"title"`
	ThumbnailImage   string               `json:"thumbnail_image"`
	ThumbnailAltText string               `json:"thumbnail_alt_text"`
	CompanyLogo      string               `json:"company_logo"`
	CompanyAltText   string               `json:"company_alt_text"`
	Language         enums.PageLanguage   `json:"language"`
	CompanyName      string               `json:"company_name"`
	CompanyDetail    string               `json:"company_detail"`
	LeadBody         string               `json:"lead_body"`
	Challenges       string               `json:"challenges"`
	Solutions        string               `json:"solutions"`
	Results          string               `json:"results"`
	AuthoredAt       time.Time            `json:"authored_at"`
	HTMLInput        string               `json:"html_input"`
	Mode             enums.PageMode       `json:"mode"`
	WorkflowStatus   enums.WorkflowStatus `json:"workflow_status"`
	PublishStatus    enums.PublishStatus  `json:"publish_status"`
	URLAlias         string               `gorm:"not null" json:"url_alias"`
	URL              string               `gorm:"not null" json:"url"`
	MetaTag          *MetaTag             `gorm:"foreignKey:MetaTagID" json:"meta_tag,omitempty"`
	MetaTagID        uuid.UUID            `gorm:"unique" json:"meta_tag_id"`
	IsRecommended    bool                 `json:"is_recommended"`
	PublishOn        time.Time            `json:"publish_on"`
	UnpublishOn      time.Time            `json:"unpublish_on"`
	AuthoredOn       time.Time            `json:"authored_on"`
	ExpiredAt        time.Time            `json:"expired_at"`
	CreatedAt        time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time            `gorm:"autoUpdateTime" json:"updated_at"`

	Page          *PartnerPage   `gorm:"foreignKey:PageID" json:"page,omitempty"`
	ApprovalEmail pq.StringArray `gorm:"type:text[]" json:"approval_email"`
	Revision      *Revision      `gorm:"foreignKey:PartnerContentID" json:"revision,omitempty"`
	Categories    []*Category    `gorm:"many2many:partner_content_categories;joinForeignKey:PartnerContentID;joinReferences:CategoryID" json:"categories,omitempty"`
	Components    []*Component   `gorm:"foreignKey:PartnerContentID" json:"components,omitempty"`
}
