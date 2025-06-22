package repositories

import (
	"open-crm/internal/app/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(user *models.User, account *models.Account) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// Register Repository
func (r *authRepository) Register(user *models.User, account *models.Account) (*models.User, error) {
	tx := r.db.Begin()

	// Create user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Account
	account.UserID = user.ID
	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, tx.Commit().Error
}
