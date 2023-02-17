package core

import (
	"context"
	"fmt"
	"sync"
	"time"
	"uwwolf/config"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/factory"
	"uwwolf/game/types"
	"uwwolf/util"

	"golang.org/x/exp/slices"
)

type game struct {
	id                    enum.GameID
	numberWerewolves      uint8
	nextTurnSignal        chan bool
	finishSignal          chan bool
	mutex                 sync.Mutex
	status                enum.GameStatus
	turnDuration          time.Duration
	discussionDuration    time.Duration
	scheduler             contract.Scheduler
	roleIDs               []enum.RoleID
	requiredRoleIDs       []enum.RoleID
	selectedRoleIDs       []enum.RoleID
	deadPlayerIDs         []enum.PlayerID
	disconnectedPlayerIDs []enum.PlayerID
	exitedPlayerIDs       []enum.PlayerID
	playedPlayerIDs       []enum.PlayerID
	fID2pIDs              map[enum.FactionID][]enum.PlayerID
	rID2pIDs              map[enum.RoleID][]enum.PlayerID
	players               map[enum.PlayerID]contract.Player
	polls                 map[enum.FactionID]contract.Poll
}

func NewGame(id enum.GameID, setting *types.GameSetting) contract.Game {
	game := game{
		id:                 id,
		numberWerewolves:   setting.NumberWerewolves,
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		status:             enum.Idle,
		turnDuration:       time.Duration(setting.TurnDuration) * time.Second,
		discussionDuration: time.Duration(setting.DiscussionDuration) * time.Second,
		roleIDs:            setting.RoleIDs,
		requiredRoleIDs:    setting.RequiredRoleIDs,
		scheduler:          NewScheduler(enum.NightPhaseID),
		deadPlayerIDs:      make([]enum.PlayerID, len(setting.PlayerIDs)),
		fID2pIDs:           make(map[enum.FactionID][]enum.PlayerID),
		rID2pIDs:           make(map[enum.RoleID][]enum.PlayerID),
		players:            make(map[enum.PlayerID]contract.Player),
		polls:              make(map[enum.FactionID]contract.Poll),
	}

	for _, id := range setting.PlayerIDs {
		playerID := enum.PlayerID(id)
		game.players[playerID] = NewPlayer(&game, playerID)
	}

	// Create polls for villagers and werewolves
	game.polls[enum.VillagerFactionID], _ = NewPoll(
		uint8(len(game.players)),
	)
	game.polls[enum.WerewolfFactionID], _ = NewPoll(
		game.numberWerewolves,
	)

	return &game
}

func (g *game) ID() enum.GameID {
	return g.id
}

func (g *game) Scheduler() contract.Scheduler {
	return g.scheduler
}

func (g *game) Poll(facitonID enum.FactionID) contract.Poll {
	return g.polls[facitonID]
}

func (g *game) Player(playerId enum.PlayerID) contract.Player {
	return g.players[playerId]
}

func (g *game) PlayerIDsByRoleID(roleID enum.RoleID) []enum.PlayerID {
	return g.rID2pIDs[roleID]
}

func (g *game) PlayerIDsByFactionID(factionID enum.FactionID) []enum.PlayerID {
	return g.fID2pIDs[factionID]
}

func (g *game) WerewolfPlayerIDs() []enum.PlayerID {
	return g.fID2pIDs[enum.WerewolfFactionID]
}

func (g *game) NonWerewolfPlayerIDs() []enum.PlayerID {
	var nonWerewolfPlayerIDs []enum.PlayerID

	for factionID, playerIDs := range g.fID2pIDs {
		if factionID != enum.WerewolfFactionID {
			nonWerewolfPlayerIDs = append(nonWerewolfPlayerIDs, playerIDs...)
		}
	}

	return nonWerewolfPlayerIDs
}

