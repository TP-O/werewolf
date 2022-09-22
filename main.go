package main

import (
	"strconv"
	"uwwolf/config"
	"uwwolf/module/game/model"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(model.Role{})
	})

	app.Listen("0.0.0.0:" + strconv.Itoa(config.App.Port))
}
