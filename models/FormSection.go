package models

import (
	"time"

	"github.com/google/uuid"
)

type FormSection struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FormID      uuid.UUID `gorm:"type:uuid;not null" json:"form_id"`
	Title       *string   `gorm:"type:varchar(255)" json:"title"` // Pointer for nullable
	Description *string   `gorm:"type:text" json:"description"`   // Pointer for nullable
	OrderIndex  int       `gorm:"not null;default:0" json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Form   *Form       `gorm:"foreignKey:FormID" json:"-"` // Back-reference, json:"-" to avoid circular dependency in basic response
	Fields []FormField `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE" json:"fields,omitempty"`
}
