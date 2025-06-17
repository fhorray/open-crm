package routes

import (
	"open-crm/core/controllers"
	"open-crm/core/middlewares"

	"github.com/gofiber/fiber/v2"
)

// User Router
func registerUserRoutes(router fiber.Router) {
	r := router.Group("/users")

	r.Post("/", middlewares.CheckRoles("superadmin,admin"), controllers.CreateUser)
	r.Get("/", middlewares.CheckRoles("superadmin,admin"), controllers.GetAllUsers)
	r.Get("/:id", middlewares.CheckRoles("superadmin,admin"), controllers.GetUserById)
	r.Patch("/:id", middlewares.CheckRoles("superadmin,admin,user"), controllers.UpdateUser)
	r.Delete("/:id", middlewares.CheckRoles("superadmin,admin"), controllers.DeleteUser)
}
