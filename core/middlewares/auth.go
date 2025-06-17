package middlewares

import (
	"net/http"
	"open-crm/core/services"
	"open-crm/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ignore auth for specific paths
		if strings.HasPrefix(c.Path(), "/v1/auth/sign-in") || strings.HasPrefix(c.Path(), "/v1/auth/sign-up") {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		rfToken := c.Cookies("myapp.refresh_token")

		// Must have either access token or refresh token
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

		c.Locals("user", &user)
		return c.Next()
	}
}
