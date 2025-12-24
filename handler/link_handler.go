package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/shortener-be/model"
	"github.com/heronhoga/shortener-be/service"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{service: linkService}
}

func (h *LinkHandler) CreateShortLink(c *fiber.Ctx) error {
	var req model.CreateLink

	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// request fields
	if req.Url == "" || req.Label == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Url field can't be empty",
		})
	}

	// get user id
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	err := h.service.CreateShortLink(c.Context(), &req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

func (h *LinkHandler) EditShortLink(c *fiber.Ctx) error {
	var req model.EditLink

	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// reqyest fields
	if req.ID == "" || req.Name == "" || req.Url == "" || req.Label == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// get user id
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// update existing link
	err := h.service.EditShortLink(c.Context(), &req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

func (h *LinkHandler) GetShortLink(c *fiber.Ctx) error {
	var req model.GetLink

	req.LinkID = c.Query("linkid")
	req.Page = c.QueryInt("page")

	// get user id
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	
    links, err := h.service.GetShortLinks(c.Context(), &req, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"data": links,
	})
}

func (h *LinkHandler) RedirectLink(c *fiber.Ctx) error {
	linkName := c.Params("name")

	url, err := h.service.RedirectLink(c.Context(), linkName)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "link not found")
	}

	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}


func (h *LinkHandler) DeleteLink(c *fiber.Ctx) error {
	var req model.DeleteLink

	// parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// get user id
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	err := h.service.DeleteLink(c.Context(), req.LinkID, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})

}