package core

import (
	"context"
	"fmt"
	"sync"
	"time"
	"uwwolf/game/contract"
	"uwwolf/game/role"
	"uwwolf/game/types"
	"uwwolf/util"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type game struct {
	id                 types.GameID
	numberWerewolves   uint8
	nextTurnSignal     chan bool
	finishSignal       chan bool
	mutex              *sync.Mutex
	status             types.GameStatus
	turnDuration       time.Duration
	discussionDuration time.Duration
	scheduler          contract.Scheduler
	roleIDs            []types.RoleID
	requiredRoleIDs    []types.RoleID
	selectedRoleIDs    []types.RoleID
	deadPlayerIDs      []types.PlayerID
	exitedPlayerIDs    []types.PlayerID
	playedPlayerIDs    []types.PlayerID
	fID2pIDs           map[types.FactionID][]types.PlayerID
	rID2pIDs           map[types.RoleID][]types.PlayerID
	players            map[types.PlayerID]contract.Player
	polls              map[types.FactionID]contract.Poll
}

var _ contract.Game = (*game)(nil)

func NewGame(id types.GameID, setting *types.GameSetting) contract.Game {
	game := game{
		id:                 id,
		numberWerewolves:   setting.NumberWerewolves,
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		status:             Idle,
		turnDuration:       time.Duration(setting.TurnDuration) * time.Second,
		discussionDuration: time.Duration(setting.DiscussionDuration) * time.Second,
		roleIDs:            setting.RoleIDs,
		requiredRoleIDs:    setting.RequiredRoleIDs,
		scheduler:          NewScheduler(role.NightPhaseID),
		deadPlayerIDs:      make([]types.PlayerID, len(setting.PlayerIDs)),
		fID2pIDs:           make(map[types.FactionID][]types.PlayerID),
		rID2pIDs:           make(map[types.RoleID][]types.PlayerID),
		players:            make(map[types.PlayerID]contract.Player),
		polls:              make(map[types.FactionID]contract.Poll),
	}

	for _, id := range setting.PlayerIDs {
		playerID := types.PlayerID(id)
		game.players[playerID] = NewPlayer(&game, playerID)
	}

	// Create polls for villagers and werewolves
	game.polls[role.VillagerFactionID] = NewPoll()
	game.polls[role.WerewolfFactionID] = NewPoll()

	return &game
}

func (g game) ID() types.GameID {
	return g.id
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

func (g game) PlayerIDsByRoleID(roleID types.RoleID) []types.PlayerID {
	return g.rID2pIDs[roleID]
}

func (g game) PlayerIDsByFactionID(factionID types.FactionID) []types.PlayerID {
	return g.fID2pIDs[factionID]
}

func (g game) WerewolfPlayerIDs() []types.PlayerID {
	return g.fID2pIDs[role.WerewolfFactionID]
}

func (g game) NonWerewolfPlayerIDs() []types.PlayerID {
	var nonWerewolfPlayerIDs []types.PlayerID

	for factionID, playerIDs := range g.fID2pIDs {
		if factionID != role.WerewolfFactionID {
			nonWerewolfPlayerIDs = append(nonWerewolfPlayerIDs, playerIDs...)
		}
	}

	return nonWerewolfPlayerIDs
}

func (g game) AlivePlayerIDs(roleID types.RoleID) []types.PlayerID {
	var playerIDs []types.PlayerID

	for _, player := range g.players {
		if slices.Contains(player.RoleIDs(), roleID) &&
			!slices.Contains(g.deadPlayerIDs, player.ID()) {
			playerIDs = append(playerIDs, player.ID())
		}
	}

	return playerIDs
}

func (g *game) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID types.RoleID) bool {
	isWerewolf := slices.Contains(
		role.FactionID2RoleIDs[role.WerewolfFactionID],
		roleID,
	)

	for i := 0; i < int(role.RoleSets[roleID]); i++ {
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
		if ok, _ := player.AssignRole(role.VillagerRoleID); ok {
			g.rID2pIDs[role.VillagerRoleID] = append(
				g.rID2pIDs[role.VillagerRoleID],
				player.ID(),
			)
			g.scheduler.AddPlayerTurn(types.NewPlayerTurn{
				PhaseID:      role.DayPhaseID,
				RoleID:       role.VillagerRoleID,
				BeginRoundID: types.RoundID(0),
				TurnID:       role.MidTurn,
				PlayerID:     player.ID(),
			})
		}

		if selectedRole == nil {
			continue
		}

		// Default werewolf faction's role
		if selectedRole.FactionID() == role.WerewolfFactionID {
			if ok, _ := player.AssignRole(role.WerewolfRoleID); ok {
				g.rID2pIDs[role.WerewolfRoleID] = append(
					g.rID2pIDs[role.WerewolfRoleID],
					player.ID(),
				)
				g.scheduler.AddPlayerTurn(types.NewPlayerTurn{
					PhaseID:      role.NightPhaseID,
					RoleID:       role.WerewolfRoleID,
					BeginRoundID: types.RoundID(0),
					TurnID:       role.MidTurn,
					PlayerID:     player.ID(),
				})
			}

		}

		// Main role
		if ok, _ := player.AssignRole(selectedRole.ID()); ok {
			g.rID2pIDs[selectedRole.ID()] = append(
				g.rID2pIDs[selectedRole.ID()],
				player.ID(),
			)

			// Add the main role's turn to the schedule
			g.scheduler.AddPlayerTurn(types.NewPlayerTurn{
				PhaseID:      selectedRole.PhaseID(),
				RoleID:       selectedRole.ID(),
				BeginRoundID: selectedRole.BeginRoundID(),
				TurnID:       selectedRole.TurnID(),
				PlayerID:     player.ID(),
			})
		}

		// Assign the main role's faction to the player
		g.fID2pIDs[selectedRole.FactionID()] = append(
			g.fID2pIDs[player.FactionID()],
			player.ID(),
		)
	}
}

func (g *game) randomRoleIDs() {
	werewolfCounter := len(g.WerewolfPlayerIDs())
	nonWerewolfCounter := len(g.NonWerewolfPlayerIDs())

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
	g.polls[role.VillagerFactionID].AddCandidates(maps.Keys(g.players)...)
	g.polls[role.WerewolfFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
}

func (g *game) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(util.Config().Game.PreparationDuration)*time.Second,
	)
	defer cancel()

	// Wait for timeout
	select {
	case <-g.finishSignal:
		// Finish game
	case <-ctx.Done():
		g.status = Running
	}
}

