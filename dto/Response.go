package dto

import "time"

type MetaTagResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

type ComponentResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Props     map[string]interface{} `json:"props"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type RevisionResponse struct {
	ID               string    `json:"id"`
	LandingContentID *string   `json:"landing_content_id,omitempty"`
	PartnerContentID *string   `json:"partner_content_id,omitempty"`
	FaqContentID     *string   `json:"faq_content_id,omitempty"`
	PublishStatus    string    `json:"publish_status"`
	UpdatedAt        time.Time `json:"updated_at"`
	Author           string    `json:"author"`
	Message          string    `json:"message"`
	Description      string    `json:"description"`
}

type LandingContentFileResponse struct {
	ID               string    `json:"id"`
	LandingContentID string    `json:"landing_content_id"`
	Name             string    `json:"name"`
	DownloadURL      string    `json:"download_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type LandingContentResponse struct {
	ID             string                       `json:"id"`
	PageID         string                       `json:"page_id"`
	Title          string                       `json:"title"`
	Language       string                       `json:"language"`
	AuthoredAt     time.Time                    `json:"authored_at"`
	HTMLInput      string                       `json:"html_input"`
	Mode           string                       `json:"mode"`
	WorkflowStatus string                       `json:"workflow_status"`
	PublishStatus  string                       `json:"publish_status"`
	CreatedAt      time.Time                    `json:"created_at"`
	UpdatedAt      time.Time                    `json:"updated_at"`
	MetaTag        MetaTagResponse              `json:"meta_tag"`
	ApprovalEmail  []string                     `json:"approval_email"`
	Files          []LandingContentFileResponse `json:"files"`
	Revision       RevisionResponse             `json:"revision"`
	Categories     []CategoryResponse           `json:"categories"`
	Components     []ComponentResponse          `json:"components"`
	URLAlias       string                       `json:"url_alias"`
	PublishOn      time.Time                    `json:"publish_on"`
	UnpublishOn    time.Time                    `json:"unpublish_on"`
	AuthoredOn     time.Time                    `json:"authored_on"`
}

type LandingPageResponse struct {
	ID        string                   `json:"id"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	Contents  []LandingContentResponse `json:"contents"`
}

type PartnerContentResponse struct {
	ID               string              `json:"id"`
	PageID           string              `json:"page_id"`
	Title            string              `json:"title"`
	ThumbnailImage   string              `json:"thumbnail_image"`
	ThumbnailAltText string              `json:"thumbnail_alt_text"`
	CompanyLogo      string              `json:"company_logo"`
	CompanyAltText   string              `json:"company_alt_text"`
	Language         string              `json:"language"`
	CompanyName      string              `json:"company_name"`
	CompanyDetail    string              `json:"company_detail"`
	LeadBody         string              `json:"lead_body"`
	Challenges       string              `json:"challenges"`
	Solutions        string              `json:"solutions"`
	Results          string              `json:"results"`
	AuthoredAt       time.Time           `json:"authored_at"`
	HTMLInput        string              `json:"html_input"`
	Mode             string              `json:"mode"`
	WorkflowStatus   string              `json:"workflow_status"`
	PublishStatus    string              `json:"publish_status"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	MetaTag          MetaTagResponse     `json:"meta_tag"`
	ApprovalEmail    []string            `json:"approval_email"`
	Revision         RevisionResponse    `json:"revision"`
	Categories       []CategoryResponse  `json:"categories"`
	Components       []ComponentResponse `json:"components"`
	PublishOn        time.Time           `json:"publish_on"`
	UnpublishOn      time.Time           `json:"unpublish_on"`
	AuthoredOn       time.Time           `json:"authored_on"`
	URLAlias         string              `json:"url_alias"`
	URL              string              `json:"url"`
	IsRecommended    bool                `json:"is_recommended"`
}

type PartnerPageResponse struct {
	ID        string                   `json:"id"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	Contents  []PartnerContentResponse `json:"contents"`
}

type FaqContentResponse struct {
	ID             string              `json:"id"`
	PageID         string              `json:"page_id"`
	Title          string              `json:"title"`
	Language       string              `json:"language"`
	AuthoredAt     time.Time           `json:"authored_at"`
	HTMLInput      string              `json:"html_input"`
	Mode           string              `json:"mode"`
	WorkflowStatus string              `json:"workflow_status"`
	PublishStatus  string              `json:"publish_status"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	MetaTag        MetaTagResponse     `json:"meta_tag"`
	Revision       RevisionResponse    `json:"revision"`
	Categories     []CategoryResponse  `json:"categories"`
	Components     []ComponentResponse `json:"components"`
	PublishOn      time.Time           `json:"publish_on"`
	UnpublishOn    time.Time           `json:"unpublish_on"`
	AuthoredOn     time.Time           `json:"authored_on"`
	URLAlias       string              `json:"url_alias"`
	URL            string              `json:"url"`
}

type FaqPageResponse struct {
	ID        string               `json:"id"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Contents  []FaqContentResponse `json:"contents"`
}

type FormSubmissionResponse struct {
	ID             string                 `json:"id"`
	FormID         string                 `json:"form_id"`
	SubmittedData  map[string]interface{} `json:"submitted_data"` // use a map to reflect JSONB object structure
	SubmittedAt    time.Time              `json:"submitted_at"`
	SubmittedEmail *string                `json:"submitted_email,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Form           *FormResponse          `json:"form,omitempty"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
