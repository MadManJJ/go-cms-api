package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
)

type Revision struct {
	ID               uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LandingContentID *uuid.UUID          `json:"landing_content_id,omitempty"`
	PartnerContentID *uuid.UUID          `json:"partner_content_id,omitempty"`
	FaqContentID     *uuid.UUID          `json:"faq_content_id,omitempty"`
	PublishStatus    enums.PublishStatus `json:"publish_status"`
	CreatedAt        time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	Author           string              `json:"author"`
	Message          string              `json:"message"`
	Description      string              `json:"description"`

	LandingContent *LandingContent `gorm:"foreignKey:LandingContentID" json:"-"`
	PartnerContent *PartnerContent `gorm:"foreignKey:PartnerContentID" json:"-"`
	FaqContent     *FaqContent     `gorm:"foreignKey:FaqContentID" json:"-"`
}
