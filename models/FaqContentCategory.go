package models

import (
	"github.com/google/uuid"
)

type FaqContentCategory struct {
	FaqContentID uuid.UUID `gorm:"type:uuid;not null;index" json:"faq_content_id"`
	CategoryID   uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`

	FaqContent *FaqContent `gorm:"foreignKey:FaqContentID" json:"faq_content,omitempty"`
	Category   *Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
