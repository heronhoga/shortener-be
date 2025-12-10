package middleware

import (
    "github.com/gofiber/fiber/v2"
    "os"
)

func AppKey() fiber.Handler {
    return func(c *fiber.Ctx) error {
        expectedKey := os.Getenv("APP_KEY")
        providedKey := c.Get("App-Key")

        if expectedKey == "" {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "APP_KEY not configured on server",
            })
        }

        if providedKey != expectedKey {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid or missing App-Key",
            })
        }

        return c.Next()
    }
}
