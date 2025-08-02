package cms

import (
	"errors"
	"strconv"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CMSFaqPageHandler struct {
	Service services.CMSFaqPageServiceInterface
}

func NewCMSFaqPageHandler(service services.CMSFaqPageServiceInterface) *CMSFaqPageHandler {
	return &CMSFaqPageHandler{Service: service}
}

// HandleCreateFaqPage handles POST requests to create a new FAQ page
// @Summary      Create a New FAQ Page
// @Description  Create a new FAQ page with optional components such as categories, components, and metatag.
// @Description  The request must include exactly one content block.
// @Description  A valid revision must be included in the request other parts are optional.
// @Tags         CMS - Faq Pages
// @Accept       json
// @Produce      json
// @Param        faq_page_data  body  dto.CreateFaqPageRequest  true  "FAQ Page payload (optional: categories, components)"
// @Success      200  {object} dto.CMSFaqPageSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/faqpages [post]
func (h *CMSFaqPageHandler) HandleCreateFaqPage(c *fiber.Ctx) error {
	var faqPage models.FaqPage
	if err := c.BodyParser(&faqPage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query JSON",
			"error":   err.Error(),
		})
	}

	helpers.SanitizeFaqPage(&faqPage)

	createdFaqPage, err := h.Service.CreateFaqPage(&faqPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create faq page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully create page",
		"item":    createdFaqPage,
	})
}

// HandleGetFaqPages handles GET requests to retrieve FAQ pages
// @Summary      List FAQ Pages
// @Description  Retrieve a list of FAQ pages with optional filtering, sorting, pagination, and language selection.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        query  query  string  false  "Filter criteria (URL-encoded JSON string), e.g., category_faq and category_keywords"
// @Param        sort   query  string  false  "Sorting fields, e.g., 'title:asc,updated_at:desc'"
// @Param        page     query  int     false  "Page number for pagination (default is 1)"
// @Param        limit    query  int     false  "Number of items per page (default is 10)"
// @Param        language  query  string  false  "Language code for localized content (e.g., en, th)"
// @Success      200  {object} dto.CMSFaqPagesSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/faqpages [get]
func (h *CMSFaqPageHandler) HandleGetFaqPages(c *fiber.Ctx) error {
	rawQuery := c.Query("query")
	sort := c.Query("sort", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	language := c.Query("language", "")

	results, totalCount, err := h.Service.FindFaqPages(rawQuery, sort, page, limit, language)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidQuery) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to decode the query",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find faq pages",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "successfully get all faq pages",
		"totalCount": totalCount,
		"page":       page,
		"limit":      limit,
		"items":      results,
	})
}

// HandleGetFaqPageByID handles GET requests to retrieve an FAQ page by its ID
// @Summary      Get FAQ Page by ID
// @Description  Retrieve a specific FAQ page using its unique identifier.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageid  path  string  true  "FAQ Page ID"
// @Success      200  {object} dto.CMSFaqPageSuccessResponse200
// @Failure      400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/faqpages/{pageid} [get]
func (h *CMSFaqPageHandler) HandleGetFaqPageById(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	faqPage, err := h.Service.FindFaqPageById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find faq page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get the faq page",
		"item":    faqPage,
	})
}

// HandleDeleteFaqPage handles DELETE requests to remove an FAQ page by its ID
// @Summary      Delete FAQ Page
// @Description  Delete an existing FAQ page using its unique identifier.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId  path  string  true  "FAQ Page ID"
// @Success      200  {object}  dto.CMSSuccessResponse
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/{pageId} [delete]
func (h *CMSFaqPageHandler) HandleDeleteFaqPage(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	err = h.Service.DeleteFaqPage(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete faq page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete faq page",
	})
}

// HandleGetContentByFaqPageId handles GET requests to retrieve FAQ page content by ID and language code.
// @Summary      Get FAQ Page Content
// @Description  Retrieve the content of a specific FAQ page by its unique identifier and language code.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId        path      string  true   "FAQ Page ID"
// @Param        languageCode  path      string  true   "Language Code (e.g., en, th)"
// @Param        mode          query     string  false   "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSFaqContentSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/{pageId}/contents/{languageCode} [get]
func (h *CMSFaqPageHandler) HandleGetContentByFaqPageId(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	language := c.Params("languageCode")
	mode := c.Query("mode", "published")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	faqContent, err := h.Service.FindContentByFaqPageId(id, language, mode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get faq page content",
		"item":    faqContent,
	})
}

