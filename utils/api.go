package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func SendResponse(c *fiber.Ctx, statusCode int, body APIResponse) error {
	body.Status = statusCode
	return c.Status(statusCode).JSON(body)
}
