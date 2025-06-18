package repositories

import (
	"open-crm/core/models"
	"open-crm/pkg/database"

	"github.com/google/uuid"
)

func CreateOrganization(org *models.Organization) error {
	return database.DB.Create(org).Error
}

func GetOrganizationByID(id string) (*models.Organization, error) {
	var org models.Organization
	err := database.DB.Preload("Users").Where("id = ?", id).First(&org).Error
	return &org, err
}

func ListOrganizations() ([]models.Organization, error) {
	var orgs []models.Organization
	err := database.DB.Find(&orgs).Error
	return orgs, err
}

// -- Invitations
func CreateInvitation(payload *models.CreateInvitationDTO) error {
	code := uuid.New().String()

	invitation := &models.Invitation{
		ID:             uuid.New().String(),
		OrganizationID: payload.OrganizationID,
		InvitedBy:      payload.InvitedBy,
		InvitedEmail:   payload.InvitedEmail,
		Code:           code,
		ExpiresAt:      payload.ExpiresAt,
	}

	return database.DB.Create(invitation).Error
}
