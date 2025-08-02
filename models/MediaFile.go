package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaFile struct { // .jpg, .pdf, ...
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	DownloadURL string    `gorm:"not null" json:"download_url"`
	CreatedAt      time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time           `gorm:"autoUpdateTime" json:"updated_at"`	
}