// HandleGetLatestContentByFaqPageId handles GET requests to retrieve the latest FAQ content by page ID.
// @Summary      Get Latest FAQ Page Content
// @Description  Retrieve the latest content for FAQ page.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId        path      string  true   "FAQ Page ID (UUID format)"
// @Param        languageCode  path      string  true   "Language Code (e.g., en, th)"
// @Success      200  {object}  dto.CMSFaqContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/{pageId}/latestcontent/{languageCode} [get]
func (h *CMSFaqPageHandler) HandleGetLatestContentByFaqPageId(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	language := c.Params("languageCode")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	faqContent, err := h.Service.FindLatestContentByPageId(pageId, language)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get the latest content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully find the latest faq content",
		"item":    faqContent,
	})
}

// HandleDeleteFaqContentByPageId handles DELETE request to delete FAQ page content.
// @Summary      Delete FAQ Page Content
// @Description  Delete the content of a specific FAQ page using its unique identifier, language code, and mode.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId        path      string  true  "FAQ Page ID (UUID)"
// @Param        languageCode  path      string  true  "Language Code (e.g., en, th)"
// @Param        mode          query     string  false  "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSSuccessResponse
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/{pageId}/contents/{languageCode} [delete]
func (h *CMSFaqPageHandler) HandleDeleteFaqContentByPageId(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	language := c.Params("languageCode")
	mode := c.Query("mode", "published")

	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	if err := h.Service.DeleteContentByFaqPageId(pageId, language, mode); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to Delete faq content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete faq content",
	})
}

// HandleDuplicateFaqPage handles POST request to duplicate a FAQ Page.
// @Summary      Duplicate FAQ Page
// @Description  Duplicate an existing FAQ page using its unique identifier and.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId        path      string  true  "FAQ Page ID (UUID)"
// @Success      200           {object}  dto.CMSFaqPageSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/duplicate/{pageId}/pages [post]
func (h *CMSFaqPageHandler) HandleDuplicateFaqPage(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	faqPage, err := h.Service.DuplicateFaqPage(pageId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to duplicate page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully duplicate page",
		"item":    faqPage,
	})
}

// HandleDuplicateFaqContentToAnotherLanguage handles POST request to duplicate a FAQ Content to another language.
// @Summary      Duplicate FAQ Content to another language
// @Description  Duplicate an existing FAQ content using its unique identifier and
// @Tags         CMS - Faq Pages
// @Accept       json
// @Produce      json
// @Param        contentId        path      string  true  "FAQ content ID (UUID)"
// @Param        revision  body  dto.CreateRevisionRequest  true  "Revision payload"
// @Success      200           {object}  dto.CMSFaqContentSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/duplicate/{contentId}/contents [post]
func (h *CMSFaqPageHandler) HandleDuplicateFaqContentToAnotherLanguage(c *fiber.Ctx) error {
	contentIdStr := c.Params("contentId")
	contentId, err := uuid.Parse(contentIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	var revision models.Revision
	if err := c.BodyParser(&revision); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the body",
			"error":   err.Error(),
		})
	}

	helpers.SanitizeRevision(&revision)

	faqContent, err := h.Service.DuplicateFaqContentToAnotherLanguage(contentId, &revision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to duplicate content to another language",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully duplicate content",
		"item":    faqContent,
	})
}

// HandleRevertFaqContent handles POST request to revert a FAQ Page to a specific revision.
// @Summary      Revert FAQ Page Content
// @Description  Revert FAQ content using a specific revision ID.
// @Tags         CMS - Faq Pages
// @Accept       json
// @Produce      json
// @Param        revisionId  path  string  true  "Revision ID (UUID)"
// @Param        revision    body  dto.CreateRevisionRequest  true  "Revision payload"
// @Success      200  {object}  dto.CMSFaqContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/{revisionId}/revisions [post]
func (h *CMSFaqPageHandler) HandleRevertFaqContent(c *fiber.Ctx) error {
	revisionIdStr := c.Params("revisionId")
	revisionId, err := uuid.Parse(revisionIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parsed revisionId",
			"error":   err.Error(),
		})
	}

	var revision models.Revision
	if err := c.BodyParser(&revision); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parsed the body",
			"error":   err.Error(),
		})
	}

	helpers.SanitizeRevision(&revision)

	faqContent, err := h.Service.RevertFaqContent(revisionId, &revision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to revert faq content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully revert faq content",
		"item":    faqContent,
	})
}

