package middlewares

import (
	"net/http"
	"open-crm/internal/app/services"
	"open-crm/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ignorar autenticação para certas rotas
		if strings.HasPrefix(c.Path(), "/v1/auth/sign-in") ||
			strings.HasPrefix(c.Path(), "/v1/auth/sign-up") ||
			strings.HasPrefix(c.Path(), "/reference") ||
			strings.HasPrefix(c.Path(), "/swagger.json") {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		rfToken := c.Cookies("myapp.refresh_token")

		// Deve ter access token ou refresh token
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

		// Parse do JWT
		parsedAcToken, err := utils.ParseJWT(acToken)
		if err != nil {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		claims, ok := parsedAcToken.Claims.(jwt.MapClaims)
		if !ok {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: invalid token claims",
			})
		}

		// Extrair user id do claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: invalid token claims (user_id missing)",
			})
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: invalid user ID format",
			})
		}

		user, err := services.GetUserById(userID)
		if err != nil {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: user not found",
			})
		}

		// Salva o usuário no contexto local
		c.Locals("user", *user)

		return c.Next()
	}
}
