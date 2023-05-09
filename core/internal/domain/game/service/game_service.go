package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	game "uwwolf/internal/app/game/logic"
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/config"
	"uwwolf/internal/domain/game/model"
	db "uwwolf/internal/infra/db/postgres"
	"uwwolf/internal/infra/db/redis"
	"uwwolf/internal/infra/server/api/dto"

	goredis "github.com/redis/go-redis/v9"
)

// RommService handles game-related business logic.
type GameService interface {
	// GameConfig returns game config of given room ID.
	// Returns the default config if doesn't exist.
	GameSettings(roomID string) model.GameSetting

	// UpdateGameConfig replaces the given room ID's old config to the new one.
	UpdateGameSettings(roomID string, config dto.UpdateGameSetting) error

	// RegisterGame create a game with the given config and player ID list.
	RegisterGame(config model.GameSetting, playerIDs []types.PlayerID) (game.Moderator, error)
}

type gameService struct {
	// config is global game config.
	config config.Game

	// rdb is redis connection.
	rdb *goredis.ClusterClient

	// pdb is postgreSQL connection.
	pdb db.Store

	// gameManger is game management instance.
	gameManager game.Manager
}

func NewGameService(
	config config.Game,
	rdb *goredis.ClusterClient,
	pdb db.Store,
	gameManager game.Manager,
) GameService {
	return &gameService{
		config,
		rdb,
		pdb,
		gameManager,
	}
}

// GameConfig returns game config of given room ID.
// Returns the default config if doesn't exist.
func (gs gameService) GameSettings(roomID string) model.GameSetting {
	var config model.GameSetting

	encodedConfig := gs.rdb.Get(
		context.Background(),
		redis.WaitingRoom+roomID,
	).Val()
	if err := json.Unmarshal([]byte(encodedConfig), &config); err != nil {
		return model.GameSetting{
			RoleIDs:            []types.RoleID{declare.SeerRoleID},
			NumberWerewolves:   1,
			TurnDuration:       20,
			DiscussionDuration: 90,
		}
	}

	return config
}

// UpdateGameConfig replaces the given room ID's old config to the new one.
func (gs gameService) UpdateGameSettings(roomID string, config dto.UpdateGameSetting) error {
	encodedConfig, _ := json.Marshal(config)

	return gs.rdb.Set(
		context.Background(),
		redis.GameSetting+roomID,
		string(encodedConfig),
		-1,
	).Err()
}

// RegisterGame create a game with the given config and player ID list.
func (gs gameService) RegisterGame(config model.GameSetting, playerIDs []types.PlayerID) (game.Moderator, error) {
	gameRecord, err := gs.pdb.CreateGame(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Something went wrong!")
	}

	mod, _ := gs.gameManager.RegisterGame(&types.GameRegistration{
		ID:                 types.GameID(gameRecord.ID),
		TurnDuration:       time.Duration(config.TurnDuration) * time.Second,
		DiscussionDuration: time.Duration(config.DiscussionDuration) * time.Second,
		GameInitialization: types.GameInitialization{
			RoleIDs:          config.RoleIDs,
			RequiredRoleIDs:  config.RequiredRoleIDs,
			NumberWerewolves: config.NumberWerewolves,
			PlayerIDs:        playerIDs,
		},
	})
	return mod, nil
}
