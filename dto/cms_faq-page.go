package dto

import (
	"time"
)

type FaqPageQuery struct {
	Title            string `form:"title" json:"title"`
	CategoryFaq      string `form:"category_faq" json:"category_faq"`
	CategoryKeywords string `form:"category_keywords" json:"category_keywords"`
	Status           string `form:"status" json:"status"`
	UrlAlias         string `form:"url_alias" json:"url_alias"`
	URL              string `form:"url" json:"url"`
}

type CMSFaqPagesSuccessResponse200 struct {
	Message    string            `json:"message" example:"Faq page retrieved successfully"`
	TotalCount int               `json:"totalCount" example:"100"`
	Page       int               `json:"page" example:"1"`
	Limit      int               `json:"limit" example:"10"`
	Items      []FaqPageResponse `json:"items"`
}

type CMSFaqPageSuccessResponse200 struct {
	Message string          `json:"message" example:"Faq page retrieved successfully"`
	Items   FaqPageResponse `json:"items"`
}

type CMSFaqContentSuccessResponse200 struct {
	Message string             `json:"message" example:"Faq content retrieved successfully"`
	Item    FaqContentResponse `json:"item"`
}

type CMSFaqCategoriesSuccessResponse200 struct {
	Message string             `json:"message" example:"Category retrieved successfully"`
	Item    []CategoryResponse `json:"item"`
}

type CMSFaqRevsionsSuccessResponse200 struct {
	Message string             `json:"message" example:"Revisions retrieved successfully"`
	Item    []RevisionResponse `json:"item"`
}

type UpdatedMetaTag struct {
	ID          string `json:"id" example:"uuid" swaggerignore:"true"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

type UpdatedFaqContent struct {
	ID             string    `json:"id" swaggerignore:"true"`
	PageID         string    `json:"page_id" swaggerignore:"true"`
	Title          string    `json:"title"`
	Language       string    `json:"language"`
	AuthoredAt     time.Time `json:"authored_at"`
	HTMLInput      string    `json:"html_input"`
	Mode           string    `json:"mode"`
	WorkflowStatus string    `json:"workflow_status"`
	PublishStatus  string    `json:"publish_status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// Revisions, Categories, Components omitted
}

type UpdatedFaqPage struct {
	ID          string            `json:"id" swaggerignore:"true"`
	URLAlias    string            `json:"url_alias"`
	URL         string            `json:"url"`
	MetaTag     UpdatedMetaTag    `json:"meta_tag"`
	PublishOn   time.Time         `json:"publish_on"`
	UnpublishOn time.Time         `json:"unpublish_on"`
	AuthoredOn  time.Time         `json:"authored_on"`
	Content     UpdatedFaqContent `json:"content"`
}
