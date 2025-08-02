package middleware

import "github.com/gofiber/fiber/v2"

func HPP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Context().QueryArgs()
		duplicateFound := false

		raw.VisitAll(func(key, val []byte) {
			if len(raw.PeekMulti(string(key))) > 1 {
				duplicateFound = true
			}
		})

		if duplicateFound {
			return c.Status(fiber.StatusBadRequest).SendString("HTTP Parameter Pollution detected")
		}

		return c.Next()
	}
}
