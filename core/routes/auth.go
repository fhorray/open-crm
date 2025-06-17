package routes

import (
	"open-crm/core/controllers"

	"github.com/gofiber/fiber/v2"
)

// Auth Router
func registerAuthRoutes(router fiber.Router) {
	r := router.Group("/auth")

	r.Get("/get-session", controllers.GetSession)
	r.Post("/sign-in", controllers.SignIn)
	r.Post("/sign-up", controllers.SignUp)
	r.Post("/sign-out", controllers.SignOut)
}
