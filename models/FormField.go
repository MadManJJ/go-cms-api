// package models

// import (
// 	"github.com/MadManJJ/cms-api/models/enums"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/datatypes"
// )

// type FormField struct {
// 	ID        uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
// 	SectionID uuid.UUID           `gorm:"type:uuid;not null" json:"section_id"`
// 	Label     string              `gorm:"type:varchar(255);not null" json:"label"`
// 	FieldKey  string              `gorm:"type:varchar(100);not null" json:"field_key"`
// 	FieldType enums.FormFieldType `gorm:"type:form_field_type_enum;not null" json:"field_type"`
// 	Placeholder  *string        `gorm:"type:varchar(255)" json:"placeholder"`
// 	IsRequired   bool           `gorm:"not null;default:false" json:"is_required"`
// 	DefaultValue *string        `gorm:"type:text" json:"default_value"`
// 	Properties   datatypes.JSON `gorm:"type:jsonb" json:"properties"`
// 	Display      datatypes.JSON `gorm:"type:jsonb" json:"display"`
// 	OrderIndex   int            `gorm:"not null;default:0" json:"order_index"`
// 	CreatedAt    time.Time      `json:"created_at"`
// 	UpdatedAt    time.Time      `json:"updated_at"`

//		Section *FormSection `gorm:"foreignKey:SectionID" json:"-"`
//	}
//
// In models/form_field.go
package models

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SwagJSON map[string]interface{}

type FormField struct {
	ID           uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SectionID    uuid.UUID           `gorm:"type:uuid;not null" json:"section_id"`
	Label        string              `gorm:"type:varchar(255);not null" json:"label"`
	FieldKey     string              `gorm:"type:varchar(100);not null" json:"field_key"`
	FieldType    enums.FormFieldType `gorm:"type:form_field_type_enum;not null" json:"field_type"`
	Placeholder  *string             `gorm:"type:varchar(255)" json:"placeholder,omitempty"`
	IsRequired   bool                `gorm:"not null;default:false" json:"is_required"`
	DefaultValue *string             `gorm:"type:text" json:"default_value,omitempty"`
	Properties   datatypes.JSON      `gorm:"type:jsonb" json:"properties,omitempty" swaggertype:"object,string"`
	Display      datatypes.JSON      `gorm:"type:jsonb" json:"display,omitempty" swaggertype:"object,string"`
	OrderIndex   int                 `gorm:"not null;default:0" json:"order_index"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	Section      *FormSection        `gorm:"foreignKey:SectionID" json:"-"`
}
