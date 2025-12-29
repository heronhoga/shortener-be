package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return fiber.ErrUnauthorized
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return fiber.ErrInternalServerError
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return fiber.ErrUnauthorized
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}


