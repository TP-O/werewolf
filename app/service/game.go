package service

import (
	"context"
	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core"
	"uwwolf/app/instance"
	"uwwolf/app/model"
	"uwwolf/app/types"
	"uwwolf/db"
)

var gameManger = core.NewManager()

func CreateGame(setting *types.GameSetting) contract.Game {
	game := &model.Game{}
	db.Client().Omit("WinningFactionId").Create(game)

	return gameManger.AddGame(game.Id, setting)
}

func AddPlayersToGame(gameId types.GameId, players map[types.PlayerId]contract.Player) {
	if gameManger.Game(gameId) == nil {
		return
	}

	var roleAssignments []*model.RoleAssignment
	redisPipe := instance.RedisClient.Pipeline()

	for _, player := range players {
		redisPipe.Set(
			context.Background(),
			enum.PlayerId2GameIdCacheNamespace+string(player.Id()),
			gameId,
			-1,
		)
		roleAssignments = append(roleAssignments, &model.RoleAssignment{
			GameId:   gameId,
			PlayerId: player.Id(),
			RoleId:   player.MainRoleId(),
		})

	}

	redisPipe.Exec(context.Background())
	db.Client().Omit("FactionId").Create(roleAssignments)
}
