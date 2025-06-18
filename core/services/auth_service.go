package services

import (
	"errors"
	"fmt"
	"open-crm/core/models"
	"open-crm/core/repositories"
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

	// Get Credential Account (provider_id = "credential")
	account, err := repositories.GetAccountByUserAndProvider(user.ID.String(), "credential")
	if err != nil {
		return models.SignInResponseDTO{}, errors.New("account not found")
	}

	// Check Password
	if ok := utils.CheckPasswordHash(account.Password, password); !ok {
		return models.SignInResponseDTO{}, errors.New("invalid credentials")
	}

	// Token payload
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
		return models.SignInResponseDTO{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return models.SignInResponseDTO{
		AccessToken:  acToken,
		RefreshToken: rfToken,
	}, nil
}
