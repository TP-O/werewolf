package game

import (
	"errors"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"uwwolf/game/contract"
	"uwwolf/game/types"
)

const (
	OfflineStatus types.PlayerStatus = iota
	OnlineStatus
	BusyStatus
	InGameStatus
)

type player struct {
	id         types.PlayerID
	factionID  types.FactionID
	mainRoleID types.RoleID
	isDead     bool
	game       contract.Game
	roles      map[types.RoleID]contract.Role
}

func NewPlayer(game contract.Game, id types.PlayerID) contract.Player {
	return &player{
		id:        id,
		game:      game,
		factionID: VillagerFactionID,
		roles:     make(map[types.RoleID]contract.Role),
	}
}

func (p *player) ID() types.PlayerID {
	return p.id
}

func (p *player) MainRoleID() types.RoleID {
	return p.mainRoleID
}

func (p *player) RoleIDs() []types.RoleID {
	return maps.Keys(p.roles)
}

func (p *player) Roles() map[types.RoleID]contract.Role {
	return p.roles
}

func (p *player) FactionID() types.FactionID {
	return p.factionID
}

func (p *player) IsDead() bool {
	return p.isDead
}

func (p *player) Die(isExited bool) bool {
	if p.isDead {
		return false
	}

	for _, role := range p.roles {
		if dead := role.BeforeDeath(); !dead && !isExited {
			return false
		}
	}

	p.isDead = true

	for _, role := range p.roles {
		role.AfterDeath()
	}

	return true
}

func (p *player) Revive() bool {
	if !p.isDead {
		return false
	}

	p.isDead = false

	return true
}

func (p *player) SetFactionID(factionID types.FactionID) {
	p.factionID = factionID
}

func (p *player) AssignRole(roleID types.RoleID) (bool, error) {
	if slices.Contains(p.RoleIDs(), roleID) {
		return false, errors.New("This role is already assigned ¯\\_(ツ)_/¯")
	}

	if newRole, err := NewRole(roleID, p.game, p.id); err != nil {
		return false, err
	} else {
		p.roles[roleID] = newRole

		if RoleIDWeights[newRole.ID()] > RoleIDWeights[p.mainRoleID] {
			p.mainRoleID = newRole.ID()
			p.factionID = newRole.FactionID()
		}
	}

	return true, nil
}

func (p *player) RevokeRole(roleID types.RoleID) (bool, error) {
	if len(p.roles) == 1 {
		return false, errors.New("Player must player at least one role ヾ(⌐■_■)ノ♪")
	} else if p.roles[roleID] == nil {
		return false, errors.New("Non-existent role ID  ¯\\_(ツ)_/¯")
	}

	delete(p.roles, roleID)
	var newMainRole contract.Role

	for _, rolee := range p.roles {
		if newMainRole == nil ||
			RoleIDWeights[rolee.ID()] > RoleIDWeights[newMainRole.ID()] {
			newMainRole = rolee
		}
	}

	p.mainRoleID = newMainRole.ID()
	p.factionID = newMainRole.FactionID()

	return true, nil
}

func (p *player) ExecuteAction(req types.ExecuteActionRequest) types.ActionResponse {
	if turn := p.game.Scheduler().Turn(); turn == nil {
		return types.ActionResponse{
			Ok:      false,
			Message: "Wait until game starts ノ(ジ)ー'",
		}
	} else if p.roles[turn.RoleID] == nil {
		return types.ActionResponse{
			Ok:      false,
			Message: "Wait for your turn, OK??",
		}
	} else {
		return p.roles[turn.RoleID].ExecuteAction(req)
	}
}
