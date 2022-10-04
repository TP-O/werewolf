package server

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"uwwolf/app/service"
	"uwwolf/app/types"
	"uwwolf/config"
)

func StartAPI() {
	app := fiber.New()

	app.Post("/api/v1/games", func(c *fiber.Ctx) error {
		setting := &types.GameSetting{}
		c.BodyParser(setting)

		err := service.CreateGame(setting)

		return c.JSON(err)
	})

	log.Fatal(app.Listen(":" + strconv.Itoa(config.App.HttpPort)))
}
