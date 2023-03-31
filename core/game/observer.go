package game

import (
	"uwwolf/game/contract"
	"uwwolf/game/types"

	"github.com/solarlune/resolv"
)

type observer struct {
	space        *resolv.Space
	playerWidth  float64
	playerHeight float64
	players      map[types.PlayerID]*resolv.Object
}

type ObserverSettings struct {
	PlayerWidth  float64
	PlayerHeight float64
	MapWidth     int
	MapHeight    int
	CellWidth    int
	CellHeight   int
}

func NewObserver(settings ObserverSettings) contract.Observer {
	return &observer{
		space: resolv.NewSpace(
			settings.MapWidth,
			settings.MapHeight,
			settings.CellWidth,
			settings.CellHeight,
		),
		playerWidth:  settings.PlayerWidth,
		playerHeight: settings.PlayerHeight,
		players:      make(map[types.PlayerID]*resolv.Object),
	}
}

func (o *observer) AddObjects(objs ...*resolv.Object) {
	for _, obj := range objs {
		if obj != nil {
			o.space.Add(obj)
		}
	}
}

func (o *observer) LoadMap(name string) {
	//
}

func (o *observer) ObservePlayers(playerIDs []types.PlayerID) {
	for _, playerID := range playerIDs {
		obj := resolv.NewObject(0, 0, float64(o.playerWidth), float64(o.playerHeight))
		obj.SetShape(resolv.NewRectangle(0, 0, o.playerWidth, o.playerHeight))
		o.space.Add(obj)
		o.players[playerID] = obj
	}
}

func (o *observer) MovePlayer(playerID types.PlayerID, x float64, y float64) bool {
	player := o.players[playerID]
	if player == nil {
		return false
	}

	if collision := player.Check(x, y); collision != nil {
		x = collision.ContactWithObject(collision.Objects[0]).X()
		y = collision.ContactWithObject(collision.Objects[0]).Y()
	}

	player.X = x
	player.Y = y
	player.Update()

	return true
}
