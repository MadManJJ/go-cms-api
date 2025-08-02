package dto

import (
	"time"

	"gorm.io/datatypes"
)

// --- MetaTag ---
type CreateMetaTagRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

// --- Component ---
type CreateComponentRequest struct {
	Type  string                 `json:"type"`
	Props map[string]interface{} `json:"props"`
}

// --- Revision ---
type CreateRevisionRequest struct {
	PublishStatus string `json:"publish_status"`
	Author        string `json:"author"`
	Message       string `json:"message"`
	Description   string `json:"description"`
}

// --- Category ---
type CreateCategoryRequest struct {
	LanguageCode   string  `json:"language_code"`
	CategoryTypeId string  `json:"category_type_id"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	Weight         int     `json:"weight"`
	PublishStatus  string  `json:"publish_status"`
}

type CreateCategoryId struct {
	ID  string  `json:"id"`
}

// --- LandingContentFile ---
type CreateLandingContentFileRequest struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

// --- LandingContent ---
type CreateLandingContentRequest struct {
	Title          string                            `json:"title"`
	Language       string                            `json:"language"`
	AuthoredAt     time.Time                         `json:"authored_at"`
	HTMLInput      string                            `json:"html_input"`
	Mode           string                            `json:"mode"`
	WorkflowStatus string                            `json:"workflow_status"`
	PublishStatus  string                            `json:"publish_status"`
	ApprovalEmail  []string                          `json:"approval_email"`
	MetaTag        CreateMetaTagRequest              `json:"meta_tag"`
	Files          []CreateLandingContentFileRequest `json:"files"`
	Revision       CreateRevisionRequest             `json:"revision"`
	Categories     []CreateCategoryId                `json:"categories"`
	Components     []CreateComponentRequest          `json:"components"`
	URLAlias       string                            `json:"url_alias"`
	PublishOn      time.Time                         `json:"publish_on"`
	UnpublishOn    time.Time                         `json:"unpublish_on"`
	AuthoredOn     time.Time                         `json:"authored_on"`
}

type CreateLandingContentPreviewRequest struct {
	Title          string                            `json:"title"`
	Language       string                            `json:"language"`
	AuthoredAt     time.Time                         `json:"authored_at"`
	HTMLInput      string                            `json:"html_input"`
	MetaTag        CreateMetaTagRequest              `json:"meta_tag"`
	Files          []CreateLandingContentFileRequest `json:"files"`
	Components     []CreateComponentRequest          `json:"components"`
	URLAlias       string                            `json:"url_alias"`
	PublishOn      time.Time                         `json:"publish_on"`
	UnpublishOn    time.Time                         `json:"unpublish_on"`
	AuthoredOn     time.Time                         `json:"authored_on"`
}

// --- LandingPage ---
type CreateLandingPageRequest struct {
	Contents []CreateLandingContentRequest `json:"contents"`
}

// --- PartnerContent ---
type CreatePartnerContentRequest struct {
	Title            string                   `json:"title"`
	ThumbnailImage   string                   `json:"thumbnail_image"`
	ThumbnailAltText string                   `json:"thumbnail_alt_text"`
	CompanyLogo      string                   `json:"company_logo"`
	CompanyAltText   string                   `json:"company_alt_text"`
	Language         string                   `json:"language"`
	CompanyName      string                   `json:"company_name"`
	CompanyDetail    string                   `json:"company_detail"`
	LeadBody         string                   `json:"lead_body"`
	Challenges       string                   `json:"challenges"`
	Solutions        string                   `json:"solutions"`
	Results          string                   `json:"results"`
	AuthoredAt       time.Time                `json:"authored_at"`
	HTMLInput        string                   `json:"html_input"`
	Mode             string                   `json:"mode"`
	WorkflowStatus   string                   `json:"workflow_status"`
	PublishStatus    string                   `json:"publish_status"`
	ApprovalEmail    []string                 `json:"approval_email"`
	MetaTag          CreateMetaTagRequest     `json:"meta_tag"`
	Revision         CreateRevisionRequest    `json:"revision"`
	Categories       []CreateCategoryId       `json:"categories"`
	Components       []CreateComponentRequest `json:"components"`
	PublishOn        time.Time                `json:"publish_on"`
	UnpublishOn      time.Time                `json:"unpublish_on"`
	AuthoredOn       time.Time                `json:"authored_on"`
	URLAlias         string                   `json:"url_alias"`
	URL              string                   `json:"url"`
}

type CreatePartnerContentPreviewRequest struct {
	Title            string                   `json:"title"`
	ThumbnailImage   string                   `json:"thumbnail_image"`
	ThumbnailAltText string                   `json:"thumbnail_alt_text"`
	CompanyLogo      string                   `json:"company_logo"`
	CompanyAltText   string                   `json:"company_alt_text"`
	Language         string                   `json:"language"`
	CompanyName      string                   `json:"company_name"`
	CompanyDetail    string                   `json:"company_detail"`
	LeadBody         string                   `json:"lead_body"`
	Challenges       string                   `json:"challenges"`
	Solutions        string                   `json:"solutions"`
	Results          string                   `json:"results"`
	AuthoredAt       time.Time                `json:"authored_at"`
	HTMLInput        string                   `json:"html_input"`
	MetaTag          CreateMetaTagRequest     `json:"meta_tag"`
	Components       []CreateComponentRequest `json:"components"`
	PublishOn        time.Time                `json:"publish_on"`
	UnpublishOn      time.Time                `json:"unpublish_on"`
	AuthoredOn       time.Time                `json:"authored_on"`
	URLAlias         string                   `json:"url_alias"`
	URL              string                   `json:"url"`	
}

// --- PartnerPage ---
type CreatePartnerPageRequest struct {
	Contents       []CreatePartnerContentRequest `json:"contents"`
}

// --- FAQ Content ---
type CreateFaqContentRequest struct {
	Title          string                   `json:"title"`
	Language       string                   `json:"language"`
	AuthoredAt     time.Time                `json:"authored_at"`
	HTMLInput      string                   `json:"html_input"`
	Mode           string                   `json:"mode"`
	WorkflowStatus string                   `json:"workflow_status"`
	PublishStatus  string                   `json:"publish_status"`
	MetaTag        CreateMetaTagRequest     `json:"meta_tag"`
	Revision       CreateRevisionRequest    `json:"revision"`
	Categories     []CreateCategoryId       `json:"categories"`
	Components     []CreateComponentRequest `json:"components"`
	PublishOn      time.Time                `json:"publish_on"`
	UnpublishOn    time.Time                `json:"unpublish_on"`
	AuthoredOn     time.Time                `json:"authored_on"`
	URLAlias       string                   `json:"url_alias"`
	URL            string                   `json:"url"`
}

type CreateFaqContentPreviewRequest struct {
	Title          string                   `json:"title"`
	Language       string                   `json:"language"`
	AuthoredAt     time.Time                `json:"authored_at"`
	HTMLInput      string                   `json:"html_input"`
	MetaTag        CreateMetaTagRequest     `json:"meta_tag"`
	Components     []CreateComponentRequest `json:"components"`
	PublishOn      time.Time                `json:"publish_on"`
	UnpublishOn    time.Time                `json:"unpublish_on"`
	AuthoredOn     time.Time                `json:"authored_on"`
	URLAlias       string                   `json:"url_alias"`
	URL            string                   `json:"url"`
}

// --- FAQ Page ---
type CreateFaqPageRequest struct {
	Contents []CreateFaqContentRequest `json:"contents"`
}

type CreateFormSubmissionRequest struct {
	SubmittedData   datatypes.JSON `gorm:"type:jsonb;not null" json:"submitted_data" swaggertype:"object,string"`
}