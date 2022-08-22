package core

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/exp/slices"

	"uwwolf/cache"
	"uwwolf/db"
	"uwwolf/module/game/contract"
	"uwwolf/module/game/factory"
	"uwwolf/module/game/model"
	"uwwolf/module/game/state"
	"uwwolf/types"
	"uwwolf/util"
)

type game struct {
	id                 types.GameId
	isStarted          bool
	capacity           int
	numberOfWerewolves int
	timeForTurn        time.Duration
	timeForDiscussion  time.Duration
	rolePool           []types.RoleId
	factions           map[types.FactionId][]types.PlayerId
	players            map[types.PlayerId]contract.Player
	deaths             []types.PlayerId
	polls              map[types.FactionId]*state.Poll
	round              *state.Round
}

type roleSplit struct {
	werewolf        []*model.Role
	others          []*model.Role
	reserveWerewolf *model.Role
	reserveVillager *model.Role
}

func NewGame(setting *types.GameSetting) contract.Game {
	game := game{
		id:                 setting.Id,
		capacity:           len(setting.PlayerIds),
		numberOfWerewolves: setting.NumberOfWerewolves,
		timeForTurn:        setting.TimeForTurn,
		timeForDiscussion:  setting.TimeForDiscussion,
		rolePool:           setting.RolePool,
		factions:           make(map[types.FactionId][]types.PlayerId),
		players:            make(map[types.PlayerId]contract.Player),
		deaths:             make([]types.PlayerId, len(setting.PlayerIds)),
		polls:              make(map[types.FactionId]*state.Poll),
		round:              state.NewRound(),
	}

	for _, id := range setting.PlayerIds {
		game.players[id] = NewPlayer(&game, id)
	}

	return &game
}

func (g *game) IsStarted() bool {
	return g.isStarted
}

func (g *game) GetCurrentRoundId() types.RoundId {
	return g.round.GetCurrentId()
}

func (g *game) GetCurrentRoleId() types.RoleId {
	return g.round.GetCurrentTurn().RoleId()
}

func (g *game) GetCurrentPhaseId() types.PhaseId {
	return g.round.GetCurrentPhaseId()
}

func (g *game) GetPlayer(playerId types.PlayerId) contract.Player {
	return g.players[playerId]
}

func (g *game) GetPoll(factionId types.FactionId) *state.Poll {
	return g.polls[factionId]
}

func (g *game) Start() bool {
	if g.IsStarted() {
		return false
	}

	roleSplit := g.getRolesByIds(g.rolePool)
	g.prepareRound(roleSplit)

	roles := g.randomRoles(roleSplit)
	g.assignRoles(roles)

	g.polls[types.VillagerFaction] = state.NewPoll(g.factions[types.VillagerFaction])
	g.polls[types.WerewolfFaction] = state.NewPoll(g.factions[types.WerewolfFaction])

	return true
}

func (g *game) randomRoles(roleSplit *roleSplit) []*model.Role {
	randomRoles := append(
		g.pickUpRoles(
			g.numberOfWerewolves,
			roleSplit.werewolf,
			roleSplit.reserveWerewolf,
		),
		g.pickUpRoles(
			g.capacity-g.numberOfWerewolves,
			roleSplit.others,
			roleSplit.reserveVillager,
		)...,
	)

	return randomRoles
}

// Return list of roles which can be duplicate because of
// property Set's value.
func (g *game) pickUpRoles(slots int, roles []*model.Role, reserveRole *model.Role) []*model.Role {
	var pickedUpRoles []*model.Role
	var randomRoles []*model.Role

	// Pick roles randomly
	for i := 0; i < slots; i++ {
		index, role := util.RandomElement(roles)

		if index == -1 {
			break
		}

		pickedUpRoles = append(pickedUpRoles, role)
		roles = slices.Delete(roles, index, index+1)
	}

	// Spread random roles based on Set property
	for i := 0; i < slots; i++ {
		index, role := util.RandomElement(pickedUpRoles)

		if index == -1 {
			randomRoles = append(randomRoles, reserveRole)
		} else {
			randomRoles = append(randomRoles, role)

			if role.Set == 1 {
				pickedUpRoles = slices.Delete(pickedUpRoles, index, index+1)
			} else {
				role.Set--
			}
		}
	}

	return randomRoles
}

