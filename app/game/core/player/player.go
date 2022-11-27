package player

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/factory"
	"uwwolf/app/game/types"
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
		id:    id,
		game:  game,
		roles: make(map[types.RoleID]contract.Role),
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
	return false
}

func (p *player) Die() {
	p.isDead = true
}

func (p *player) Revive() {
	p.isDead = false
}

func (p *player) SetFactionID(factionID types.FactionID) {
	p.factionID = factionID
}

func (p *player) AssignRoles(roleIDs ...types.RoleID) {
	for _, roleID := range roleIDs {
		if !slices.Contains(p.RoleIDs(), roleID) {
			newRole := factory.NewRole(roleID, p.game, p.id)
			p.roles[roleID] = newRole

			if config.RolePriorities[newRole.ID()] > config.RolePriorities[p.mainRoleID] {
				p.mainRoleID = newRole.ID()
				p.factionID = newRole.FactionID()
			}
		}
	}
}

func (p *player) RevokeRoles(roleIDs ...types.RoleID) bool {
	if len(p.roles) == 0 || len(roleIDs) == len(p.roles) {
		return false
	}

	for _, roleID := range roleIDs {
		delete(p.roles, roleID)
	}

	var newMainRole contract.Role

	for _, role := range p.roles {
		if newMainRole == nil ||
			config.RolePriorities[role.ID()] > config.RolePriorities[newMainRole.ID()] {

			newMainRole = role
		}
	}

	p.mainRoleID = newMainRole.ID()
	p.factionID = newMainRole.FactionID()

	return true
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
