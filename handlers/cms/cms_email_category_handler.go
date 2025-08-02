package cms

import (
	"errors"
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EmailCategoryHandler struct {
	Service  services.EmailCategoryServiceInterface
	validate *validator.Validate
}

func NewEmailCategoryHandler(service services.EmailCategoryServiceInterface) *EmailCategoryHandler {
	return &EmailCategoryHandler{
		Service:  service,
		validate: validator.New(),
	}
}

// CreateEmailCategory creates a new email category.
// @Summary Create Email Category
// @Description Creates a new email category.
// @Tags CMS - Email Management
// @Accept json
// @Produce json
// @Param category body dto.CreateEmailCategoryRequest true "Email Category Data"
// @Success 201 {object} dto.EmailCategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Title already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-categories [post]
func (h *EmailCategoryHandler) HandleCreateEmailCategory(c *fiber.Ctx) error {
	var req dto.CreateEmailCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}

	category, err := h.Service.CreateCategory(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to create email category", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(category)
}

// ListEmailCategories lists all email categories.
// @Summary List Email Categories
// @Description Retrieves a list of all email categories.
// @Tags CMS - Email Management
// @Produce json
// @Success 200 {array} dto.EmailCategoryResponse
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-categories [get]
func (h *EmailCategoryHandler) HandleListEmailCategories(c *fiber.Ctx) error {
	categories, err := h.Service.ListCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to list email categories", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(categories)
}

// GetEmailCategory retrieves a specific email category by ID.
// @Summary Get Email Category
// @Description Retrieves a specific email category by its UUID.
// @Tags CMS - Email Management
// @Produce json
// @Param id path string true "Email Category ID (UUID)"
// @Success 200 {object} dto.EmailCategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid ID format"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-categories/{id} [get]
func (h *EmailCategoryHandler) HandleGetEmailCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	category, err := h.Service.GetCategoryByID(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Email category not found"})
		}
		if strings.Contains(err.Error(), "invalid category ID format") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get email category", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(category)
}

// UpdateEmailCategory updates an existing email category.
// @Summary Update Email Category
// @Description Updates an existing email category by its UUID.
// @Tags CMS - Email Management
// @Accept json
// @Produce json
// @Param id path string true "Email Category ID (UUID)"
// @Param category body dto.UpdateEmailCategoryRequest true "Email Category Update Data"
// @Success 200 {object} dto.EmailCategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Title already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-categories/{id} [patch]
func (h *EmailCategoryHandler) HandleUpdateEmailCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var req dto.UpdateEmailCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}

	category, err := h.Service.UpdateCategory(idStr, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Email category not found for update"})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to update email category", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(category)
}

// DeleteEmailCategory deletes an email category.
// @Summary Delete Email Category
// @Description Deletes an email category by its UUID, if not in use.
// @Tags CMS - Email Management
// @Param id path string true "Email Category ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID format"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Category in use"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/email-categories/{id} [delete]
func (h *EmailCategoryHandler) HandleDeleteEmailCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	err := h.Service.DeleteCategory(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Email category not found for delete"})
		}
		if strings.Contains(err.Error(), "in use") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to delete email category", Message: err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
