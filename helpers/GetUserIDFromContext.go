package helpers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	claimsRaw := c.Locals("user")
	if claimsRaw == nil {
		return uuid.Nil, fmt.Errorf("no user info found in token")
	}

	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims format")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token payload")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}