func (g *game) runScheduler() {
	// Wait a little bit for the player to prepare
	g.waitForPreparation()
	g.scheduler.NextTurn()

	fmt.Println("Starttttttttttttt")

	for g.status == Running {
		g.playedPlayerIDs = make([]types.PlayerID, 0)

		func() {
			var duration time.Duration

			fmt.Println("round:", g.scheduler.RoundID())
			fmt.Println("turn:", g.scheduler.TurnID())

			if g.scheduler.PhaseID() == role.DayPhaseID &&
				g.scheduler.TurnID() == role.MidTurn {
				duration = g.discussionDuration

				g.Poll(role.VillagerFactionID).Open()
				defer g.Poll(role.VillagerFactionID).Close()
			} else {
				duration = g.turnDuration

				if g.scheduler.PhaseID() == role.NightPhaseID &&
					g.scheduler.TurnID() == role.MidTurn {
					g.Poll(role.WerewolfFactionID).Open()
					defer g.Poll(role.WerewolfFactionID).Close()
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			// Wait for signal or timeout
			select {
			case <-g.nextTurnSignal:
				g.scheduler.NextTurn()
				fmt.Println("Time up")
			case <-ctx.Done():
				g.mutex.Lock()
				defer g.mutex.Unlock()
				g.scheduler.NextTurn()
				fmt.Println("Done")
			case <-g.finishSignal:
				fmt.Println("Finished")
			}
		}()
	}
}

func (g *game) Start() int64 {
	if g.status != Idle {
		return -1
	}

	g.randomRoleIDs()
	g.addCandidatesToPolls()
	go g.runScheduler()
	g.status = Waiting

	return time.Now().Unix()
}

func (g *game) Finish() bool {
	if g.status == Finished {
		return false
	}

	g.finishSignal <- true
	g.status = Finished
	close(g.nextTurnSignal)
	close(g.finishSignal)

	return true
}

func (g *game) Play(playerID types.PlayerID, req types.ActivateAbilityRequest) types.ActionResponse {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Validate
	var msg string
	if g.status != Running {
		msg = "Wait until the game starts (╥﹏╥)"
	} else if g.Player(playerID) == nil {
		msg = "Unable to play this game (╥﹏╥)"
	} else if slices.Contains(g.playedPlayerIDs, playerID) {
		msg = "Wait for next turn (╥﹏╥)"
	} else if slices.Contains(g.deadPlayerIDs, playerID) {
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

	res := g.players[playerID].ActivateAbility(req)
	if res.Ok {
		g.playedPlayerIDs = append(g.playedPlayerIDs, playerID)
		// Move to the next turn if all players have finished their job
		if len(g.playedPlayerIDs) == len(g.scheduler.PlayablePlayerIDs()) {
			g.nextTurnSignal <- true
		}
	}

	return res
}

func (g *game) ExitPlayer(playerID types.PlayerID) bool {
	if g.players[playerID] == nil ||
		slices.Contains(g.exitedPlayerIDs, playerID) {
		return false
	}

	g.KillPlayer(playerID, true)
	g.exitedPlayerIDs = append(g.exitedPlayerIDs, playerID)

	return true
}

func (g *game) KillPlayer(playerID types.PlayerID, isExited bool) contract.Player {
	player := g.players[playerID]
	if player == nil ||
		slices.Contains(g.deadPlayerIDs, playerID) ||
		!player.Die(isExited) {
		return nil
	}

	g.polls[role.VillagerFactionID].RemoveElector(player.ID())
	g.polls[role.VillagerFactionID].RemoveCandidate(player.ID())
	if player.FactionID() == role.WerewolfFactionID {
		g.polls[role.WerewolfFactionID].RemoveElector(player.ID())
	} else {
		g.polls[role.WerewolfFactionID].RemoveCandidate(player.ID())
	}

	g.deadPlayerIDs = append(g.deadPlayerIDs, playerID)

	return player
}
