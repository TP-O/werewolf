package game

// type gameManager struct {
// 	games map[types.GameID]contract.Game
// }

// var manger *gameManager

// func Manager() contract.GameManger {
// 	if manger == nil {
// 		manger = &gameManager{
// 			games: make(map[types.GameID]contract.Game),
// 		}
// 	}

// 	return manger
// }

// func (g *gameManager) Game(gameID types.GameID) contract.Game {
// 	return g.games[gameID]
// }

// func (g *gameManager) CreateGame(setting *types.GameSetting) (contract.Game, error) {
// 	game :=

// 	// if err := db.Client().Query(
// 	// 	`INSERT INTO games (id) VALUES (?)`,
// 	// 	game.ID(),
// 	// ).Exec(); err != nil {
// 	// 	game.Finish()

// 	// 	return nil, errors.New("Unable to create game (╯°□°)╯︵ ┻━┻")
// 	// }

// 	m.games[game.ID()] = game

// 	return game, nil
// }

// func (g *gameManager) RemoveGame(gameID types.GameID) (bool, error) {
// 	if removedGame := m.games[gameID]; removedGame == nil {
// 		return false, errors.New("Game does not exist (• ε •)")
// 	} else {
// 		removedGame.Finish()
// 		delete(m.games, gameID)

// 		return true, nil
// 	}
// }
