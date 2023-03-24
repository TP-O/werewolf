package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"uwwolf/app/data"
	"uwwolf/app/dto"
	"uwwolf/app/enum"
	"uwwolf/config"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	"uwwolf/storage/postgres"

	"github.com/redis/go-redis/v9"
)

// RommService handles game-related business logic.
type GameService interface {
	// GameConfig returns game config of given room ID.
	// Returns the default config if doesn't exist.
	GameConfig(roomID string) data.GameConfig

	// UpdateGameConfig replaces the given room ID's old config to the new one.
	UpdateGameConfig(roomID string, config dto.ReplaceGameConfigDto) error

	// CheckBeforeRegistration checks the combination of room and game config before
	// registering a game. This makes sure the game runs properly without any unexpectation.
	CheckBeforeRegistration(room data.WaitingRoom, cfg data.GameConfig) error

	// RegisterGame create a game with the given config and player ID list.
	RegisterGame(config data.GameConfig, playerIDs []types.PlayerID) (contract.Moderator, error)
}

type gameService struct {
	// rdb is redis connection.
	rdb *redis.ClusterClient
	// pdb is postgreSQL connection.
	pdb postgres.Store

	// gameManger is game management instance.
	gameManager contract.Manager
}

func NewGameService(rdb *redis.ClusterClient, pdb postgres.Store, gameManager contract.Manager) GameService {
	return &gameService{
		rdb,
		pdb,
		gameManager,
	}
}

// GameConfig returns game config of given room ID.
// Returns the default config if doesn't exist.
func (gs gameService) GameConfig(roomID string) data.GameConfig {
	var config data.GameConfig

	encodedConfig := gs.rdb.Get(
		context.Background(),
		enum.RoomID2GameConfigRNs+roomID,
	).Val()
	if err := json.Unmarshal([]byte(encodedConfig), &config); err != nil {
		return data.GameConfig{
			RoleIDs:            []types.RoleID{vars.SeerRoleID},
			NumberWerewolves:   1,
			TurnDuration:       20,
			DiscussionDuration: 90,
		}
	}

	return config
}

// UpdateGameConfig replaces the given room ID's old config to the new one.
func (gs gameService) UpdateGameConfig(roomID string, config dto.ReplaceGameConfigDto) error {
	encodedConfig, _ := json.Marshal(config)

	return gs.rdb.Set(
		context.Background(),
		enum.RoomID2GameConfigRNs+roomID,
		string(encodedConfig),
		-1,
	).Err()
}

// CheckBeforeRegistration checks the combination of room and game config before
// registering a game. This makes sure the game runs properly without any unexpectation.
func (gs gameService) CheckBeforeRegistration(room data.WaitingRoom, cfg data.GameConfig) error {
	if len(room.PlayerIDs) < int(config.Game().MinCapacity) {
		return fmt.Errorf("Invite more players to play!")
	} else if len(room.PlayerIDs) > int(config.Game().MaxCapacity) {
		return fmt.Errorf("Too many players!")
	}

	numberOfPlayers := len(room.PlayerIDs)
	if (numberOfPlayers%2 == 0 && numberOfPlayers/2 <= int(cfg.NumberWerewolves)) ||
		(numberOfPlayers%2 != 0 && numberOfPlayers/2 < int(cfg.NumberWerewolves)) {
		return fmt.Errorf("Unblanced number of werewolves!")
	}

	return nil
}

// RegisterGame create a game with the given config and player ID list.
func (gs gameService) RegisterGame(config data.GameConfig, playerIDs []types.PlayerID) (contract.Moderator, error) {
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
