package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
)

type Category struct {
	ID             uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CategoryTypeID uuid.UUID           `gorm:"type:uuid;not null;index" json:"category_type_id"`
	LanguageCode   enums.PageLanguage  `gorm:"type:varchar(10);not null;index" json:"language_code"`
	Name           string              `gorm:"type:varchar(255);not null" json:"name"`
	Description    *string             `gorm:"type:text" json:"description,omitempty"`
	Weight         int                 `gorm:"not null;default:0" json:"weight"`
	PublishStatus  enums.PublishStatus `gorm:"type:varchar(50);not null;default:'Unpublished'" json:"publish_status"`
	CreatedAt      time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time           `gorm:"autoUpdateTime" json:"updated_at"`

	CategoryType *CategoryType `gorm:"foreignKey:CategoryTypeID" json:"category_type,omitempty"`
}
