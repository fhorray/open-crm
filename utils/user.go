package utils

import "open-crm/core/models"

// Converte um único User em UserResponseDTO
func ToUserResponseDTO(user models.User) models.UserResponseDTO {
	return models.UserResponseDTO{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Email:          user.Email,
		Roles:          user.Roles,
		Image:          user.Image,
		EmailVerified:  user.EmailVerified,
		IsActive:       user.IsActive,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

// Converte um slice de User em um slice de UserResponseDTO
func ToUserResponseDTOs(users []models.User) []models.UserResponseDTO {
	result := make([]models.UserResponseDTO, 0, len(users))
	for _, user := range users {
		result = append(result, ToUserResponseDTO(user))
	}
	return result
}
