package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CreateFormRequest struct {
	Name            string               `json:"name" validate:"required,min=1,max=255"`
	Description     *string              `json:"description,omitempty" validate:"omitempty,max=1000"`
	Sections        []FormSectionRequest `json:"sections"`
	EmailCategoryID *string              `json:"email_category_id,omitempty" validate:"omitempty,uuid"`
	Language        *string              `json:"language,omitempty" validate:"omitempty,oneof=th en"`
}

type FormSectionRequest struct {
	Title       *string            `json:"title,omitempty" validate:"omitempty,max=255"`
	Description *string            `json:"description,omitempty" validate:"omitempty,max=1000"`
	OrderIndex  int                `json:"order_index" validate:"required,min=0"`
	Fields      []FormFieldRequest `json:"fields" validate:"required,dive"`
}

type FormFieldRequest struct {
	Label        string         `json:"label" validate:"required,min=1,max=255"`
	FieldKey     string         `json:"field_key" validate:"required,min=1,max=100,fieldkey"`
	FieldType    string         `json:"field_type" validate:"required,oneof=text textarea select checkbox radio date time email number file"`
	IsRequired   bool           `json:"is_required" validate:"required"`
	Placeholder  *string        `json:"placeholder,omitempty" validate:"omitempty,max=255"`
	DefaultValue *string        `json:"default_value,omitempty" validate:"omitempty,max=1000"`
	OrderIndex   int            `json:"order_index" validate:"required,min=0"`
	Properties   datatypes.JSON `json:"properties,omitempty" swaggertype:"object,string"`
	Display      datatypes.JSON `json:"display,omitempty" swaggertype:"object,string"`
}

type UpdateFormRequest struct {
	Name            string                     `json:"name" validate:"required,min=1,max=255"`
	Description     *string                    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Sections        []UpdateFormSectionRequest `json:"sections"`
	EmailCategoryID *string                    `json:"email_category_id,omitempty" validate:"omitempty,uuid"`
	Language        *string                    `json:"language,omitempty" validate:"omitempty,oneof=th en"`
}
type UpdateFormSectionRequest struct {
	Title       *string                  `json:"title,omitempty" validate:"omitempty,max=255"`
	Description *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	OrderIndex  int                      `json:"order_index" validate:"required,min=0"`
	Fields      []UpdateFormFieldRequest `json:"fields" validate:"required,dive"`
}

type UpdateFormFieldRequest struct {
	Label        string         `json:"label" validate:"required,min=1,max=255"`
	FieldKey     string         `json:"field_key" validate:"required,min=1,max=100,fieldkey"`
	FieldType    string         `json:"field_type" validate:"required,oneof=text textarea select checkbox radio date time email number file"`
	IsRequired   bool           `json:"is_required" validate:"required"`
	Placeholder  *string        `json:"placeholder,omitempty" validate:"omitempty,max=255"`
	DefaultValue *string        `json:"default_value,omitempty" validate:"omitempty,max=1000"`
	OrderIndex   int            `json:"order_index" validate:"required,min=0"`
	Properties   datatypes.JSON `json:"properties,omitempty" swaggertype:"object,string"`
	Display      datatypes.JSON `json:"display,omitempty" swaggertype:"object,string"`
}

// Response ของ GET /cms/forms/{formId}
// Response ของ POST /cms/forms
// Response ของ PUT /cms/forms/{formId}
type FormFieldResponse struct {
	ID           uuid.UUID      `json:"id"`
	Label        string         `json:"label"`
	FieldKey     string         `json:"field_key"`
	FieldType    string         `json:"field_type"`
	IsRequired   bool           `json:"is_required"`
	Placeholder  *string        `json:"placeholder,omitempty"`
	DefaultValue *string        `json:"default_value,omitempty"`
	OrderIndex   int            `json:"order_index"`
	Properties   datatypes.JSON `json:"properties,omitempty" swaggertype:"object,string"`
	Display      datatypes.JSON `json:"display,omitempty" swaggertype:"object,string"`
}
type FormSectionResponse struct {
	ID          uuid.UUID           `json:"id"`
	Title       *string             `json:"title,omitempty"`
	Description *string             `json:"description,omitempty"`
	OrderIndex  int                 `json:"order_index"`
	Fields      []FormFieldResponse `json:"fields"`
}

