package services

import (
	"open-crm/core/models"
	"open-crm/core/repositories"

	"github.com/google/uuid"
)

// CreateAccount
func CreateAccount(payload *models.CreateAccountDTO) error {
	account := &models.Account{
		ID:       uuid.New(),
		UserID:   uuid.MustParse(payload.UserID),
		Password: payload.Password,
	}

	return repositories.CreateAccount(account)
}
