package logic

import (
	"time"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type world struct {
	isLoaded bool

	// numberWerewolves is required number of players belongs to
	// the werewolf faction.
	numberWerewolves uint8

	// scheduler is turn manager.
	moderator contract.Moderator

	// roleIDs is the possible role IDs in the game.
	roleIDs []types.RoleId

	// requiredRoleIDs is the required role IDs in the game.
	requiredRoleIDs []types.RoleId

	// selectRoleIDs is the selected role IDs from `roleIDs` and `requiredRoleIDs`.
	selectedRoleIDs []types.RoleId

	// players contains all players playing the game.
	players map[types.PlayerId]contract.Player

	// polls contains the polls of villager and werewolf factions.
	polls map[types.FactionId]contract.Poll

	gameMap contract.Map
}

func NewWorld(moderator contract.Moderator, init *types.GameInitialization) contract.World {
	world := world{
		numberWerewolves: init.NumberWerewolves,
		roleIDs:          init.RoleIds,
		requiredRoleIDs:  init.RequiredRoleIds,
		moderator:        moderator,
		players:          make(map[types.PlayerId]contract.Player),
		polls:            make(map[types.FactionId]contract.Poll),
		gameMap:          NewMap(),
	}

	for i, id := range init.PlayerIDs {
		world.players[id] = NewPlayer(moderator, id)
		world.gameMap.AddEntity(string(id), contract.EntitySettings{
			Type:    contract.PlayerEntity,
			X:       float64(64*i) + 200,
			Y:       50,
			Width:   64,
			Height:  64,
			IsSolid: true,
			Speed:   1,
		})
	}

	// Create polls for villagers and werewolves
	world.polls[constants.VillagerFactionId] = NewPoll()
	world.polls[constants.WerewolfFactionId] = NewPoll()

	return &world
}

func (w world) Map() contract.Map {
	return w.gameMap
}

// Scheduler returns turn manager.
func (w world) Scheduler() contract.Scheduler {
	return w.moderator.Scheduler()
}

// Poll returns the in-game poll management state.
// Each specific faction has different poll to interact with.
func (w world) Poll(facitonId types.FactionId) contract.Poll {
	return w.polls[facitonId]
}

// Player returns the player by given player ID.
func (w world) Player(playerId types.PlayerId) contract.Player {
	return w.players[playerId]
}

// Players returns the player list.
func (w world) Players() map[types.PlayerId]contract.Player {
	return w.players
}

// AlivePlayerIDsWithRoleID returns the alive player ID list having the
// givent role ID.
func (w world) AlivePlayerIdsWithRoleId(roleId types.RoleId) []types.PlayerId {
	playerIDs := make([]types.PlayerId, 0)
	for playerID, player := range w.players {
		if !player.IsDead() && slices.Contains(player.RoleIds(), roleId) {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// AlivePlayerIDsWithFactionID returns the alive player ID list having the
// given faction ID.
func (w world) AlivePlayerIdsWithFactionId(factionId types.FactionId) []types.PlayerId {
	playerIDs := make([]types.PlayerId, 0)
	for playerID, player := range w.players {
		if !player.IsDead() && player.FactionId() == factionId {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// AlivePlayerIDsWithoutFactionID returns the alive player ID list not having
// the given faction ID.
func (w world) AlivePlayerIdsWithoutFactionId(factionId types.FactionId) []types.PlayerId {
	playerIDs := make([]types.PlayerId, 0)
	for playerID, player := range w.players {
		if !player.IsDead() && player.FactionId() != factionId {
			playerIDs = append(playerIDs, playerID)
		}
	}

	return playerIDs
}

// selectRoleID selects set of the given role ID. Return false if the full set
// can't be selected.
//
// Note: Don't work with unlimited set.
func (w *world) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID types.RoleId) bool {
	isWerewolf := slices.Contains(
		constants.FactionId2RoleIds.BlindGet(constants.WerewolfFactionId),
		roleID,
	)

	for i := 0; i < int(constants.RoleSets.BlindGet(roleID)); i++ {
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
	roleIDs := lo.Filter(w.roleIDs, func(roleID types.RoleId, _ int) bool {
		return !slices.Contains(w.requiredRoleIDs, roleID)
	})
	for {
		i, randomRoleID := util.RandomElement(roleIDs)
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
		w.selectedRoleIDs = append(w.selectedRoleIDs, constants.WerewolfRoleId)
		werewolfCounter++
	}
}

// assignRoles assigns the selected roles to the players randomly.
func (w *world) assignRoles() {
	selectedRoleIDs := slices.Clone(w.selectedRoleIDs)

	for _, player := range w.players {
		i, selectedRoleID := util.RandomElement(selectedRoleIDs)
		// Remove the assigned role
		if i != -1 {
			selectedRoleIDs = slices.Delete(selectedRoleIDs, i, i+1)
		}

		// Assign default role
		player.AssignRole(constants.VillagerRoleId) // nolint: errcheck

		selectedRole, _ := w.moderator.RoleFactory().CreateById(selectedRoleID, w.moderator, player.Id())
		if selectedRole == nil {
			continue
		}

		// Assign default werewolf faction's role
		if selectedRole.FactionId() == constants.WerewolfFactionId &&
			selectedRole.Id() != constants.WerewolfRoleId {
			player.AssignRole(constants.WerewolfRoleId) // nolint: errcheck
		}

		// Assign main role
		player.AssignRole(selectedRole.Id()) // nolint: errcheck
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
