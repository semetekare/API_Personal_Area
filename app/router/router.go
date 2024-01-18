package router

import (
	Handlers "api_sotr/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// /
func RegisterHTTPEndpoints(router fiber.Router) {
	router.Post("/login", Handlers.Login)
	router.Get("/Ok", Handlers.Ok)

}
