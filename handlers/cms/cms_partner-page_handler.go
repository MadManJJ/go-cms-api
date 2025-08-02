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

type CMSPartnerPageHandler struct {
	Service services.CMSPartnerPageServiceInterface
}

func NewCMSPartnerPageHandler(service services.CMSPartnerPageServiceInterface) *CMSPartnerPageHandler {
	return &CMSPartnerPageHandler{Service: service}
}

// HandleCreatePartnerPage handles POST requests to create a new Partner page
// @Summary      Create a New Partner Page
// @Description  Create a new Partner page with optional components such as categories, revisions, and content blocks.
// @Tags         CMS - Partner Pages
// @Accept       json
// @Produce      json
// @Param        Partner_page_data  body  dto.CreatePartnerPageRequest  true  "Partner Page payload (optional: categories, components)"
// @Success      200  {object} dto.CMSPartnerPageSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/partnerpages [post]
func (h *CMSPartnerPageHandler) HandleCreatePartnerPage(c *fiber.Ctx) error {
	var partnerPage models.PartnerPage
	if err := c.BodyParser(&partnerPage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query JSON",
			"error":   err.Error(),
		})
	}

	helpers.SanitizePartnerPage(&partnerPage)

	createdPartnerPage, err := h.Service.CreatePartnerPage(&partnerPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Partner page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully create page",
		"item":    createdPartnerPage,
	})
}

// HandleGetPartnerPages handles GET requests to retrieve Partner pages
// @Summary      List Partner Pages
// @Description  Retrieve a list of Partner pages with optional filtering, sorting, pagination, and language selection.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        query  query  string  false  "Filter criteria (URL-encoded JSON string), e.g., category_Partner and category_keywords"
// @Param        sort   query  string  false  "Sorting fields, e.g., 'title:asc,updated_at:desc'"
// @Param        page     query  int     false  "Page number for pagination (default is 1)"
// @Param        limit    query  int     false  "Number of items per page (default is 10)"
// @Param        language  query  string  false  "Language code for localized content (e.g., en, th)"
// @Success      200  {object} dto.CMSPartnerPageSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/partnerpages [get]
func (h *CMSPartnerPageHandler) HandleGetPartnerPages(c *fiber.Ctx) error {
	rawQuery := c.Query("query")
	sort := c.Query("sort", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	language := c.Query("language", "")

	results, totalCount, err := h.Service.FindPartnerPages(rawQuery, sort, page, limit, language)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidQuery) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to decode the query",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find Partner pages",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "successfully get all Partner pages",
		"totalCount": totalCount,
		"page":       page,
		"limit":      limit,
		"items":      results,
	})
}

// HandleGetPartnerPageByID handles GET requests to retrieve an Partner page by its ID
// @Summary      Get Partner Page by ID
// @Description  Retrieve a specific Partner page using its unique identifier.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageid  path  string  true  "Partner Page ID"
// @Success      200  {object} dto.CMSPartnerPageSuccessResponse200
// @Failure      400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/partnerpages/{pageid} [get]
func (h *CMSPartnerPageHandler) HandleGetPartnerPageById(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	partnerPage, err := h.Service.FindPartnerPageById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find Partner page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get the Partner page",
		"item":    partnerPage,
	})
}

// HandleDeletePartnerPage handles DELETE requests to remove an Partner page by its ID
// @Summary      Delete Partner Page
// @Description  Delete an existing Partner page using its unique identifier.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId  path  string  true  "Partner Page ID"
// @Success      200  {object}  dto.CMSSuccessResponse
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/{pageId} [delete]
func (h *CMSPartnerPageHandler) HandleDeletePartnerPage(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	err = h.Service.DeletePartnerPage(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete Partner page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete Partner page",
	})
}

