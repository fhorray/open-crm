package main

import (
	"open-crm/core/middlewares"
	"open-crm/core/routes"
	"open-crm/pkg/config"
	"open-crm/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Load .env Variables
	config.Load()

	// Conect to db
	database.Connect()

	app := fiber.New()

	// Middlewares
	middlewares.Middlewares(app)

	// Routing
	routes.RegisterRoutes(app)

	app.Listen(":8787")
}
