package game

// import (
// 	"context"
// 	"time"
// 	"uwwolf/app/game/config"
// 	"uwwolf/app/game/contract"
// 	"uwwolf/app/game/core/player"
// 	"uwwolf/app/game/core/role"
// 	"uwwolf/app/game/factory"
// 	"uwwolf/app/game/types"
// 	"uwwolf/util"

// 	"golang.org/x/exp/slices"
// )

// type game struct {
// 	id                    types.GameID
// 	numberOfWerewolves    int
// 	nextTurnSignal        chan bool
// 	status                types.GameStatus
// 	turnDuration          time.Duration
// 	discussionDuration    time.Duration
// 	scheduler             contract.Scheduler
// 	roleIDs               []types.RoleID
// 	requiredRoleIDs       []types.RoleID
// 	selectedRoleIDs       []types.RoleID
// 	deadPlayerIDs         []types.PlayerID
// 	disconnectedPlayerIDs []types.PlayerID
// 	exitedPlayerIDs       []types.PlayerID
// 	playedPlayerIDs       []types.PlayerID
// 	fID2pIDs              map[types.FactionID][]types.PlayerID
// 	rID2pIDs              map[types.RoleID][]types.PlayerID
// 	players               map[types.PlayerID]contract.Player
// 	polls                 map[types.FactionID]contract.Poll
// }

// func NewGame(id types.GameID, setting *types.GameSetting) contract.Game {
// 	game := game{
// 		id:                 id,
// 		numberOfWerewolves: setting.NumberOfWerewolves,
// 		nextTurnSignal:     make(chan bool),
// 		status:             config.Idle,
// 		turnDuration:       setting.TurnDuration * time.Second,
// 		discussionDuration: setting.DiscussionDuration * time.Second,
// 		roleIDs:            setting.RoleIDs,
// 		requiredRoleIDs:    setting.RequiredRoleIDs,
// 		scheduler:          NewScheduler(config.NightPhaseID),
// 		deadPlayerIDs:      make([]types.PlayerID, len(setting.PlayerIDs)),
// 		fID2pIDs:           make(map[types.FactionID][]types.PlayerID),
// 		rID2pIDs:           make(map[types.RoleID][]types.PlayerID),
// 		players:            make(map[types.PlayerID]contract.Player),
// 		polls:              make(map[types.FactionID]contract.Poll),
// 	}

// 	for _, id := range setting.PlayerIDs {
// 		game.players[id] = player.NewPlayer(&game, id)
// 	}

// 	// Create polls for villagers and werewolves
// 	game.polls[config.VillagerFactionID], _ = NewPoll(
// 		uint(len(game.players)),
// 	)
// 	game.polls[config.WerewolfFactionID], _ = NewPoll(
// 		uint(len(game.fID2pIDs[config.WerewolfFactionID])),
// 	)

// 	return &game
// }

// func (g *game) ID() types.GameID {
// 	return g.id
// }

// func (g *game) Scheduler() contract.Scheduler {
// 	return g.scheduler
// }

// func (g *game) Poll(facitonID types.FactionID) contract.Poll {
// 	return g.polls[facitonID]
// }

// func (g *game) Player(playerId types.PlayerID) contract.Player {
// 	return g.players[playerId]
// }

// func (g *game) PlayerIDsByRoleID(roleID types.RoleID) []types.PlayerID {
// 	return g.rID2pIDs[roleID]
// }

// func (g *game) PlayerIDsByFactionID(factionID types.FactionID) []types.PlayerID {
// 	return g.fID2pIDs[factionID]
// }

// func (g *game) WerewolfPlayerIDs() []types.PlayerID {
// 	return g.fID2pIDs[config.WerewolfFactionID]
// }

// func (g *game) NonWerewolfPlayerIDs() []types.PlayerID {
// 	var nonWerewolfPlayerIDs []types.PlayerID

// 	for factionID, playerIDs := range g.fID2pIDs {
// 		if factionID != config.WerewolfFactionID {
// 			nonWerewolfPlayerIDs = append(nonWerewolfPlayerIDs, playerIDs...)
// 		}
// 	}

// 	return nonWerewolfPlayerIDs
// }

// func (g *game) AlivePlayerIDs(roleID types.RoleID) []types.PlayerID {
// 	var playerIDs []types.PlayerID

// 	for _, player := range g.players {
// 		if slices.Contains(player.RoleIDs(), roleID) &&
// 			!slices.Contains(g.deadPlayerIDs, player.ID()) {
// 			playerIDs = append(playerIDs, player.ID())
// 		}
// 	}

// 	return playerIDs
// }

// func (g *game) selectRoleID(werewolfCounter *int, nonWerewolfCounter *int, roleID types.RoleID) bool {
// 	isWerewolf := slices.Contains(
// 		config.RolesByFaction[config.WerewolfFactionID],
// 		roleID,
// 	)

