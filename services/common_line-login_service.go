package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"

	"gorm.io/gorm"
)

type LineLoginServiceInterface interface {
	GetLoginLink() (string, error)
	Authenticate(autorizationCode string) (*dto.LineTokenResponse, *models.User, error)
	RefreshToken(refreshToken string) (*dto.LineTokenResponse, error)
}

type LineLoginService struct{
	cfg *config.Config
	repo repositories.CMSAuthRepositoryInterface
}

func NewLineLoginService(cfg *config.Config, repo repositories.CMSAuthRepositoryInterface) LineLoginServiceInterface {
	return &LineLoginService{cfg: cfg, repo: repo}
}

func (s *LineLoginService) GetLoginLink() (string, error) {
	clientId, redirectUri := s.cfg.Line.ClientId, s.cfg.Line.RedirectUri

	if clientId == "" || redirectUri == "" {
		return "", errs.ErrNotFound
	}

	state := helpers.RandomState()

	baseURL := s.cfg.Line.AuthorizeUrl

	// Construct the query parameters
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientId)
	params.Add("redirect_uri", redirectUri)
	params.Add("scope", "profile openid")
	params.Add("state", state)	

	// Construct the final login link
	loginLink := fmt.Sprintf("%s?%s", baseURL, params.Encode())	
	fmt.Println(loginLink)

	return loginLink, nil
}

func (s *LineLoginService) Authenticate(autorizationCode string) (*dto.LineTokenResponse, *models.User, error) {
	if autorizationCode == "" {
		return nil, nil, errs.ErrBadRequest
	}

	clientId, redirectUri, clientSecret, tokenUrl := s.cfg.Line.ClientId, s.cfg.Line.RedirectUri, s.cfg.Line.ClientSecret, s.cfg.Line.TokenUrl

	if tokenUrl == "" || clientId == "" || redirectUri == "" || clientSecret == "" {
		return nil, nil, errs.ErrInternalServerError
	}

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", autorizationCode)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectUri)

	// Create HTTP request
	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Failed to create request: ", err)
		return nil, nil, errs.ErrInternalServerError
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed: ", err)
		return nil, nil, errs.ErrInternalServerError
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body: ", err)
		return nil, nil, errs.ErrInternalServerError
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Response status: ", resp.StatusCode)
		fmt.Println("Response body: ", string(body)) // Add this to see the error details
		return nil, nil, errs.ErrInternalServerError
	}

	// Parse JSON response
	var tokenResponse dto.LineTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Failed to parse JSON response: ", err)
		fmt.Println("Response body: ", string(body)) // Add this to see what you're trying to parse
		return nil, nil, errs.ErrInternalServerError
	}

	// Get user info from ID token
	idToken := tokenResponse.IDToken
	claims, err := helpers.ParseJWTWithKey(idToken, s.cfg.SecretKey.LineKey)
	if err != nil {
		return nil, nil, errs.ErrInvalidCredentials
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, nil, errors.New("sub claim missing or invalid")
	}
	
	id := helpers.UUIDFromSub(sub)
	findUser, err := s.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// user not found, register new user
			user := &models.User{
				ID: id,
				Provider: enums.ProviderLine,
			}			
			createdUser, err := s.repo.RegisterUser(user)
			if err != nil {
				return nil, nil, err
			}

			return &tokenResponse, createdUser, nil
		} else {
			// some other DB error occurred
			return nil, nil, err
		}
	} 

	// user found
	return &tokenResponse, findUser, nil
}

func (s *LineLoginService) RefreshToken(refreshToken string) (*dto.LineTokenResponse, error) {
	if refreshToken == "" {
		return nil, errs.ErrBadRequest
	}
	clientId, clientSecret, tokenUrl := s.cfg.Line.ClientId, s.cfg.Line.ClientSecret, s.cfg.Line.TokenUrl

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)	

	// Create HTTP request
	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Failed to create request: ", err)
		return nil, errs.ErrInternalServerError
	}	

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed: ", err)
		return nil, errs.ErrInternalServerError
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body: ", err)
		return nil, errs.ErrInternalServerError
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Response status: ", resp.StatusCode)
		fmt.Println("Response body: ", string(body))
		return nil, errs.ErrInternalServerError
	}

	// Parse JSON response
	var tokenResponse dto.LineTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Failed to parse JSON response: ", err)
		fmt.Println("Response body: ", string(body))
		return nil, errs.ErrInternalServerError
	}

	return &tokenResponse, nil	
}