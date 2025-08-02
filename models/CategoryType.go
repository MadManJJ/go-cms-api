package models

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TypeCode  string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"type_code"`
	Name      string       `gorm:"type:varchar(255)" json:"name"`
	IsActive  bool         `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`

	Categories []*Category `gorm:"foreignKey:CategoryTypeID" json:"categories,omitempty"`
}
