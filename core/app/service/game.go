package service

import (
	"context"
	"encoding/json"
	"fmt"
	"uwwolf/app/data"
	"uwwolf/app/dto"
	"uwwolf/game"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	"uwwolf/storage/postgres"

	"github.com/redis/go-redis/v9"
)

type GameService interface {
	RegisterGame(config *types.GameConfig) (contract.Moderator, error)
	UpdateGameConfig(roomID string, config dto.UpdateGameConfigDto) error
	GameConfig(roomID string) *data.GameConfig
}

type gameService struct {
	rdb *redis.ClusterClient
	pdb *postgres.Store
}

func NewGameService(rdb *redis.ClusterClient, pdb *postgres.Store) GameService {
	return &gameService{
		rdb,
		pdb,
	}
}

// RegisterGame creates a moderator and assigns a game for it.
func (gs gameService) RegisterGame(config *types.GameConfig) (contract.Moderator, error) {
	mod, err := game.NewModerator(config)
	if err != nil {
		return nil, err
	}

	gameRecord, err := gs.pdb.CreateGame(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to create game")
	}

	// Return the old moderator if the game ID has been assigned to it
	if ok, _ := game.Manager().AddModerator(types.GameID(gameRecord.ID), mod); !ok {
		return game.Manager().Moderator(types.GameID(gameRecord.ID)), nil
	}

	return mod, nil
}

func (gs gameService) UpdateGameConfig(roomID string, config dto.UpdateGameConfigDto) error {
	gs.rdb.Set(
		context.Background(),
		"dsds",
		config,
		-1,
	)

	return nil
}

func (gs gameService) GameConfig(roomID string) *data.GameConfig {
	var config *data.GameConfig
	configJson := gs.rdb.Get(context.Background(), "").String()
	if err := json.Unmarshal([]byte(configJson), config); err != nil {
		return &data.GameConfig{
			RoleIDs:            []types.RoleID{vars.SeerRoleID},
			NumberWerewolves:   1,
			TurnDuration:       20,
			DiscussionDuration: 90,
		}
	}

	return config
}
