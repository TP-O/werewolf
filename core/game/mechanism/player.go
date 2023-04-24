package mechanism

import (
	"errors"
	"fmt"
	"uwwolf/game/declare"
	"uwwolf/game/mechanism/contract"
	"uwwolf/game/mechanism/role"
	"uwwolf/game/tool"
	"uwwolf/game/types"

	"github.com/paulmach/orb"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// player represents the player in a game.
type player struct {
	id         types.PlayerID
	factionID  types.FactionID
	mainRoleID types.RoleID
	isDead     bool
	world      contract.World
	roles      map[types.RoleID]contract.Role
}

func NewPlayer(world contract.World, id types.PlayerID) contract.Player {
	return &player{
		id:        id,
		world:     world,
		factionID: declare.VillagerFactionID,
		roles:     make(map[types.RoleID]contract.Role),
	}
}

// ID returns player's ID.
func (p player) ID() types.PlayerID {
	return p.id
}

// MainRoleID returns player's main role id.
func (p player) MainRoleID() types.RoleID {
	return p.mainRoleID
}

// RoleIDs returns player's assigned role ids.
func (p player) RoleIDs() []types.RoleID {
	return maps.Keys(p.roles)
}

// Roles returns player's assigned roles.
func (p player) Roles() map[types.RoleID]contract.Role {
	return p.roles
}

// FactionID returns player's faction ID.
func (p player) FactionID() types.FactionID {
	return p.factionID
}

// IsDead checks if player is dead.
func (p player) IsDead() bool {
	return p.isDead
}

func (p player) Location() (float64, float64) {
	// entitiy := p.world.Player(p.ID())
	// return entitiy.X, entitiy.Y
	return 1, 1
}

// SetFactionID assigns this player to the new faction.
func (p *player) SetFactionID(factionID types.FactionID) {
	p.factionID = factionID
}

// Die marks this player as dead and triggers roles events.
// If `isExited` is true, any trigger preventing death is ignored.
func (p *player) Die(isExited bool) bool {
	if p.isDead {
		return false
	}

	for _, role := range p.roles {
		if isDead := role.OnBeforeDeath(); !isDead && !isExited {
			return false
		}
	}

	p.isDead = true
	for _, role := range p.roles {
		role.OnAfterDeath()
		role.OnRevoke()
	}
	p.world.Map().RemoveEntity(tool.EntityID(fmt.Sprintf("%v_%v", tool.PlayerEntity, p.ID())))

	return true
}

// AssignRole assigns the role to the player, and the faction can
// be updated based on this role.
func (p *player) AssignRole(roleID types.RoleID) (bool, error) {
	if slices.Contains(p.RoleIDs(), roleID) {
		return false, fmt.Errorf("This role is already assigned ¯\\_(ツ)_/¯")
	}

	if newRole, err := role.NewRole(roleID, p.world, p.id); err != nil {
		return false, err
	} else {
		p.roles[roleID] = newRole
		if declare.RoleWeights.BindGet(newRole.ID()) > declare.RoleWeights.BindGet(p.mainRoleID) {
			p.mainRoleID = newRole.ID()
			p.factionID = newRole.FactionID()
		}
		newRole.OnAssign()
	}

	return true, nil
}

// RevokeRole removes the role from the player, and the faction can
// be updated based on removed role.
func (p *player) RevokeRole(roleID types.RoleID) (bool, error) {
	if len(p.roles) == 1 {
		return false, errors.New("Player must player at least one role ヾ(⌐■_■)ノ♪")
	} else if p.roles[roleID] == nil {
		return false, errors.New("Non-existent role ID  ¯\\_(ツ)_/¯")
	}

	p.roles[roleID].OnRevoke()
	delete(p.roles, roleID)

	if roleID == p.mainRoleID {
		var newMainRole contract.Role

		for _, role := range p.roles {
			if newMainRole == nil ||
				declare.RoleWeights.BindGet(role.ID()) > declare.RoleWeights.BindGet(newMainRole.ID()) {
				newMainRole = role
			}
		}

		p.mainRoleID = newMainRole.ID()
		p.factionID = newMainRole.FactionID()
	}

	return true, nil
}

// ActivateAbility executes one of player's available ability.
// The executed ability is selected based on the requested
// action.
func (p *player) ActivateAbility(req *types.ActivateAbilityRequest) *types.ActionResponse {
	if p.isDead {
		return &types.ActionResponse{
			Ok:      false,
			Message: "You're died (╥﹏╥)",
		}
	} else if !p.world.Scheduler().CanPlay(p.id) {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Wait for your turn, OK??",
		}
	} else {
		turn := p.world.Scheduler().Turn()
		return p.roles[turn[p.id].RoleID].ActivateAbility(req)
	}
}

func (p *player) Move(position orb.Point) (bool, error) {
	return p.world.Map().MoveEntity(tool.EntityID(fmt.Sprintf("%v_%v", tool.PlayerEntity, p.ID())), position)
}
