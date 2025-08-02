package app

import (
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AppLandingPageHandler struct {
	Service services.AppLandingPageServiceInterface
}

func NewAppLandingPageHandler(service services.AppLandingPageServiceInterface) *AppLandingPageHandler {
	return &AppLandingPageHandler{Service: service}
}

// HandleGetLandingPageByUrlAlias handles GET requests to retrieve a Landing Page by its UrlAlias
// @Summary      Get Landing Page by UrlAlias
// @Description  Retrieves a Landing Page by its UrlAlias
// @Tags         App - Landing Pages
// @Produce      json
// @Param        languageCode  path  string  true  "Language"
// @Param        url_alias  query  string  true  "Landing Page UrlAlias"
// @Param        select  query     string  false  "Comma-separated list of fields to preload (e.g. files, revisions, categories, components, metatag). If blank, all fields will be preloaded."
// @Success      200  {object} dto.LandingPageSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /app/landingpages/{languageCode}/by-alias [get]
func(h *AppLandingPageHandler) HandleGetLandingPageByUrlAlias(c *fiber.Ctx) error {
	urlAlias := c.Query("url_alias")
	language := c.Params("languageCode")
	selectParam := c.Query("select")

	landingPage, err := h.Service.GetLandingPageByUrlAlias(urlAlias, selectParam, language)

	if err != nil {
		switch err {
			case errs.ErrNotFound:
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Landing page not found",
					"error": err.Error(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get Landing page",
					"error": err.Error(),
				})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Landing page retrieved successfully",
		"data":    landingPage,	
	})
}

// HandleGetLandingContentPreview handles GET requests to retrieve a Landing Content Preview
// @Summary      Get Landing Content by content id
// @Description  Retrieves a Landing Content by its id
// @Tags         App - Landing Pages
// @Produce      json
// @Param        id  path  string  true  "Content Id"
// @Success      200  {object} dto.LandingContentSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /app/landingpages/previews/{id} [get]
func (h *AppLandingPageHandler) HandleGetLandingContentPreview(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the id",
			"error":   err.Error(),
		})
	}	

	landingContent, err := h.Service.GetLandingContentPreview(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to preview the content",
			"error":   err.Error(),
		})		
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get preview content",
		"data": landingContent,
	})	
}