// Get roles by ids then split them to 2 faction Werewolf and the rest
func (g *game) getRolesByIds(ids []types.RoleId) *roleSplit {
	var roles []*model.Role
	cacheRoles := cache.Local().Get(types.RoleCacheKey)

	// Cache roles
	if cacheRoles == nil || cacheRoles.IsExpired() {
		db.Client().Order("id").Find(&roles)
		cache.Local().Set(types.RoleCacheKey, roles, ttlcache.DefaultTTL)
	} else {
		roles = cacheRoles.Value().([]*model.Role)
	}

	roleSplit := &roleSplit{
		reserveWerewolf: roles[types.WerewolfRole-1],
		reserveVillager: roles[types.VillagerRole-1],
	}

	// Add role to split if id is valid, also skip 2 reserve roles
	for _, role := range roles {
		if role.ID != types.WerewolfRole &&
			role.ID != types.VillagerRole &&
			slices.Contains(ids, role.ID) {

			if role.FactionID == types.WerewolfFaction {
				roleSplit.werewolf = append(roleSplit.werewolf, role)
			} else {
				roleSplit.others = append(roleSplit.others, role)
			}
		}

		// Break if enough roles
		if len(roleSplit.werewolf)+len(roleSplit.others) == len(ids) {
			break
		}
	}

	return roleSplit
}

func (g *game) prepareRound(roleSplit *roleSplit) {
	g.round.AddTurn(&types.TurnSetting{
		PhaseId:    roleSplit.reserveWerewolf.PhaseID,
		RoleId:     roleSplit.reserveWerewolf.ID,
		BeginRound: roleSplit.reserveWerewolf.BeginRound,
		Priority:   roleSplit.reserveWerewolf.Priority,
		Expiration: roleSplit.reserveWerewolf.Expiration,
		Position:   types.SortedPosition,
	})
	g.round.AddTurn(&types.TurnSetting{
		PhaseId:    roleSplit.reserveVillager.PhaseID,
		RoleId:     roleSplit.reserveVillager.ID,
		BeginRound: roleSplit.reserveVillager.BeginRound,
		Priority:   roleSplit.reserveVillager.Priority,
		Expiration: roleSplit.reserveVillager.Expiration,
		Position:   types.SortedPosition,
	})

	for _, role := range roleSplit.werewolf {
		g.round.AddTurn(&types.TurnSetting{
			PhaseId:    role.PhaseID,
			RoleId:     role.ID,
			BeginRound: role.BeginRound,
			Priority:   role.Priority,
			Expiration: role.Expiration,
			Position:   types.SortedPosition,
		})
	}
	for _, role := range roleSplit.others {
		g.round.AddTurn(&types.TurnSetting{
			PhaseId:    role.PhaseID,
			RoleId:     role.ID,
			BeginRound: role.BeginRound,
			Priority:   role.Priority,
			Expiration: role.Expiration,
			Position:   types.SortedPosition,
		})
	}
}

func (g *game) assignRoles(roles []*model.Role) {
	for _, player := range g.players {
		index, role := util.RandomElement(roles)

		player.AssignRoles(
			factory.Role(role.ID, g, &types.RoleSetting{
				OwnerId:    player.GetId(),
				FactionId:  role.FactionID,
				BeginRound: role.BeginRound,
				Expiration: role.Expiration,
			}),
			factory.Role(types.VillagerRole, g, &types.RoleSetting{
				OwnerId:    player.GetId(),
				FactionId:  types.VillagerFaction,
				BeginRound: types.FirstRound,
				Expiration: types.UnlimitedTimes,
			}),
		)

		if role.FactionID == types.WerewolfFaction {
			g.factions[types.WerewolfFaction] = append(g.factions[types.WerewolfFaction], player.GetId())

			player.AssignRoles(factory.Role(types.WerewolfRole, g, &types.RoleSetting{
				OwnerId:    player.GetId(),
				FactionId:  types.WerewolfFaction,
				BeginRound: types.FirstRound,
				Expiration: types.UnlimitedTimes,
			}))
		} else {
			g.factions[types.VillagerFaction] = append(g.factions[types.VillagerFaction], player.GetId())
		}

		roles = slices.Delete(roles, index, index+1)
	}
}

func (g *game) KillPlayer(playerId types.PlayerId) contract.Player {
	if player := g.players[playerId]; player == nil {
		return nil
	} else {
		g.deaths = append(g.deaths, playerId)

		return player
	}
}

func (g *game) RequestAction(playerId types.PlayerId, req *types.ActionRequest) *types.ActionResponse {
	if playerId != req.Actor ||
		slices.Contains(g.deaths, playerId) ||
		!g.round.IsAllowed(playerId) {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.UnauthorizedErrorTag,
				Msg: map[string]string{
					types.AlertErrorField: "Not your turn or you're died!",
				},
			},
		}
	}

	return g.GetPlayer(playerId).UseSkill(req)
}
