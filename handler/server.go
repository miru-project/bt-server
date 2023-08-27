package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/miru-project/bt-server/config"
)

func Hello(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Miru BT server(%s) is running ✨", config.VERSION))
}

func Version(c *fiber.Ctx) error {
	return c.SendString(config.VERSION)
}
