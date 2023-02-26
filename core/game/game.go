package game

import (
	"context"
	"fmt"
	"sync"
	"time"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/util"

	"golang.org/x/exp/slices"
)

const (
	Idle types.GameStatus = iota
	Waiting
	Running
	Finished
)

type game struct {
	id                    types.GameID
	numberWerewolves      uint8
	nextTurnSignal        chan bool
	finishSignal          chan bool
	mutex                 sync.Mutex
	status                types.GameStatus
	turnDuration          time.Duration
	discussionDuration    time.Duration
	scheduler             contract.Scheduler
	roleIDs               []types.RoleID
	requiredRoleIDs       []types.RoleID
	selectedRoleIDs       []types.RoleID
	deadPlayerIDs         []types.PlayerID
	disconnectedPlayerIDs []types.PlayerID
	exitedPlayerIDs       []types.PlayerID
	playedPlayerIDs       []types.PlayerID
	fID2pIDs              map[types.FactionID][]types.PlayerID
	rID2pIDs              map[types.RoleID][]types.PlayerID
	players               map[types.PlayerID]contract.Player
	polls                 map[types.FactionID]contract.Poll
}

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
		scheduler:          NewScheduler(NightPhaseID),
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
	game.polls[VillagerFactionID], _ = NewPoll(
		uint8(len(game.players)),
	)
	game.polls[WerewolfFactionID], _ = NewPoll(
		game.numberWerewolves,
	)

	return &game
}

func (g *game) ID() types.GameID {
	return g.id
}

func (g *game) Scheduler() contract.Scheduler {
	return g.scheduler
}

func (g *game) Poll(facitonID types.FactionID) contract.Poll {
	return g.polls[facitonID]
}

func (g *game) Player(playerId types.PlayerID) contract.Player {
	return g.players[playerId]
}

func (g *game) PlayerIDsByRoleID(roleID types.RoleID) []types.PlayerID {
	return g.rID2pIDs[roleID]
}

func (g *game) PlayerIDsByFactionID(factionID types.FactionID) []types.PlayerID {
	return g.fID2pIDs[factionID]
}

func (g *game) WerewolfPlayerIDs() []types.PlayerID {
	return g.fID2pIDs[WerewolfFactionID]
}

func (g *game) NonWerewolfPlayerIDs() []types.PlayerID {
	var nonWerewolfPlayerIDs []types.PlayerID

	for factionID, playerIDs := range g.fID2pIDs {
		if factionID != WerewolfFactionID {
			nonWerewolfPlayerIDs = append(nonWerewolfPlayerIDs, playerIDs...)
		}
	}

	return nonWerewolfPlayerIDs
}