// HandleGetContentByPartnerPageId handles GET requests to retrieve Partner page content by ID and language code.
// @Summary      Get Partner Page Content
// @Description  Retrieve the content of a specific Partner page by its unique identifier and language code.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId        path      string  true   "Partner Page ID"
// @Param        languageCode  path      string  true   "Language Code (e.g., en, th)"
// @Param        mode          query     string  false   "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSPartnerContentSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/{pageId}/contents/{languageCode} [get]
func (h *CMSPartnerPageHandler) HandleGetContentByPartnerPageId(c *fiber.Ctx) error {
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

	PartnerContent, err := h.Service.FindContentByPartnerPageId(id, language, mode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get Partner page content",
		"item":    PartnerContent,
	})
}

// HandleGetLatestContentByPartnerPageId handles GET requests to retrieve the latest Partner content by page ID.
// @Summary      Get Latest Partner Page Content
// @Description  Retrieve the latest content for Partner page.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId        path      string  true   "Partner Page ID (UUID format)"
// @Param        languageCode  path      string  true   "Language Code (e.g., en, th)"
// @Success      200  {object}  dto.CMSPartnerContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/{pageId}/latestcontent/{languageCode} [get]
func (h *CMSPartnerPageHandler) HandleGetLatestContentByPartnerPageId(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	language := c.Params("languageCode")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	PartnerContent, err := h.Service.FindLatestContentByPageId(pageId, language)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get the latest content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully find the latest Partner content",
		"item":    PartnerContent,
	})
}

// HandleDeletePartnerContentByPageId handles DELETE request to delete Partner page content.
// @Summary      Delete Partner Page Content
// @Description  Delete the content of a specific Partner page using its unique identifier, language code, and mode.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId        path      string  true  "Partner Page ID (UUID)"
// @Param        languageCode  path      string  true  "Language Code (e.g., en, th)"
// @Param        mode          query     string  false  "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSSuccessResponse
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/{pageId}/contents/{languageCode} [delete]
func (h *CMSPartnerPageHandler) HandleDeletePartnerContentByPageId(c *fiber.Ctx) error {
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

	if err := h.Service.DeleteContentByPartnerPageId(pageId, language, mode); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to Delete Partner content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete Partner content",
	})
}

// HandleDuplicatePartnerPage handles POST request to duplicate a Partner Page.
// @Summary      Duplicate Partner Page
// @Description  Duplicate an existing Partner page using its unique identifier and language code.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId        path      string  true  "Partner Page ID (UUID)"
// @Success      200           {object}  dto.CMSSuccessResponse
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/duplicate/{pageId} [post]
func (h *CMSPartnerPageHandler) HandleDuplicatePartnerPage(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	partnerPage, err := h.Service.DuplicatePartnerPage(pageId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to duplicate page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully duplicate page",
		"item":    partnerPage,
	})
}

// HandleDuplicatePartnerContentToAnotherLanguage handles POST request to duplicate a Partner Content to another language.
// @Summary      Duplicate Partner Content to another language
// @Description  Duplicate an existing Partner content using its unique identifier and
// @Tags         CMS - Partner Pages
// @Accept       json
// @Produce      json
// @Param        contentId        path      string  true  "Partner content ID (UUID)"
// @Param        revision  body  dto.CreateRevisionRequest  true  "Revision payload"
// @Success      200           {object}  dto.CMSPartnerContentSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/duplicate/{contentId}/contents [post]
func (h *CMSPartnerPageHandler) HandleDuplicatePartnerContentToAnotherLanguage(c *fiber.Ctx) error {
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

	partnerContent, err := h.Service.DuplicatePartnerContentToAnotherLanguage(contentId, &revision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to duplicate content to another language",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully duplicate content",
		"item":    partnerContent,
	})
}

// HandleRevertPartnerContent handles POST request to revert a Partner Page to a specific revision.
// @Summary      Revert Partner Page Content
// @Description  Revert Partner content using a specific revision ID.
// @Tags         CMS - Partner Pages
// @Accept       json
// @Produce      json
// @Param        revisionId  path  string  true  "Revision ID (UUID)"
// @Param        revision    body  dto.CreateRevisionRequest  true  "Revision payload"
// @Success      200  {object}  dto.CMSPartnerContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/{revisionId}/revisions [post]
func (h *CMSPartnerPageHandler) HandleRevertPartnerContent(c *fiber.Ctx) error {
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

	PartnerContent, err := h.Service.RevertPartnerContent(revisionId, &revision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to revert Partner content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully revert Partner content",
		"item":    PartnerContent,
	})
}