// HandleUpdateFaqContent handles PUT request to update FAQ content.
// @Summary      Update FAQ Content
// @Description  Update an existing FAQ content by its content ID.
// @Description  This operation does not merge with previous content — it overwrites everything.
// @Description  A valid revision must be included in the request.
// @Tags         CMS - Faq Pages
// @Accept       json
// @Produce      json
// @Param        contentId      path  string           true  "FAQ Content ID (UUID)"
// @Param        faqContent     body  dto.CreateFaqContentRequest  true  "Updated FAQ Content"
// @Success      200  {object}  dto.CMSFaqContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/{contentId}/contents [put]
func (h *CMSFaqPageHandler) HandleUpdateFaqContent(c *fiber.Ctx) error {
	contentIdStr := c.Params("contentId")
	contentId, err := uuid.Parse(contentIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	var updatedContent models.FaqContent
	if err := c.BodyParser(&updatedContent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the body",
			"error":   err.Error(),
		})
	}

	helpers.SanitizeFaqContent(&updatedContent)

	faqContent, err := h.Service.UpdateFaqContent(&updatedContent, contentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update faq content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully update faq content",
		"item":    faqContent,
	})
}

// HandleGetCategory handles GET request to retrieve Category for that FAQ Page
// @Summary      Get FAQ Page Categories
// @Description  Retrieve the categories of a specific FAQ page using its unique identifier, category type code, language code, and mode.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId        path      string  true  "FAQ Page ID (UUID)"
// @Param        categoryTypeCode  path      string  true  "Category Type Code"
// @Param        languageCode      path      string  true  "Language Code (e.g., en, th)"
// @Param        mode              query     string  false "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSFaqCategoriesSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/category/{categoryTypeCode}/{pageId}/{languageCode} [get]
func (h *CMSFaqPageHandler) HandleGetCategory(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	categoryTypeCode := c.Params("categoryTypeCode")
	language := c.Params("languageCode")
	mode := c.Query("mode", "published")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}
	categories, err := h.Service.FindCategories(pageId, categoryTypeCode, language, mode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get categories",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get categories",
		"item":    categories,
	})
}

// HandleGetRevisions handles GET request to retrieve revisions of a specific FAQ page
// @Summary      Get FAQ Page Revisions
// @Description  Retrieve a list of revisions for a specific FAQ page based on its unique identifier and language code. Each revision includes metadata such as author, status, and timestamps.
// @Tags         CMS - Faq Pages
// @Produce      json
// @Param        pageId        path      string  true  "FAQ Page ID (UUID)"
// @Param        languageCode      path      string  true  "Language Code (e.g., en, th)"
// @Success      200           {object}  dto.CMSFaqRevsionsSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/revisions/{languageCode}/{pageId} [get]
func (h *CMSFaqPageHandler) HandleGetRevisions(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	language := c.Params("languageCode")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	revisions, err := h.Service.FindRevisions(pageId, language)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get revisions",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get revisions",
		"item":    revisions,
	})
}

// HandlePreviewFaqContent handles POST request to preview FAQ content.
// @Summary      Preview FAQ Content
// @Description  Preview an existing FAQ content by its content ID.
// @Description  This operation does not merge with previous content — it overwrites everything.
// @Tags         CMS - Faq Pages
// @Accept       json
// @Produce      json
// @Param        pageId      path  string           true  "FAQ Page ID (UUID)"
// @Param        faqContent     body  dto.CreateFaqContentPreviewRequest  true  "Preview FAQ Content"
// @Success      200  {object}  dto.CMSFaqContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/faqpages/previews/{pageId} [post]
func (h *CMSFaqPageHandler) HandlePreviewFaqContent(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	var faqContentPreview models.FaqContent
	if err := c.BodyParser(&faqContentPreview); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid faq content JSON",
			"error":   err.Error(),
		})
	}	

	helpers.SanitizeFaqContent(&faqContentPreview)

	url, err := h.Service.PreviewFaqContent(pageId, &faqContentPreview)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to preview content",
			"error":   err.Error(),
		})
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully preview faq content",
		"url":    url,
	})
}