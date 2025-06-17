package repositories

import (
	"open-crm/core/models"
	"open-crm/pkg/database"
)

// Create User
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// Gell All Users
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}

// Get User by ID
func GetUserById(id string) (models.User, error) {
	var user models.User

	err := database.DB.First(&user, "id = ?", id).Error
	return user, err
}

// Delete User by ID
func DeleteUser(id string) (models.User, error) {
	var user models.User

	// Get user
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return user, err
	}

	// Delete user from db
	if err := database.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
