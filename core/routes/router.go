package routes

import (
	"github.com/gofiber/fiber/v2"
)

// Function to Register Routes
func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	// Register Routes
	registerUserRoutes(v1)
	registerAuthRoutes(v1)
}
