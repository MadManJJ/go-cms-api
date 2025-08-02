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

type CMSLandingPageHandler struct {
	Service services.CMSLandingPageServiceInterface
}

func NewCMSLandingPageHandler(service services.CMSLandingPageServiceInterface) *CMSLandingPageHandler {
	return &CMSLandingPageHandler{Service: service}
}

// HandleCreateLandingPage handles POST requests to create a new Landing page
// @Summary      Create a New Landing Page
// @Description  Create a new Landing page with optional components such as categories, revisions, and content blocks.
// @Tags         CMS - Landing Pages
// @Accept       json
// @Produce      json
// @Param        Landing_page_data  body  dto.CreateLandingPageRequest  true  "Landing Page payload (optional: categories, components)"
// @Success      200  {object} dto.CMSLandingPageSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/landingpages [post]
func (h *CMSLandingPageHandler) HandleCreateLandingPage(c *fiber.Ctx) error {
	var landingPage models.LandingPage
	if err := c.BodyParser(&landingPage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query JSON",
			"error":   err.Error(),
		})
	}

	helpers.SanitizeLandingPage(&landingPage)

	createdLandingPage, err := h.Service.CreateLandingPage(&landingPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Landing page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully create page",
		"item":    createdLandingPage,
	})
}

// HandleGetLandingPages handles GET requests to retrieve Landing pages
// @Summary      List Landing Pages
// @Description  Retrieve a list of Landing pages with optional filtering, sorting, pagination, and language selection.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        query  query  string  false  "Filter criteria (URL-encoded JSON string), e.g., category_Landing and category_keywords"
// @Param        sort   query  string  false  "Sorting fields, e.g., 'title:asc,updated_at:desc'"
// @Param        page     query  int     false  "Page number for pagination (default is 1)"
// @Param        limit    query  int     false  "Number of items per page (default is 10)"
// @Param        language  query  string  false  "Language code for localized content (e.g., en, th)"
// @Success      200  {object} dto.CMSLandingPageSuccessResponse200
// @Failure 		 400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/landingpages [get]
func (h *CMSLandingPageHandler) HandleGetLandingPages(c *fiber.Ctx) error {
	rawQuery := c.Query("query")
	sort := c.Query("sort", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	language := c.Query("language", "")

	results, totalCount, err := h.Service.FindLandingPages(rawQuery, sort, page, limit, language)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidQuery) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to decode the query",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find Landing pages",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "successfully get all Landing pages",
		"totalCount": totalCount,
		"page":       page,
		"limit":      limit,
		"items":      results,
	})
}

// HandleGetLandingPageByID handles GET requests to retrieve an Landing page by its ID
// @Summary      Get Landing Page by ID
// @Description  Retrieve a specific Landing page using its unique identifier.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageid  path  string  true  "Landing Page ID"
// @Success      200  {object} dto.CMSLandingPageSuccessResponse200
// @Failure      400  {object} dto.ErrorResponse400
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /cms/landingpages/{pageid} [get]
func (h *CMSLandingPageHandler) HandleGetLandingPageById(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	LandingPage, err := h.Service.FindLandingPageById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to find Landing page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get the Landing page",
		"item":    LandingPage,
	})
}

// HandleDeleteLandingPage handles DELETE requests to remove an Landing page by its ID
// @Summary      Delete Landing Page
// @Description  Delete an existing Landing page using its unique identifier.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId  path  string  true  "Landing Page ID"
// @Success      200  {object}  dto.CMSSuccessResponse
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/{pageId} [delete]
func (h *CMSLandingPageHandler) HandleDeleteLandingPage(c *fiber.Ctx) error {
	idStr := c.Params("pageId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	err = h.Service.DeleteLandingPage(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete Landing page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete Landing page",
	})
}

// HandleGetContentByLandingPageId handles GET requests to retrieve Landing page content by ID and language code.
// @Summary      Get Landing Page Content
// @Description  Retrieve the content of a specific Landing page by its unique identifier and language code.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId        path      string  true   "Landing Page ID"
// @Param        languageCode  path      string  true   "Language Code (e.g., en, th)"
// @Param        mode          query     string  false   "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSLandingContentSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/{pageId}/contents/{languageCode} [get]
func (h *CMSLandingPageHandler) HandleGetContentByLandingPageId(c *fiber.Ctx) error {
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

	LandingContent, err := h.Service.FindContentByLandingPageId(id, language, mode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get Landing page content",
		"item":    LandingContent,
	})
}

