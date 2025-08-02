package models

import (
	"time"

	"github.com/google/uuid"
)

type PartnerPage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Contents []*PartnerContent `gorm:"foreignKey:PageID" json:"contents,omitempty"`
}
