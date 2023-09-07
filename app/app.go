package app

import (
	"net/http"

	"api_sotr/app/models"
	"api_sotr/app/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type App struct {
	httpServer *http.Server
}

func (a *App) Run(port string) {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "AuthorizationService",
		AppName:       "AuthorizationService IrGUPS v 1.0",
	})
	models.ConnectDatabase()
	app.Use(logger.New(), recover.New(), cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Route(
		"/",
		router.RegisterHTTPEndpoints,
		"abit.",
	)
	log.Fatal(app.Listen(":" + port))
}
