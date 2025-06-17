package utils

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success *bool  `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func SendResponse(c *fiber.Ctx, body APIResponse) error {
	if body.Success == nil {
		success := body.Status >= 200 && body.Status < 300
		body.Success = &success
	}

	return c.Status(body.Status).JSON(body)
}
