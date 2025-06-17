package services

import (
	"open-crm/core/models"
	"open-crm/core/repositories"
	"open-crm/utils"

	"github.com/lucsky/cuid"
)

// Create User
func CreateUser(payload models.CreateUserDTO) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       cuid.New(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Roles:    payload.Roles,
	}

	if err := repositories.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Get All Users
func GetAllUsers() ([]models.User, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return []models.User{}, err
	}

	// Montar array de resposta
	var usersData []models.User
	for _, user := range users {
		usersData = append(usersData, models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Roles: user.Roles,
			Image: user.Image,
		})
	}

	return usersData, nil

}

// Get User By Id
func GetUserById(id string) (models.User, error) {
	user, err := repositories.GetUserById(id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Get User By Email
func GetUserByEmail(email string, includePassword bool) (models.User, error) {
	var usersData models.User
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	if includePassword {
		usersData = models.User{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		}
	} else {
		usersData = models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return usersData, nil

}

// Update User
func UpdateUser(id string, values models.User) (models.User, error) {
	user, err := repositories.UpdateUser(id, values)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Delete User
func DeleteUser(id string) (models.User, error) {
	return repositories.DeleteUser(id)
}
