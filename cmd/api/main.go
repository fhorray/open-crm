package main

import (
	"open-crm/config"
	"open-crm/internal/app/handlers"
	"open-crm/internal/app/middlewares"
	"open-crm/internal/app/routes"

	"open-crm/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load .env Variables
	config.Load()

	// Connect to DB
	database.Connect()

	app := fiber.New()

	// Serve Scalar API Reference using Fiber
	app.Get("/docs", handlers.ScalarHandler())

	// Middlewares
	middlewares.Middlewares(app)

	// Routes
	routes.RegisterRoutes(app)

	app.Listen(":8787")
}
