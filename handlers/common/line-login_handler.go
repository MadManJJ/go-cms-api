package common

import (
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
)

type LineLoginHandler struct {
	Service services.LineLoginServiceInterface
}

func NewLineLoginHandler(service services.LineLoginServiceInterface) *LineLoginHandler {
	return &LineLoginHandler{Service: service}
}

// HandleGetLoginLink handles GET requests to generate a LINE login URL
// @Summary      Generate LINE Login Link
// @Description  Returns a LINE OAuth 2.0 login URL for redirect-based authentication.
// @Tags         Common - Line Login
// @Produce      json
// @Success      200  {object} dto.SuccessLinkResponse
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /common/login-link [get]
func (h *LineLoginHandler) HandleGetLoginLink(c *fiber.Ctx) error {
	loginLink, err := h.Service.GetLoginLink()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get login link",
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get login link",
		"item": loginLink,
	})
}

// HandleAuthenticate handles POST requests to initiate LINE authentication
// @Summary      Initiate LINE Authentication
// @Description  Accepts authentication data and returns a LINE login URL for OAuth 2.0 authentication.
// @Tags         Common - Line Login
// @Accept       json
// @Produce      json
// @Param        authenticate_data  body  dto.AuthenticateDto  true  "Authentication payload"
// @Success      200  {object} dto.SuccessTokenResponse  "Successfully generated login link"
// @Failure      400  {object} dto.ErrorResponse500      "Invalid input"
// @Failure      500  {object} dto.ErrorResponse500      "Internal server error"
// @Router       /common/authenticate [post]
func (h *LineLoginHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authenticateDto dto.AuthenticateDto
	if err := c.BodyParser(&authenticateDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the body",
			"error": err.Error(),
		})		
	}

	tokenResponse, user, err := h.Service.Authenticate(authenticateDto.AutorizationCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to authenticate",
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully authenticate",
		"item": tokenResponse,
		"user": user,
	})
}

// HandleRefreshToken handles POST requests to refresh LINE authentication tokens
// @Summary      Refresh LINE Authentication Token
// @Description  Accepts a refresh token and returns new access and refresh tokens for LINE OAuth 2.0.
// @Tags         Common - Line Login
// @Accept       json
// @Produce      json
// @Param        refresh_toke_data  body  dto.RefreshTokenDto  true  "Refresh token payload"
// @Success      200  {object} dto.SuccessTokenResponse  "Successfully refresh"
// @Failure      400  {object} dto.ErrorResponse500      "Invalid input"
// @Failure      500  {object} dto.ErrorResponse500      "Internal server error"
// @Router       /common/refresh-token [post]
func (h *LineLoginHandler) HandleRefreshToken(c *fiber.Ctx) error {
	var refreshTokenDto dto.RefreshTokenDto
	if err := c.BodyParser(&refreshTokenDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the body",
			"error": err.Error(),
		})		
	}	

	refreshToken, err := h.Service.RefreshToken(refreshTokenDto.RefreshToken)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to refresh the token",
			"error": err.Error(),
		})
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully refresh the token",
		"item": refreshToken,
	})	
}