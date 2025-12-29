package handler

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/shortener-be/model"
	"github.com/heronhoga/shortener-be/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

func (h *UserHandler) RegisterNewUser(c *fiber.Ctx) error {
	var req model.RegisterUser

	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// check if request part is nil
	if req.Email == "" || req.Username == "" || req.Password == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// call service layer
	if err := h.service.RegisterNewUser(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
    var req model.LoginUser

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if req.Provider == "" {
        req.Provider = "local"
    }

    if req.Provider == "local" {
        if req.Email == "" || req.Password == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Email and password are required",
            })
        }
    }

    if req.Provider == "google" {
        if req.Token == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Google token is required",
            })
        }
    }

    token, err := h.service.LoginUser(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    //get cookie configuration
    cookieSecurity, err := strconv.ParseBool(os.Getenv("COOKIE_SECURITY"))

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Error parsing boolean value",
        })
    }

    //set httpOnly token
    c.Cookie(&fiber.Cookie{
        Name: "access_token",
        Value: token,
        HTTPOnly: true,
        Secure: cookieSecurity,
        SameSite: fiber.CookieSameSiteLaxMode,
        Path: "/",
        MaxAge: 60 * 60 * 24,
    })

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "ok",
    })
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
    userID := c.Locals("user_id")
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "user_id": userID,
    })
}
