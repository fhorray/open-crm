package handlers

import (
	"net/http"
	"open-crm/internal/app/models"
	"open-crm/internal/app/services"
	"open-crm/pkg/utils"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

// Create User
func CreateUser(c *fiber.Ctx) error {
	var payload models.CreateUserDTO

	// Parse body
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid body",
		})
	}

	// Validate payload
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = "Invalid field: " + e.Tag()
		}

		return utils.SendResponse(c, utils.APIResponse{
			Status: http.StatusBadRequest,
			Data:   errors,
		})
	}

	// Call service to create user
	user, err := services.CreateUser(payload)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Message: "Error creating user",
			Status:  http.StatusInternalServerError,
		})
	}

	// response without password
	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusCreated,
		Data: models.User{
			Name:  user.Name,
			Email: user.Email,
			Image: user.Image,
			Roles: user.Roles,
		},
	})
}

// Get All Users
func GetAllUsers(c *fiber.Ctx) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error trying to get all users",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   users,
	})
}

// Get User By Id
func GetUserById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(strings.TrimSpace(idStr))
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	user, err := services.GetUserById(id)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   user,
	})
}

// Update User
func UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(strings.TrimSpace(idStr))
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	var payload models.User
	// Parse body
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid body",
		})
	}

	// Validate payload
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = "Invalid field: " + e.Tag()
		}

		return utils.SendResponse(c, utils.APIResponse{
			Status: http.StatusBadRequest,
			Data:   errors,
		})
	}

	user, err := services.UpdateUser(id, &payload)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   user,
	})
}

// Delete User
func DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(strings.TrimSpace(idStr))
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	user, err := services.DeleteUser(id)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Error trying to delete user",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   user,
	})
}
