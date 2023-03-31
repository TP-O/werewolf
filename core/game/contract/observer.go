package contract

import (
	"uwwolf/game/types"

	"github.com/solarlune/resolv"
)

type Observer interface {
	LoadMap(name string)

	AddObjects(objs ...*resolv.Object)

	ObservePlayers(playerIDs []types.PlayerID)

	MovePlayer(playerID types.PlayerID, x float64, y float64) bool
}
