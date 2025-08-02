package common

import (
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EmailSendingHandler struct {
	Service  services.EmailSendingServiceInterface
	validate *validator.Validate
}

func NewEmailSendingHandler(service services.EmailSendingServiceInterface) *EmailSendingHandler {
	return &EmailSendingHandler{
		Service:  service,
		validate: validator.New(),
	}
}

// SendEmail sends an email based on a template and provided data.
// @Summary Send Email
// @Description Sends an email using a pre-configured template identified by category, label, and language.
// @Tags Email Sending
// @Accept json
// @Produce json
// @Param emailRequest body dto.SendEmailRequest true "Email Sending Request"
// @Success 202 {object} map[string]string "Accepted - Email queued for sending"
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 404 {object} dto.ErrorResponse "Email category or content template not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error or Email Sending Failed"
// @Router /emails/send [post]
func (h *EmailSendingHandler) HandleSendEmail(c *fiber.Ctx) error {
	var req dto.SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}

	err := h.Service.SendEmail(req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "email sending is not configured") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(dto.ErrorResponse{Error: "Service Unavailable", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to send email", Message: err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Email queued for sending"})
}
