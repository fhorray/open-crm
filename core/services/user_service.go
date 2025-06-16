package services

import (
	"open-crm/core/models"
	"open-crm/core/repositories"
)

// Create User
func CreateUser(user *models.User) error {
	return repositories.CreateUser(user)
}

// Get All Users
func GetAllUsers() ([]models.User, error) {
	return repositories.GetAllUsers()
}

// Get User By Id
func GetUserById(id string) (models.User, error) {
	return repositories.GetUserById(id)
}

// Delete User
func DeleteUser(id string) (models.User, error) {
	return repositories.DeleteUser(id)
}
