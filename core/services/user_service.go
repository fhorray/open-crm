package services

import (
	"open-crm/core/models"
	"open-crm/core/repositories"
	"open-crm/utils"

	"github.com/google/uuid"
)

// Create User
func CreateUser(payload models.CreateUserDTO) (*models.UserResponseDTO, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:    uuid.New(),
		Name:  payload.Name,
		Email: payload.Email,
		Roles: payload.Roles,
	}

	if err := repositories.CreateUser(user); err != nil {
		return nil, err
	}

	account := &models.Account{
		ID:         uuid.New(),
		UserID:     user.ID,
		ProviderID: "credential",
		Password:   hashedPassword,
	}

	if err := repositories.CreateAccount(account); err != nil {
		return nil, err
	}

	dto := utils.ToUserResponseDTO(*user)
	return &dto, nil
}

// Get All Users
func GetAllUsers() ([]models.UserResponseDTO, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return utils.ToUserResponseDTOs(users), nil
}

// Get User By Id
func GetUserById(id string) (*models.UserResponseDTO, error) {
	user, err := repositories.GetUserById(id)
	if err != nil {
		return nil, err
	}

	dto := utils.ToUserResponseDTO(user)
	return &dto, nil
}

// Get User By Email
func GetUserByEmail(email string, includePassword bool) (models.User, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	if includePassword {
		return user, nil
	}

	// Retorna apenas dados públicos
	return models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Update User
func UpdateUser(id string, values models.User) (*models.UserResponseDTO, error) {
	user, err := repositories.UpdateUser(id, values)
	if err != nil {
		return nil, err
	}

	dto := utils.ToUserResponseDTO(user)
	return &dto, nil
}

// Delete User
func DeleteUser(id string) (*models.UserResponseDTO, error) {
	user, err := repositories.DeleteUser(id)
	if err != nil {
		return nil, err
	}

	dto := utils.ToUserResponseDTO(user)
	return &dto, nil
}
