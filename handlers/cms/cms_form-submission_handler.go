package cms

import (
	"strconv"

	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CMSFormSubmissionHandler struct {
	Service services.CMSFormSubmissionServiceInterface
}

func NewCMSFormSubmissionHandler(service services.CMSFormSubmissionServiceInterface) *CMSFormSubmissionHandler {
	return &CMSFormSubmissionHandler{Service: service}
}

// HandleCreateFormSubmission handles POST requests to create a new form submission
// @Summary      Create a New Form Submission
// @Description  Create a new form submission with optional components such as categories, components, and metatag.
// @Security     BearerAuth
// @Tags         CMS - Form Submissions
// @Accept       json
// @Produce      json
// @Param        formId  path  string  true  "Form ID"
// @Param        form_submission_data  body  dto.CreateFormSubmissionRequest  true  "Form Submission payload"
// @Success      200  {object} dto.CMSFormSubmissionSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/forms/{formId}/submissions [post]
func (h *CMSFormSubmissionHandler) HandleCreateFormSubmission(c *fiber.Ctx) error {
	formIdStr := c.Params("formId")
	formId, err := uuid.Parse(formIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the formId",
			"error":   err.Error(),
		})
	}

	var formSubmission models.FormSubmission
	if err := c.BodyParser(&formSubmission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse the formSubmission",
		})
	}

	createdFormSubmission, err := h.Service.CreateFormSubmission(formId, &formSubmission)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create the formSubmission",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "successfully created the formSubmission",
		"item":    createdFormSubmission,
	})
}

// HandleGetFormSubmissions handles GET requests to create a new form submission
// @Summary      List Submissions
// @Description  Get form submissions with optional components such as categories, components, and metatag.
// @Security     BearerAuth
// @Tags         CMS - Form Submissions
// @Produce      json
// @Param        formId  path  string  true  "Form ID"
// @Param        sort   query  string  false  "Sorting fields, e.g., 'title:asc,updated_at:desc'"
// @Param        page     query  int     false  "Page number for pagination (default is 1)"
// @Param        limit    query  int     false  "Number of items per page (default is 10)"
// @Success      200  {object} dto.CMSFormSubmissionsSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/forms/{formId}/submissions [get]
func (h *CMSFormSubmissionHandler) HandleGetFormSubmissions(c *fiber.Ctx) error {
	sort := c.Query("sort", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))		
	formIdStr := c.Params("formId")
	formId, err := uuid.Parse(formIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the formId",
			"error":   err.Error(),
		})
	}

	formSubmissions, totalCount, err := h.Service.GetFormSubmissions(formId, sort, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get the formSubmissions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully got the formSubmissions",
		"totalCount": totalCount,
		"page": page,
		"limit": limit,
		"items":   formSubmissions,
	})
}

// HandleGetFormSubmissions handles GET requests to get a form submission
// @Summary      Get Form Submission
// @Description  Retrieve a specific form submission using its unique identifier.
// @Security     BearerAuth
// @Tags         CMS - Form Submissions
// @Produce      json
// @Param        submissionId  path  string  true  "Form Submission ID"
// @Success      200  {object} dto.CMSFormSubmissionSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/forms/submissions/{submissionId} [get]
func (h *CMSFormSubmissionHandler) HandleGetFormSubmission(c *fiber.Ctx) error {
	submissionIdStr := c.Params("submissionId")
	submissionId, err := uuid.Parse(submissionIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the submissionId",
			"error":   err.Error(),
		})
	}

	formSubmission, err := h.Service.GetFormSubmission(submissionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get the formSubmission",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully got the formSubmission",
		"item":   formSubmission,
	})
}