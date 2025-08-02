package app

import (
	_ "github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/gofiber/fiber/v2"
)

// AppHandler handles HTTP requests for the app domain
type AppHandler struct {
	Service services.AppServiceInterface
}

// NewAppHandler creates a new instance of AppHandler
func NewAppHandler(service services.AppServiceInterface) *AppHandler {
	return &AppHandler{
		Service: service,
	}
}

// HandleTest handles the test endpoint
// @Summary      Test app endpoint
// @Description  Returns test data from App repository
// @Tags         App
// @Produce      json
// @Success      200  {object} dto.AppData "Success response"
// @Failure      500  {object} map[string]string "Internal Server Error"
// @Router       /app/test [get]
func (h *AppHandler) HandleTest(c *fiber.Ctx) error {
	data, err := h.Service.GetTestData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get test data",
		})
	}

	return c.JSON(data)
}

// HandleAdditional handles the additional endpoint
func (h *AppHandler) HandleAdditional(c *fiber.Ctx) error {
	data, err := h.Service.GetAdditionalData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get additional data",
		})
	}

	return c.JSON(data)
}