// HandleGetLatestContentBylandingPageId handles GET requests to retrieve the latest landing content by page ID.
// @Summary      Get Latest landing Page Content
// @Description  Retrieve the latest content for landing page.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId        path      string  true   "landing Page ID (UUID format)"
// @Param        languageCode  path      string  true   "Language Code (e.g., en, th)"
// @Success      200  {object}  dto.CMSLandingContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/{pageId}/latestcontents/{languageCode} [get]
func (h *CMSLandingPageHandler) HandleGetLatestContentByLandingPageId(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	language := c.Params("languageCode")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	LandingContent, err := h.Service.FindLatestContentByPageId(pageId, language)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get the latest content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully find the latest Landing content",
		"item":    LandingContent,
	})
}

// HandleDeleteLandingContentByPageId handles DELETE request to delete Landing page content.
// @Summary      Delete Landing Page Content
// @Description  Delete the content of a specific Landing page using its unique identifier, language code, and mode.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId        path      string  true  "Landing Page ID (UUID)"
// @Param        languageCode  path      string  true  "Language Code (e.g., en, th)"
// @Param        mode          query     string  false  "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSSuccessResponse
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/{pageId}/contents/{languageCode} [delete]
func (h *CMSLandingPageHandler) HandleDeleteLandingContentByPageId(c *fiber.Ctx) error {
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

	if err := h.Service.DeleteContentByLandingPageId(pageId, language, mode); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to Delete Landing content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete Landing content",
	})
}

// HandleDuplicateLandingPage handles POST request to duplicate a Landing Page.
// @Summary      Duplicate Landing Page
// @Description  Duplicate an existing Landing page using its unique identifier and language code.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId        path      string  true  "Landing Page ID (UUID)"
// @Success      200           {object}  dto.CMSSuccessResponse
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/duplicate/{pageId} [post]
func (h *CMSLandingPageHandler) HandleDuplicateLandingPage(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	LandingPage, err := h.Service.DuplicateLandingPage(pageId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to duplicate page",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully duplicate page",
		"item":    LandingPage,
	})
}

// HandleDuplicateLandingContentToAnotherLanguage handles POST request to duplicate a Landing Content to another language.
// @Summary      Duplicate Landing Content to another language
// @Description  Duplicate an existing Landing content using its unique identifier and
// @Tags         CMS - Landing Pages
// @Accept       json
// @Produce      json
// @Param        contentId        path      string  true  "Landing content ID (UUID)"
// @Param        revision  body  dto.CreateRevisionRequest  true  "Revision payload"
// @Success      200           {object}  dto.CMSLandingContentSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/duplicate/{contentId}/contents [post]
func (h *CMSLandingPageHandler) HandleDuplicateLandingContentToAnotherLanguage(c *fiber.Ctx) error {
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

	landingContent, err := h.Service.DuplicateLandingContentToAnotherLanguage(contentId, &revision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to duplicate content to another language",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully duplicate content",
		"item":    landingContent,
	})
}

// HandleRevertLandingContent handles POST request to revert a Landing Page to a specific revision.
// @Summary      Revert Landing Page Content
// @Description  Revert Landing content using a specific revision ID.
// @Tags         CMS - Landing Pages
// @Accept       json
// @Produce      json
// @Param        revisionId  path  string  true  "Revision ID (UUID)"
// @Param        revision    body  dto.CreateRevisionRequest  true  "Revision payload"
// @Success      200  {object}  dto.CMSLandingContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/{revisionId}/revisions [post]
func (h *CMSLandingPageHandler) HandleRevertLandingContent(c *fiber.Ctx) error {
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

	LandingContent, err := h.Service.RevertLandingContent(revisionId, &revision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to revert Landing content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully revert Landing content",
		"item":    LandingContent,
	})
}

