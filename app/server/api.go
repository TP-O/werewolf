package server

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"uwwolf/app/model"
	"uwwolf/app/service"
	"uwwolf/app/types"
	"uwwolf/config"
	"uwwolf/db"
)

func StartAPI() {
	app := fiber.New()

	app.Post("/api/v1/start", func(c *fiber.Ctx) error {
		setting := &types.GameSetting{}
		c.BodyParser(setting)

		if service.ArePlayersReadyToPlay(setting.PlayerIds...) {
			return errors.New("All players must be online and not in any game!")
		}

		game := service.CreateGame(setting)
		if players, err := game.Start(); err != nil {
			return c.JSON(fiber.Map{
				"ok":    false,
				"error": err,
			})
		} else {
			var roleAssignments []*model.RoleAssignment

			for _, player := range players {
				roleAssignments = append(roleAssignments, &model.RoleAssignment{
					GameId:   game.Id(),
					PlayerId: player.Id(),
					RoleId:   player.MainRoleId(),
				})
			}

			db.Client().Omit("FactionId").Create(roleAssignments)
		}

		return c.JSON(fiber.Map{
			"ok": true,
			"data": fiber.Map{
				"gameId": game.Id(),
			},
		})
	})

	log.Fatal(app.Listen(":" + strconv.Itoa(config.App.HttpPort)))
}
