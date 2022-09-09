package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Hello(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Hello, user at %s", c.Context().RemoteAddr().String()))
}
