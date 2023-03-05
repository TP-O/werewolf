package core

import (
	"time"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	"uwwolf/util"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type game struct {
	id               types.GameID
	numberWerewolves uint8
	statusID         types.GameStatusID
	scheduler        contract.Scheduler
	roleIDs          []types.RoleID
	requiredRoleIDs  []types.RoleID
	selectedRoleIDs  []types.RoleID
	exitedPlayerIDs  []types.PlayerID
	players          map[types.PlayerID]contract.Player
	polls            map[types.FactionID]contract.Poll
}

var _ contract.Game = (*game)(nil)

func NewGame(scheduler contract.Scheduler, setting types.GameSetting) contract.Game {
	game := game{
		id:               setting.GameID,
		numberWerewolves: setting.NumberWerewolves,
		statusID:         vars.Idle,
		roleIDs:          setting.RoleIDs,
		requiredRoleIDs:  setting.RequiredRoleIDs,
		scheduler:        scheduler,
		players:          make(map[types.PlayerID]contract.Player),
		polls:            make(map[types.FactionID]contract.Poll),
	}

	for _, id := range setting.PlayerIDs {
		playerID := types.PlayerID(id)
		game.players[playerID] = NewPlayer(&game, playerID)
	}

	// Create polls for villagers and werewolves
	game.polls[vars.VillagerFactionID] = NewPoll()
	game.polls[vars.WerewolfFactionID] = NewPoll()

	return &game
}

func (g game) ID() types.GameID {
	return g.id
}

func (g game) StatusID() types.GameStatusID {
	return g.statusID
}

func (g game) Scheduler() contract.Scheduler {
	return g.scheduler
}

func (g game) Poll(facitonID types.FactionID) contract.Poll {
	return g.polls[facitonID]
}

func (g game) Player(playerId types.PlayerID) contract.Player {
	return g.players[playerId]
}

func (g game) Players() map[types.PlayerID]contract.Player {
	return g.players
}

func (g game) PlayerIDsWithRoleID(roleID types.RoleID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range g.players {
		if slices.Contains(player.RoleIDs(), roleID) {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

func (g game) PlayerIDsWithFactionID(
	factionID types.FactionID,
	onlyAlive bool,
) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range g.players {
		if player.FactionID() == factionID &&
			(!onlyAlive || (onlyAlive && !player.IsDead())) {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

func (g game) PlayerIDsWithoutFactionID(
	factionID types.FactionID,
	onlyAlive bool,
) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range g.players {
		if player.FactionID() != factionID &&
			(!onlyAlive || (onlyAlive && !player.IsDead())) {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

func (g *game) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID types.RoleID) bool {
	isWerewolf := slices.Contains(
		vars.FactionID2RoleIDs[vars.WerewolfFactionID],
		roleID,
	)

	for i := 0; i < int(vars.RoleSets[roleID]); i++ {
		isMissingWerewolf := *werewolfCounter < int(g.numberWerewolves)
		isMissingNonWerewolf := *nonWerewolfCounter < len(g.players)-int(g.numberWerewolves)

		if !isMissingWerewolf && !isMissingNonWerewolf {
			return false
		}

		if isMissingWerewolf && isWerewolf {
			g.selectedRoleIDs = append(g.selectedRoleIDs, roleID)
			*werewolfCounter++
		} else if isMissingNonWerewolf && !isWerewolf {
			g.selectedRoleIDs = append(g.selectedRoleIDs, roleID)
			*nonWerewolfCounter++
		}
	}

	return true
}

func (g *game) assignRoles() {
	selectedRoleIDs := slices.Clone(g.selectedRoleIDs)

	for _, player := range g.players {
		i, selectedRoleID := util.RandomElement(selectedRoleIDs)
		selectedRole, _ := NewRole(selectedRoleID, g, player.ID())

		// Remove the random role from array
		if i != -1 {
			selectedRoleIDs = slices.Delete(selectedRoleIDs, i, i+1)
		}

		// Default role
		if ok, _ := player.AssignRole(vars.VillagerRoleID); ok {
			g.scheduler.AddPlayerTurn(types.NewPlayerTurn{
				PhaseID:      vars.DayPhaseID,
				RoleID:       vars.VillagerRoleID,
				BeginRoundID: types.RoundID(0),
				TurnID:       vars.MidTurn,
				PlayerID:     player.ID(),
			})
		}

		if selectedRole == nil {
			continue
		}

		// Default werewolf faction's role
		if selectedRole.FactionID() == vars.WerewolfFactionID {
			if ok, _ := player.AssignRole(vars.WerewolfRoleID); ok {
				g.scheduler.AddPlayerTurn(types.NewPlayerTurn{
					PhaseID:      vars.NightPhaseID,
					RoleID:       vars.WerewolfRoleID,
					BeginRoundID: types.RoundID(0),
					TurnID:       vars.MidTurn,
					PlayerID:     player.ID(),
				})
			}

		}

		// Main role
		if ok, _ := player.AssignRole(selectedRole.ID()); ok {
			// Add the main role's turn to the schedule
			g.scheduler.AddPlayerTurn(types.NewPlayerTurn{
				PhaseID:      selectedRole.PhaseID(),
				RoleID:       selectedRole.ID(),
				BeginRoundID: selectedRole.BeginRoundID(),
				TurnID:       selectedRole.TurnID(),
				PlayerID:     player.ID(),
			})
		}
	}
}

func (g *game) randomRoleIDs() {
	werewolfCounter := len(g.PlayerIDsWithFactionID(vars.WerewolfFactionID, true))
	nonWerewolfCounter := len(g.PlayerIDsWithoutFactionID(vars.WerewolfFactionID, true))

	// Select required roles
	for _, requiredRoleID := range g.requiredRoleIDs {
		if !g.selectRoleID(&werewolfCounter, &nonWerewolfCounter, requiredRoleID) {
			break
		}
	}

	roleIDs := slices.Clone(g.roleIDs)

	// Select random roles
	for {
		i, randomRoleID := util.RandomElement(roleIDs)

		if i == -1 ||
			!g.selectRoleID(&werewolfCounter, &nonWerewolfCounter, randomRoleID) {
			break
		}

		roleIDs = slices.Delete(roleIDs, i, i+1)
	}

	g.assignRoles()
}

func (g *game) addCandidatesToPolls() {
	g.polls[vars.VillagerFactionID].AddCandidates(maps.Keys(g.players)...)
	g.polls[vars.WerewolfFactionID].AddCandidates(
		g.PlayerIDsWithoutFactionID(vars.WerewolfFactionID, true)...,
	)
}

func (g *game) Prepare() int64 {
	if g.statusID != vars.Idle {
		return -1
	}

	g.randomRoleIDs()
	g.addCandidatesToPolls()
	g.statusID = vars.Waiting

	return time.Now().Unix()
}

func (g *game) Start() bool {
	if g.statusID != vars.Waiting {
		return false
	}

	g.statusID = vars.Starting
	return true
}

func (g *game) Finish() bool {
	if g.statusID == vars.Finished {
		return false
	}

	g.statusID = vars.Finished
	return true
}

func (g *game) Play(playerID types.PlayerID, req types.ActivateAbilityRequest) types.ActionResponse {
	// Validate
	var msg string
	if g.statusID != vars.Starting {
		msg = "Wait until the game starts (╥﹏╥)"
	} else if g.Player(playerID) == nil {
		msg = "Unable to play this game (╥﹏╥)"
	} else if g.players[playerID].IsDead() {
		msg = "You're died (╥﹏╥)"
	} else if !slices.Contains(maps.Keys(g.scheduler.Turn()), playerID) {
		msg = "Not your turn (╥﹏╥)"
	}
	if msg != "" {
		return types.ActionResponse{
			Ok:      false,
			Message: msg,
		}
	}

	return g.players[playerID].ActivateAbility(req)
}

func (g *game) KillPlayer(playerID types.PlayerID, isExited bool) contract.Player {
	player := g.players[playerID]
	if player == nil ||
		g.players[playerID].IsDead() ||
		!player.Die(isExited) {
		return nil
	}

	g.polls[vars.VillagerFactionID].RemoveElector(player.ID())
	g.polls[vars.VillagerFactionID].RemoveCandidate(player.ID())
	if player.FactionID() == vars.WerewolfFactionID {
		g.polls[vars.WerewolfFactionID].RemoveElector(player.ID())
	} else {
		g.polls[vars.WerewolfFactionID].RemoveCandidate(player.ID())
	}

	return player
}
