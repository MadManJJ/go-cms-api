package dto

import (
	"time"

	"github.com/MadManJJ/cms-api/models/enums"
)

// CategoryDetailRequest

type CategoryResponse struct {
	ID             string                `json:"id"`
	CategoryTypeID string                `json:"category_type_id"`
	CategoryType   *CategoryTypeResponse `json:"category_type,omitempty"` // Optional
	LanguageCode   enums.PageLanguage    `json:"language_code"`
	Name           string                `json:"name"`
	Description    *string               `json:"description,omitempty"`
	Weight         int                   `json:"weight"`
	PublishStatus  enums.PublishStatus   `json:"publish_status"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}
type CategoryCreateRequest struct {
	CategoryTypeID string              `json:"category_type_id" validate:"required,uuid"`
	LanguageCode   enums.PageLanguage  `json:"language_code" validate:"required,oneof=th en"`
	Name           string              `json:"name" validate:"required,min=1,max=255"`
	Description    *string             `json:"description,omitempty"`
	Weight         *int                `json:"weight,omitempty"`
	PublishStatus  enums.PublishStatus `json:"publish_status" validate:"required,oneof=Published UnPublished"`
}

type CategoryUpdateRequest struct {
	// CategoryTypeID และ LanguageCode โดยทั่วไปจะไม่เปลี่ยนสำหรับ Category (Detail) ที่มีอยู่แล้ว
	Name          *string              `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string              `json:"description,omitempty"`
	Weight        *int                 `json:"weight,omitempty"`
	PublishStatus *enums.PublishStatus `json:"publish_status,omitempty" validate:"omitempty,oneof=Published UnPublished"`
}
type CategoryFilter struct {
	CategoryTypeID *string              `query:"category_type_id" validate:"omitempty,uuid"` // Filter by CategoryType ID (UUID)
	LanguageCode   *enums.PageLanguage  `query:"lang" validate:"omitempty,oneof=th en"`      // Filter details by language
	Name           *string              `query:"name" validate:"omitempty,max=255"`
	PublishStatus  *enums.PublishStatus `query:"publish_status" validate:"omitempty,oneof=Published UnPublished"`
}
