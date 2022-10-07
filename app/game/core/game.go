package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/role"
	"uwwolf/app/model"
	"uwwolf/app/types"
	"uwwolf/app/validator"
	"uwwolf/util"
)

type game struct {
	id                 types.GameId
	isStarted          bool
	capacity           int
	numberOfWerewolves int
	switchTurnSignal   chan bool
	turnDuration       time.Duration
	discussionDuration time.Duration
	round              contract.Round
	werewolfRoleIds    []types.RoleId
	nonWerewolfRoleIds []types.RoleId
	selectedRoleIds    []types.RoleId
	deadPlayerIds      []types.PlayerId
	fId2pIds           map[types.FactionId][]types.PlayerId
	rId2pIds           map[types.RoleId][]types.PlayerId
	players            map[types.PlayerId]contract.Player
	polls              map[types.FactionId]contract.Poll
}

type roleSplit struct {
	werewolfFaction []*model.Role
	otherFactions   []*model.Role
	reserveWerewolf *model.Role
	reserveVillager *model.Role
}

func NewGame(id types.GameId, setting *types.GameSetting) contract.Game {
	game := game{
		id:                 id,
		capacity:           len(setting.PlayerIds),
		numberOfWerewolves: setting.NumberOfWerewolves,
		switchTurnSignal:   make(chan bool),
		turnDuration:       setting.TurnDuration,
		discussionDuration: setting.DiscussionDuration,
		werewolfRoleIds:    setting.WerewolfRoleIds,
		nonWerewolfRoleIds: setting.NonWerewolfRoleIds,
		round:              NewRound(),
		deadPlayerIds:      make([]types.PlayerId, len(setting.PlayerIds)),
		fId2pIds:           make(map[types.FactionId][]types.PlayerId),
		rId2pIds:           make(map[types.RoleId][]types.PlayerId),
		players:            make(map[types.PlayerId]contract.Player),
		polls:              make(map[types.FactionId]contract.Poll),
	}

	for _, id := range setting.PlayerIds {
		game.players[id] = NewPlayer(&game, id)
	}

	// Create polls for villagers and werewolves
	game.polls[enum.VillagerFactionId] = NewPoll()
	game.polls[enum.WerewolfFactionId] = NewPoll()

	return &game
}

func (g *game) IsStarted() bool {
	return g.isStarted
}

func (g *game) Id() types.GameId {
	return g.id
}

func (g *game) Round() contract.Round {
	return g.round
}

