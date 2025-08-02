package cms

import (
	"github.com/MadManJJ/cms-api/dto"  // หรือ path ที่ถูกต้องสำหรับ DTOs ของคุณ
	"github.com/MadManJJ/cms-api/errs" // Custom errors

	// ถ้าคุณตัดสินใจใช้ Models เป็น Request Body โดยตรง
	"errors"
	"fmt"
	"strings"

	"github.com/MadManJJ/cms-api/services"

	// สำหรับ ParseInt
	"github.com/go-playground/validator/v10" // Validator
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CMSFormHandler struct {
	service  services.CMSFormServiceInterface
	validate *validator.Validate
}

func NewCMSFormHandler(service services.CMSFormServiceInterface) *CMSFormHandler {
	return &CMSFormHandler{
		service:  service,
		validate: validator.New(),
	}
}

// --- Handler Methods ---

// POST /api/v1/cms/forms
// HandleCreateForm สร้าง Form Template ใหม่
// @Summary Create Form Template
// @Description Creates a new Form Template with its sections and fields.
// @Tags CMS - Forms
// @Accept json
// @Produce json
// @Param form_data body dto.CreateFormRequest true "Form Template Data" // หรือ models.Form ถ้าใช้ Model โดยตรง
// @Success 201 {object} dto.FormResponse "Successfully created form template"
// @Failure 400 {object} dto.ErrorResponse "Bad Request - Validation error or invalid input"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /cms/forms [post]
// handlers/cms/cms_form_handler.go
func (h *CMSFormHandler) HandleCreateForm(c *fiber.Ctx) error {

	var createFormReq dto.CreateFormRequest
	if err := c.BodyParser(&createFormReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Cannot parse JSON request body into DTO", Error: err.Error()})
	}
	if createFormReq.EmailCategoryID != nil && *createFormReq.EmailCategoryID == "" {
		createFormReq.EmailCategoryID = nil
	}

	// Validate DTO
	if err := h.validate.Struct(createFormReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	createdFormResponse, err := h.service.CreateNewForm(createFormReq)

	if err != nil {
		if errors.Is(err, errs.ErrBadRequest) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Failed to create form", Error: err.Error()})
		}
		fmt.Printf("Error creating form: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "Internal server error while creating form"})
	}

	return c.Status(fiber.StatusCreated).JSON(createdFormResponse)
}

// GET /api/v1/cms/forms
// HandleListForms ดึงรายการ Form Templates ทั้งหมด
// @Summary List Form Templates
// @Description Retrieves a list of all form templates with pagination and filtering.
// @Tags CMS - Forms
// @Produce json
// @Param name query string false "Filter by form name (partial match)"
// @Param created_by_user_id query string false "Filter by creator User ID (UUID)"
// @Param created_at query string false "Filter by creation date (YYYY-MM-DD)"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param items_per_page query int false "Number of items per page (default: 10, max: 100)"
// @Param sort_by query string false "Sort by field (e.g., name, created_at, updated_at)" Enums(name, created_at, updated_at)
// @Param sort_order query string false "Sort order (asc, desc)" Enums(asc, desc)
// @Success 200 {object} dto.PaginatedFormListResponse "List of form templates"
// @Failure 400 {object} dto.ErrorResponse "Bad Request - Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /cms/forms [get]
func (h *CMSFormHandler) HandleListForms(c *fiber.Ctx) error {
	var filter dto.FormListFilter

	// Parse query parameters into the filter struct
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid query parameters", Error: err.Error()})
	}

	if err := h.validate.Struct(filter); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Validation failed for query parameters", Error: err.Error()})
	}

	// Set default pagination and sorting if not provided or handle in service
	if filter.Page == nil {
		defaultPage := 1
		filter.Page = &defaultPage
	}
	if filter.ItemsPerPage == nil {
		defaultLimit := 10
		filter.ItemsPerPage = &defaultLimit
	}

	paginatedResponse, err := h.service.GetAllForms(filter)
	if err != nil {
		if errors.Is(err, errs.ErrBadRequest) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Failed to list forms due to bad filter", Error: err.Error()})
		}
		fmt.Printf("Error listing forms: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "Internal server error while listing forms"})
	}

	return c.Status(fiber.StatusOK).JSON(paginatedResponse)
}

// GET /api/v1/cms/forms/{formId}
// HandleGetForm ดึงข้อมูล Form Template เดียวตาม ID
// @Summary Get Form Template by ID
// @Description Retrieves a specific form template by its UUID, including sections and fields.
// @Tags CMS - Forms
// @Produce json
// @Param formId path string true "Form Template ID (UUID)"
// @Success 200 {object} dto.FormResponse "Details of the form template"
// @Failure 400 {object} dto.ErrorResponse "Bad Request - Invalid Form ID format"
// @Failure 404 {object} dto.ErrorResponse "Not Found - Form template not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /cms/forms/{formId} [get]
func (h *CMSFormHandler) HandleGetForm(c *fiber.Ctx) error {
	formIDStr := c.Params("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid form ID format", Error: err.Error()})
	}

	formResponse, err := h.service.GetFormDetails(formID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Message: fmt.Sprintf("Form with ID %s not found", formIDStr)})
		}
		fmt.Printf("Error getting form %s: %v\n", formIDStr, err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "Internal server error while getting form details"})
	}

	return c.Status(fiber.StatusOK).JSON(formResponse)
}

