package server

import (
	"errors"
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

		if service.ArePlayersReadyToPlay(setting.PlayerIds...) {
			return errors.New("All players must be online and not in any game!")
		}

		if err := service.CreateGame(setting); err != nil {
			return c.JSON(fiber.Map{
				"ok":    false,
				"error": err,
			})
		}

		return c.JSON(fiber.Map{
			"ok": true,
		})
	})

	log.Fatal(app.Listen(":" + strconv.Itoa(config.App.HttpPort)))
}
