package dto

type LineTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthenticateDto struct {
	AutorizationCode string `json:"authorization_code"`
}

type RefreshTokenDto struct {
	RefreshToken string `json:"refresh_token"`
}