// GET /api/v1/forms/{formId}/structure
// HandleGetForm ดึงข้อมูล Form Structure เดียวตาม ID
// @Summary Get Form Structure by ID
// @Description Retrieves a specific form structure by its UUID, including sections and fields.
// @Tags App - Forms
// @Produce json
// @Param formId path string true "Form Structure ID (UUID)"
// @Success 200 {object} dto.FormResponse "Details of the form structure"
// @Failure 400 {object} dto.ErrorResponse "Bad Request - Invalid Form ID format"
// @Failure 404 {object} dto.ErrorResponse "Not Found - Form structure not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /app/forms/{formId}/structure [get]
func (h *CMSFormHandler) HandleGetFormStructure(c *fiber.Ctx) error {
	formIDStr := c.Params("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid form ID format", Error: err.Error()})
	}

	formResponse, err := h.service.GetFormStructure(formID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Message: fmt.Sprintf("Form with ID %s not found", formIDStr)})
		}
		fmt.Printf("Error getting form %s: %v\n", formIDStr, err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "Internal server error while getting form details"})
	}

	return c.Status(fiber.StatusOK).JSON(formResponse)
}

// HandleUpdateForm อัปเดต Form Template ที่มีอยู่
// @Summary Update Form Template
// @Description Updates an existing Form Template, including its sections and fields.
// @Tags CMS - Forms
// @Accept json
// @Produce json
// @Param formId path string true "Form Template ID (UUID)"
// @Param body body dto.UpdateFormRequest true "Updated Form Template Data"
// @Success 200 {object} dto.FormResponse "Successfully updated form template"
// @Failure 400 {object} dto.ErrorResponse "Bad Request - Validation error or invalid input"
// @Failure 404 {object} dto.ErrorResponse "Not Found - Form template not found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /cms/forms/{formId} [put]
func (h *CMSFormHandler) HandleUpdateForm(c *fiber.Ctx) error {
	formIDStr := c.Params("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid form ID format", Error: err.Error()})
	}

	var updateFormReq dto.UpdateFormRequest
	if err := c.BodyParser(&updateFormReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Cannot parse JSON request body", Error: err.Error()})
	}
	if updateFormReq.EmailCategoryID != nil && *updateFormReq.EmailCategoryID == "" {
		updateFormReq.EmailCategoryID = nil
	}

	// Validate DTO
	if err := h.validate.Struct(updateFormReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	updatedFormResponse, err := h.service.UpdateExistingForm(formID, updateFormReq)

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Message: fmt.Sprintf("Form with ID %s not found for update", formIDStr)})
		}
		if errors.Is(err, errs.ErrBadRequest) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Failed to update form due to invalid data", Error: err.Error()})
		}
		if strings.Contains(err.Error(), "url_send field is read-only") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Message: "Invalid request: URL send field cannot be updated",
				Error:   err.Error(),
			})
		}
		fmt.Printf("Error updating form %s: %v\n", formIDStr, err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "Internal server error while updating form"})
	}

	return c.Status(fiber.StatusOK).JSON(updatedFormResponse)
}

// DELETE /api/v1/cms/forms/{formId}
// HandleDeleteForm ลบ Form Template
// @Summary Delete Form Template
// @Description Deletes a specific form template by its UUID (soft delete).
// @Tags CMS - Forms
// @Produce json
// @Param formId path string true "Form Template ID (UUID)"
// @Success 204 "No Content - Form template deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad Request - Invalid Form ID format"
// @Failure 404 {object} dto.ErrorResponse "Not Found - Form template not found"
// @Failure 409 {object} dto.ErrorResponse "Conflict - Cannot delete form (e.g., has submissions and ON DELETE RESTRICT)"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /cms/forms/{formId} [delete]
func (h *CMSFormHandler) HandleDeleteForm(c *fiber.Ctx) error {
	formIDStr := c.Params("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid form ID format", Error: err.Error()})
	}

	err = h.service.DeleteExistingForm(formID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Message: fmt.Sprintf("Form with ID %s not found for deletion", formIDStr)})
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") && strings.Contains(err.Error(), "form_submissions_form_id_fkey") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Message: "Cannot delete form: it has existing submissions."})
		}
		fmt.Printf("Error deleting form %s: %v\n", formIDStr, err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: "Internal server error while deleting form"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
