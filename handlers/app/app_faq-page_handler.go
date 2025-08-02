package app

import (
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AppFaqPageHandler struct {
	Service services.AppFaqPageServiceInterface
}

func NewAppFaqPageHandler(service services.AppFaqPageServiceInterface) *AppFaqPageHandler {
	return &AppFaqPageHandler{Service: service}
}

// HandleGetFaqPage handles GET requests to retrieve a Faq Page by its UrlAlias
// @Summary      Get Faq Page by UrlAlias
// @Description  Retrieves a Faq Page by its UrlAlias
// @Tags         App - Faq Pages
// @Produce      json
// @Param        languageCode  path  string  true  "Language"
// @Param        url_alias  query  string  true  "Faq Page UrlAlias"
// @Param        select  query     string  false  "Comma-separated list of fields to preload (e.g.revisions, categories, components, metatag). If blank, all fields will be preloaded."
// @Success      200  {object} dto.FaqPageSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /app/faqpages/{languageCode}/by-alias [get]
func (h *AppFaqPageHandler) HandleGetFaqPageByAlias(c *fiber.Ctx) error {
	slug := c.Query("url_alias")
	language := c.Params("languageCode")
	isAlias := true	

	selectParam := c.Query("select")

	faqPage, err := h.Service.GetFaqPage(slug, isAlias, selectParam, language)	

	if err != nil {
		switch err {
			case errs.ErrNotFound:
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Faq page not found",
					"error": err.Error(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get Faq page",
					"error": err.Error(),
				})
		}
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Faq page retrieved successfully",
		"data": faqPage,	
	})	
}

// HandleGetFaqPage handles GET requests to retrieve a Faq Page by its Url
// @Summary      Get Faq Page by Url
// @Description  Retrieves a Faq Page by its Url
// @Tags         App - Faq Pages
// @Produce      json
// @Param        languageCode  path  string  true  "Language"
// @Param        url  query  string  true  "Faq Page Url"
// @Param        select  query     string  false  "Comma-separated list of fields to preload (e.g.revisions, categories, components, metatag). If blank, all fields will be preloaded."
// @Success      200  {object} dto.FaqPageSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /app/faqpages/{languageCode}/by-url [get]
func (h *AppFaqPageHandler) HandleGetFaqPageByUrl(c *fiber.Ctx) error {
	slug := c.Query("url")
	language := c.Params("languageCode")
	isAlias := false	

	selectParam := c.Query("select")

	faqPage, err := h.Service.GetFaqPage(slug, isAlias, selectParam, language)	

	if err != nil {
		switch err {
			case errs.ErrNotFound:
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Faq page not found",
					"error": err.Error(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get Faq page",
					"error": err.Error(),
				})
		}
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Faq page retrieved successfully",
		"data": faqPage,	
	})	
}

// HandleGetFaqContentPreview handles GET requests to retrieve a Faq Content Preview
// @Summary      Get Faq Content by content id
// @Description  Retrieves a Faq Content by its id
// @Tags         App - Faq Pages
// @Produce      json
// @Param        id  path  string  true  "Content Id"
// @Success      200  {object} dto.FaqContentSuccessResponse200
// @Failure 		 404  {object} dto.ErrorResponse404
// @Failure      500  {object} dto.ErrorResponse500 
// @Router       /app/faqpages/previews/{id} [get]
func (h *AppFaqPageHandler) HandleGetFaqContentPreview(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse the id",
			"error":   err.Error(),
		})
	}	

	faqContent, err := h.Service.GetFaqContentPreview(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to preview the content",
			"error":   err.Error(),
		})		
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully get preview content",
		"data": faqContent,
	})	
}