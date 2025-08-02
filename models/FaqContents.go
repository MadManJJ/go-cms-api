package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
)

type FaqContents struct {
	ID             uuid.UUID            `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PageID         uuid.UUID            `gorm:"not null" json:"page_id"`
	Title          string               `json:"title"`
	Language       enums.PageLanguage   `json:"language"`
	AuthoredAt     time.Time            `json:"authored_at"`
	HTMLInput      string               `json:"html_input"`
	Mode           enums.PageMode       `json:"mode"`
	WorkflowStatus enums.WorkflowStatus `json:"workflow_status"`
	PublishStatus  enums.PublishStatus  `json:"publish_status"`      
	PublishOn      time.Time            `json:"publish_on"`
	UnpublishOn    time.Time            `json:"unpublish_on"`
	AuthoredOn     time.Time            `json:"authored_on"`
	URLAlias       string               `gorm:"not null" json:"url_alias"`
	URL            string               `gorm:"not null" json:"url"`
	MetaTagID      uuid.UUID            `gorm:"unique" json:"meta_tag_id"`
	MetaTag        *MetaTag             `gorm:"foreignKey:MetaTagID" json:"meta_tag,omitempty"`
	ExpiredAt      time.Time            `json:"expired_at"`
	CreatedAt      time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time            `gorm:"autoUpdateTime" json:"updated_at"`

	Page       *FaqPage     `gorm:"foreignKey:PageID" json:"page,omitempty"`
	Revision   *Revision    `gorm:"foreignKey:FaqContentID" json:"revision,omitempty"`
	Categories []*Category  `gorm:"many2many:faq_content_categories;joinForeignKey:FaqContentID;joinReferences:CategoryID" json:"categories,omitempty"`
	Components []*Component `gorm:"foreignKey:FaqContentID" json:"components,omitempty"`
}
