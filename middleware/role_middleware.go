package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// CheckRoleMiddleware checks if the user has one of the allowed roles.
// Use after CheckTokenMiddleware to ensure the user is authenticated.
func CheckRoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(jwt.MapClaims)
		fmt.Println("User in CheckRoleMiddleware:", user)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		role, exists := user["role"].(string)
		if !exists {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role not found",
			})
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions",
		})
	}
}