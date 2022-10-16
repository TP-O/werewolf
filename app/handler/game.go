package handler

import (
	"uwwolf/app/enum"
	"uwwolf/app/service"
	"uwwolf/app/types"

	"github.com/gofiber/fiber/v2"
)

func CreateGame(c *fiber.Ctx) error {
	payload := c.Locals(enum.FiberLocalPayloadKey).(*types.GameSetting)

	if service.ArePlayersReadyToPlay(payload.PlayerIds...) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
			"error": fiber.Map{
				"tag":     enum.InvalidInputErrorTag,
				"message": "All players must be online and not in any game!",
			},
		})
	}

	game := service.CreateGame(payload)

	if players, err := game.Init(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
			"error": fiber.Map{
				"tag":     enum.SystemErrorTag,
				"message": err.Error(),
			},
		})
	} else {
		service.AddPlayersToGame(game.Id(), players)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok": true,
		"data": fiber.Map{
			"gameId": game.Id(),
		},
	})
}
