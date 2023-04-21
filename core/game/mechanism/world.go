package mechanism

import (
	"time"
	"uwwolf/game/declare"
	"uwwolf/game/mechanism/contract"
	"uwwolf/game/mechanism/role"
	"uwwolf/game/tool"
	"uwwolf/game/types"
	"uwwolf/util/helper"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type world struct {
	isLoaded bool

	// numberWerewolves is required number of players belongs to
	// the werewolf faction.
	numberWerewolves uint8

	// scheduler is turn manager.
	scheduler tool.Scheduler

	// roleIDs is the possible role IDs in the game.
	roleIDs []types.RoleID

	// requiredRoleIDs is the required role IDs in the game.
	requiredRoleIDs []types.RoleID

	// selectRoleIDs is the selected role IDs from `roleIDs` and `requiredRoleIDs`.
	selectedRoleIDs []types.RoleID

	// players contains all players playing the game.
	players map[types.PlayerID]contract.Player

	// polls contains the polls of villager and werewolf factions.
	polls map[types.FactionID]tool.Poll

	gameMap tool.Map
}

func NewWorld(scheduler tool.Scheduler, init *types.GameInitialization) contract.World {
	world := world{
		numberWerewolves: init.NumberWerewolves,
		roleIDs:          init.RoleIDs,
		requiredRoleIDs:  init.RequiredRoleIDs,
		scheduler:        scheduler,
		players:          make(map[types.PlayerID]contract.Player),
		polls:            make(map[types.FactionID]tool.Poll),
	}

	for _, id := range init.PlayerIDs {
		world.players[id] = NewPlayer(&world, id)
	}

	// Create polls for villagers and werewolves
	world.polls[declare.VillagerFactionID] = tool.NewPoll()
	world.polls[declare.WerewolfFactionID] = tool.NewPoll()

	return &world
}

func (w world) Map() tool.Map {
	return w.gameMap
}

// Scheduler returns turn manager.
func (w world) Scheduler() tool.Scheduler {
	return w.scheduler
}

// Poll returns the in-game poll management state.
// Each specific faction has different poll to interact with.
func (w world) Poll(facitonID types.FactionID) tool.Poll {
	return w.polls[facitonID]
}

// Player returns the player by given player ID.
func (w world) Player(playerId types.PlayerID) contract.Player {
	return w.players[playerId]
}

// Players returns the player list.
func (w world) Players() map[types.PlayerID]contract.Player {
	return w.players
}

// AlivePlayerIDsWithRoleID returns the alive player ID list having the
// givent role ID.
func (w world) AlivePlayerIDsWithRoleID(roleID types.RoleID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range w.players {
		if !player.IsDead() && slices.Contains(player.RoleIDs(), roleID) {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// AlivePlayerIDsWithFactionID returns the alive player ID list having the
// given faction ID.
func (w world) AlivePlayerIDsWithFactionID(factionID types.FactionID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range w.players {
		if !player.IsDead() && player.FactionID() == factionID {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// AlivePlayerIDsWithoutFactionID returns the alive player ID list not having
// the given faction ID.
func (w world) AlivePlayerIDsWithoutFactionID(factionID types.FactionID) []types.PlayerID {
	playerIDs := make([]types.PlayerID, 0)
	for playerID, player := range w.players {
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
func (w *world) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID types.RoleID) bool {
	isWerewolf := slices.Contains(
		declare.FactionID2RoleIDs.BindGet(declare.WerewolfFactionID),
		roleID,
	)

	for i := 0; i < int(declare.RoleSets.BindGet(roleID)); i++ {
		isMissingWerewolf := *werewolfCounter < int(w.numberWerewolves)
		isMissingNonWerewolf := *nonWerewolfCounter < len(w.players)-int(w.numberWerewolves)

		// Break if number of selectedRoleIDs is enough
		if !isMissingWerewolf && !isMissingNonWerewolf {
			return false
		}

		if isMissingWerewolf && isWerewolf {
			// Select werewolf role
			w.selectedRoleIDs = append(w.selectedRoleIDs, roleID)
			*werewolfCounter++
		} else if isMissingNonWerewolf && !isWerewolf {
			// Select non-werewolf role
			w.selectedRoleIDs = append(w.selectedRoleIDs, roleID)
			*nonWerewolfCounter++
		}
	}

	return true
}

// selectRoleIDs selects the required role IDs. If the selectedRoleIDs isn't enough,
// continue to select role IDs in roleIDs.
func (w *world) selectRoleIDs() {
	werewolfCounter := 0
	nonWerewolfCounter := 0

	// Select required roles
	for _, requiredRoleID := range w.requiredRoleIDs {
		// Stop if selectedRoleIDs is enough
		if !w.selectRoleID(&werewolfCounter, &nonWerewolfCounter, requiredRoleID) {
			break
		}
	}

	// Select random roles
	roleIDs := lo.Filter(w.roleIDs, func(roleID types.RoleID, _ int) bool {
		return !slices.Contains(w.requiredRoleIDs, roleID)
	})
	for {
		i, randomRoleID := helper.RandomElement(roleIDs)
		if i == -1 ||
			!w.selectRoleID(&werewolfCounter, &nonWerewolfCounter, randomRoleID) {
			// Stop if selectedRoleIDs is enough or roleIDs is fully checked
			break
		} else {
			// Remove selected roleID
			roleIDs = slices.Delete(roleIDs, i, i+1)
		}
	}

	// Add missing werewolf roles
	for werewolfCounter < int(w.numberWerewolves) {
		w.selectedRoleIDs = append(w.selectedRoleIDs, declare.WerewolfRoleID)
		werewolfCounter++
	}
}

// assignRoles assigns the selected roles to the players randomly.
func (w *world) assignRoles() {
	selectedRoleIDs := slices.Clone(w.selectedRoleIDs)

	for _, player := range w.players {
		i, selectedRoleID := helper.RandomElement(selectedRoleIDs)
		// Remove the assigned role
		if i != -1 {
			selectedRoleIDs = slices.Delete(selectedRoleIDs, i, i+1)
		}

		// Assign default role
		player.AssignRole(declare.VillagerRoleID) // nolint: errcheck

		selectedRole, _ := role.NewRole(selectedRoleID, w, player.ID())
		if selectedRole == nil {
			continue
		}

		// Assign default werewolf faction's role
		if selectedRole.FactionID() == declare.WerewolfFactionID &&
			selectedRole.ID() != declare.WerewolfRoleID {
			player.AssignRole(declare.WerewolfRoleID) // nolint: errcheck
		}

		// Assign main role
		player.AssignRole(selectedRole.ID()) // nolint: errcheck
	}
}

// Prepare sets up the game and returns completion time in milisecond.
func (w *world) Load() int64 {
	if w.isLoaded {
		return -1
	}

	w.selectRoleIDs()
	w.assignRoles()
	w.isLoaded = true
	return time.Now().Unix()
}
