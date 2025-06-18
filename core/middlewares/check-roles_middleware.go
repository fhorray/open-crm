package middlewares

import (
	"slices"
	"strings"

	"open-crm/core/models"
	"open-crm/utils"

	"github.com/gofiber/fiber/v2"
)

func CheckRoles(roles string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := c.Locals("user")
		if userCtx == nil {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "Unauthorized: user not found in context",
			})
		}

		user, ok := userCtx.(models.UserResponseDTO)
		if !ok {
			return utils.SendResponse(c, utils.APIResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Internal error: invalid user type",
			})
		}

		allowedRoles := strings.Split(roles, ",")
		for i := range allowedRoles {
			allowedRoles[i] = strings.TrimSpace(allowedRoles[i])
		}

		userRoles := strings.SplitSeq(user.Roles, ",")
		for userRole := range userRoles {
			userRole = strings.TrimSpace(userRole)
			if slices.Contains(allowedRoles, userRole) {
				// user is permitted to continue
				return c.Next()
			}
		}

		// no role registered on user profile
		return utils.SendResponse(c, utils.APIResponse{
			Status:  fiber.StatusForbidden,
			Message: "Forbidden: insufficient permissions",
		})
	}
}
