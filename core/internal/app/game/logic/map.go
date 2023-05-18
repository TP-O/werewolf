package logic

import (
	"errors"
	"fmt"
	"math"
	"time"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/pkg/util"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"github.com/samber/lo"
)

type mapData struct {
	TileWidth     int
	TileHeight    int
	TilePositions []orb.Point
	Obstacles     []orb.Bound
}

// gameMap is the instance managing entities position of the game.
type gameMap struct {
	// tree is the quad tree to manage object positions.
	tree *quadtree.Quadtree

	// data is the map data.
	data mapData

	// entityPositions is the position of all entities in the map.
	entityPositions map[contract.EntityID]*orb.Point

	// entityByPosition is entity at all stored positions in the map.
	entityByPosition map[*orb.Point]*contract.Entity

	// timePerEachMove is the max time to move from the current position to the new one.
	timePerEachMove time.Duration

	// searchExpansion is the added area bounding around the searched area to make sure
	// all objects in it are detected. We need this because the quad tree only stores
	// the coordinator of object's top left corner, so the expansion will help us to reach
	// that point.
	searchExpansion float64
}

func NewMap() contract.Map {
	data := mapData{
		TileWidth:     64,
		TileHeight:    64,
		TilePositions: []orb.Point{},
		Obstacles:     []orb.Bound{},
	}
	m := &gameMap{
		tree: quadtree.New(orb.Bound{
			Min: orb.Point{0, 0},
			Max: orb.Point{float64(data.TileWidth), float64(data.TileHeight)},
		}),
		data:             data,
		entityPositions:  make(map[contract.EntityID]*orb.Point),
		entityByPosition: make(map[*orb.Point]*contract.Entity),
		timePerEachMove:  500 * time.Microsecond,
	}
	for i, tile := range data.TilePositions {
		m.AddEntity(fmt.Sprintf("%d", i), contract.EntitySettings{
			Type:    contract.ObstacleEntity,
			X:       tile.X(),
			Y:       tile.Y(),
			Width:   data.TileWidth,
			Height:  data.TileHeight,
			IsSolid: true,
		})
	}
	for i, obstacle := range data.Obstacles {
		m.AddEntity(fmt.Sprintf("%d", i), contract.EntitySettings{
			Type:    contract.TileEntity,
			X:       obstacle.Left(),
			Y:       obstacle.Top(),
			Width:   int(obstacle.Right()),
			Height:  int(obstacle.Bottom()),
			IsSolid: true,
		})
	}

	return m
}

func (m *gameMap) Entity(ID contract.EntityID) *contract.Entity {
	position := m.entityPositions[ID]
	if position == nil {
		return nil
	}
	return m.entityByPosition[position]
}

func (m gameMap) EntityInArea(area orb.Bound) []*contract.Entity {
	// Expand search area
	expandedArea := area.Extend(orb.Point{
		area.Right() + m.searchExpansion, area.Bottom() + m.searchExpansion,
	})
	expandedArea = expandedArea.Extend(orb.Point{
		area.Left() - m.searchExpansion, area.Top() - m.searchExpansion,
	})

	candidateEntities := lo.Map(
		m.tree.InBound([]orb.Pointer{}, expandedArea),
		func(p orb.Pointer, _ int) *contract.Entity {
			return m.entityByPosition[p.(*orb.Point)]
		})

	// Filter out-of-range entiies
	return lo.Filter(candidateEntities, func(e *contract.Entity, _ int) bool {
		return area.Intersects(orb.Bound{
			Min: *e.Position,
			Max: orb.Point{
				e.Position.X() + float64(e.Width),
				e.Position.Y() + float64(e.Height),
			},
		})
	})
}

func (m *gameMap) AddEntity(IDPerType string, settings contract.EntitySettings) (bool, error) {
	ID := contract.EntityID(fmt.Sprintf("%v_%v", settings.Type, IDPerType))
	if m.Entity(ID) != nil {
		return false, errors.New("Entity already existed!")
	} else if settings.X < 0 || settings.Y < 0 ||
		settings.Width > int(m.tree.Bound().Max.X()) || settings.Height > int(m.tree.Bound().Max.Y()) {
		return false, errors.New("Invalid postion ore size")
	}

	entity := &contract.Entity{
		Position:   &orb.Point{settings.X, settings.Y},
		Width:      settings.Width,
		Height:     settings.Height,
		IsSolid:    settings.IsSolid,
		Speed:      settings.Speed,
		LastMoveAt: 0,
	}
	m.entityPositions[ID] = entity.Position
	m.entityByPosition[entity.Position] = entity
	m.tree.Add(entity.Position)

	diagonal := util.CalculateDiagonal(float64(entity.Width), float64(entity.Height))
	if diagonal > m.searchExpansion {
		m.searchExpansion = diagonal
	}

	return true, nil
}

func (m *gameMap) RemoveEntity(ID contract.EntityID) bool {
	entity := m.Entity(ID)
	if entity == nil {
		return false
	}

	m.tree.Remove(entity.Position, func(p orb.Pointer) bool {
		return p == entity.Position
	})
	delete(m.entityByPosition, entity.Position)
	delete(m.entityPositions, ID)

	diagonal := util.CalculateDiagonal(float64(entity.Width), float64(entity.Height))
	if diagonal == m.searchExpansion {
		diagonal = 0
		for _, e := range m.entityByPosition {
			if d := util.CalculateDiagonal(float64(e.Width), float64(e.Height)); d > diagonal {
				diagonal = d
			}
		}
		m.searchExpansion = diagonal
	}

	return true
}

func (m *gameMap) MoveEntity(ID contract.EntityID, position orb.Point) (bool, error) {
	entity := m.Entity(ID)
	if entity == nil {
		return false, errors.New("Entity doesn't exist!")
	}

	now := time.Now().UnixMilli()
	if now-entity.LastMoveAt.Milliseconds() < int64(m.timePerEachMove) {
		return false, errors.New("Move too fast!")
	}
	dx := position.X() - entity.Position.X()
	dy := position.Y() - entity.Position.Y()
	if math.Sqrt(dx*dx+dy*dy)/m.timePerEachMove.Seconds() > entity.Speed {
		return false, errors.New("Invalid position!")
	}

	if m.checkCollision(*entity, position) {
		return false, errors.New("Collided!")
	}

	m.tree.Remove(entity.Position, func(p orb.Pointer) bool {
		return p == entity.Position
	})
	entity.Position = &position
	m.tree.Add(entity.Position)
	entity.LastMoveAt = time.Duration(now * int64(time.Millisecond))

	return true, nil
}

func (m gameMap) checkCollision(entity contract.Entity, position orb.Point) bool {
	entities := m.EntityInArea(orb.Bound{
		Min: orb.Point{position.X(), position.Y() + float64(entity.Height)},
		Max: orb.Point{position.X() + float64(entity.Width), position.Y()},
	})
	for i := 0; i < len(entities); i++ {
		if entities[i].IsSolid {
			return false
		}
	}

	return true
}
