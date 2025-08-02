package cms

import (
	_ "github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
)

// CMSAuthHandler handles HTTP requests related to CMS authentication.
type CMSAuthHandler struct {
	Service services.CMSAuthServiceInterface
}

// NewCMSAuthHandler creates a new instance of CMSAuthHandler
func NewAuthCMSHandler(service services.CMSAuthServiceInterface) *CMSAuthHandler {
	return &CMSAuthHandler{
		Service: service,
	}
}

// Handle register
// @Summary      Register
// @Description  Registers a new user
// @Tags         CMS - Email Login
// @Accept       json
// @Produce      json
// @Param 			 User body dto.AuthUserRequest true "User"
// @Success      200  {object} dto.RegisterResponse
// @Failure 		 400  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse 
// @Router       /cms/auth/register [post]
func (h *CMSAuthHandler) HandleRegister(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse user data",
			"error": err.Error(),
		})
	}

	createdUser, err := h.Service.RegisterUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register user",
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Register successful",
		"user": createdUser,
	})
}

// Handle login
// @Summary      Login
// @Description  Logs in a user and returns a JWT token
// @Tags         CMS - Email Login
// @Accept       json
// @Produce      json
// @Param 			 User body dto.AuthUserRequest true "User"
// @Success      200  {object} dto.LoginResponse
// @Failure 		 401  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse 
// @Router       /cms/auth/login [post]
func (h *CMSAuthHandler) HandleLogin(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse user data",
			"error": err.Error(),
		})
	}

	// Validate user input
	authenticatedUser, token, err := h.Service.LoginUser(&user)
	if err != nil {
		switch err {
		case errs.ErrInvalidCredentials: // This is a custom error defined in errs/errs.go
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid email or password",
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
				"error": err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": authenticatedUser,
	})
}
