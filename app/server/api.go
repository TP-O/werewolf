package server

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"uwwolf/config"
)

func StartAPI() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"aa": 1, "bb": "aa"})
	})

	log.Fatal(app.Listen(":" + strconv.Itoa(config.App.HttpPort)))
}
