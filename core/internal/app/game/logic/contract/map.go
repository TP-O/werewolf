package contract

import (
	"time"

	"github.com/paulmach/orb"
)

type EntityID string

type EntityType = string

const (
	PlayerEntity   = "P"
	ObstacleEntity = "O"
	TileEntity     = "T"
)

type EntitySettings struct {
	Type    EntityType
	X       float64
	Y       float64
	Width   int
	Height  int
	IsSolid bool
	Speed   float64
}

type Entity struct {
	Position   *orb.Point
	Width      int
	Height     int
	IsSolid    bool
	Speed      float64
	LastMoveAt time.Duration
}

type Map interface {
	Entity(ID EntityID) *Entity

	EntityInArea(area orb.Bound) []*Entity

	AddEntity(IDPerType string, settings EntitySettings) (bool, error)

	RemoveEntity(ID EntityID) bool

	MoveEntity(ID EntityID, position orb.Point) (bool, error)
}
