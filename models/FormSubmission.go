package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type FormSubmission struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FormID         uuid.UUID      `gorm:"type:uuid;not null" json:"form_id"`
	SubmittedData  datatypes.JSON `gorm:"type:jsonb;not null" json:"submitted_data" swaggertype:"object,string"`
	SubmittedAt    time.Time      `gorm:"not null;default:now()" json:"submitted_at"`
	SubmittedEmail *string        `gorm:"type:varchar(255)" json:"submitted_email,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Form           *Form          `gorm:"foreignKey:FormID;constraint:OnDelete:RESTRICT" json:"form,omitempty"`
}
