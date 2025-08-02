// api/handlers/cms/cms_category_handler.go
package cms

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSCategoryHandler struct {
	Service  services.CMSCategoryServiceInterface
	validate *validator.Validate
}

func NewCMSCategoryHandler(service services.CMSCategoryServiceInterface) *CMSCategoryHandler {
	return &CMSCategoryHandler{
		Service:  service,
		validate: validator.New(),
	}
}

// POST /api/v1/cms/categories
// @Summary Create Category (Detail)
// @Description Creates a new category detail item.
// @Tags CMS - Categories
// @Accept json
// @Produce json
// @Param category_data body dto.CategoryCreateRequest true "Category (Detail) Data"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request (e.g., invalid category_type_id)"
// @Failure 404 {object} dto.ErrorResponse "Category Type not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Category with this name, language, and type already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/categories [post]
func (h *CMSCategoryHandler) HandleCreateCategory(c *fiber.Ctx) error {
	var req dto.CategoryCreateRequest // DTO นี้คือสำหรับสร้าง Category (Detail) หนึ่งรายการ
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}

	// การ default weight จะถูกจัดการใน Service Layer
	// ลบ loop นี้ออก:
	// for i := range req.Details {
	// 	if req.Details[i].Weight == nil {
	// 		defaultWeight := 0
	// 		req.Details[i].Weight = &defaultWeight
	// 	}
	// }

	categoryResponse, err := h.Service.CreateCategory(req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") { // e.g., category_type_id not found
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "invalid") { // e.g., invalid category_type_id format
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "already exist") ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to create category", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(categoryResponse)
}

// GET /api/v1/cms/categories
// @Summary List Categories (Details)
// @Description Lists all category details, optionally filtered by category_type_id, language, name, or publish_status.
// @Tags CMS - Categories
// @Produce json
// @Param category_type_id query string false "Filter by CategoryType ID (UUID)"
// @Param lang query string false "Filter by language (th, en)" Enums(th, en)
// @Param name query string false "Filter by name (partial match)"
// @Param publish_status query string false "Filter by publish status (Published, Unpublished)" Enums(Published, Unpublished)
// @Success 200 {array} dto.CategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid filter parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/categories [get]
func (h *CMSCategoryHandler) HandleListAllCategories(c *fiber.Ctx) error {
	var filter dto.CategoryFilter

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid query parameters", Message: err.Error()})
	}

	if err := h.validate.Struct(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed for filter", Message: err.Error()})
	}

	categories, err := h.Service.ListAllCategories(filter)
	if err != nil {
		if strings.Contains(err.Error(), "invalid category_type_id") { // ควรมาจาก service layer ถ้ามีการ parse ที่นั่น
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to list categories", Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(categories)
}

// GET /api/v1/cms/categories/{categoryUuid}
// @Summary Get Category (Detail) by UUID
// @Description Retrieves a specific category detail by its UUID.
// @Tags CMS - Categories
// @Produce json
// @Param categoryUuid path string true "Category (Detail) UUID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid UUID format"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/categories/{categoryUuid} [get]
func (h *CMSCategoryHandler) HandleGetCategoryByUUID(c *fiber.Ctx) error {
	uuidStr := c.Params("categoryUuid")
	if _, err := uuid.Parse(uuidStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid category UUID format", Message: err.Error()})
	}

	categoryResponse, err := h.Service.GetCategoryByUUID(uuidStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Category not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get category", Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(categoryResponse)
}

// PATCH /api/v1/cms/categories/{categoryUuid}
// @Summary Update Category (Detail)
// @Description Updates an existing category detail by its UUID.
// @Tags CMS - Categories
// @Accept json
// @Produce json
// @Param categoryUuid path string true "Category (Detail) UUID"
// @Param category_data body dto.CategoryUpdateRequest true "Category (Detail) Update Data"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Category with this name, language, and type already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/categories/{categoryUuid} [patch]
func (h *CMSCategoryHandler) HandleUpdateCategory(c *fiber.Ctx) error {
	uuidStr := c.Params("categoryUuid")
	if _, err := uuid.Parse(uuidStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid category UUID format", Message: err.Error()})
	}

	var req dto.CategoryUpdateRequest // DTO นี้สำหรับอัปเดต Category (Detail) หนึ่งรายการ
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil { // Validate optional fields
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed for update request", Message: err.Error()})
	}

	categoryResponse, err := h.Service.UpdateCategoryByUUID(uuidStr, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		// "invalid" error for category_type_id should not happen here as CategoryUpdateRequest doesn't include it.
		// However, if name uniqueness check fails in service:
		if strings.Contains(err.Error(), "already exist") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to update category", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(categoryResponse)
}

// DELETE /api/v1/cms/categories/{categoryUuid}
// @Summary Delete Category (Detail)
// @Description Deletes a specific category detail by its UUID.
// @Tags CMS - Categories
// @Produce json
// @Param categoryUuid path string true "Category (Detail) UUID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Invalid UUID format"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /cms/categories/{categoryUuid} [delete]
func (h *CMSCategoryHandler) HandleDeleteCategory(c *fiber.Ctx) error {
	uuidStr := c.Params("categoryUuid")
	if _, err := uuid.Parse(uuidStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid category UUID format", Message: err.Error()})
	}

	err := h.Service.DeleteCategoryByUUID(uuidStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: fmt.Sprintf("Category with ID '%s' not found for delete", uuidStr)})
		}
		// "cannot delete category with children" logic is likely in CategoryType deletion now.
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to delete category", Message: err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

