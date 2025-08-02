package models

import (
	"github.com/google/uuid"
)

type PartnerContentCategory struct {
	PartnerContentID uuid.UUID      `gorm:"type:uuid;not null;index" json:"partner_content_id"`
	CategoryID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
  
	PartnerContent   *PartnerContent `gorm:"foreignKey:PartnerContentID" json:"partner_content,omitempty"`
	Category         *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

