package common

import "github.com/gofiber/fiber/v2"

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HandleHealthCheck handles the health check endpoint
func (h *HealthHandler) HandleHealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
