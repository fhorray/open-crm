package controllers

import (
	"net/http"
	"open-crm/core/models"
	"open-crm/core/services"
	"open-crm/pkg/config"
	"open-crm/utils"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Sign In
func SignIn(c *fiber.Ctx) error {
	var body models.SignInRequestDTO

	if err := c.BodyParser(&body); err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Body",
		})
	}

	tokens, err := services.SignIn(body.Email, body.Password)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	}

	// Set Access Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(config.Cfg.AUTH.JWT_EXPIRES_IN),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	// Set Refresh Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
	})
}

// Sign Up
func SignUp(c *fiber.Ctx) error {
	var payload models.CreateUserDTO

	if err := c.BodyParser(&payload); err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Body",
		})
	}

	//Validate Payload
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

	// Call service to create user
	user, err := services.CreateUser(models.CreateUserDTO{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Roles:    config.Cfg.AUTH.DEFAULT_ROLE,
	})
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Message: "Error creating user",
			Status:  http.StatusInternalServerError,
		})
	}

	// Sign In user to get tokens
	tokens, err := services.SignIn(payload.Email, payload.Password)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	}

	// Set Access Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(config.Cfg.AUTH.JWT_EXPIRES_IN),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	// Set Refresh Token Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "myapp.refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(config.Cfg.AUTH.JWT_REFRESH_TOKEN_EXPIRES_IN),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
	})

	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusCreated,
		Data: models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Image: user.Image,
			Roles: user.Roles,
		},
	})
}

// Sign Out
func SignOut(c *fiber.Ctx) error {
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
	authHeader := c.Get("Authorization")
	rfToken := c.Cookies("myapp.refresh_token")

	// Deve ter token de acesso ou refresh token
	if !strings.HasPrefix(authHeader, "Bearer ") && rfToken == "" {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: not authenticated",
		})
	}

	acToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if acToken == "" {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: missing access token",
		})
	}

	// Parse JWT
	parsedAcToken, err := utils.ParseJWT(acToken)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	claims := parsedAcToken.Claims.(jwt.MapClaims)

	userId, ok := claims["id"].(string)
	if !ok {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: invalid token claims",
		})
	}

	user, err := services.GetUserById(userId)
	if err != nil {
		return utils.SendResponse(c, utils.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized: user not found",
		})
	}

	// Retorna os dados do usuário na resposta
	return utils.SendResponse(c, utils.APIResponse{
		Status: http.StatusOK,
		Data: models.GetSessionResponseDTO{
			Session: claims,
			User:    user,
		},
	})
}
