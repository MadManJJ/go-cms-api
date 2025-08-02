package cms

import (
	_ "github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
)

// CMSHandler handles HTTP requests for the CMS domain
type CMSHandler struct {
	Service services.CMSServiceInterface
}

// NewCMSHandler creates a new instance of CMSHandler
func NewCMSHandler(service services.CMSServiceInterface) *CMSHandler {
	return &CMSHandler{
		Service: service,
	}
}

// HandleTest handles the test endpoint
// @Summary      Test cms endpoint
// @Description  Returns test data from CMS repository
// @Tags         CMS
// @Produce      json
// @Success      200  {object} dto.CMSData "Success response"
// @Failure      500  {object} map[string]string "Internal Server Error" 
// @Router       /cms/test [get]
func (h *CMSHandler) HandleTest(c *fiber.Ctx) error {
	data, err := h.Service.GetTestData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get test data",
		})
	}

	return c.JSON(data)
}

// HandleAdditional handles the additional endpoint
func (h *CMSHandler) HandleAdditional(c *fiber.Ctx) error {
	data, err := h.Service.GetAdditionalData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get additional data",
		})
	}

	return c.JSON(data)
}
