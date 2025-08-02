package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     *string            `json:"email,omitempty" validate:"omitempty,email" gorm:"uniqueIndex"`
	Password  *string            `json:"password,omitempty" validate:"omitempty,min=6"`
	Provider  enums.ProviderType `json:"provider" gorm:"type:varchar(20);default:'normal'"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}
