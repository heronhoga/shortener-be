package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/shortener-be/handler"
	"github.com/heronhoga/shortener-be/middleware"
)

type UserRoute struct {
	Handler *handler.UserHandler
}

func NewUserRoute(handler *handler.UserHandler) *UserRoute {
	return &UserRoute{Handler: handler}
}

func (r *UserRoute) Register(router fiber.Router) {
	user := router.Group("/users").Use(middleware.AppKey())
	user.Post("/register", r.Handler.RegisterNewUser)
	user.Post("/login", r.Handler.LoginUser)
	user.Post("/logout", middleware.VerifyToken(), r.Handler.LogoutUser)

	// check existing access token - middleware
	user.Get("/me", middleware.VerifyToken(), r.Handler.Me)
}