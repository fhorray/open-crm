package main

import (
	"open-crm/core/controllers"
	"open-crm/core/middlewares"
	"open-crm/core/routes"
	"open-crm/pkg/config"
	"open-crm/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load .env Variables
	config.Load()

	// Connect to DB
	database.Connect()

	app := fiber.New()

	// Middlewares
	middlewares.Middlewares(app)

	// Routes
	routes.RegisterRoutes(app)

	// Serve Scalar API Reference using Fiber
	app.Get("/reference", controllers.ScalarHandler())

	app.Listen(":8787")
}
