package server

import "github.com/gofiber/fiber/v2"

func (h *handler) healthCheck(c *fiber.Ctx) error {
	err := h.dbc.Ping()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
