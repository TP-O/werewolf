package game

import (
	"errors"
	"uwwolf/game/contract"
	"uwwolf/game/types"

	"github.com/google/uuid"
)

type manager struct {
	games map[types.GameID]contract.Game
}

var mangerInstance *manager

func Manager() contract.GameManger {
	if mangerInstance == nil {
		mangerInstance = &manager{
			games: make(map[types.GameID]contract.Game),
		}
	}

	return mangerInstance
}

func (m *manager) Game(gameID types.GameID) contract.Game {
	return m.games[gameID]
}

func (m *manager) NewGame(setting *types.GameSetting) (contract.Game, error) {
	game := NewGame(uuid.NewString(), setting)

	// if err := db.Client().Query(
	// 	`INSERT INTO games (id) VALUES (?)`,
	// 	game.ID(),
	// ).Exec(); err != nil {
	// 	game.Finish()

	// 	return nil, errors.New("Unable to create game (╯°□°)╯︵ ┻━┻")
	// }

	m.games[game.ID()] = game

	return game, nil
}

func (m *manager) RemoveGame(gameID types.GameID) (bool, error) {
	if removedGame := m.games[gameID]; removedGame == nil {
		return false, errors.New("Game does not exist (• ε •)")
	} else {
		removedGame.Finish()
		delete(m.games, gameID)

		return true, nil
	}
}
