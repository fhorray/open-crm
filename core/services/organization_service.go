package services

import (
	"open-crm/core/models"
	"open-crm/core/repositories"
	"open-crm/pkg/database"

	"github.com/google/uuid"
)

func CreateOrganization(payload *models.CreateOrganizationDTO) (models.Organization, error) {
	org := &models.Organization{
		ID:       uuid.New(),
		Name:     payload.Name,
		Domain:   payload.Domain,
		IsActive: false,
	}

	err := repositories.CreateOrganization(org)
	return *org, err
}

// -- Invitations
func CreateInvitation(payload *models.CreateInvitationDTO) error {

	invitation := &models.CreateInvitationDTO{
		OrganizationID: payload.OrganizationID,
		InvitedBy:      payload.InvitedBy,
		InvitedEmail:   payload.InvitedEmail,
		ExpiresAt:      payload.ExpiresAt,
	}

	return repositories.CreateInvitation(invitation)
}

// Buscar usuários de uma organização
func GetUsersByOrgID(orgID string) ([]models.UserResponseDTO, error) {
	var users []models.UserResponseDTO

	err := database.DB.
		Table("core.users").
		Where("organization_id = ?", orgID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// Buscar organização de um usuário (inverso)
func GetOrgByUserID(userID string) (*models.OrganizationResponseDTO, error) {
	var user models.User
	err := database.DB.
		Table("core.users").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	var org models.OrganizationResponseDTO
	err = database.DB.
		Table("core.organizations").
		Where("id = ?", user.OrganizationID).
		First(&org).Error

	if err != nil {
		return nil, err
	}

	return &org, nil
}
