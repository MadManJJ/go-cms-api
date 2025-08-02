package dto

type CMSFormSubmissionSuccessResponse200 struct {
	Message string                 `json:"message" example:"Form submission retrieved successfully"`
	Item    FormSubmissionResponse `json:"item"`
}

type CMSFormSubmissionsSuccessResponse200 struct {
	Message string                   `json:"message" example:"Form submission retrieved successfully"`
	Item    []FormSubmissionResponse `json:"item"`
}