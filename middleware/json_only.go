package middleware

import (
    "github.com/gofiber/fiber/v2"
)

func JSONOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        method := c.Method()

        if method == fiber.MethodPost || 
           method == fiber.MethodPut || 
           method == fiber.MethodPatch {

            contentType := c.Get("Content-Type")
            if contentType != "application/json" {
                return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
                    "error": "Content-Type must be application/json",
                })
            }
        }

        accept := c.Get("Accept")
        if accept != "" && accept != "application/json" && accept != "*/*" {
            return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
                "error": "Only application/json responses are supported",
            })
        }

        return c.Next()
    }
}
