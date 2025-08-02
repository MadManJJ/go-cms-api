package dto

import (
	"time"
)

type PartnerPageQuery struct {
	Title             string `form:"title" json:"title"`
	CategoryPartner   string `form:"category_Partner" json:"category_Partner"`
	CategoryKeywords  string `form:"category_keywords" json:"category_keywords"`
	CategoryScale     string `form:"category_scale" json:"category_scale"`
	CategoryIndustry  string `form:"category_industry" json:"category_industry"`
	CategoryGoal      string `form:"category_goal" json:"category_goal"`
	CategoryFunctions string `form:"category_functions" json:"category_functions"`
	UrlAlias          string `form:"url_alias" json:"url_alias"`
	URL               string `form:"url" json:"url"`

	Status string `form:"status" json:"status"`
}

type CMSPartnerPageSuccessResponse200 struct {
	Message    string                `json:"message" example:"Partner page retrieved successfully"`
	TotalCount int                   `json:"totalCount" example:"100"`
	Page       int                   `json:"page" example:"1"`
	Limit      int                   `json:"limit" example:"10"`
	Items      []PartnerPageResponse `json:"items"`
}

type CMSPartnerContentSuccessResponse200 struct {
	Message string                 `json:"message" example:"Partner content retrieved successfully"`
	Item    PartnerContentResponse `json:"item"`
}

type CMSPartnerCategoriesSuccessResponse200 struct {
	Message string             `json:"message" example:"Category retrieved successfully"`
	Item    []CategoryResponse `json:"item"`
}

type UpdatedPartnerContent struct {
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

type UpdatedPartnerPage struct {
	ID          string                `json:"id" swaggerignore:"true"`
	URLAlias    string                `json:"url_alias"`
	URL         string                `json:"url"`
	MetaTag     UpdatedMetaTag        `json:"meta_tag"`
	PublishOn   time.Time             `json:"publish_on"`
	UnpublishOn time.Time             `json:"unpublish_on"`
	AuthoredOn  time.Time             `json:"authored_on"`
	Content     UpdatedPartnerContent `json:"content"`
}

type CMSPartnerRevsionsSuccessResponse200 struct {
	Message string             `json:"message" example:"Revisions retrieved successfully"`
	Item    []RevisionResponse `json:"item"`
}