func (g *game) AlivePlayerIDs(roleID types.RoleID) []types.PlayerID {
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
		FactionID2RoleIDs[WerewolfFactionID],
		roleID,
	)

	for i := 0; i < int(RoleSets[roleID]); i++ {
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
		if ok, _ := player.AssignRole(VillagerRoleID); ok {
			g.rID2pIDs[VillagerRoleID] = append(
				g.rID2pIDs[VillagerRoleID],
				player.ID(),
			)
		}

		if selectedRole == nil {
			continue
		}

		// Default werewolf faction's role
		if selectedRole.FactionID() == WerewolfFactionID {
			if ok, _ := player.AssignRole(WerewolfRoleID); ok {
				g.rID2pIDs[WerewolfRoleID] = append(
					g.rID2pIDs[WerewolfRoleID],
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
				// Priority:   selectedRole.Priority(),
				Limit:    selectedRole.ActiveLimit(0),
				Position: SortedPosition,
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
		PhaseID:    DayPhaseID,
		RoleID:     VillagerRoleID,
		BeginRound: FirstRound,
		// Priority:   VillagerTurnPriority,
		Limit:    Unlimited,
		Position: SortedPosition,
	})
	g.scheduler.AddTurn(&types.TurnSetting{
		PhaseID:    NightPhaseID,
		RoleID:     WerewolfRoleID,
		BeginRound: FirstRound,
		// Priority:   WerewolfTurnPriority,
		Limit:    Unlimited,
		Position: SortedPosition,
	})
}

func (g *game) addCandidatesToPolls() {
	g.polls[VillagerFactionID].AddCandidates(g.WerewolfPlayerIDs()...)
	g.polls[VillagerFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
	g.polls[WerewolfFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
}

func (g *game) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(util.Config().Game.PreparationDuration)*time.Second,
	)
	defer cancel()

	fmt.Println("Preparing...")

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
	g.scheduler.NextTurn(false)

	fmt.Println("Starttttttttttttt")

	for g.status == Running {
		g.playedPlayerIDs = make([]types.PlayerID, 0)

		func() {
			var duration time.Duration

			fmt.Println("turn of", g.scheduler.Turn().RoleID)

			// if g.

			if g.scheduler.Turn().RoleID == VillagerRoleID {
				duration = g.discussionDuration

				g.Poll(VillagerFactionID).Open()
				defer g.Poll(VillagerFactionID).Close()
			} else {
				duration = g.turnDuration

				if g.scheduler.Turn().RoleID == WerewolfRoleID {
					g.Poll(WerewolfFactionID).Open()
					defer g.Poll(WerewolfFactionID).Close()
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
	if g.status != Idle {
		return -1
	}

	g.randomRoleIDs()
	g.addCandidatesToPolls()
	g.addDefaultTurnsToSchedule()

	go g.runScheduler()

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

func (g *game) UsePlayerRole(playerID types.PlayerID, req types.ExecuteActionRequest) types.ActionResponse {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Check actor + target
	if g.status != Running ||
		g.Player(playerID) == nil ||
		slices.Contains(g.playedPlayerIDs, playerID) ||
		slices.Contains(g.deadPlayerIDs, playerID) ||
		slices.Contains(g.disconnectedPlayerIDs, playerID) ||
		!slices.Contains(g.rID2pIDs[g.scheduler.Turn().RoleID], playerID) {

		return types.ActionResponse{
			Ok:      false,
			Message: "Not your turn or you're died (╥﹏╥)",
		}
	}

	res := g.Player(playerID).ExecuteAction(req)

	if res.Ok {
		g.playedPlayerIDs = append(g.playedPlayerIDs, playerID)

		if len(g.playedPlayerIDs) == len(g.AlivePlayerIDs(g.scheduler.Turn().RoleID)) {
			g.nextTurnSignal <- true
		}
	}

	return res
}

func (g *game) ConnectPlayer(playerID types.PlayerID, isConnected bool) bool {
	if g.status != Running || g.players[playerID] == nil {
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

func (g *game) ExitPlayer(playerID types.PlayerID) bool {
	if g.status != Running ||
		g.players[playerID] == nil ||
		slices.Contains(g.exitedPlayerIDs, playerID) {
		return false
	}

	// g.KillPlayer(playerID, true)
	g.KillPlayer(playerID)
	g.exitedPlayerIDs = append(g.exitedPlayerIDs, playerID)

	return true
}

func (g *game) KillPlayer(playerID types.PlayerID) contract.Player {
	// Out game thi noi tai se ko kich hoat
	player := g.players[playerID]

	if g.status != Running ||
		player == nil ||
		slices.Contains(g.deadPlayerIDs, playerID) {
		// !player.Die(isExited) {
		return nil
	}

	g.polls[VillagerFactionID].RemoveElector(player.ID())
	g.polls[VillagerFactionID].RemoveCandidate(player.ID())

	if player.FactionID() == WerewolfFactionID {
		g.polls[WerewolfFactionID].RemoveElector(player.ID())
	} else {
		g.polls[WerewolfFactionID].RemoveCandidate(player.ID())
	}

	g.deadPlayerIDs = append(g.deadPlayerIDs, playerID)

	return player
}
