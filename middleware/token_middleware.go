package middleware

import (
	"fmt"
	"strings"

	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/gofiber/fiber/v2"
)

// CheckTokenMiddleware checks if the request has a valid JWT token in the Authorization header.
// Use before CheckRoleMiddleware to ensure the user is authenticated.
func CheckAnyTokenMiddleware(lineKey, normalKey string, repo repositories.CMSAuthRepositoryInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			fmt.Println("Missing or invalid Authorization header:", token)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		actualToken := strings.TrimPrefix(token, "Bearer ")

		// Try LineKey first
		claims, err := helpers.ParseJWTWithKey(actualToken, lineKey)
		if err == nil {
			id := helpers.UUIDFromSub(claims["sub"].(string))
			user, err := repo.FindUserById(id)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			claims["user_id"] = user.ID.String()
			c.Locals("user", claims)
			c.Locals("token_type", "line")
			return c.Next()
		}
		fmt.Println("line login failed: ", err)

		// Try NormalKey if LineKey fails
		claims, err = helpers.ParseJWTWithKey(actualToken, normalKey)
		if err == nil {
			c.Locals("user", claims)
			c.Locals("token_type", "normal")
			return c.Next()
		}
		fmt.Println("normal login failed: ", err)

		// Both failed
		fmt.Println("JWT validation failed for all keys")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
}
