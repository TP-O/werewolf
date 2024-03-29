package logic

import (
	"errors"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"

	"github.com/paulmach/orb"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// player represents the player in a game.
type player struct {
	// id is the player ID.
	id types.PlayerId

	// mainRoleId is the role ID of main role.
	mainRoleId types.RoleId

	// factionId is the faction ID of the main role.
	factionId types.FactionId

	// isDead indicates whether player is dead or not.
	isDead bool

	// moderator is the game moderator.
	moderator contract.Moderator

	// roles is the assinged roles.
	roles map[types.RoleId]contract.Role

	// records the the play records.
	records []types.PlayerRecord
}

func NewPlayer(moderator contract.Moderator, id types.PlayerId) contract.Player {
	return &player{
		id:        id,
		moderator: moderator,
		factionId: constants.VillagerFactionId,
		roles:     make(map[types.RoleId]contract.Role),
		records:   make([]types.PlayerRecord, 0),
	}
}

// Id returns player's ID.
func (p player) Id() types.PlayerId {
	return p.id
}

// MainRoleId returns player's main role ID.
func (p player) MainRoleId() types.RoleId {
	return p.mainRoleId
}

// RoleIds returns player's assigned role IDs.
func (p player) RoleIds() []types.RoleId {
	return maps.Keys(p.roles)
}

// Roles returns player's assigned roles.
func (p player) Roles() map[types.RoleId]contract.Role {
	return p.roles
}

// PlayRecords returns play records of the player.
func (p player) PlayRecords() []types.PlayerRecord {
	return p.records
}

// FactionId returns player's faction ID.
func (p player) FactionId() types.FactionId {
	return p.factionId
}

// IsDead checks if player is dead.
func (p player) IsDead() bool {
	return p.isDead
}

// Position returns the curernt location of the player.
func (p player) Location() (float64, float64) {
	// entitiy := p.world.Player(p.Id())
	// return entitiy.X, entitiy.Y
	return 1, 1
}

// SetFactionId assigns the player to the new faction.
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
	// p.moderator.World().Map().RemoveEntity(contract.EntityID(fmt.Sprintf("%v_%v", contract.PlayerEntity, p.Id())))

	return true
}

// AssignRole assigns role to the player, and the faction can
// be updated based on the role.
func (p *player) AssignRole(roleID types.RoleId) (bool, error) {
	if slices.Contains(p.RoleIds(), roleID) {
		return false, errors.New("This role is already assigned ¯\\_(ツ)_/¯")
	}

	if newRole, err := p.moderator.RoleFactory().CreateById(roleID, p.moderator, p.id); err != nil {
		return false, err
	} else {
		p.roles[roleID] = newRole
		if constants.RoleWeights.BlindGet(newRole.Id()) > constants.RoleWeights.BlindGet(p.mainRoleId) {
			p.mainRoleId = newRole.Id()
			p.factionId = newRole.FactionId()
		}
		newRole.OnAfterAssign()
	}

	return true, nil
}

// RevokeRole removes the role from the player, and updates faction
// if needed
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
				constants.RoleWeights.BlindGet(role.Id()) > constants.RoleWeights.BlindGet(newMainRole.Id()) {
				newMainRole = role
			}
		}

		p.mainRoleId = newMainRole.Id()
		p.factionId = newMainRole.FactionId()
	}

	return true, nil
}

// UseRole uses one of player's available ability.
func (p *player) UseRole(req types.RoleRequest) types.RoleResponse {
	if p.isDead {
		return types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "You're died (╥﹏╥)",
			},
		}
	} else if !p.moderator.Scheduler().CanPlay(p.id) {
		return types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "Wait for your turn, OK??",
			},
		}
	} else {
		turn := p.moderator.Scheduler().TurnSlots()
		res := p.roles[turn[p.id].RoleId].Use(req)
		p.records = append(p.records, types.PlayerRecord{
			Round:     res.Round,
			Turn:      res.Turn,
			PhaseId:   res.PhaseId,
			RoleId:    res.RoleId,
			ActionId:  res.ActionId,
			IsSkipped: res.IsSkipped,
			TargetId:  res.TargetId,
		})

		return res
	}
}

// Move moves the player to the given location.
func (p *player) Move(position orb.Point) (bool, error) {
	// return p.moderator.World().Map().MoveEntity(contract.EntityID(fmt.Sprintf("%v_%v", contract.PlayerEntity, p.Id())), position)
	return true, nil
}
