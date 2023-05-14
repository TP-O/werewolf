package logic

import (
	"errors"
	"fmt"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/role"
	"uwwolf/internal/app/game/logic/types"

	"github.com/paulmach/orb"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// type playerRecord struct {
// 	RoundId  types.RoundID
// 	TurnId   types.TurnId
// 	RoleId   types.RoleId
// 	TargetId types.PlayerId
// }

// player represents the player in a game.
type player struct {
	id         types.PlayerId
	factionId  types.FactionId
	mainRoleId types.RoleId
	isDead     bool
	moderator  contract.Moderator
	roles      map[types.RoleId]contract.Role
	// records    []playerRecord
}

func NewPlayer(moderator contract.Moderator, id types.PlayerId) contract.Player {
	return &player{
		id:        id,
		moderator: moderator,
		factionId: constants.VillagerFactionId,
		roles:     make(map[types.RoleId]contract.Role),
		// records:   make([]playerRecord, 0),
	}
}

// ID returns player's ID.
func (p player) Id() types.PlayerId {
	return p.id
}

// MainRoleID returns player's main role id.
func (p player) MainRoleId() types.RoleId {
	return p.mainRoleId
}

// RoleIDs returns player's assigned role ids.
func (p player) RoleIds() []types.RoleId {
	return maps.Keys(p.roles)
}

// Roles returns player's assigned roles.
func (p player) Roles() map[types.RoleId]contract.Role {
	return p.roles
}

// FactionID returns player's faction ID.
func (p player) FactionId() types.FactionId {
	return p.factionId
}

// IsDead checks if player is dead.
func (p player) IsDead() bool {
	return p.isDead
}

func (p player) Location() (float64, float64) {
	// entitiy := p.world.Player(p.Id())
	// return entitiy.X, entitiy.Y
	return 1, 1
}

// SetFactionID assigns this player to the new faction.
func (p *player) SetFactionId(factionID types.FactionId) {
	p.factionId = factionID
}

// Die kills the player and triggers roles events.
func (p *player) Die() bool {
	return p.die(false)
}

// Exit kills the player and ignores any trigger preventing death.
func (p *player) Exit() bool {
	return p.die(true)
}

func (p *player) die(isExited bool) bool {
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
		role.OnAfterRevoke()
	}
	p.moderator.World().Map().RemoveEntity(contract.EntityID(fmt.Sprintf("%v_%v", contract.PlayerEntity, p.Id())))

	return true
}

// AssignRole assigns the role to the player, and the faction can
// be updated based on this role.
func (p *player) AssignRole(roleID types.RoleId) (bool, error) {
	if slices.Contains(p.RoleIds(), roleID) {
		return false, fmt.Errorf("This role is already assigned ¯\\_(ツ)_/¯")
	}

	if newRole, err := role.NewRole(roleID, p.moderator, p.id); err != nil {
		return false, err
	} else {
		p.roles[roleID] = newRole
		if constants.RoleWeights.BindGet(newRole.Id()) > constants.RoleWeights.BindGet(p.mainRoleId) {
			p.mainRoleId = newRole.Id()
			p.factionId = newRole.FactionId()
		}
		newRole.OnAfterAssign()
	}

	return true, nil
}

// RevokeRole removes the role from the player, and the faction can
// be updated based on removed role.
func (p *player) RevokeRole(roleID types.RoleId) (bool, error) {
	if len(p.roles) == 1 {
		return false, errors.New("Player must player at least one role ヾ(⌐■_■)ノ♪")
	} else if p.roles[roleID] == nil {
		return false, errors.New("Non-existent role ID  ¯\\_(ツ)_/¯")
	}

	p.roles[roleID].OnAfterRevoke()
	delete(p.roles, roleID)

	if roleID == p.mainRoleId {
		var newMainRole contract.Role

		for _, role := range p.roles {
			if newMainRole == nil ||
				constants.RoleWeights.BindGet(role.Id()) > constants.RoleWeights.BindGet(newMainRole.Id()) {
				newMainRole = role
			}
		}

		p.mainRoleId = newMainRole.Id()
		p.factionId = newMainRole.FactionId()
	}

	return true, nil
}

// ActivateAbility executes one of player's available ability.
// The executed ability is selected based on the requested
// action.
func (p *player) ActivateAbility(req *types.RoleRequest) *types.RoleResponse {
	if p.isDead {
		return &types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "You're died (╥﹏╥)",
			},
		}
	} else if !p.moderator.World().Scheduler().CanPlay(p.id) {
		return &types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "Wait for your turn, OK??",
			},
		}
	} else {
		turn := p.moderator.World().Scheduler().TurnSlots()
		res := p.roles[turn[p.id].RoleId].Use(*req)
		// p.records = append(p.records, playerRecord{
		// 	RoundId:  p.moderator.Scheduler().RoundID(),
		// 	TurnId:   p.moderator.World().Scheduler().TurnID(),
		// 	RoleId:   1,
		// 	TargetId: res.TargetId,
		// })

		return &res
	}
}

func (p *player) Move(position orb.Point) (bool, error) {
	return p.moderator.World().Map().MoveEntity(contract.EntityID(fmt.Sprintf("%v_%v", contract.PlayerEntity, p.Id())), position)
}
