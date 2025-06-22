package handlers

import (
	"net/http"
	"open-crm/internal/app/models"
	"open-crm/internal/app/repositories"
	"open-crm/internal/app/services"
	"open-crm/pkg/database"
	"open-crm/pkg/utils"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

// Register Handler
func Register(c *fiber.Ctx) error {
	// Create repositories instances
	userRepo := repositories.NewUserRepository(database.DB)
	sessionRepo := repositories.NewSessionRepository(database.DB)
	accountRepo := repositories.NewAccountRepository(database.DB)

	var payload models.RegisterRequestDTO
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid body",
		})
	}

	// validation
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = "Invalid field: " + e.Tag()
		}
		return utils.SendResponse(c, utils.APIResponse{
			Status: http.StatusBadRequest,
			Data:   errors,
		})
	}

	// Create Service
	authService := services.NewAuthService(c, userRepo, sessionRepo, accountRepo)
	sessionData, err := authService.Register(&payload)

	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Set cookien on client
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.access_token",
		Value:    *sessionData.Session.AccessToken,
		Expires:  sessionData.Session.AccessTokenExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.refresh_token",
		Value:    sessionData.Session.RefreshToken,
		Expires:  sessionData.Session.RefreshTokenExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	// Retorna os dados da sessão criada
	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusCreated,
		Data:   sessionData,
	})
}

// Login
func Login(c *fiber.Ctx) error {
	// Create repositories instances
	userRepo := repositories.NewUserRepository(database.DB)
	sessionRepo := repositories.NewSessionRepository(database.DB)
	accountRepo := repositories.NewAccountRepository(database.DB)

	// Crate Auth Service instance
	authService := services.NewAuthService(c, userRepo, sessionRepo, accountRepo)

	var body models.LoginRequestDTO
	if err := c.BodyParser(&body); err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid body",
		})
	}

	// Validação (opcional, se quiser usar validate)
	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = "Invalid field: " + e.Tag()
		}
		return utils.SendResponse(c, utils.APIResponse{
			Status: http.StatusBadRequest,
			Data:   errors,
		})
	}

	sessionData, err := authService.Login(&body)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	// Setar cookies com tokens e expiração vindos do service
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.access_token",
		Value:    *sessionData.Session.AccessToken,
		Expires:  sessionData.Session.AccessTokenExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "myapp.refresh_token",
		Value:    sessionData.Session.RefreshToken,
		Expires:  sessionData.Session.RefreshTokenExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   sessionData,
	})
}

// Sign Out
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	// Set Refresh Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
	})
}

// Get Session
func GetSession(c *fiber.Ctx) error {
	sessionRepo := repositories.NewSessionRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)

	authHeader := c.Get("Authorization")
	rfToken := c.Cookies("myapp.refresh_token")

	// must have the access token or refresh token
	if !strings.HasPrefix(authHeader, "Bearer ") && rfToken == "" {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: not authenticated",
		})
	}

	acToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if acToken == "" && rfToken == "" {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: missing access token and refresh token",
		})
	}

	// tries to get session by access token, if not found tries the refresh token
	var session *models.Session
	var err error

	if rfToken != "" {
		session, err = sessionRepo.FindByToken(rfToken)

		if err != nil || session == nil {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "session not found",
			})
		}
	} else {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "session not found",
		})
	}

	// parse JWT to validate and get claims
	tokenToParse := acToken
	if tokenToParse == "" {
		tokenToParse = acToken
	}
	parsedToken, err := utils.ParseJWT(tokenToParse)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid token: " + err.Error(),
		})
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: invalid token",
		})
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: invalid token claims",
		})
	}

	// get user by id
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
		})
	}

	user, err := userRepo.FindByID(userUUID)
	if err != nil || user == nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not found",
		})
	}

	// response
	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data: models.GetSessionResponseDTO{
			User: utils.ToUserResponseDTO(*user),
			Session: &models.SessionDTO{
				ID:                    session.ID,
				UserID:                user.ID,
				UserAgent:             session.UserAgent,
				IPAddress:             session.IPAddress,
				RefreshToken:          session.RefreshToken,
				RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
			},
		},
	})
}

// Refresh Token
func RefreshToken(c *fiber.Ctx) error {
	// sessionRepo := repositories.NewSessionRepository(database.DB)
	// userRepo := repositories.NewUserRepository(database.DB)

	authHeader := c.Get("Authorization")
	rfToken := c.Cookies("myapp.refresh_token")

	// must have the access token or refresh token
	if !strings.HasPrefix(authHeader, "Bearer ") && rfToken == "" {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: not authenticated",
		})
	}

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data:   "Success",
	})

}
