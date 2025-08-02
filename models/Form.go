package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Form struct {
	ID              uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name            string              `gorm:"type:varchar(255);not null" json:"name"`
	Slug            string              `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Description     *string             `gorm:"type:text" json:"description"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	DeletedAt       gorm.DeletedAt      `gorm:"index" json:"-" swaggerignore:"true"`
	EmailCategoryID *uuid.UUID          `gorm:"type:uuid" json:"email_category_id,omitempty"`
	EmailCategory   *EmailCategory      `gorm:"foreignKey:EmailCategoryID" json:"email_category,omitempty"`
	Language        *enums.PageLanguage `gorm:"type:varchar(10)" json:"language,omitempty"`
	Sections        []FormSection       `gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE" json:"sections,omitempty"`
	Submissions     []FormSubmission    `gorm:"foreignKey:FormID;constraint:OnDelete:RESTRICT" json:"-"`
}
