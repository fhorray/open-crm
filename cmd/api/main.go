package main

import (
	"open-crm/config"
	"open-crm/core/routes"
	"open-crm/internal/app/handlers"
	"open-crm/internal/app/middlewares"

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
	app.Get("/reference", handlers.ScalarHandler())

	app.Listen(":8787")
}
