package helpers

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

func ParseJWTWithKey(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	// Check expiration here
	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims format")
	}

	return claims, nil
}