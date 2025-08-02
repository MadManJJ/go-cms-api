package dto

import "time"

// CategoryTypeResponse matches OpenAPI components.schemas.CategoryType
type CategoryTypeResponse struct {
	ID            string          `json:"id"`
	TypeCode      string          `json:"type_code"`
	Name          *string         `json:"name,omitempty"` // OpenAPI บอก nullable
	IsActive      bool            `json:"is_active"`
	ChildrenCount *map[string]int `json:"children_count,omitempty"` // e.g., {"th": 2, "en": 1}
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type CreateCategoryTypeRequest struct {
	TypeCode string  `json:"type_code" validate:"required"`
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"` // Default true ใน model
}

// UpdateCategoryTypeRequest
type UpdateCategoryTypeRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type CategoryTypeWithDetailsResponse struct {
	ID         string             `json:"id"` // CategoryType ID
	TypeCode   string             `json:"type_code"`
	Name       *string            `json:"name,omitempty"`
	IsActive   bool               `json:"is_active"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Categories []CategoryResponse `json:"categories"` // List of Category (Detail) responses for the specified language
}
