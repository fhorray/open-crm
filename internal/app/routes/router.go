package routes

import (
	"open-crm/internal/app/handlers"
	"open-crm/internal/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Auth
func registerAuthRoutes(router fiber.Router) {
	r := router.Group("/auth")

	r.Get("/get-session", handlers.GetSession)
	r.Post("/sign-in", handlers.Login)
	r.Post("/sign-up", handlers.Register)
	r.Post("/sign-out", handlers.Logout)
	r.Post("/refresh", handlers.RefreshToken)
}

// Users
func registerUserRoutes(router fiber.Router) {
	r := router.Group("/users")

	r.Post("/", middlewares.CheckRoles("superadmin,admin"), handlers.CreateUser)
	r.Get("/", middlewares.CheckRoles("superadmin,admin"), handlers.GetAllUsers)
	r.Get("/:id", middlewares.CheckRoles("superadmin,admin"), handlers.GetUserById)
	r.Patch("/:id", middlewares.CheckRoles("superadmin,admin,user"), handlers.UpdateUser)
	r.Delete("/:id", middlewares.CheckRoles("superadmin,admin"), handlers.DeleteUser)
}

// Organizations
func registerOrganizationsRoutes(router fiber.Router) {
	r := router.Group("/orgs")

	r.Post("/", middlewares.CheckRoles("superadmin,admin"), handlers.CreateOrganization)
	r.Get("/:id", middlewares.CheckRoles("superadmin,admin,user,owner"), handlers.GetOrganizationByID)
}

// Relationshp Routes
func registerRelationshipRoutes(router fiber.Router) {
	// Rota: /users/:id/orgs
	router.Get("/users/:id/orgs", middlewares.CheckRoles("superadmin,admin,user"), handlers.GetOrgsByUserID)

	// Rota: /orgs/:id/users
	router.Get("/orgs/:id/users", middlewares.CheckRoles("superadmin,admin,user"), handlers.GetUsersByOrgID)
}

// Function to Register Routes
func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	// Register Routes
	registerUserRoutes(v1)
	registerAuthRoutes(v1)

	registerOrganizationsRoutes(v1)

	registerRelationshipRoutes(v1)

	// Static Files
	// v1.Static("/docs/openapi.yaml", "./docs/openapi.yaml")
}
