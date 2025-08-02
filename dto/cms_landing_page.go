package dto

import (
	"time"

	"github.com/lib/pq"
)

// --- Request DTOs for Landing Page Creation & Initial Content ---

type LandingPageQuery struct {
	Title            string `form:"title" json:"title"`
	CategoryKeywords string `form:"category_keywords" json:"category_keywords"`
	Status           string `form:"status" json:"status"`
	UrlAlias         string `form:"url_alias" json:"url_alias"`
}

type CMSLandingPageSuccessResponse200 struct {
	Message    string            `json:"message" example:"Faq page retrieved successfully"`
	TotalCount int               `json:"totalCount" example:"100"`
	Page       int               `json:"page" example:"1"`
	Limit      int               `json:"limit" example:"10"`
	Items      []FaqPageResponse `json:"items"`
}
type CMSLandingContentSuccessResponse200 struct {
	Message string             `json:"message" example:"Faq content retrieved successfully"`
	Item    FaqContentResponse `json:"item"`
}

type CMSLandingCategoriesSuccessResponse200 struct {
	Message string             `json:"message" example:"Category retrieved successfully"`
	Item    []CategoryResponse `json:"item"`
}

type UpdatedLandingContent struct {
	ID             string                `json:"id" swaggerignore:"true"`
	PageID         string                `json:"page_id" swaggerignore:"true"`
	Title          string                `json:"title"`
	Language       string                `json:"language"`
	AuthoredAt     time.Time             `json:"authored_at"`
	HTMLInput      string                `json:"html_input"`
	Mode           string                `json:"mode"`
	WorkflowStatus string                `json:"workflow_status"`
	PublishStatus  string                `json:"publish_status"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	PublishOn      time.Time             `json:"publish_on"`
	UnpublishOn    time.Time             `json:"unpublish_on"`
	AuthoredOn     time.Time             `json:"authored_on"`
	URLAlias       string                `json:"url_alias"`
	URL            string                `json:"url"`
	ApprovalEmail  pq.StringArray        `gorm:"type:text[]" json:"approval_email"`
	Files          []*LandingContentFile `gorm:"foreignKey:LandingContentID" json:"files,omitempty"`
	// Revisions, Categories, Components omitted
}
type LandingContentFile struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}
type UpdatedLandingPage struct {
	ID      string            `json:"id" swaggerignore:"true"`
	MetaTag UpdatedMetaTag    `json:"meta_tag"`
	Content UpdatedFaqContent `json:"content"`
}

type RevsionsSuccessResponse200 struct {
	Message string             `json:"message" example:"Revisions retrieved successfully"`
	Item    []RevisionResponse `json:"item"`
}
