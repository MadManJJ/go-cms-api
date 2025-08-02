package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailCategory struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title         string          `gorm:"type:varchar(255);uniqueIndex;not null" json:"title"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	EmailContents []*EmailContent `gorm:"foreignKey:EmailCategoryID" json:"email_contents,omitempty"`
}
