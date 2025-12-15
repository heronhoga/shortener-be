package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/shortener-be/handler"
	"github.com/heronhoga/shortener-be/middleware"
)

type LinkRoute struct {
	Handler *handler.LinkHandler
}

func NewLinkRoute(handler *handler.LinkHandler) *LinkRoute {
	return &LinkRoute{Handler: handler}
}

func (r *LinkRoute) Register(router fiber.Router) {
	user := router.Group("/links")
	user.Post("/create", middleware.VerifyToken(), r.Handler.CreateShortLink)
}