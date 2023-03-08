package game

// import (
// 	"fmt"
// 	"uwwolf/game/contract"
// 	"uwwolf/game/types"
// 	"uwwolf/game/vars"

// 	"golang.org/x/exp/maps"
// 	"golang.org/x/exp/slices"
// )

// // player represents the player in a game.
// type player struct {
// 	id         types.PlayerID
// 	factionID  types.FactionID
// 	mainRoleID types.RoleID
// 	isDead     bool
// 	game       contract.Game
// 	roles      map[types.RoleID]contract.Role
// }

// var _ contract.Player = (*player)(nil)

// func NewPlayer(game contract.Game, id types.PlayerID) contract.Player {
// 	return &player{
// 		id:        id,
// 		game:      game,
// 		factionID: vars.VillagerFactionID,
// 		roles:     make(map[types.RoleID]contract.Role),
// 	}
// }

// func (p player) ID() types.PlayerID {
// 	return p.id
// }

// func (p player) MainRoleID() types.RoleID {
// 	return p.mainRoleID
// }

// func (p player) RoleIDs() []types.RoleID {
// 	return maps.Keys(p.roles)
// }

// func (p player) Roles() map[types.RoleID]contract.Role {
// 	return p.roles
// }

// func (p player) FactionID() types.FactionID {
// 	return p.factionID
// }

// func (p player) IsDead() bool {
// 	return p.isDead
// }

// func (p *player) Die(isExited bool) bool {
// 	if p.isDead {
// 		return false
// 	}

// 	for _, role := range p.roles {
// 		if dead := role.BeforeDeath(); !dead && !isExited {
// 			return false
// 		}
// 	}

// 	p.isDead = true
// 	for _, role := range p.roles {
// 		role.AfterDeath()
// 	}

// 	return true
// }

// func (p *player) Revive() bool {
// 	if !p.isDead {
// 		return false
// 	}

// 	p.isDead = false
// 	for _, role := range p.roles {
// 		role.AfterSaved()
// 	}

// 	return true
// }

// func (p *player) SetFactionID(factionID types.FactionID) {
// 	p.factionID = factionID
// }

// func (p *player) AssignRole(roleID types.RoleID) (bool, error) {
// 	if slices.Contains(p.RoleIDs(), roleID) {
// 		return false, fmt.Errorf("This role is already assigned ¯\\_(ツ)_/¯")
// 	}

// 	if newRole, err := NewRole(roleID, p.game, p.id); err != nil {
// 		return false, err
// 	} else {
// 		p.roles[roleID] = newRole
// 		if vars.RoleWeights[newRole.ID()] > vars.RoleWeights[p.mainRoleID] {
// 			p.mainRoleID = newRole.ID()
// 			p.factionID = newRole.FactionID()
// 		}
// 	}

// 	return true, nil
// }

// func (p *player) ActivateAbility(req types.ActivateAbilityRequest) types.ActionResponse {
// 	if turn := p.game.Scheduler().Turn(); turn == nil {
// 		return types.ActionResponse{
// 			Ok:      false,
// 			Message: "The game is about to start ノ(ジ)ー'",
// 		}
// 	} else if playerTurn := turn[p.id]; playerTurn == nil ||
// 		playerTurn.FrozenLimit != vars.OutOfTimes ||
// 		p.game.Scheduler().RoundID() < playerTurn.BeginRoundID {
// 		return types.ActionResponse{
// 			Ok:      false,
// 			Message: "Wait for your turn, OK??",
// 		}
// 	} else {
// 		return p.roles[playerTurn.RoleID].ActivateAbility(req)
// 	}
// }
