package api

import (
	"open-crm/api/handlers"

	"github.com/gofiber/fiber/v2"
)

// User Router
func registerUserRoutes(router fiber.Router) {
	user := router.Group("/users")
	user.Post("/", handlers.CreateUser)
	user.Get("/", handlers.GetAllUsers)
	user.Get("/:id", handlers.GetUserById)
	user.Delete("/:id", handlers.DeleteUser)
}

// Function to Register Routes
func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	// Register Routes
	registerUserRoutes(v1)
}
