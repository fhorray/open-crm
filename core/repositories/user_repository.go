package repositories

import (
	"open-crm/core/models"
	"open-crm/pkg/db"
)

// Create User
func CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

// Gell All Users
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Find(&users).Error
	return users, err
}

// Get User by ID
func GetUserById(id string) (models.User, error) {
	var user models.User
	err := db.DB.First(&user, id).Error
	return user, err
}

// Delete User by ID
func DeleteUser(id string) (models.User, error) {
	var user models.User

	// Get user
	if err := db.DB.First(&user, "id =?", id).Error; err != nil {
		return user, err
	}

	// Delete user from db
	if err := db.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
