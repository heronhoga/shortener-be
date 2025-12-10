package handler

import (
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

	// call service layer
	if err := h.service.RegisterNewUser(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "ok",
	})
}
