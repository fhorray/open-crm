package repositories

import (
	"open-crm/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(acc *models.Account) (*models.Account, error)
	FindByUserAndProvider(userID uuid.UUID, provider string) (*models.Account, error)
	DeleteAccount(id uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}

// Constructor
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

// Create Account
func (r *accountRepository) Create(acc *models.Account) (*models.Account, error) {
	if err := r.db.Create(acc).Error; err != nil {
		return nil, err
	}
	return acc, nil
}

// Find account by user ID and provider
func (r *accountRepository) FindByUserAndProvider(userID uuid.UUID, provider string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("user_id = ? AND provider_id = ?", userID, provider).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Delete Account
func (r *accountRepository) DeleteAccount(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Account{}).Error
}
