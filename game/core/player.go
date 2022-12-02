package core

import (
	"errors"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/factory"
	"uwwolf/game/types"
)

type player struct {
	id         enum.PlayerID
	factionID  enum.FactionID
	mainRoleID enum.RoleID
	isDead     bool
	game       contract.Game
	roles      map[enum.RoleID]contract.Role
}

func NewPlayer(game contract.Game, id enum.PlayerID) contract.Player {
	return &player{
		id:        id,
		game:      game,
		factionID: enum.VillagerFactionID,
		roles:     make(map[enum.RoleID]contract.Role),
	}
}

func (p *player) ID() enum.PlayerID {
	return p.id
}

func (p *player) MainRoleID() enum.RoleID {
	return p.mainRoleID
}

func (p *player) RoleIDs() []enum.RoleID {
	return maps.Keys(p.roles)
}

func (p *player) Roles() map[enum.RoleID]contract.Role {
	return p.roles
}

func (p *player) FactionID() enum.FactionID {
	return p.factionID
}

func (p *player) IsDead() bool {
	return false
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

func (p *player) SetFactionID(factionID enum.FactionID) {
	p.factionID = factionID
}

func (p *player) AssignRole(roleID enum.RoleID) (bool, error) {
	if slices.Contains(p.RoleIDs(), roleID) {
		return false, errors.New("Non-existent role ID ¯\\_(ツ)_/¯")
	}

	if newRole, err := factory.NewRole(roleID, p.game, p.id); err != nil {
		return false, err
	} else {
		p.roles[roleID] = newRole

		if types.RoleIDRanks[newRole.ID()] > types.RoleIDRanks[p.mainRoleID] {
			p.mainRoleID = newRole.ID()
			p.factionID = newRole.FactionID()
		}
	}

	return true, nil
}

func (p *player) RevokeRole(roleID enum.RoleID) (bool, error) {
	if len(p.roles) == 1 {
		return false, errors.New("Player must player at least one role ヾ(⌐■_■)ノ♪")
	} else if p.roles[roleID] == nil {
		return false, errors.New("Non-existent role ID  ¯\\_(ツ)_/¯")
	}

	delete(p.roles, roleID)
	var newMainRole contract.Role

	for _, role := range p.roles {
		if newMainRole == nil ||
			types.RoleIDRanks[role.ID()] > types.RoleIDRanks[newMainRole.ID()] {
			newMainRole = role
		}
	}

	p.mainRoleID = newMainRole.ID()
	p.factionID = newMainRole.FactionID()

	return true, nil
}

func (p *player) UseAbility(req *types.UseRoleRequest) *types.ActionResponse {
	if turn := p.game.Scheduler().Turn(); turn != nil && p.roles[turn.RoleID] != nil {
		return p.roles[turn.RoleID].UseAbility(req)
	}

	return &types.ActionResponse{
		Ok:      false,
		Message: "Wait for your turn, OK??",
	}
}
