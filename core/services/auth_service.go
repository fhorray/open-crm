package services

import (
	"errors"
	"fmt"
	"open-crm/core/models"
	"open-crm/pkg/config"
	"open-crm/utils"
)

// Sign in
func SignIn(email, password string) (models.SignInResponseDTO, error) {
	// Get User
	user, err := GetUserByEmail(email, true)
	if err != nil {
		return models.SignInResponseDTO{}, errors.New("user not found")
	}

	// Check Password
	matchPassword := utils.CheckPasswordHash(user.Password, password)

	// if password matches generate token
	if !matchPassword {
		return models.SignInResponseDTO{}, errors.New("invalid credentials")
	}

	// generate jwt-token
	tokenData := map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	// Access Token
	acToken, err := utils.GenerateJWT(tokenData, config.Cfg.AUTH.JWT_EXPIRES_IN, "", "")
	if err != nil {
		return models.SignInResponseDTO{}, fmt.Errorf("error generating access token: %w", err)
	}

	// Refresh Token
	rfToken, err := utils.GenerateJWT(tokenData, config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN, "", "")
	if err != nil {
		return models.SignInResponseDTO{}, fmt.Errorf("error generating access token: %w", err)
	}

	return models.SignInResponseDTO{
		AccessToken:  acToken,
		RefreshToken: rfToken,
	}, nil
}
