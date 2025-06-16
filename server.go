package main

import (
	"open-crm/api"
	"open-crm/pkg/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Conect to db
	db.Connect()

	app := fiber.New()

	// Routing
	api.RegisterRoutes(app)

	app.Listen(":8787")
}
