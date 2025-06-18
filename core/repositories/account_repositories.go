package repositories

import (
	"open-crm/core/models"
	"open-crm/pkg/database"
)

func CreateAccount(acc *models.Account) error {
	return database.DB.Create(acc).Error
}

func GetAccountByID(id string) (*models.Account, error) {
	var acc models.Account
	err := database.DB.Where("id = ?", id).First(&acc).Error
	return &acc, err
}

func ListAccountsByUser(userID string) ([]models.Account, error) {
	var accounts []models.Account
	err := database.DB.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func GetAccountByUserAndProvider(userID string, providerID string) (*models.Account, error) {
	var acc models.Account
	err := database.DB.
		Where("user_id = ? AND provider_id = ?", userID, providerID).
		First(&acc).Error

	return &acc, err
}