// HandleUpdateLandingContent handles PUT request to update Landing content.
// @Summary      Update Landing Content
// @Description  Update an existing Landing content by its content ID.
// @Description  This operation does not merge with previous content — it overwrites everything.
// @Description  A valid revision must be included in the request.
// @Tags         CMS - Landing Pages
// @Accept       json
// @Produce      json
// @Param        contentId      path  string           true  "Landing Content ID (UUID)"
// @Param        landingContent     body  dto.CreateLandingContentRequest  true  "Updated Landing Content"
// @Success      200  {object}  dto.CMSLandingContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/{contentId}/contents [put]
func (h *CMSLandingPageHandler) HandleUpdateLandingContent(c *fiber.Ctx) error {
	contentIdStr := c.Params("contentId")
	contentId, err := uuid.Parse(contentIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse pageId",
			"error":   err.Error(),
		})
	}

	var updatedContent models.LandingContent
	if err := c.BodyParser(&updatedContent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the body",
			"error":   err.Error(),
		})
	}

	helpers.SanitizeLandingContent(&updatedContent)

	LandingContent, err := h.Service.UpdateLandingContent(&updatedContent, contentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update Landing content",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully update Landing content",
		"item":    LandingContent,
	})
}

// HandleGetCategory handles GET request to retrieve Category for that Landing Page
// @Summary      Get Landing Page Categories
// @Description  Retrieve the categories of a specific Landing page using its unique identifier, category type code, language code, and mode.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId        path      string  true  "Landing Page ID (UUID)"
// @Param        categoryTypeCode  path      string  true  "Category Type Code"
// @Param        languageCode      path      string  true  "Language Code (e.g., en, th)"
// @Param        mode              query     string  false "Mode (e.g., draft, published, histories, preview). Defaults to 'published'."
// @Success      200           {object}  dto.CMSLandingCategoriesSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/category/{categoryTypeCode}/{pageId}/{languageCode} [get]
func (h *CMSLandingPageHandler) HandleGetCategory(c *fiber.Ctx) error {
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

// TODO: Add handler for FindLatestContentByPageId and UpdateLandingContent
// TODO: Test all previous handler
// HandleGetRevisions handles GET request to retrieve revisions of a specific landing page
// @Summary      Get landing Page Revisions
// @Description  Retrieve a list of revisions for a specific landing page based on its unique identifier and language code. Each revision includes metadata such as author, status, and timestamps.
// @Tags         CMS - Landing Pages
// @Produce      json
// @Param        pageId        path      string  true  "landing Page ID (UUID)"
// @Param        languageCode      path      string  true  "Language Code (e.g., en, th)"
// @Success      200           {object}  dto.RevsionsSuccessResponse200
// @Failure      400           {object}  dto.ErrorResponse400
// @Failure      500           {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/revisions/{languageCode}/{pageId} [get]
func (h *CMSLandingPageHandler) HandleGetRevisions(c *fiber.Ctx) error {
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

// HandlePreviewLandingContent handles POST request to preview Landing content.
// @Summary      Preview Landing Content
// @Description  Preview an existing Landing content by its content ID.
// @Description  This operation does not merge with previous content — it overwrites everything.
// @Tags         CMS - Landing Pages
// @Accept       json
// @Produce      json
// @Param        pageId      path  string           true  "Landing Page ID (UUID)"
// @Param        landingContent     body  dto.CreateLandingContentPreviewRequest  true  "Preview Landing Content"
// @Success      200  {object}  dto.CMSLandingContentSuccessResponse200
// @Failure      400  {object}  dto.ErrorResponse400
// @Failure      500  {object}  dto.ErrorResponse500
// @Router       /cms/landingpages/previews/{pageId} [post]
func (h *CMSLandingPageHandler) HandlePreviewLandingContent(c *fiber.Ctx) error {
	pageIdStr := c.Params("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the pageId",
			"error":   err.Error(),
		})
	}

	var landingContentPreview models.LandingContent
	if err := c.BodyParser(&landingContentPreview); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query JSON",
			"error":   err.Error(),
		})
	}	

	helpers.SanitizeLandingContent(&landingContentPreview)

	url, err := h.Service.PreviewLandingContent(pageId, &landingContentPreview)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to preview content",
			"error":   err.Error(),
		})
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully preview landing content",
		"url":    url,
	})
}