type FormResponse struct {
	ID              uuid.UUID             `json:"id"`
	Name            string                `json:"name"`
	Slug            string                `json:"slug"`
	Description     *string               `json:"description,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	EmailCategoryID *string               `json:"email_category_id,omitempty"`
	Language        *string               `json:"language,omitempty"`
	Sections        []FormSectionResponse `json:"sections"`
}

type FormListItemResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PaginationMeta struct {
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
}

// Response ของ GET /cms/forms
type PaginatedFormListResponse struct {
	Data []FormListItemResponse `json:"data"`
	Meta PaginationMeta         `json:"meta"`
}

type FormListFilter struct {
	Name         *string    `form:"name" json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	CreatedAt    *time.Time `form:"created_at" json:"created_at,omitempty" layout:"2006-01-02"`
	Page         *int       `form:"page" json:"page,omitempty" validate:"omitempty,min=1"`
	ItemsPerPage *int       `form:"items_per_page" json:"items_per_page,omitempty" validate:"omitempty,min=1,max=100"`
	Sort         *string    `form:"sort" json:"sort,omitempty" validate:"omitempty,oneof=name_asc name_desc updated_at_asc updated_at_desc"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// for  /forms/{formId}/structure
type PublicFormStructureResponse struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	Description *string               `json:"description,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Sections    []FormSectionResponse `json:"sections"`
}

// for POST /forms/{formId}/submissions
type SubmitPublicFormRequest struct {
	Data map[string]interface{} `json:"data" validate:"required"`
}

type SubmissionSuccessResponse struct {
	Message      string    `json:"message"`
	SubmissionID uuid.UUID `json:"submission_id"`
}

type SubmissionListFilter struct {
	FormID         *uuid.UUID `form:"form_id" json:"form_id,omitempty" validate:"omitempty,uuid"`
	SubmitterEmail *string    `form:"submitter_email" json:"submitter_email,omitempty" validate:"omitempty,email"`
	SubmittedAt    *time.Time `form:"submitted_at" json:"submitted_at,omitempty" layout:"2006-01-02"`
	Page           *int       `form:"page" json:"page,omitempty" validate:"omitempty,min=1"`
	ItemsPerPage   *int       `form:"items_per_page" json:"items_per_page,omitempty" validate:"omitempty,min=1,max=100"`
	SortBy         *string    `form:"sort_by" json:"sort_by,omitempty" validate:"omitempty,oneof=submitted_at status form_name submitter_email"`
	SortOrder      *string    `form:"sort_order" json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

type FormSubmissionAdminView struct {
	ID               uuid.UUID            `json:"id"`
	FormID           uuid.UUID            `json:"form_id"`
	FormName         string               `json:"form_name"`
	SubmittedAt      time.Time            `json:"submitted_at"`
	SubmitterEmail   *string              `json:"submitter_email,omitempty"`
	DecoratedData    []DecoratedFieldData `json:"decorated_data"`
	RawSubmittedData datatypes.JSON       `json:"raw_submitted_data,omitempty"`
}

type DecoratedFieldData struct {
	Label     string      `json:"label"`
	Value     interface{} `json:"value"`
	FieldKey  string      `json:"field_key"`
	FieldType string      `json:"field_type"`
}

// Response ของ GET /cms/forms/{formId}/submissions
type PaginatedSubmissionResponse struct {
	Data []FormSubmissionAdminView `json:"data"`
	Meta PaginationMeta            `json:"meta"`
}
