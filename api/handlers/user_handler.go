package handlers

import (
	"open-crm/core/models"
	"open-crm/core/services"
	"open-crm/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lucsky/cuid"
)

var validate = validator.New()

// Create User
func CreateUser(c *fiber.Ctx) error {
	var payload models.CreateUserDTO

	// Verify body JSON
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Validate Data
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)

		for _, e := range validationErrors {
			errors[e.Field()] = "Invalid field: " + e.Tag()
		}

		return utils.SendResponse(c, 422, utils.APIResponse{
			Data: errors,
		})
	}

	// Map to real user
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return c.Status(500).JSON("Error trying to create user")
	}
	user := models.User{
		ID:       cuid.New(),
		Email:    payload.Email,
		Name:     payload.Name,
		Password: hashedPassword,
	}

	if err := services.CreateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error trying to create user."})
	}

	return c.Status(201).JSON(user)
}

// Get All Users
func GetAllUsers(c *fiber.Ctx) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error trying to get all users"})
	}

	return c.JSON(users)
}

// Get User By Id
func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user, err := services.GetUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// Delete User
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user, err := services.DeleteUser(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Error trying to delete user"})
	}

	return c.JSON(user)
}
