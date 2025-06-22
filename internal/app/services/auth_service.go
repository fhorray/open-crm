package services

import (
	"errors"
	"fmt"
	"open-crm/config"
	"open-crm/internal/app/models"

	"open-crm/internal/app/repositories"
	"open-crm/pkg/database"
	"open-crm/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create Interface
type AuthService interface {
	Register(user *models.RegisterRequestDTO) (*models.GetSessionResponseDTO, error)
	Login(credentials *models.LoginRequestDTO) (*models.GetSessionResponseDTO, error)
	RefreshToken(refreshToken string) (*models.GetSessionResponseDTO, error)
}

type authService struct {
	fiberContext *fiber.Ctx
	userRepo     repositories.UserRepository
	sessionRepo  repositories.SessionRepository
	accountRepo  repositories.AccountRepository
}

func NewAuthService(
	fiberContext *fiber.Ctx,
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	accountRepo repositories.AccountRepository,
) AuthService {
	return &authService{
		fiberContext: fiberContext,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		accountRepo:  accountRepo,
	}
}

// Register
func (s *authService) Register(req *models.RegisterRequestDTO) (*models.GetSessionResponseDTO, error) {
	sessionRepo := repositories.NewSessionRepository(database.DB)
	authRepo := repositories.NewAuthRepository(database.DB)

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

	// Create user data
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Roles: config.Cfg.AUTH.DEFAULT_ROLE,
	}

	// Create account data
	account := &models.Account{
		ProviderID: "credential",
		Password:   hashedPassword,
	}

	// Call the Register method inside auth repository to create user and account
	user, err = authRepo.Register(user, account)
	if err != nil {
		return nil, err
	}

	// Upsert token payload
	tokenData := &utils.Claims{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	// Create Access Tk
	accessToken, accessClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_EXPIRES_IN, "", "")
	if err != nil {
		return nil, err
	}

	// Create Refresh Tk
	refreshToken, refreshClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN, "", "")
	if err != nil {
		return nil, err
	}

	// Create Session inside DB
	session := &models.Session{
		UserID:                user.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		UserAgent:             string(s.fiberContext.Request().Header.UserAgent()),
		IPAddress:             s.fiberContext.IP(),
	}
	createdSession, err := sessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return &models.GetSessionResponseDTO{
		User: utils.ToUserResponseDTO(*user),
		Session: &models.SessionDTO{
			ID:                    createdSession.ID,
			AccessToken:           &accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		},
	}, nil
}

// Login
func (s *authService) Login(req *models.LoginRequestDTO) (*models.GetSessionResponseDTO, error) {
	// Create Repositories
	sessionRepo := repositories.NewSessionRepository(database.DB)

	// Get user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Get users' account by provider "credential"
	account, err := s.accountRepo.FindByUserAndProvider(user.ID, "credential")
	if err != nil || account == nil {
		return nil, errors.New("account not found")
	}

	// Delete users sessions from DB
	err = s.sessionRepo.DeleteUserSessions(user.ID)
	if err != nil {
		return nil, err
	}

	// Verify password
	if !utils.CheckPasswordHash(account.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Upsert token Payload
	tokenData := &utils.Claims{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	// Generate Access Tk
	accessToken, accessClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_EXPIRES_IN, "", "")
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// Generate Refresh Tk
	refreshToken, refreshClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN, "", "")
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Create session inside DB
	session := &models.Session{
		UserID:                user.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
	}

	createdSession, err := sessionRepo.Create(session)
	if err != nil {
		return nil, fmt.Errorf("error saving session: %w", err)
	}

	// Response
	return &models.GetSessionResponseDTO{
		User: utils.ToUserResponseDTO(*user),
		Session: &models.SessionDTO{
			ID:                    createdSession.ID,
			UserID:                user.ID,
			AccessToken:           &accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
			UserAgent:             string(s.fiberContext.Context().UserAgent()),
			IPAddress:             s.fiberContext.IP(),
		},
	}, nil
}

// Refresh Token
func (s *authService) RefreshToken(refreshToken string) (*models.GetSessionResponseDTO, error) {
	// Parse e valida o refresh token
	token, err := utils.ParseJWT(refreshToken)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id missing in token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID in token claims")
	}

	// Busca o usuário
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verifica se existe sessão ativa com esse refresh token
	session, err := s.sessionRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, errors.New("session not found or expired")
	}

	// Gera novos tokens
	tokenData := &utils.Claims{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	accessToken, accessClaims, err := utils.GenerateToken(tokenData, config.Cfg.AUTH.JWT_EXPIRES_IN, "", "")
	if err != nil {
		return nil, errors.New("failed to generate new access token")
	}

	// session.AccessToken = accessToken
	// session.AccessTokenExpiresAt = accessClaims.ExpiresAt.Time

	if err := s.sessionRepo.Update(session); err != nil {
		return nil, errors.New("failed to update session")
	}

	return &models.GetSessionResponseDTO{
		User: utils.ToUserResponseDTO(*user),
		Session: &models.SessionDTO{
			ID:                    session.ID,
			AccessToken:           &accessToken,
			RefreshToken:          session.RefreshToken,
			AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
		},
	}, nil
}
