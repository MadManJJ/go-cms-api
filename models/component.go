package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Component struct {
	ID               uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LandingContentID *uuid.UUID          `json:"landing_content_id,omitempty"`
	LandingContent   *LandingContent     `gorm:"foreignKey:LandingContentID" json:"landing_content,omitempty"`
	PartnerContentID *uuid.UUID          `json:"partner_content_id,omitempty"`
	PartnerContent   *PartnerContent     `gorm:"foreignKey:PartnerContentID" json:"partner_content,omitempty"`
	FaqContentID     *uuid.UUID          `json:"faq_content_id,omitempty"`
	FaqContent       *FaqContent         `gorm:"foreignKey:FaqContentID" json:"faq_content,omitempty"`
	Type             enums.ComponentType `json:"type,omitempty"`
	Props            datatypes.JSON      `gorm:"type:jsonb" json:"props,omitempty"`
	CreatedAt        time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
}
