package dto

import (
	"time"
)

type UploadMediaFileRequest struct {
	Path    *string `form:"path"`
	Replace *bool   `form:"replace"`
}

type MediaFileListFilter struct {
	Search   *string `query:"search"`
	Page     int     `query:"page" validate:"omitempty,min=1"`
	PageSize int     `query:"pageSize" validate:"omitempty,min=1,max=100"`
	SortBy   *string `query:"sortBy" validate:"omitempty,oneof=name created_at"`
	Order    *string `query:"order" validate:"omitempty,oneof=asc desc"`
}

type MediaFileResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MediaFileListItemResponse is a simplified version for list views.
type MediaFileListItemResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type MediaFilesListResponse struct {
	Data     []MediaFileListItemResponse `json:"data"`
	Total    int64                       `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"pageSize"`
}
