package cms

import (
	"errors" // Import errors package
	"strconv"
	"strings" // Import strings package

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models/enums" // ต้อง import enums
	"github.com/MadManJJ/cms-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm" // Import gorm
)

type CMSCategoryTypeHandler struct {
	Service  services.CMSCategoryTypeServiceInterface
	validate *validator.Validate
}

func NewCMSCategoryTypeHandler(service services.CMSCategoryTypeServiceInterface) *CMSCategoryTypeHandler {
	return &CMSCategoryTypeHandler{
		Service:  service,
		validate: validator.New(),
	}
}

// POST /api/v1/cms/category-types
// @Summary      Create Category Type
// @Description  Creates a new category type
// @Tags         CMS - Category Types
// @Accept       json
// @Produce      json
// @Param 		 category_type_data body dto.CreateCategoryTypeRequest true "Category Type Data"
// @Success      201  {object} dto.CategoryTypeResponse
// @Failure 	 400  {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure 	 409  {object} dto.ErrorResponse "Conflict - TypeCode already exists"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/category-types [post]
func (h *CMSCategoryTypeHandler) HandleCreateCategoryType(c *fiber.Ctx) error {
	var req dto.CreateCategoryTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}
	resp, err := h.Service.CreateCategoryType(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to create category type", Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GET /api/v1/cms/category-types
// @Summary      List Category Types
// @Description  Lists all category types, optionally filtered by isActive status. Includes count of categories per language.
// @Tags         CMS - Category Types
// @Produce      json
// @Param        isActive  query     bool  false  "Filter by active status"
// @Success      200  {array} dto.CategoryTypeResponse
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/category-types [get]
func (h *CMSCategoryTypeHandler) HandleListCategoryTypes(c *fiber.Ctx) error {
	var isActive *bool
	if qIsActive := c.Query("isActive"); qIsActive != "" {
		val, err := strconv.ParseBool(qIsActive)
		if err == nil {
			isActive = &val
		} else {
			// Optional: return bad request if isActive is not a valid boolean
			// return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid isActive query parameter", Message: "isActive must be true or false"})
		}
	}
	resp, err := h.Service.ListCategoryTypes(isActive) // This service method should now populate children_count
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to list category types", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GET /api/v1/cms/category-types/{id}
// @Summary      Get Category Type by ID
// @Description  Retrieves a specific category type by its ID, including count of categories per language.
// @Tags         CMS - Category Types
// @Produce      json
// @Param        id path string true "Category Type ID (UUID)"
// @Success      200  {object} dto.CategoryTypeResponse
// @Failure      400  {object} dto.ErrorResponse "Invalid ID format"
// @Failure      404  {object} dto.ErrorResponse "Category type not found"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/category-types/{id} [get]
func (h *CMSCategoryTypeHandler) HandleGetCategoryType(c *fiber.Ctx) error {
	idStr := c.Params("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid ID format", Message: "ID must be a valid UUID"})
	}

	resp, err := h.Service.GetCategoryTypeByID(idStr) // This service method should now populate children_count
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Category type not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get category type", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// PATCH /api/v1/cms/category-types/{id}
// @Summary      Update Category Type
// @Description  Updates an existing category type by its ID.
// @Tags         CMS - Category Types
// @Accept       json
// @Produce      json
// @Param        id path string true "Category Type ID (UUID)"
// @Param 		 category_type_data body dto.UpdateCategoryTypeRequest true "Category Type Update Data"
// @Success      200  {object} dto.CategoryTypeResponse
// @Failure      400  {object} dto.ErrorResponse "Validation Error or Bad Request"
// @Failure      404  {object} dto.ErrorResponse "Category type not found"
// @Failure      409  {object} dto.ErrorResponse "Conflict - TypeCode already exists"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/category-types/{id} [patch]
func (h *CMSCategoryTypeHandler) HandleUpdateCategoryType(c *fiber.Ctx) error {
	idStr := c.Params("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid ID format", Message: "ID must be a valid UUID"})
	}
	var req dto.UpdateCategoryTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Cannot parse JSON", Message: err.Error()})
	}
	if err := h.validate.Struct(req); err != nil { // Validate even if fields are optional
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed", Message: err.Error()})
	}
	resp, err := h.Service.UpdateCategoryType(idStr, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to update category type", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// DELETE /api/v1/cms/category-types/{id}
// @Summary      Delete Category Type
// @Description  Deletes a category type by its ID. Fails if categories are associated with it.
// @Tags         CMS - Category Types
// @Produce      json
// @Param        id path string true "Category Type ID (UUID)"
// @Success      204  "No Content"
// @Failure      400  {object} dto.ErrorResponse "Invalid ID format"
// @Failure      404  {object} dto.ErrorResponse "Category type not found"
// @Failure      409  {object} dto.ErrorResponse "Conflict - Category type in use"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/category-types/{id} [delete]
func (h *CMSCategoryTypeHandler) HandleDeleteCategoryType(c *fiber.Ctx) error {
	idStr := c.Params("id")
	_, err := uuid.Parse(idStr) // Validate UUID format first
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid ID format", Message: "ID must be a valid UUID"})
	}

	err = h.Service.DeleteCategoryType(idStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		if strings.Contains(err.Error(), "in use") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "Conflict", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to delete category type", Message: err.Error()})
	}

	// For DELETE, a 204 No Content is more standard than returning a JSON body
	return c.SendStatus(fiber.StatusNoContent)
}

// GET /api/v1/cms/category-types/{categoryTypeId}/categories
// @Summary      List Categories for a Specific Type
// @Description  Retrieves all categories (details) associated with a specific category type, filtered by language.
// @Tags         CMS - Category Types
// @Produce      json
// @Param        categoryTypeId path string true "Category Type ID (UUID)"
// @Param        lang query string true "Language code (th or en)" Enums(th, en)
// @Success      200  {object} dto.CategoryTypeWithDetailsResponse
// @Failure      400  {object} dto.ErrorResponse "Invalid ID or language format"
// @Failure      404  {object} dto.ErrorResponse "Category type not found"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/category-types/{categoryTypeId}/categories [get]
func (h *CMSCategoryTypeHandler) HandleListCategoriesForType(c *fiber.Ctx) error {
	categoryTypeIDStr := c.Params("categoryTypeId")
	langQuery := c.Query("lang")

	// Validate categoryTypeId from path
	if _, err := uuid.Parse(categoryTypeIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid categoryTypeId format", Message: "Category Type ID must be a valid UUID."})
	}

	// Validate lang query parameter
	if langQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Missing language query parameter", Message: "'lang' query parameter is required (e.g., th, en)."})
	}
	var lang enums.PageLanguage
	switch strings.ToLower(langQuery) {
	case string(enums.PageLanguageTH):
		lang = enums.PageLanguageTH
	case string(enums.PageLanguageEN):
		lang = enums.PageLanguageEN
	default:
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid language query parameter", Message: "Supported languages are 'th' or 'en'."})
	}

	// Call the service method (assuming it's created in CMSCategoryTypeServiceInterface and its implementation)
	// The service method GetCategoryTypeWithDetails should handle fetching the CategoryType
	// and then fetching its associated Categories based on the language.
	response, err := h.Service.GetCategoryTypeWithDetails(categoryTypeIDStr, string(lang))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: err.Error()})
		}
		// Handle other potential errors from the service, e.g., validation errors for language if service handles it
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to list categories for the specified type", Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
