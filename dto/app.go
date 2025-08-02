package dto

// AppData represents data from the app domain
type AppData struct {
	ID      string `json:"id" example:"123"`
	Message string `json:"message" example:"Hello, World!"`
}

type SuccessResponseLandingPage struct {
	Message string `json:"message" example:"Landing page retrieved successfully"`
}

type LandingPageSuccessResponse200 struct {
	Message string              `json:"message" example:"Landing page retrieved successfully"`
	Data    LandingPageResponse `json:"data"`
}

type LandingContentSuccessResponse200 struct {
	Message string                 `json:"message" example:"Landing content retrieved successfully"`
	Data    LandingContentResponse `json:"data"`
}

type PartnerPageSuccessResponse200 struct {
	Message string              `json:"message" example:"Partner page retrieved successfully"`
	Data    PartnerPageResponse `json:"data"`
}

type PartnerContentSuccessResponse200 struct {
	Message string                 `json:"message" example:"Partner page retrieved successfully"`
	Data    PartnerContentResponse `json:"data"`
}

type FaqPageSuccessResponse200 struct {
	Message string          `json:"message" example:"Faq page retrieved successfully"`
	Data    FaqPageResponse `json:"data"`
}

type FaqContentSuccessResponse200 struct {
	Message string             `json:"message" example:"Faq content retrieved successfully"`
	Data    FaqContentResponse `json:"data"`
}

type ErrorResponse400 struct {
	Message string `json:"message" example:"Bad Request"`
	Error   string `json:"error" example:"Bad Request"`
}

type ErrorResponse404 struct {
	Message string `json:"message" example:"Page not found"`
	Error   string `json:"error" example:"Not Found"`
}

type ErrorResponse500 struct {
	Message string `json:"message" example:"An unexpected error occurred on the server"`
	Error   string `json:"error" example:"Internal Server Error"`
}
