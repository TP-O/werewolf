package service

import (
	"context"
	"fmt"
	"uwwolf/db"
	"uwwolf/game"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

type GameService interface {
	RegisterGame(modInit *types.ModeratorInit) (contract.Moderator, error)
}

type gameService struct {
	store *db.Store
}

func NewGameService(db *db.Store) GameService {
	return &gameService{
		store: db,
	}
}

// RegisterGame creates a moderator and assigns a game for it.
func (gs gameService) RegisterGame(modInit *types.ModeratorInit) (contract.Moderator, error) {
	mod, err := game.NewModerator(modInit)
	if err != nil {
		return nil, err
	}

	gameRecord, err := db.DB().CreateGame(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to create game")
	}

	// Return the old moderator if the game ID has been assigned to it
	if ok, _ := game.Manager().AddModerator(types.GameID(gameRecord.ID), mod); !ok {
		return game.Manager().Moderator(types.GameID(gameRecord.ID)), nil
	}

	return mod, nil
}
