package dto

type SuccessTokenResponse struct {
	Item    LineTokenResponse `json:"item"`
	Message string            `json:"message" example:"Success"`
}

type SuccessLinkResponse struct {
	Item    string `json:"item"`
	Message string `json:"message" example:"Success"`
}