func (g *game) AlivePlayerIDs(roleID enum.RoleID) []enum.PlayerID {
	var playerIDs []enum.PlayerID

	for _, player := range g.players {
		if slices.Contains(player.RoleIDs(), roleID) &&
			!slices.Contains(g.deadPlayerIDs, player.ID()) {
			playerIDs = append(playerIDs, player.ID())
		}
	}

	return playerIDs
}

func (g *game) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID enum.RoleID) bool {
	isWerewolf := slices.Contains(
		types.RoleIDsByFactionID[enum.WerewolfFactionID],
		roleID,
	)

	for i := 0; i < int(types.RoleIDSets[roleID]); i++ {
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
		selectedRole, _ := factory.NewRole(selectedRoleID, g, player.ID())

		// Remove the random role from array
		if i != -1 {
			selectedRoleIDs = slices.Delete(selectedRoleIDs, i, i+1)
		}

		// Default role
		if ok, _ := player.AssignRole(enum.VillagerRoleID); ok {
			g.rID2pIDs[enum.VillagerRoleID] = append(
				g.rID2pIDs[enum.VillagerRoleID],
				player.ID(),
			)
		}

		if selectedRole == nil {
			continue
		}

		// Default werewolf faction's role
		if selectedRole.FactionID() == enum.WerewolfFactionID {
			if ok, _ := player.AssignRole(enum.WerewolfRoleID); ok {
				g.rID2pIDs[enum.WerewolfRoleID] = append(
					g.rID2pIDs[enum.WerewolfRoleID],
					player.ID(),
				)
			}

		}

		// Main role
		if ok, _ := player.AssignRole(selectedRole.ID()); ok {
			g.rID2pIDs[selectedRole.ID()] = append(
				g.rID2pIDs[selectedRole.ID()],
				player.ID(),
			)

			// Add the main role's turn to the schedule
			g.scheduler.AddTurn(&types.TurnSetting{
				PhaseID:    selectedRole.PhaseID(),
				RoleID:     selectedRole.ID(),
				BeginRound: selectedRole.BeginRound(),
				Priority:   selectedRole.Priority(),
				Limit:      selectedRole.ActiveLimit(0),
				Position:   enum.SortedPosition,
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

func (g *game) addDefaultTurnsToSchedule() {
	g.scheduler.AddTurn(&types.TurnSetting{
		PhaseID:    enum.DayPhaseID,
		RoleID:     enum.VillagerRoleID,
		BeginRound: enum.FirstRound,
		Priority:   enum.VillagerTurnPriority,
		Limit:      enum.Unlimited,
		Position:   enum.SortedPosition,
	})
	g.scheduler.AddTurn(&types.TurnSetting{
		PhaseID:    enum.NightPhaseID,
		RoleID:     enum.WerewolfRoleID,
		BeginRound: enum.FirstRound,
		Priority:   enum.WerewolfTurnPriority,
		Limit:      enum.Unlimited,
		Position:   enum.SortedPosition,
	})
}

func (g *game) addCandidatesToPolls() {
	g.polls[enum.VillagerFactionID].AddCandidates(g.WerewolfPlayerIDs()...)
	g.polls[enum.VillagerFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
	g.polls[enum.WerewolfFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
}

func (g *game) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.Game().PreparationTime)*time.Second,
	)
	defer cancel()

	fmt.Println("Preparing...")

	// Wait for timeout
	select {
	case <-g.finishSignal:
		// Finish game
	case <-ctx.Done():
		g.status = enum.Running
	}
}

func (g *game) runScheduler() {
	// Wait a little bit for the player to prepare
	g.waitForPreparation()
	g.scheduler.NextTurn(false)

	fmt.Println("Starttttttttttttt")

	for g.status == enum.Running {
		g.playedPlayerIDs = make([]enum.PlayerID, 0)

		func() {
			var duration time.Duration

			fmt.Println("turn of", g.scheduler.Turn().RoleID)

			// if g.

			if g.scheduler.Turn().RoleID == enum.VillagerRoleID {
				duration = g.discussionDuration

				g.Poll(enum.VillagerFactionID).Open()
				defer g.Poll(enum.VillagerFactionID).Close()
			} else {
				duration = g.turnDuration

				if g.scheduler.Turn().RoleID == enum.WerewolfRoleID {
					g.Poll(enum.WerewolfFactionID).Open()
					defer g.Poll(enum.WerewolfFactionID).Close()
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			// Wait for signal or timeout
			select {
			case <-g.nextTurnSignal:
				g.scheduler.NextTurn(false)
				fmt.Println("Time up")
			case <-ctx.Done():
				g.scheduler.NextTurn(false)
				fmt.Println("Done")
			case <-g.finishSignal:
				fmt.Println("Finished")
			}
		}()
	}
}

func (g *game) Start() int64 {
	if g.status != enum.Idle {
		return -1
	}

	g.randomRoleIDs()
	g.addCandidatesToPolls()
	g.addDefaultTurnsToSchedule()

	go g.runScheduler()

	return time.Now().Unix()
}

func (g *game) Finish() bool {
	if g.status == enum.Finished {
		return false
	}

	g.finishSignal <- true
	g.status = enum.Finished
	close(g.nextTurnSignal)
	close(g.finishSignal)

	return true
}

func (g *game) UsePlayerRole(playerID enum.PlayerID, req *types.UseRoleRequest) *types.ActionResponse {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Check actor + target
	if g.status != enum.Running ||
		g.Player(playerID) == nil ||
		slices.Contains(g.playedPlayerIDs, playerID) ||
		slices.Contains(g.deadPlayerIDs, playerID) ||
		slices.Contains(g.disconnectedPlayerIDs, playerID) ||
		!slices.Contains(g.rID2pIDs[g.scheduler.Turn().RoleID], playerID) {

		return &types.ActionResponse{
			Ok:      false,
			Message: "Not your turn or you're died (╥﹏╥)",
		}
	}

	res := g.Player(playerID).UseAbility(req)

	if res.Ok {
		g.playedPlayerIDs = append(g.playedPlayerIDs, playerID)

		if len(g.playedPlayerIDs) == len(g.AlivePlayerIDs(g.scheduler.Turn().RoleID)) {
			g.nextTurnSignal <- true
		}
	}

	return res
}

func (g *game) ConnectPlayer(playerID enum.PlayerID, isConnected bool) bool {
	if g.status != enum.Running || g.players[playerID] == nil {
		return false
	}

	disconnectedIndex := slices.Index(g.disconnectedPlayerIDs, playerID)

	if isConnected && disconnectedIndex != -1 {
		g.disconnectedPlayerIDs = slices.Delete(
			g.disconnectedPlayerIDs,
			disconnectedIndex,
			disconnectedIndex+1,
		)
	} else if !isConnected && disconnectedIndex == -1 {
		g.disconnectedPlayerIDs = append(g.disconnectedPlayerIDs, playerID)
	} else {
		return false
	}

	return true
}

func (g *game) ExitPlayer(playerID enum.PlayerID) bool {
	if g.status != enum.Running ||
		g.players[playerID] == nil ||
		slices.Contains(g.exitedPlayerIDs, playerID) {
		return false
	}

	g.KillPlayer(playerID, true)
	g.exitedPlayerIDs = append(g.exitedPlayerIDs, playerID)

	return true
}

func (g *game) KillPlayer(playerID enum.PlayerID, isExited bool) contract.Player {
	player := g.players[playerID]

	if g.status != enum.Running ||
		player == nil ||
		slices.Contains(g.deadPlayerIDs, playerID) ||
		!player.Die(isExited) {
		return nil
	}

	g.polls[enum.VillagerFactionID].RemoveElector(player.ID())
	g.polls[enum.VillagerFactionID].RemoveCandidate(player.ID())

	if player.FactionID() == enum.WerewolfFactionID {
		g.polls[enum.WerewolfFactionID].RemoveElector(player.ID())
	} else {
		g.polls[enum.WerewolfFactionID].RemoveCandidate(player.ID())
	}

	g.deadPlayerIDs = append(g.deadPlayerIDs, playerID)

	return player
}
