package game

import (
	"time"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	"uwwolf/util"

	"golang.org/x/exp/slices"
)

type game struct {
	numberWerewolves uint8
	statusID         types.GameStatusID
	scheduler        contract.Scheduler
	roleIDs          []types.RoleID
	requiredRoleIDs  []types.RoleID
	selectedRoleIDs  []types.RoleID
	players          map[types.PlayerID]contract.Player
	polls            map[types.FactionID]contract.Poll
}

func NewGame(scheduler contract.Scheduler, setting *types.GameSetting) contract.Game {
	game := game{
		numberWerewolves: setting.NumberWerewolves,
		statusID:         vars.Idle,
		roleIDs:          setting.RoleIDs,
		requiredRoleIDs:  setting.RequiredRoleIDs,
		scheduler:        scheduler,
		players:          make(map[types.PlayerID]contract.Player),
		polls:            make(map[types.FactionID]contract.Poll),
	}

	for _, id := range setting.PlayerIDs {
		game.players[id] = NewPlayer(&game, id)
	}

	// Create polls for villagers and werewolves
	game.polls[vars.VillagerFactionID] = NewPoll()
	game.polls[vars.WerewolfFactionID] = NewPoll()

	return &game
}

// StatusID retusn current game status ID.
func (g game) StatusID() types.GameStatusID {
	return g.statusID
}

// Scheduler returns turn manager.
func (g game) Scheduler() contract.Scheduler {
	return g.scheduler
}

// Poll returns the in-game poll management state.
// Each specific faction has different poll to interact with.
func (g game) Poll(facitonID types.FactionID) contract.Poll {
	return g.polls[facitonID]
}

// Player returns the player by given player ID.
func (g game) Player(playerId types.PlayerID) contract.Player {
	return g.players[playerId]
}

// Players returns the player list.
func (g game) Players() map[types.PlayerID]contract.Player {
	return g.players
}

// AlivePlayerIDsWithRoleID returns the alive player ID list having the
// givent role ID.
func (g game) AlivePlayerIDsWithRoleID(roleID types.RoleID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range g.players {
		if !player.IsDead() && slices.Contains(player.RoleIDs(), roleID) {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// AlivePlayerIDsWithFactionID returns the alive player ID list having the
// given faction ID.
func (g game) AlivePlayerIDsWithFactionID(factionID types.FactionID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range g.players {
		if !player.IsDead() && player.FactionID() == factionID {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// AlivePlayerIDsWithoutFactionID returns the alive player ID list not having
// the given faction ID.
func (g game) AlivePlayerIDsWithoutFactionID(factionID types.FactionID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range g.players {
		if !player.IsDead() && player.FactionID() != factionID {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// selectRoleID selects set of the given role ID. Return false if the full set
// can't be selected.
//
// Note: Don't work with unlimited set.
func (g *game) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID types.RoleID) bool {
	isWerewolf := slices.Contains(
		vars.FactionID2RoleIDs[vars.WerewolfFactionID],
		roleID,
	)

	for i := 0; i < int(vars.RoleSets[roleID]); i++ {
		isMissingWerewolf := *werewolfCounter < int(g.numberWerewolves)
		isMissingNonWerewolf := *nonWerewolfCounter < len(g.players)-int(g.numberWerewolves)

		// Break if number of selectedRoleIDs is enough
		if !isMissingWerewolf && !isMissingNonWerewolf {
			return false
		}

		if isMissingWerewolf && isWerewolf {
			// Select werewolf role
			g.selectedRoleIDs = append(g.selectedRoleIDs, roleID)
			*werewolfCounter++
		} else if isMissingNonWerewolf && !isWerewolf {
			// Select non-werewolf role
			g.selectedRoleIDs = append(g.selectedRoleIDs, roleID)
			*nonWerewolfCounter++
		}
	}

	return true
}

// selectRoleIDs selects the required role IDs. If the selectedRoleIDs isn't enough,
// continue to select role IDs in roleIDs.
func (g *game) selectRoleIDs() {
	werewolfCounter := 0
	nonWerewolfCounter := 0

	// Select required roles
	for _, requiredRoleID := range g.requiredRoleIDs {
		// Stop if selectedRoleIDs is enough
		if !g.selectRoleID(&werewolfCounter, &nonWerewolfCounter, requiredRoleID) {
			return
		}
	}

	// Select random roles
	roleIDs := util.FilterSlice(g.roleIDs, func(roleID types.RoleID) bool {
		return !slices.Contains(g.requiredRoleIDs, roleID)
	})
	for {
		i, randomRoleID := util.RandomElement(roleIDs)
		if i == -1 ||
			!g.selectRoleID(&werewolfCounter, &nonWerewolfCounter, randomRoleID) {
			// Stop if selectedRoleIDs is enough or roleIDs is fully checked
			return
		} else {
			// Remove selected roleID
			roleIDs = slices.Delete(roleIDs, i, i+1)
		}
	}
}

// assignRoles assigns the selected roles to the players randomly.
func (g *game) assignRoles() {
	selectedRoleIDs := slices.Clone(g.selectedRoleIDs)

	for _, player := range g.players {
		i, selectedRoleID := util.RandomElement(selectedRoleIDs)
		// Remove the assigned role
		if i != -1 {
			selectedRoleIDs = slices.Delete(selectedRoleIDs, i, i+1)
		}

		// Assign default role
		player.AssignRole(vars.VillagerRoleID)

		selectedRole, _ := NewRole(selectedRoleID, g, player.ID())
		if selectedRole == nil {
			continue
		}

		// Assign default werewolf faction's role
		if selectedRole.FactionID() == vars.WerewolfFactionID {
			player.AssignRole(vars.WerewolfRoleID)
		}

		// Assign main role
		player.AssignRole(selectedRole.ID())
	}
}

// Prepare sets up the game and returns completion time in milisecond.
func (g *game) Prepare() int64 {
	if g.statusID != vars.Idle {
		return -1
	}

	g.selectRoleIDs()
	g.assignRoles()
	g.statusID = vars.Waiting

	return time.Now().Unix()
}

// Start starts the game.
func (g *game) Start() bool {
	if g.statusID != vars.Waiting {
		return false
	}

	g.statusID = vars.Starting
	return true
}

// Finish fishes the game.
func (g *game) Finish() bool {
	if g.statusID == vars.Finished {
		return false
	}

	g.statusID = vars.Finished
	return true
}

// Play activates the player's ability.
func (g *game) Play(playerID types.PlayerID, req *types.ActivateAbilityRequest) *types.ActionResponse {
	var msg string
	if g.statusID != vars.Starting {
		msg = "The game is about to start ノ(ジ)ー'"
	} else if g.Player(playerID) == nil {
		msg = "Unable to play this game (╥﹏╥)"
	}
	if msg != "" {
		return &types.ActionResponse{
			Ok:      false,
			Message: msg,
		}
	}

	return g.players[playerID].ActivateAbility(req)
}

// KillPlayer kills the player by the given player ID.
// If `isExited` is true, any trigger preventing death is ignored.
func (g *game) KillPlayer(playerID types.PlayerID, isExited bool) contract.Player {
	if player := g.players[playerID]; player == nil ||
		player.IsDead() ||
		!player.Die(isExited) {
		return nil
	} else {
		return player
	}
}