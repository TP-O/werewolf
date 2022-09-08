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
	"uwwolf/validator"
)

type game struct {
	id                 types.GameId
	isStarted          bool
	capacity           int
	numberOfWerewolves int
	turnDuration       time.Duration
	discussionDuration time.Duration
	round              *state.Round
	rolePool           []types.RoleId
	deadPlayerIds      []types.PlayerId
	fId2pIds           map[types.FactionId][]types.PlayerId
	rId2pIds           map[types.RoleId][]types.PlayerId
	players            map[types.PlayerId]contract.Player
	polls              map[types.FactionId]*state.Poll
}

type roleSplit struct {
	werewolfFaction []*model.Role
	otherFactions   []*model.Role
	reserveWerewolf *model.Role
	reserveVillager *model.Role
}

func NewGame(setting *types.GameSetting) contract.Game {
	game := game{
		id:                 setting.Id,
		capacity:           len(setting.PlayerIds),
		numberOfWerewolves: setting.NumberOfWerewolves,
		turnDuration:       setting.TurnDuration,
		discussionDuration: setting.DiscussionDuration,
		rolePool:           setting.RolePool,
		round:              state.NewRound(),
		deadPlayerIds:      make([]types.PlayerId, len(setting.PlayerIds)),
		fId2pIds:           make(map[types.FactionId][]types.PlayerId),
		rId2pIds:           make(map[types.RoleId][]types.PlayerId),
		players:            make(map[types.PlayerId]contract.Player),
		polls:              make(map[types.FactionId]*state.Poll),
	}

	for _, id := range setting.PlayerIds {
		game.players[id] = NewPlayer(&game, id)
	}

	// Create polls for villagers and werewolves
	game.polls[types.VillagerFaction] = state.NewPoll()
	game.polls[types.WerewolfFaction] = state.NewPoll()

	return &game
}

func (g *game) IsStarted() bool {
	return g.isStarted
}

func (g *game) Player(playerId types.PlayerId) contract.Player {
	return g.players[playerId]
}

func (g *game) PlayerIdsWithRole(roleId types.RoleId) []types.PlayerId {
	return g.rId2pIds[roleId]
}

func (g *game) PlayerIdsWithFaction(factionId types.FactionId) []types.PlayerId {
	return g.fId2pIds[factionId]
}

func (g *game) Round() *state.Round {
	return g.round
}

func (g *game) Poll(factionId types.FactionId) *state.Poll {
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

	// Create polls for villagers and werewolves
	g.polls[types.VillagerFaction].AddElectors(g.fId2pIds[types.VillagerFaction])
	g.polls[types.VillagerFaction].AddElectors(g.fId2pIds[types.WerewolfFaction])
	g.polls[types.WerewolfFaction].AddElectors(g.fId2pIds[types.WerewolfFaction])

	g.isStarted = true

	return true
}

func (g *game) randomRoles(roleSplit *roleSplit) []*model.Role {
	randomRoles := append(
		g.pickUpRoles(
			g.numberOfWerewolves,
			roleSplit.werewolfFaction,
			roleSplit.reserveWerewolf,
		),
		g.pickUpRoles(
			g.capacity-g.numberOfWerewolves,
			roleSplit.otherFactions,
			roleSplit.reserveVillager,
		)...,
	)

	return randomRoles
}

// Return list of roles which can be duplicate because
// one role can be assigned to many players.
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
				roleSplit.werewolfFaction = append(roleSplit.werewolfFaction, role)
			} else {
				roleSplit.otherFactions = append(roleSplit.otherFactions, role)
			}
		}

		// Break if enough roles
		if len(roleSplit.werewolfFaction)+len(roleSplit.otherFactions) == len(ids) {
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

	for _, role := range roleSplit.werewolfFaction {
		g.round.AddTurn(&types.TurnSetting{
			PhaseId:    role.PhaseID,
			RoleId:     role.ID,
			BeginRound: role.BeginRound,
			Priority:   role.Priority,
			Expiration: role.Expiration,
			Position:   types.SortedPosition,
		})
	}
	for _, role := range roleSplit.otherFactions {
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

		g.rId2pIds[types.VillagerRole] = append(g.rId2pIds[types.VillagerRole], player.Id())
		g.rId2pIds[role.ID] = append(g.rId2pIds[role.ID], player.Id())

		player.AssignRoles(
			factory.Role(role.ID, g, &types.RoleSetting{
				Id:         role.ID,
				OwnerId:    player.Id(),
				PhaseId:    role.PhaseID,
				FactionId:  role.FactionID,
				BeginRound: role.BeginRound,
				Expiration: role.Expiration,
			}),
			factory.Role(types.VillagerRole, g, &types.RoleSetting{
				Id:         types.VillagerRole,
				OwnerId:    player.Id(),
				PhaseId:    types.DayPhase,
				FactionId:  types.VillagerFaction,
				BeginRound: types.FirstRound,
				Expiration: types.UnlimitedTimes,
			}),
		)

		g.round.AddPlayer(player.Id(), role.ID)
		g.round.AddPlayer(player.Id(), types.VillagerRole)

		if role.FactionID == types.WerewolfFaction {
			g.rId2pIds[types.WerewolfRole] = append(g.rId2pIds[types.WerewolfRole], player.Id())
			g.fId2pIds[types.WerewolfFaction] = append(g.fId2pIds[types.WerewolfFaction], player.Id())

			g.round.AddPlayer(player.Id(), types.WerewolfRole)

			player.AssignRoles(factory.Role(types.WerewolfRole, g, &types.RoleSetting{
				Id:         types.WerewolfRole,
				OwnerId:    player.Id(),
				PhaseId:    types.NightPhase,
				FactionId:  types.WerewolfFaction,
				BeginRound: types.FirstRound,
				Expiration: types.UnlimitedTimes,
			}))
		} else {
			g.fId2pIds[types.VillagerFaction] = append(g.fId2pIds[types.VillagerFaction], player.Id())
		}

		roles = slices.Delete(roles, index, index+1)
	}
}

func (g *game) KillPlayer(playerId types.PlayerId) contract.Player {
	if player := g.players[playerId]; player == nil {
		return nil
	} else {
		g.Round().DeletePlayerFromAllTurns(player.Id())
		g.polls[types.VillagerFaction].RemoveElector(player.Id())
		g.polls[types.WerewolfFaction].RemoveElector(player.Id())
		g.deadPlayerIds = append(g.deadPlayerIds, playerId)

		return player
	}
}

func (g *game) RequestAction(req *types.ActionRequest) *types.ActionResponse {
	if errs := validator.ValidateStruct(req); errs != nil {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.InvalidInputErrorTag,
				Msg: errs,
			},
		}
	}

	if slices.Contains(g.deadPlayerIds, req.ActorId) ||
		!g.round.IsAllowed(req.ActorId) {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.UnauthorizedErrorTag,
				Alert: "Not your turn or you're died!",
			},
		}
	}

	return g.Player(req.ActorId).UseSkill(req)
}
