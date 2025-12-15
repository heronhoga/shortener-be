package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/heronhoga/shortener-be/config"
	"github.com/heronhoga/shortener-be/handler"
	"github.com/heronhoga/shortener-be/middleware"
	"github.com/heronhoga/shortener-be/repository"
	"github.com/heronhoga/shortener-be/route"
	"github.com/heronhoga/shortener-be/service"
	"github.com/heronhoga/shortener-be/util"
)

func main() {
	// load env
	util.LoadEnv()

	// create fiber app
	app := fiber.New()

	// db connect
	db := config.ConnectDB()
	defer db.Close(context.Background())

	// front-end app
	frontEndApp := os.Getenv("FRONTEND_APP")

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontEndApp,
		AllowHeaders:     "Origin, Content-Type, Accept, App-Key",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowCredentials: true,
	}))

	//api route
	api := app.Group("/api/v1", middleware.JSONOnly(), middleware.AppKey())
	
	// dependencies
	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoute := route.NewUserRoute(userHandler)
	userRoute.Register(api)

	// link
	linkRepo := repository.NewLinkRepository(db)
	linkService := service.NewLinkService(linkRepo)
	linkHandler := handler.NewLinkHandler(linkService)
	linkRoute := route.NewLinkRoute(linkHandler)
	linkRoute.Register(api)

	// listen
	app.Listen(":8000")
}