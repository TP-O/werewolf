package handler

import (
	"uwwolf/app/enum"
	"uwwolf/app/model"
	"uwwolf/app/service"
	"uwwolf/app/types"
	"uwwolf/db"

	"github.com/gofiber/fiber/v2"
)

func StartGame(c *fiber.Ctx) error {
	payload := c.Locals(enum.FiberLocalDataKey).(*types.GameSetting)

	if service.ArePlayersReadyToPlay(payload.PlayerIds...) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": enum.InvalidInputErrorTag,
			"error":   "All players must be online and not in any game!",
		})
	}

	game := service.CreateGame(payload)
	if players, err := game.Start(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": enum.SystemErrorTag,
			"error":   err,
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok": true,
		"data": fiber.Map{
			"gameId": game.Id(),
		},
	})
}
