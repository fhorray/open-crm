package repositories

import (
	"open-crm/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Declaring interface
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetAll() (*[]models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(id uuid.UUID, values *models.User) (*models.User, error)
	Delete(id uuid.UUID) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// Functions to create a Repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create User
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Gell All Users
func (r *userRepository) GetAll() (*[]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return &users, err
}

// Get User by ID
func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, "id = ?", id).Error

	return &user, err
}

// Get User by Email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update User
func (r *userRepository) Update(id uuid.UUID, values *models.User) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return &user, err
	}

	// Update user inside DB
	if err := r.db.Model(&user).Where("id = ?", id).Updates(values).Error; err != nil {
		return &user, err
	}

	return &user, nil
}

// Delete User by ID
func (r *userRepository) Delete(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}
