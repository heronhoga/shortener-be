package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/shortener-be/handler"
)

type UserRoute struct {
	Handler *handler.UserHandler
}

func NewUserRoute(handler *handler.UserHandler) *UserRoute {
	return &UserRoute{Handler: handler}
}

func (r *UserRoute) Register(router fiber.Router) {
	user := router.Group("/users")
	user.Post("/register", r.Handler.RegisterNewUser)
	user.Post("/login", r.Handler.LoginUser)
}