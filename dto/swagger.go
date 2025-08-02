package dto

import (
	"time"

	enum "github.com/MadManJJ/cms-api/models/enums"
)

type ComponentSwagger struct {
	ID        string      `json:"id"`
	ContentID string      `json:"contentId"`
	Type      string      `json:"type"`
	Props     interface{} `json:"props"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type ContentSwagger struct {
	ID             string               `json:"id"`
	PageID         string               `json:"pageId"`
	Language       enum.PageLanguage    `json:"language,omitempty"`
	HtmlInput      string               `json:"htmlInput,omitempty"`
	Mode           enum.PageMode        `json:"mode,omitempty"`
	WorkflowStatus enum.WorkflowStatus  `json:"workflowStatus,omitempty"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	Files          []*FileSwagger       `json:"files,omitempty"`
	Revisions      []*RevisionSwagger   `json:"revisions,omitempty"`
	Components     []*ComponentSwagger  `json:"components,omitempty"`
}

type PageSwagger struct {
	ID          string            `json:"id"`
	UrlAlias    string            `json:"urlAlias"`
	PublishOn   *time.Time        `json:"publishOn,omitempty"`
	UnpublishOn *time.Time        `json:"unpublishOn,omitempty"`
	AuthoredOn  *time.Time        `json:"authoredOn,omitempty"`
	MetaTag     *MetaTagSwagger          `json:"metaTag,omitempty"`
	Content     *ContentSwagger   `json:"content,omitempty"`
}

type MetaTagSwagger struct {
	PageID      string `json:"pageId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
}

type ContentResponse struct {
	ID             string              `json:"id"`
	PageID         string              `json:"pageId"`
	Page           *PageSwagger        `json:"page,omitempty"`
	Language       enum.PageLanguage   `json:"language,omitempty"`
	HtmlInput      string              `json:"htmlInput,omitempty"`
	Mode           enum.PageMode       `json:"mode,omitempty"`
	WorkflowStatus enum.WorkflowStatus `json:"workflowStatus,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`

	Files      []*FileSwagger      `json:"files,omitempty"`
	Revisions  []*RevisionSwagger  `json:"revisions,omitempty"`
	Components []*ComponentSwagger `json:"components,omitempty"`
}

type FileSwagger struct {
	ID          string    `json:"id"`
	ContentID   string    `json:"contentId"`
	Position    int       `json:"position"`
	FileType    string    `json:"fileType"`
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	DownloadURL string    `json:"downloadUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RevisionSwagger struct {
	ID            string    `json:"id"`
	ContentID     string    `json:"contentId"`
	PublishStatus string    `json:"publishStatus"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Author        string    `json:"author"`
	Message       string    `json:"message"`
	Description   string    `json:"description"`
}
