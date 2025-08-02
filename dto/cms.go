package dto

import "github.com/MadManJJ/cms-api/models"

// CMSData represents data from the CMS domain
type CMSData struct {
	ID      string `json:"id" example:"123"`
	Message string `json:"message" example:"Hello, World!"`
}

type AuthUserRequest struct {
	Email    string `json:"email" example:"johnabc@example.com"`
	Password string `json:"password" example:"123456"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"User registered successfully"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Failed to parse user data"`
	Message string `json:"message" example:"Something went wrong"`
}

type LoginResponse struct {
	Message string      `json:"message" example:"Login successful"`
	Token   string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User    models.User `json:"user"`
}

type RegisterResponse struct {
	Message string      `json:"message" example:"Register successful"`
	User    models.User `json:"user"`
}

type CMSErrorResponse struct {
	Error string `json:"error,omitempty" example:"error message here"`
}

type CMSSuccessResponse struct {
	Error string `json:"error,omitempty" example:"success message here"`
}