// HandleUpdatePartnerContent handles PUT request to update Partner content.
// @Summary      Update Partner Content
// @Description  Update an existing Partner content by its content ID.
// @Description  This operation does not merge with previous content — it overwrites everything.
// @Description  A valid revision must be included in the request.
// @Tags         CMS - Partner Pages
// @Accept       json
// @Produce      json
// @Param        contentId      path  string           true  "Partner Content ID (UUID)"
// @Param        partnerContent     body  dto.CreatePartnerContentRequest  true  "Updated Partner Content"
// @Success      200  {object}  dto.CMSPartnerContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/{contentId}/contents [put]
func (h *CMSPartnerPageHandler) HandleUpdatePartnerContent(c *fiber.Ctx) error {
	contentIdStr := c.Params("contentId")
	contentId, err := uuid.Parse(contentIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	var updatedContent models.PartnerContent
	if err := c.BodyParser(&updatedContent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the body",
			"error":   err.Error(),
		})
	}

	helpers.SanitizePartnerContent(&updatedContent)

	PartnerContent, err := h.Service.UpdatePartnerContent(&updatedContent, contentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update Partner content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully update Partner content",
		"item":    PartnerContent,
	})
}

// HandleGetCategory handles GET request to retrieve Category for that Partner Page
// @Summary      Get Partner Page Categories
// @Description  Retrieve the categories of a specific Partner page using its unique identifier, category type code, language code, and mode.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId        path      string  true  "Partner Page ID (UUID)"
// @Param        categoryTypeCode  path      string  true  "Category Type Code"
// @Param        languageCode      path      string  true  "Language Code (e.g., en, th)"
// @Param        mode              query     string  false "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSPartnerCategoriesSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/category/{categoryTypeCode}/{pageId}/{languageCode} [get]
func (h *CMSPartnerPageHandler) HandleGetCategory(c *fiber.Ctx) error {
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
	categories, err := h.Service.GetCategory(pageId, categoryTypeCode, language, mode)
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

// HandleGetRevisions handles GET request to retrieve revisions of a specific Partner page
// @Summary      Get Partner Page Revisions
// @Description  Retrieve a list of revisions for a specific Partner page based on its unique identifier and language code. Each revision includes metadata such as author, status, and timestamps.
// @Tags         CMS - Partner Pages
// @Produce      json
// @Param        pageId        path      string  true  "Partner Page ID (UUID)"
// @Param        languageCode      path      string  true  "Language Code (e.g., en, th)"
// @Success      200           {object}  dto.CMSPartnerRevsionsSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/revisions/{languageCode}/{pageId} [get]
func (h *CMSPartnerPageHandler) HandleGetRevisions(c *fiber.Ctx) error {
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

// HandlePreviewPartnerContent handles POST request to preview Partner content.
// @Summary      Preview Partner Content
// @Description  Preview an existing Partner content by its content ID.
// @Description  This operation does not merge with previous content — it overwrites everything.
// @Tags         CMS - Partner Pages
// @Accept       json
// @Produce      json
// @Param        pageId      path  string           true  "Partner Page ID (UUID)"
// @Param        partnerContent     body  dto.CreatePartnerContentPreviewRequest  true  "Preview Partner Content"
// @Success      200  {object}  dto.CMSPartnerContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/partnerpages/previews/{pageId} [post]
func (h *CMSPartnerPageHandler) HandlePreviewPartnerContent(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	var partnerContentPreview models.PartnerContent
	if err := c.BodyParser(&partnerContentPreview); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query JSON",
			"error":   err.Error(),
		})
	}	

	helpers.SanitizePartnerContent(&partnerContentPreview)

	url, err := h.Service.PreviewPartnerContent(pageId, &partnerContentPreview)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to preview content",
			"error":   err.Error(),
		})
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully preview partner content",
		"url":    url,
	})
}