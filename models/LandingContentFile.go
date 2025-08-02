package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
)

type LandingContentFile struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LandingContentID uuid.UUID       `gorm:"not null" json:"landing_content_id"`
	LandingContent   *LandingContent `gorm:"foreignKey:LandingContentID" json:"landing_content,omitempty"`
	Name             string          `gorm:"not null" json:"name"`
	DownloadURL      string          `gorm:"not null" json:"download_url"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`	
	FileType         enums.FileType  `gorm:"not null" json:"file_type"` // e.g., "image", "video", "document"
}