// 	for i := 0; i < int(config.RoleSets[roleID]); i++ {
// 		isMissingWerewolf := *werewolfCounter < g.numberOfWerewolves
// 		isMissingNonWerewolf := *nonWerewolfCounter < len(g.players)-g.numberOfWerewolves

// 		if !isMissingWerewolf && !isMissingNonWerewolf {
// 			return false
// 		}

// 		if isMissingWerewolf && isWerewolf {
// 			g.selectedRoleIDs = append(g.selectedRoleIDs, roleID)
// 			*werewolfCounter++
// 		} else if isMissingNonWerewolf && !isWerewolf {
// 			g.selectedRoleIDs = append(g.selectedRoleIDs, roleID)
// 			*nonWerewolfCounter++
// 		}
// 	}

// 	return true
// }

// func (g *game) assignRoles() {
// 	selectedRoleIDs := slices.Clone(g.selectedRoleIDs)

// 	for _, player := range g.players {
// 		i, selectedRoleID := util.RandomElement(selectedRoleIDs)
// 		selectedRole, _ := factory.NewRole(selectedRoleID, g, player.ID())

// 		// Remove the random role from array
// 		if i != -1 {
// 			selectedRoleIDs = slices.Delete(selectedRoleIDs, i, i+1)
// 		}

// 		// Default role
// 		g.rID2pIDs[config.VillagerRoleID] = append(
// 			g.rID2pIDs[config.VillagerRoleID],
// 			player.ID(),
// 		)
// 		player.AssignRoles(config.VillagerRoleID)

// 		if selectedRole == nil {
// 			continue
// 		}

// 		// Main role
// 		g.rID2pIDs[selectedRole.ID()] = append(
// 			g.rID2pIDs[selectedRole.ID()],
// 			player.ID(),
// 		)
// 		player.AssignRoles(selectedRole.ID())

// 		// Default werewolf faction's role
// 		if selectedRole.FactionID() == config.WerewolfFactionID {
// 			g.rID2pIDs[config.WerewolfRoleID] = append(
// 				g.rID2pIDs[config.WerewolfRoleID],
// 				player.ID(),
// 			)
// 			player.AssignRoles(config.WerewolfRoleID)
// 		}

// 		// Assign the main role's faction to the player
// 		g.fID2pIDs[selectedRole.FactionID()] = append(
// 			g.fID2pIDs[selectedRole.FactionID()],
// 			player.ID(),
// 		)

// 		// Add the main role's turn to the schedule
// 		g.scheduler.AddTurn(&types.TurnSetting{
// 			PhaseID:    selectedRole.PhaseID(),
// 			RoleID:     selectedRole.ID(),
// 			BeginRound: selectedRole.BeginRound(),
// 			Priority:   selectedRole.Priority(),
// 			Position:   config.SortedPosition,
// 		})
// 	}
// }

// func (g *game) randomRoleIDs() {
// 	werewolfCounter := len(g.WerewolfPlayerIDs())
// 	nonWerewolfCounter := len(g.NonWerewolfPlayerIDs())

// 	// Select required roles
// 	for _, requiredRoleID := range g.requiredRoleIDs {
// 		if !g.selectRoleID(&werewolfCounter, &nonWerewolfCounter, requiredRoleID) {
// 			break
// 		}
// 	}

// 	roleIDs := slices.Clone(g.roleIDs)

// 	// Select random roles
// 	for {
// 		i, randomRoleID := util.RandomElement(roleIDs)
// 		roleIDs = slices.Delete(roleIDs, i, i+1)

// 		if i == -1 ||
// 			!g.selectRoleID(&werewolfCounter, &nonWerewolfCounter, randomRoleID) {
// 			break
// 		}
// 	}

// 	g.assignRoles()
// }

// func (g *game) addDefaultTurnsToSchedule() {
// 	villager, _ := role.NewVillager(nil, types.PlayerID(""))
// 	werewolf, _ := role.NewWerewolf(nil, types.PlayerID(""))

// 	g.scheduler.AddTurn(&types.TurnSetting{
// 		PhaseID:    villager.PhaseID(),
// 		RoleID:     villager.ID(),
// 		BeginRound: villager.BeginRound(),
// 		Priority:   villager.Priority(),
// 		Position:   config.SortedPosition,
// 	})
// 	g.scheduler.AddTurn(&types.TurnSetting{
// 		PhaseID:    werewolf.PhaseID(),
// 		RoleID:     werewolf.ID(),
// 		BeginRound: werewolf.BeginRound(),
// 		Priority:   werewolf.Priority(),
// 		Position:   config.SortedPosition,
// 	})
// }

// func (g *game) addCandidatesToPolls() {
// 	g.polls[config.VillagerFactionID].AddCandidates(g.WerewolfPlayerIDs()...)
// 	g.polls[config.VillagerFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
// 	g.polls[config.WerewolfFactionID].AddCandidates(g.NonWerewolfPlayerIDs()...)
// }

// func (g *game) waitForPreparation() {
// 	ctx, cancel := context.WithTimeout(
// 		context.Background(),
// 		time.Duration(time.Duration(config.PreparationTime).Seconds()),
// 	)
// 	defer cancel()

