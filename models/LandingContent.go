package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type LandingContent struct {
	ID             uuid.UUID            `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PageID         uuid.UUID            `gorm:"not null" json:"page_id"`
	Title          string               `json:"title"`
	Language       enums.PageLanguage   `json:"language"`
	AuthoredAt     time.Time            `json:"authored_at"`
	HTMLInput      string               `json:"html_input"`
	Mode           enums.PageMode       `json:"mode"`
	WorkflowStatus enums.WorkflowStatus `json:"workflow_status"`
	PublishStatus  enums.PublishStatus  `json:"publish_status"`
	UrlAlias       string               `gorm:"not null" json:"url_alias"` // Ensure unique UrlAlias
	MetaTagID      uuid.UUID            `gorm:"unique" json:"meta_tag_id"` // MetaTag is 1-to-1
	MetaTag        *MetaTag             `gorm:"foreignKey:MetaTagID" json:"meta_tag,omitempty"`
	PublishOn      *time.Time           `json:"publish_on,omitempty"`
	UnpublishOn    *time.Time           `json:"unpublish_on,omitempty"`
	AuthoredOn     *time.Time           `json:"authored_on,omitempty"`
	ExpiredAt      time.Time            `json:"expired_at"`
	CreatedAt      time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time            `gorm:"autoUpdateTime" json:"updated_at"`

	Page          *LandingPage          `gorm:"foreignKey:PageID" json:"page,omitempty"`
	ApprovalEmail pq.StringArray        `gorm:"type:text[]" json:"approval_email"`
	Files         []*LandingContentFile `gorm:"foreignKey:LandingContentID" json:"files,omitempty"`
	Revision      *Revision             `gorm:"foreignKey:LandingContentID" json:"revision,omitempty"`
	Categories    []*Category           `gorm:"many2many:landing_content_categories;joinForeignKey:LandingContentID;joinReferences:CategoryID" json:"categories,omitempty"`
	Components    []*Component          `gorm:"foreignKey:LandingContentID" json:"components,omitempty"`
}
