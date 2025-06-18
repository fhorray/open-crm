package routes

import (
	"open-crm/core/controllers"
	"open-crm/core/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Auth
func registerAuthRoutes(router fiber.Router) {
	r := router.Group("/auth")

	r.Get("/get-session", controllers.GetSession)
	r.Post("/sign-in", controllers.SignIn)
	r.Post("/sign-up", controllers.SignUp)
	r.Post("/sign-out", controllers.SignOut)
}

// Users
func registerUserRoutes(router fiber.Router) {
	r := router.Group("/users")

	r.Post("/", middlewares.CheckRoles("superadmin,admin"), controllers.CreateUser)
	r.Get("/", middlewares.CheckRoles("superadmin,admin"), controllers.GetAllUsers)
	r.Get("/:id", middlewares.CheckRoles("superadmin,admin"), controllers.GetUserById)
	r.Patch("/:id", middlewares.CheckRoles("superadmin,admin,user"), controllers.UpdateUser)
	r.Delete("/:id", middlewares.CheckRoles("superadmin,admin"), controllers.DeleteUser)
}

// Organizations
func registerOrganizationsRoutes(router fiber.Router) {
	r := router.Group("/orgs")

	r.Post("/", middlewares.CheckRoles("superadmin,admin"), controllers.CreateOrganization)
	r.Get("/:id", middlewares.CheckRoles("superadmin,admin,user,owner"), controllers.GetOrganizationByID)
}

// Relationshp Routes
func registerRelationshipRoutes(router fiber.Router) {
	// Rota: /users/:id/orgs
	router.Get("/users/:id/orgs", middlewares.CheckRoles("superadmin,admin,user"), controllers.GetOrgsByUserID)

	// Rota: /orgs/:id/users
	router.Get("/orgs/:id/users", middlewares.CheckRoles("superadmin,admin,user"), controllers.GetUsersByOrgID)
}

// Function to Register Routes
func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	// Register Routes
	registerUserRoutes(v1)
	registerAuthRoutes(v1)
	registerOrganizationsRoutes(v1)

	registerRelationshipRoutes(v1)
}
