package dto

import "time"

type CreateEmailCategoryRequest struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
}

type UpdateEmailCategoryRequest struct {
	Title string `json:"title" validate:"omitempty,min=3,max=255"`
}

type EmailCategoryResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
