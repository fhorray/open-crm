package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Logger adiciona middlewares globais ao app
func Middlewares(app *fiber.App) {
	// Recover panics
	app.Use(recover.New())

	// Request Loggers
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// Cors Middleware
	app.Use(Cors())

	// Auth Middleware
	app.Use(AuthMiddleware())

}