// 	// Wait for timeout
// 	select {
// 	case <-ctx.Done():
// 	}

// 	g.status = config.Starting
// }

// func (g *game) runScheduler() {
// 	// Wait a little bit for the player to prepare
// 	g.waitForPreparation()

// 	for g.status == config.Starting {
// 		g.playedPlayerIDs = make([]types.PlayerID, 0)

// 		func() {
// 			var duration time.Duration

// 			if g.scheduler.Turn().RoleID == config.VillagerRoleID {
// 				duration = g.discussionDuration

// 				g.Poll(config.VillagerFactionID).Open()
// 				defer g.Poll(config.VillagerFactionID).Close()
// 			} else {
// 				duration = g.turnDuration

// 				if g.scheduler.Turn().RoleID == config.WerewolfRoleID {
// 					g.Poll(config.WerewolfFactionID).Open()
// 					defer g.Poll(config.WerewolfFactionID).Close()
// 				}
// 			}

// 			ctx, cancel := context.WithTimeout(context.Background(), duration)
// 			defer cancel()

// 			// Wait for signal or timeout
// 			select {
// 			case <-g.nextTurnSignal:
// 			case <-ctx.Done():
// 			}

// 			g.scheduler.NextTurn(false)
// 		}()
// 	}
// }

// func (g *game) Start() int64 {
// 	if g.status != config.Idle {
// 		return -1
// 	}

// 	g.randomRoleIDs()
// 	g.addDefaultTurnsToSchedule()
// 	g.addCandidatesToPolls()

// 	go g.runScheduler()

// 	return time.Now().UnixMilli()
// }

// func (g *game) Finish() bool {
// 	if g.status != config.Starting {
// 		return false
// 	}

// 	g.status = config.Finished
// 	close(g.nextTurnSignal)

// 	return true
// }

// func (g *game) UsePlayerRole(playerID types.PlayerID, req *types.UseRoleRequest) *types.ActionResponse {
// 	if g.status != config.Starting ||
// 		g.Player(playerID) == nil ||
// 		slices.Contains(g.playedPlayerIDs, playerID) ||
// 		slices.Contains(g.deadPlayerIDs, playerID) ||
// 		slices.Contains(g.disconnectedPlayerIDs, playerID) ||
// 		!slices.Contains(g.rID2pIDs[g.scheduler.Turn().RoleID], playerID) {

// 		return &types.ActionResponse{
// 			Ok:      false,
// 			Message: "Not your turn or you're died (╥﹏╥)",
// 		}
// 	}

// 	res := g.Player(playerID).UseAbility(req)

// 	if res.Ok {
// 		g.playedPlayerIDs = append(g.playedPlayerIDs, playerID)

// 		if len(g.playedPlayerIDs) == len(g.AlivePlayerIDs(g.scheduler.Turn().RoleID)) {
// 			g.nextTurnSignal <- true
// 		}
// 	}

// 	return res
// }

// func (g *game) ConnectPlayer(playerID types.PlayerID, isConnected bool) bool {
// 	if g.status != config.Starting || g.players[playerID] == nil {
// 		return false
// 	}

// 	disconnectedIndex := slices.Index(g.disconnectedPlayerIDs, playerID)

// 	if isConnected && disconnectedIndex != -1 {
// 		g.disconnectedPlayerIDs = slices.Delete(
// 			g.disconnectedPlayerIDs,
// 			disconnectedIndex,
// 			disconnectedIndex+1,
// 		)
// 	} else if !isConnected && disconnectedIndex == -1 {
// 		g.disconnectedPlayerIDs = append(g.disconnectedPlayerIDs, playerID)
// 	} else {
// 		return false
// 	}

// 	return true
// }

// func (g *game) ExitPlayer(playerID types.PlayerID) bool {
// 	if g.status != config.Starting ||
// 		g.players[playerID] == nil ||
// 		slices.Contains(g.exitedPlayerIDs, playerID) {
// 		return false
// 	}

// 	g.KillPlayer(playerID)
// 	g.exitedPlayerIDs = append(g.exitedPlayerIDs, playerID)

// 	return true
// }

// func (g *game) KillPlayer(playerID types.PlayerID) contract.Player {
// 	player := g.players[playerID]

// 	if g.status != config.Starting ||
// 		player == nil ||
// 		slices.Contains(g.deadPlayerIDs, playerID) {
// 		return nil
// 	}

// 	g.polls[config.VillagerFactionID].RemoveElector(player.ID())
// 	g.polls[config.VillagerFactionID].RemoveCandidate(player.ID())

// 	if player.FactionID() == config.WerewolfFactionID {
// 		g.polls[config.WerewolfFactionID].RemoveElector(player.ID())
// 	} else {
// 		g.polls[config.WerewolfFactionID].RemoveCandidate(player.ID())
// 	}

// 	g.deadPlayerIDs = append(g.deadPlayerIDs, playerID)

// 	return player
// }
