package models

import (
	"github.com/google/uuid"
)

type LandingContentCategory struct {
	LandingContentID uuid.UUID `gorm:"type:uuid;primaryKey" json:"landing_content_id"`
	CategoryID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"category_id"`

	LandingContent *LandingContent `gorm:"foreignKey:LandingContentID" json:"landing_content,omitempty"`
	Category       *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
