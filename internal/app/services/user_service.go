package services

import (
	"open-crm/internal/app/models"
	"open-crm/internal/app/repositories"
	"open-crm/pkg/database"
	"open-crm/pkg/utils"

	"github.com/google/uuid"
)

// Create User
func CreateUser(payload models.CreateUserDTO) (*models.UserResponseDTO, error) {
	userRepo := repositories.NewUserRepository(database.DB)
	accountRepo := repositories.NewAccountRepository(database.DB)

	// Hash da senha
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	// Criar usuário
	user := &models.User{
		Name:  payload.Name,
		Email: payload.Email,
		Roles: payload.Roles,
	}
	user, err = userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Criar conta
	account := &models.Account{
		UserID:     user.ID,
		ProviderID: "credential",
		Password:   hashedPassword,
	}
	_, err = accountRepo.Create(account)
	if err != nil {
		return nil, err
	}

	dto := utils.ToUserResponseDTO(*user)
	return dto, nil
}

// Get All Users
func GetAllUsers() ([]models.UserResponseDTO, error) {
	userRepo := repositories.NewUserRepository(database.DB)

	users, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return *utils.ToUserResponseDTOs(*users), nil
}

// Get User By Id
func GetUserById(id uuid.UUID) (*models.UserResponseDTO, error) {
	userRepo := repositories.NewUserRepository(database.DB)

	user, err := userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return utils.ToUserResponseDTO(*user), nil
}

// Get User By Email
func GetUserByEmail(email string, includePassword bool) (*models.User, error) {
	userRepo := repositories.NewUserRepository(database.DB)

	user, err := userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if includePassword {
		return user, nil
	}

	// Retorna dados públicos (sem senha)
	publicUser := &models.User{
		Name:  user.Name,
		Email: user.Email,
	}
	return publicUser, nil
}

// Update User
func UpdateUser(id uuid.UUID, values *models.User) (*models.UserResponseDTO, error) {
	userRepo := repositories.NewUserRepository(database.DB)

	updatedUser, err := userRepo.Update(id, values)
	if err != nil {
		return nil, err
	}

	return utils.ToUserResponseDTO(*updatedUser), nil
}

// Delete User
func DeleteUser(id uuid.UUID) (*models.UserResponseDTO, error) {
	userRepo := repositories.NewUserRepository(database.DB)

	user, err := userRepo.Delete(id)
	if err != nil {
		return nil, err
	}

	return utils.ToUserResponseDTO(*user), nil
}
