package repositories

import (
	"open-crm/internal/app/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	// Usuário e conta
	Register(user *models.User, account *models.Account) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindAccountByUserAndProvider(userID string, providerID string) (*models.Account, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// Auth
func (r *authRepository) Register(user *models.User, account *models.Account) (*models.User, error) {
	tx := r.db.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	account.UserID = user.ID
	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, tx.Commit().Error
}

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindAccountByUserAndProvider(userID string, providerID string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("user_id = ? AND provider_id = ?", userID, providerID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
