package services

import (
	"open-crm/internal/app/models"
	"open-crm/internal/app/repositories"
	"open-crm/pkg/database"

	"github.com/google/uuid"
)

// CreateAccount
func CreateAccount(payload *models.CreateAccountDTO) error {
	accRepo := repositories.NewAccountRepository(database.DB)

	account := &models.Account{
		UserID:   uuid.MustParse(payload.UserID),
		Password: payload.Password,
	}

	_, err := accRepo.Create(account)
	return err
}
