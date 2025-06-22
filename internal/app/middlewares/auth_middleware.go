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
		// Rotas públicas (pode extrair para uma função se quiser)
		if strings.HasPrefix(c.Path(), "/v1/auth") ||
			strings.HasPrefix(c.Path(), "/reference") ||
			strings.HasPrefix(c.Path(), "/swagger.json") {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: missing bearer token",
			})
		}

		accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		parsedToken, err := utils.ParseJWT(accessToken)
		if err != nil || !parsedToken.Valid {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: invalid or expired token",
			})
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: invalid token claims",
			})
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: missing user_id",
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

		c.Locals("user", *user)
		return c.Next()
	}
}