func (g *game) Poll(facitonId types.FactionId) contract.Poll {
	return g.polls[facitonId]
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

func (g *game) Start() (map[types.PlayerId]contract.Player, error) {
	if g.IsStarted() {
		return nil, errors.New("Game is starting!")
	}

	g.selectRoleIds()
	g.assignRoles()
	g.initRound()

	// Create polls for villagers and werewolves
	g.polls[enum.VillagerFactionId].AddElectors(g.fId2pIds[enum.VillagerFactionId])
	g.polls[enum.WerewolfFactionId].AddElectors(g.fId2pIds[enum.WerewolfFactionId])

	time.AfterFunc(10*time.Second, g.listenTurnSwitching)

	g.isStarted = true

	return maps.Clone(g.players), nil
}

func (g *game) selectRoleIds() {
	var selectedWerewolfRoleCounter int
	var selectedNonWerewolfRoleCounter int
	werewolfRoleIds := slices.Clone(g.werewolfRoleIds)
	nonWerewolfRoleIds := slices.Clone(g.nonWerewolfRoleIds)

	for selectedWerewolfRoleCounter < g.numberOfWerewolves {
		if len(werewolfRoleIds) == 0 {
			g.selectedRoleIds = append(g.selectedRoleIds, enum.WerewolfRoleId)
			selectedWerewolfRoleCounter++

			continue
		}

		i, roleId := util.RandomElement(werewolfRoleIds)
		werewolfRoleIds = slices.Delete(werewolfRoleIds, i, i+1)

		for i := 0; i < enum.RoleSets[roleId]; i++ {
			g.selectedRoleIds = append(g.selectedRoleIds, roleId)
			selectedWerewolfRoleCounter++
		}
	}

	for selectedNonWerewolfRoleCounter < len(g.players)-g.numberOfWerewolves {
		if len(nonWerewolfRoleIds) == 0 {
			g.selectedRoleIds = append(g.selectedRoleIds, enum.VillagerRoleId)
			selectedNonWerewolfRoleCounter++

			continue
		}

		i, roleId := util.RandomElement(nonWerewolfRoleIds)
		nonWerewolfRoleIds = slices.Delete(nonWerewolfRoleIds, i, i+1)

		for i := 0; i < enum.RoleSets[roleId]; i++ {
			g.selectedRoleIds = append(g.selectedRoleIds, roleId)
			selectedNonWerewolfRoleCounter++
		}
	}
}

func (g *game) assignRoles() {
	selectedRoleIds := slices.Clone(g.selectedRoleIds)

	for _, player := range g.players {
		i, selectedRoleId := util.RandomElement(selectedRoleIds)
		selectedRoleIds = slices.Delete(selectedRoleIds, i, i+1)
		selectedRole := role.New(selectedRoleId, g, player.Id())

		g.rId2pIds[selectedRoleId] = append(
			g.rId2pIds[selectedRoleId],
			player.Id(),
		)
		g.rId2pIds[enum.VillagerRoleId] = append(
			g.rId2pIds[enum.VillagerRoleId],
			player.Id(),
		)

		// Villager role is always assigned to the player
		player.AssignRoles(
			selectedRole,
			role.New(enum.VillagerRoleId, g, player.Id()),
		)

		// Werewolf role is always assigned to the player belongs
		// to werewolf faction
		if selectedRole.FactionId() == enum.WerewolfFactionId {
			g.rId2pIds[enum.WerewolfRoleId] = append(
				g.rId2pIds[enum.WerewolfRoleId],
				player.Id(),
			)

			player.AssignRoles(role.New(enum.WerewolfRoleId, g, player.Id()))
		}

		g.fId2pIds[selectedRole.FactionId()] = append(
			g.fId2pIds[selectedRole.FactionId()],
			player.Id(),
		)
	}
}

func (g *game) initRound() {
	g.round.AddTurn(&types.TurnSetting{
		PhaseId:    enum.DayPhaseId,
		RoleId:     enum.VillagerRoleId,
		BeginRound: enum.FirstRound,
		Priority:   enum.VillagerPriority,
		Expiration: enum.UnlimitedTimes,
		Position:   enum.SortedPosition,
	})
	g.round.AddTurn(&types.TurnSetting{
		PhaseId:    enum.NightPhaseId,
		RoleId:     enum.WerewolfRoleId,
		BeginRound: enum.FirstRound,
		Priority:   enum.WerewolfPriority,
		Expiration: enum.UnlimitedTimes,
		Position:   enum.SortedPosition,
	})

	for _, player := range g.players {
		for _, assignedRole := range player.Roles() {
			g.round.AddTurn(&types.TurnSetting{
				PhaseId:    assignedRole.PhaseId(),
				RoleId:     assignedRole.Id(),
				BeginRound: assignedRole.BeginRound(),
				Priority:   assignedRole.Priority(),
				Expiration: assignedRole.Expiration(),
				Position:   enum.SortedPosition,
			})
			g.round.AddPlayer(player.Id(), assignedRole.Id())
		}
	}
}

func (g *game) listenTurnSwitching() {
	for {
		func() {
			var duration time.Duration

			if g.Round().CurrentTurn().RoleId == enum.VillagerRoleId {
				duration = g.discussionDuration

				g.Poll(enum.VillagerFactionId).Open()
				defer g.Poll(enum.VillagerFactionId).Close()
			} else {
				duration = g.turnDuration

				if g.Round().CurrentTurn().RoleId == enum.WerewolfRoleId {
					g.Poll(enum.WerewolfFactionId).Open()
					defer g.Poll(enum.WerewolfFactionId).Close()
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), duration*time.Second)
			defer cancel()

			select {
			case isSkipped := <-g.switchTurnSignal:
				// fmt.Println("next turn!")
				g.Round().NextTurn(isSkipped)

			case <-ctx.Done():
				// fmt.Println("timeout!")
				g.Round().NextTurn(true)
			}
		}()
	}
}

func (g *game) KillPlayer(playerId types.PlayerId) contract.Player {
	if player := g.players[playerId]; player == nil {
		return nil
	} else {
		g.Round().DeletePlayerFromAllTurns(player.Id())
		g.polls[enum.VillagerFactionId].RemoveElector(player.Id())
		g.polls[enum.WerewolfFactionId].RemoveElector(player.Id())
		g.deadPlayerIds = append(g.deadPlayerIds, playerId)

		return player
	}
}

func (g *game) RequestAction(req *types.ActionRequest) *types.ActionResponse {
	if errs := validator.ValidateStruct(req); errs != nil {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: enum.InvalidInputErrorTag,
				Msg: errs,
			},
		}
	}

	if slices.Contains(g.deadPlayerIds, req.ActorId) ||
		!g.round.IsAllowed(req.ActorId) {

		fmt.Println(g.round.CurrentTurn().PlayerIds)

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   enum.ForbiddenErrorTag,
				Alert: "Not your turn or you're died!",
			},
		}
	}

	res := g.Player(req.ActorId).UseSkill(req)

	if res.Ok {
		g.switchTurnSignal <- req.IsSkipped
	}

	return res
}
