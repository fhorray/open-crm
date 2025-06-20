package handlers

import (
	"net/http"
	"open-crm/internal/app/models"
	"open-crm/internal/app/repositories"
	"open-crm/internal/app/services"
	"open-crm/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Create Organization
func CreateOrganization(c *fiber.Ctx) error {
	var payload models.CreateOrganizationDTO

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

	// Call service toc create Organization
	org, err := services.CreateOrganization(&payload)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Message: "Error creating organization",
			Status:  http.StatusInternalServerError,
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusCreated,
		Data:   org,
	})
}

func GetOrganizationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	org, err := repositories.GetOrganizationByID(id)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusNotFound,
			Message: "Organization not found",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   org,
	})
}

// GET USER ORGANIZATION
func GetOrgsByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")

	org, err := services.GetOrgByUserID(userID)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Erro ao buscar organização do usuário",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: fiber.StatusOK,
		Data:   org,
	})
}

// GET ORGANIZATION's USERS
func GetUsersByOrgID(c *fiber.Ctx) error {
	orgID := c.Params("id")

	users, err := services.GetUsersByOrgID(orgID)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Erro ao buscar usuários da organização",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: fiber.StatusOK,
		Data:   users,
	})
}
