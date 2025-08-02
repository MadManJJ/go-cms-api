package app

import (
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AppPartnerPageHandler struct {
	Service services.AppPartnerPageServiceInterface
}

func NewAppPartnerPageHandler(service services.AppPartnerPageServiceInterface) *AppPartnerPageHandler {
	return &AppPartnerPageHandler{Service: service}
}

// HandleGetPartnerPageByAlias handles GET requests to retrieve a Partner Page by its UrlAlias
// @Summary      Get Partner Page by UrlAlias
// @Description  Retrieves a Partner Page by its UrlAlias
// @Tags         App - Partner Pages
// @Produce      json
// @Param        languageCode  path  string  true  "Language"
// @Param        url_alias  query  string  true  "Partner Page UrlAlias"
// @Param        select  query     string  false  "Comma-separated list of fields to preload (e.g. page, files, components, revisions, metatag)"
// @Success      200  {object} dto.PartnerPageSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /app/partnerpages/{languageCode}/by-alias [get]
func (h *AppPartnerPageHandler) HandleGetPartnerPageByAlias(c *fiber.Ctx) error {
	slug := c.Query("url_alias")
	language := c.Params("languageCode")
	isAlias := true	

	selectParam := c.Query("select")

	partnerPage, err := h.Service.GetPartnerPage(slug, isAlias, selectParam, language)

	if err != nil {
		switch err {
		case errs.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Partner page not found",
				"error":   err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get Partner page",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Partner page retrieved successfully",
		"data":    partnerPage,
	})
}

// HandleGetPartnerPageByUrl handles GET requests to retrieve a Partner Page by its UrlAlias
// @Summary      Get Partner Page by UrlAlias
// @Description  Retrieves a Partner Page by its UrlAlias
// @Tags         App - Partner Pages
// @Produce      json
// @Param        languageCode  path  string  true  "Language"
// @Param        url  query  string  true  "Partner Page Url"
// @Param        select  query     string  false  "Comma-separated list of fields to preload (e.g. page, files, components, revisions, metatag)"
// @Success      200  {object} dto.PartnerPageSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /app/partnerpages/{languageCode}/by-url [get]
func (h *AppPartnerPageHandler) HandleGetPartnerPageByUrl(c *fiber.Ctx) error {
	slug := c.Query("url")
	language := c.Params("languageCode")
	isAlias := false	

	selectParam := c.Query("select")

	partnerPage, err := h.Service.GetPartnerPage(slug, isAlias, selectParam, language)

	if err != nil {
		switch err {
		case errs.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Partner page not found",
				"error":   err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get Partner page",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Partner page retrieved successfully",
		"data":    partnerPage,
	})
}

// HandleGetPartnerContentPreview handles GET requests to retrieve a Partner Content Preview
// @Summary      Get Partner Content by content id
// @Description  Retrieves a Partner Content by its id
// @Tags         App - Partner Pages
// @Produce      json
// @Param        id  path  string  true  "Content Id"
// @Success      200  {object} dto.PartnerContentSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /app/partnerpages/previews/{id} [get]
func (h *AppPartnerPageHandler) HandleGetPartnerContentPreview(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the id",
			"error":   err.Error(),
		})
	}	

	partnerContent, err := h.Service.GetPartnerContentPreview(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to preview the content",
			"error":   err.Error(),
		})		
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get preview content",
		"data": partnerContent,
	})	
}