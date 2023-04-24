package tool

import (
	"errors"
	"fmt"
	"math"
	"time"
	"uwwolf/util/helper"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"github.com/samber/lo"
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
	position   *orb.Point
	width      int
	height     int
	isSolid    bool
	speed      float64
	lastMoveAt time.Duration
}

type MapData struct {
	TileWidth     int
	TileHeight    int
	TilePositions []orb.Point
	Obstacles     []orb.Bound
}

var mapData = MapData{
	TileWidth:     64,
	TileHeight:    64,
	TilePositions: []orb.Point{},
	Obstacles:     []orb.Bound{},
}

type gameMap struct {
	tree             *quadtree.Quadtree
	entityPositions  map[EntityID]*orb.Point
	entityByPosition map[*orb.Point]*Entity
	timePerEachMove  time.Duration
	searchExpansion  float64
}

type Map interface {
	Entity(ID EntityID) *Entity

	EntityInArea(area orb.Bound) []*Entity

	AddEntity(IDPerType string, settings EntitySettings) (bool, error)

	RemoveEntity(ID EntityID) bool

	MoveEntity(ID EntityID, position orb.Point) (bool, error)
}

func NewMap() Map {
	m := &gameMap{
		tree: quadtree.New(orb.Bound{
			Min: orb.Point{0, 0},
			Max: orb.Point{float64(mapData.TileWidth), float64(mapData.TileHeight)},
		}),
		entityPositions:  make(map[EntityID]*orb.Point),
		entityByPosition: make(map[*orb.Point]*Entity),
		timePerEachMove:  500 * time.Microsecond,
	}
	for i, tile := range mapData.TilePositions {
		m.AddEntity(fmt.Sprintf("%d", i), EntitySettings{
			Type:    ObstacleEntity,
			X:       tile.X(),
			Y:       tile.Y(),
			Width:   mapData.TileWidth,
			Height:  mapData.TileHeight,
			IsSolid: true,
		})
	}
	for i, obstacle := range mapData.Obstacles {
		m.AddEntity(fmt.Sprintf("%d", i), EntitySettings{
			Type:    TileEntity,
			X:       obstacle.Left(),
			Y:       obstacle.Top(),
			Width:   int(obstacle.Right()),
			Height:  int(obstacle.Bottom()),
			IsSolid: true,
		})
	}

	return m
}

func (m *gameMap) Entity(ID EntityID) *Entity {
	position := m.entityPositions[ID]
	if position == nil {
		return nil
	}
	return m.entityByPosition[position]
}

func (m gameMap) EntityInArea(area orb.Bound) []*Entity {
	// Expand search area
	expandedArea := area.Extend(orb.Point{
		area.Right() + m.searchExpansion, area.Bottom() + m.searchExpansion,
	})
	expandedArea = expandedArea.Extend(orb.Point{
		area.Left() - m.searchExpansion, area.Top() - m.searchExpansion,
	})

	candidateEntities := lo.Map(
		m.tree.InBound([]orb.Pointer{}, expandedArea),
		func(p orb.Pointer, _ int) *Entity {
			return m.entityByPosition[p.(*orb.Point)]
		})

	// Filter out-of-range entiies
	return lo.Filter(candidateEntities, func(e *Entity, _ int) bool {
		return area.Intersects(orb.Bound{
			Min: *e.position,
			Max: orb.Point{
				e.position.X() + float64(e.width),
				e.position.Y() + float64(e.height),
			},
		})
	})
}

func (m *gameMap) AddEntity(IDPerType string, settings EntitySettings) (bool, error) {
	ID := EntityID(fmt.Sprintf("%v_%v", settings.Type, IDPerType))
	if m.Entity(ID) != nil {
		return false, errors.New("Entity already existed!")
	} else if settings.X < 0 || settings.Y < 0 ||
		settings.Width > int(m.tree.Bound().Max.X()) || settings.Height > int(m.tree.Bound().Max.Y()) {
		return false, errors.New("Invalid postion ore size")
	}

	entity := &Entity{
		position:   &orb.Point{settings.X, settings.Y},
		width:      settings.Width,
		height:     settings.Height,
		isSolid:    settings.IsSolid,
		speed:      settings.Speed,
		lastMoveAt: 0,
	}
	m.entityPositions[ID] = entity.position
	m.entityByPosition[entity.position] = entity
	m.tree.Add(entity.position)

	diagonal := helper.CalculateDiagonal(float64(entity.width), float64(entity.height))
	if diagonal > m.searchExpansion {
		m.searchExpansion = diagonal
	}

	return true, nil
}

func (m *gameMap) RemoveEntity(ID EntityID) bool {
	entity := m.Entity(ID)
	if entity == nil {
		return false
	}

	m.tree.Remove(entity.position, func(p orb.Pointer) bool {
		return p == entity.position
	})
	delete(m.entityByPosition, entity.position)
	delete(m.entityPositions, ID)

	diagonal := helper.CalculateDiagonal(float64(entity.width), float64(entity.height))
	if diagonal == m.searchExpansion {
		diagonal = 0
		for _, e := range m.entityByPosition {
			if d := helper.CalculateDiagonal(float64(e.width), float64(e.height)); d > diagonal {
				diagonal = d
			}
		}
		m.searchExpansion = diagonal
	}

	return true
}

func (m *gameMap) MoveEntity(ID EntityID, position orb.Point) (bool, error) {
	entity := m.Entity(ID)
	if entity == nil {
		return false, errors.New("Entity doesn't exist!")
	}

	now := time.Now().UnixMilli()
	if now-entity.lastMoveAt.Milliseconds() < int64(m.timePerEachMove) {
		return false, errors.New("Move too fast!")
	}

	dx := position.X() - entity.position.X()
	dy := position.Y() - entity.position.Y()
	if math.Sqrt(dx*dx+dy*dy)/m.timePerEachMove.Seconds() > entity.speed {
		return false, errors.New("Invalid position!")
	}

	if m.checkCollision(*entity, position) {
		return false, errors.New("Collided!")
	}

	m.tree.Remove(entity.position, func(p orb.Pointer) bool {
		return p == entity.position
	})
	entity.position = &position
	m.tree.Add(entity.position)
	entity.lastMoveAt = time.Duration(now * int64(time.Millisecond))

	return true, nil
}

func (m gameMap) checkCollision(entity Entity, position orb.Point) bool {
	entities := m.EntityInArea(orb.Bound{
		Min: orb.Point{position.X(), position.Y() + float64(entity.height)},
		Max: orb.Point{position.X() + float64(entity.width), position.Y()},
	})
	for i := 0; i < len(entities); i++ {
		if entities[i].isSolid {
			return false
		}
	}

	return true
}
