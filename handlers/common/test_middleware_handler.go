package common

import "github.com/gofiber/fiber/v2"

type TestMiddlewareHandler struct{}

func NewTestMiddlewareHandler() *TestMiddlewareHandler {
	return &TestMiddlewareHandler{}
}

// HandleTestMiddleware test
// @Summary      Test middleware
// @Description  Test middleware
// @Tags         test
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} dto.SuccessResponse
// @Failure      401  {object} dto.ErrorResponse500
// @Failure      500  {object} dto.ErrorResponse500
// @Router       /middleware/test [get]
func (h *TestMiddlewareHandler) HandleTestMiddleware(c *fiber.Ctx) error {
	// This is a test middleware handler that simply returns a success message
	return c.JSON(fiber.Map{
		"message": "Test middleware executed successfully",
	})
}