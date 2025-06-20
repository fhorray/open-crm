package services

import (
	"errors"
	"fmt"
	"open-crm/config"
	"open-crm/internal/app/models"
	"open-crm/internal/app/repositories"
	"open-crm/pkg/database"
	"open-crm/pkg/utils"

	"gorm.io/gorm"
)

// Create Interface
type AuthService interface {
	Register(user *models.RegisterRequestDTO) (*models.GetSessionResponseDTO, error)
	Login(credentials *models.LoginRequestDTO) (*models.GetSessionResponseDTO, error)
	// ValidateToken(token string) (uuid.UUID, error)
	// Logout() error
}

type authService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	accountRepo repositories.AccountRepository
}

func NewAuthService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	accountRepo repositories.AccountRepository,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		accountRepo: accountRepo,
	}
}

// Register
func (s *authService) Register(req *models.RegisterRequestDTO) (*models.GetSessionResponseDTO, error) {
	sessionRepo := repositories.NewSessionRepository(database.DB)

	// Verifica se já existe o usuário
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already taken")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Gera hash da senha
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Gera o UUID manualmente
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Roles: config.Cfg.AUTH.DEFAULT_ROLE,
	}

	// Salva usuário
	user, err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Cria conta
	account := &models.Account{
		UserID:     user.ID,
		ProviderID: "credential",
		Password:   hashedPassword,
	}
	if _, err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	// Gera tokens
	tokenData := map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	accessToken, accessClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_EXPIRES_IN, "", "")
	if err != nil {
		return nil, err
	}

	refreshToken, refreshClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN, "", "")
	if err != nil {
		return nil, err
	}

	// Cria sessão
	session := &models.Session{
		UserID:                user.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
	}
	if err := sessionRepo.Create(session); err != nil {
		return nil, err
	}

	// Retorna a sessão real com ID preenchido corretamente
	return &models.GetSessionResponseDTO{
		User:    utils.ToUserResponseDTO(*user),
		Session: session,
	}, nil
}

// Sign in
func (s *authService) Login(req *models.LoginRequestDTO) (*models.GetSessionResponseDTO, error) {
	// Buscar o usuário por email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Buscar a conta do usuário com provider "credential"
	account, err := s.accountRepo.FindByUserAndProvider(user.ID, "credential")
	if err != nil || account == nil {
		return nil, errors.New("account not found")
	}

	// Verificar a senha
	if !utils.CheckPasswordHash(req.Password, account.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Payload do token
	tokenData := map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	// Gerar tokens
	accessToken, accessClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_EXPIRES_IN, "", "")
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, refreshClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN, "", "")
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Criar ou atualizar sessão no banco
	session := &models.Session{
		UserID:                user.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
	}

	sessionRepo := repositories.NewSessionRepository(database.DB)

	if err := sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("error saving session: %w", err)
	}

	// Retornar dados para o cliente
	return &models.GetSessionResponseDTO{
		User: utils.ToUserResponseDTO(*user),
		Session: &models.Session{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		},
	}, nil
}
