package cms

import (
	"errors"
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EmailContentHandler struct {
	Service  services.EmailContentServiceInterface
	validate *validator.Validate
}

func NewEmailContentHandler(service services.EmailContentServiceInterface) *EmailContentHandler {
	// TODO: Register custom validators like 'email_list' and 'alphanumdash'
	// validate.RegisterValidation("email_list", customEmailListValidator)
	// validate.RegisterValidation("alphanumdash", customAlphanumDashValidator)
	return &EmailContentHandler{
		Service:  service,
		validate: validator.New(),
	}
}

// CreateEmailContent creates a new email content for a category.
// @Summary Create Email Content
// @Description Creates a new email content associated with an email category.
// @Tags CMS - Email Management
// @Accept json
// @Produce json
// @Param content body dto.CreateEmailContentRequest true "Email Content Data"
// @Success 201 {object} dto.EmailContentResponse
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 404 {object} dto.ErrorResponse "Email Category not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Content with this category, language, and label already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-contents [post]
func (h *EmailContentHandler) HandleCreateEmailContent(c *fiber.Ctx) error {
	var req dto.CreateEmailContentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}

	content, err := h.Service.CreateContent(req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to create email content", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(content)
}

// ListEmailContents lists email contents, optionally filtered.
// @Summary List Email Contents
// @Description Retrieves a list of email contents, can be filtered by category_id, language, or label.
// @Tags CMS - Email Management
// @Produce json
// @Param email_category_id query string false "Filter by Email Category ID (UUID)"
// @Param language query string false "Filter by language (th, en)" Enums(th, en)
// @Param label query string false "Filter by label (e.g., customer-welcome)"
// @Success 200 {array} dto.EmailContentResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid filter parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-contents [get]
func (h *EmailContentHandler) HandleListEmailContents(c *fiber.Ctx) error {
	var filters dto.EmailContentFilter
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid query parameters", Message: err.Error()})
	}
	if err := h.validate.Struct(filters); err != nil { // Assuming you might add validation tags to EmailContentFilter
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed for filters", Message: err.Error()})
	}

	contents, err := h.Service.ListContents(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to list email contents", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(contents)
}

// GetEmailContent retrieves a specific email content by ID.
// @Summary Get Email Content
// @Description Retrieves a specific email content by its UUID.
// @Tags CMS - Email Management
// @Produce json
// @Param id path string true "Email Content ID (UUID)"
// @Success 200 {object} dto.EmailContentResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid ID format"
// @Failure 404 {object} dto.ErrorResponse "Content not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-contents/{id} [get]
func (h *EmailContentHandler) HandleGetEmailContent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	content, err := h.Service.GetContentByID(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Email content not found"})
		}
		if strings.Contains(err.Error(), "invalid content ID format") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get email content", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(content)
}

// UpdateEmailContent updates an existing email content.
// @Summary Update Email Content
// @Description Updates an existing email content by its UUID.
// @Tags CMS - Email Management
// @Accept json
// @Produce json
// @Param id path string true "Email Content ID (UUID)"
// @Param content body dto.UpdateEmailContentRequest true "Email Content Update Data"
// @Success 200 {object} dto.EmailContentResponse
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 404 {object} dto.ErrorResponse "Content not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Another content with this category, language, and label already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-contents/{id} [patch]
func (h *EmailContentHandler) HandleUpdateEmailContent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var req dto.UpdateEmailContentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}

	content, err := h.Service.UpdateContent(idStr, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Email content not found for update"})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to update email content", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(content)
}

// DeleteEmailContent deletes an email content.
// @Summary Delete Email Content
// @Description Deletes an email content by its UUID.
// @Tags CMS - Email Management
// @Param id path string true "Email Content ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID format"
// @Failure 404 {object} dto.ErrorResponse "Content not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-contents/{id} [delete]
func (h *EmailContentHandler) HandleDeleteEmailContent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	err := h.Service.DeleteContent(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Email content not found for delete"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to delete email content", Message: err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetEmailContentByCategoryAndLanguage retrieves email contents by category ID and language.
// @Summary Get Email Content by Category and Language
// @Description Retrieves email contents for a specific category and language.
// @Tags CMS - Email Management
// @Produce json
// @Param email_category_id path string true "Email Category ID (UUID)"
// @Param language path string true "Language (th, en)"
// @Success 200 {array} dto.EmailContentResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid parameters"
// @Failure 404 {object} dto.ErrorResponse "Category not found or no contents available"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-contents/category/{email_category_id}/language/{language} [get]
func (h *EmailContentHandler) HandleGetEmailContentByCategoryAndLanguage(c *fiber.Ctx) error {
	categoryIDStr := c.Params("email_category_id")
	langStr := c.Params("language")

	language, err := h.prepareLanguageEnum(langStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid Language", Message: err.Error()})
	}

	contents, err := h.Service.GetEmailContentByCategoryIDAndLanguage(categoryIDStr, language)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "No email content found for the given category and language"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get email content", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(contents)
}

func (h *EmailContentHandler) prepareLanguageEnum(languageStr string) (enums.PageLanguage, error) {
	switch strings.ToLower(languageStr) {
	case "th":
		return enums.PageLanguageTH, nil
	case "en":
		return enums.PageLanguageEN, nil
	default:
		return "", errors.New("invalid language parameter")
	}
}
