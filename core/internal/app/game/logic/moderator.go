package logic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/config"
	"uwwolf/pkg/util"

	"golang.org/x/exp/slices"
)

// moderator controlls a game.
type moderator struct {
	gameID types.GameID

	// gameStatus is the current game status ID.
	gameStatus types.GameStatusID

	world               contract.World
	config              config.Game
	scheduler           contract.Scheduler
	mutex               *sync.Mutex
	nextTurnSignal      chan bool
	finishSignal        chan bool
	turnDuration        time.Duration
	discussionDuration  time.Duration
	playedPlayerID      []types.PlayerId
	winningFaction      types.FactionId
	onPhaseChanged      func(mod contract.Moderator)
	actionRegistrations []types.ExecuteActionRegistration
}

func NewModerator(config config.Game, reg *types.GameRegistration) contract.Moderator {
	m := &moderator{
		gameID:              reg.ID,
		gameStatus:          constants.Idle,
		config:              config,
		nextTurnSignal:      make(chan bool),
		finishSignal:        make(chan bool),
		mutex:               new(sync.Mutex),
		turnDuration:        reg.TurnDuration,
		discussionDuration:  reg.DiscussionDuration,
		scheduler:           NewScheduler(constants.NightPhaseId),
		actionRegistrations: make([]types.ExecuteActionRegistration, 0),
	}
	m.world = NewWorld(m, &types.GameInitialization{
		RoleIds:          reg.RoleIds,
		RequiredRoleIds:  reg.RequiredRoleIds,
		NumberWerewolves: reg.NumberWerewolves,
		PlayerIDs:        reg.PlayerIDs,
	})

	return m
}

func (m *moderator) RegisterActionExecution(regis types.ExecuteActionRegistration) {
	m.actionRegistrations = append(m.actionRegistrations, regis)
}

func (m *moderator) OnPhaseChanged(fn func(mod contract.Moderator)) {
	m.onPhaseChanged = fn
}

func (m moderator) GameID() types.GameID {
	return m.gameID
}

func (m moderator) Scheduler() contract.Scheduler {
	return m.scheduler
}

func (m moderator) World() contract.World {
	return m.world
}

// StatusID retusn current world status ID.
func (m moderator) GameStatus() types.GameStatusID {
	return m.gameStatus
}

func (m moderator) Player(ID types.PlayerId) contract.Player {
	return m.world.Player(ID)
}

// checkWinConditions checks if any faction satisfies its win condition,
// if any, finish the game.
func (m *moderator) checkWinConditions() {
	m.mutex.Lock()
	if len(m.world.AlivePlayerIdsWithFactionId(constants.WerewolfFactionId)) == 0 {
		// Villager wins if all werewolves are dead
		m.winningFaction = constants.VillagerFactionId
	} else if len(m.world.AlivePlayerIdsWithFactionId(constants.WerewolfFactionId)) >=
		len(m.world.AlivePlayerIdsWithoutFactionId(constants.WerewolfFactionId)) {
		// Werewolf wins if the number is overwhelming or equal to villager
		m.winningFaction = constants.WerewolfFactionId
	}
	m.mutex.Unlock()

	if !util.IsZero(m.winningFaction) {
		m.FinishGame()
	}
}

// handlePoll handles poll result of each faction.
func (m moderator) handlePoll(factionID types.FactionId) {
	if poll := m.world.Poll(factionID); poll != nil && poll.Close() {
		if record := poll.Record(constants.ZeroRound); !util.IsZero(record.WinnerId) {
			if player := m.world.Player(record.WinnerId); player != nil {
				player.Die()
			}
		}
	}
}

// runScheduler switches turns automatically.
func (m *moderator) runScheduler() {
	for m.GameStatus() == constants.Starting {
		m.mutex.Lock()
		m.playedPlayerID = make([]types.PlayerId, 0)
		m.scheduler.NextTurn()

		for i, regis := range m.actionRegistrations {
			if regis.CanExecute() {
				regis.Exec()
				// Check this carefully whether for loop skip the next element or not
				m.actionRegistrations = slices.Delete(m.actionRegistrations, i, 1)
			}
		}

		func() {
			var duration time.Duration

			if m.scheduler.PhaseId() == constants.DayPhaseId &&
				m.scheduler.Turn() == constants.MidTurn {
				duration = m.discussionDuration

				m.world.Poll(constants.VillagerFactionId).Open() // nolint: errcheck
				defer m.handlePoll(constants.VillagerFactionId)
			} else {
				duration = m.turnDuration

				if m.scheduler.PhaseId() == constants.NightPhaseId &&
					m.scheduler.Turn() == constants.MidTurn {
					m.world.Poll(constants.WerewolfFactionId).Open() // nolint: errcheck
					defer m.handlePoll(constants.WerewolfFactionId)
				}
			}

			// Notify new turn is started and its duration
			fmt.Println("New turn!!!")
			m.mutex.Unlock()

			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			select {
			case <-m.nextTurnSignal:
				m.checkWinConditions()
			case <-ctx.Done():
				m.checkWinConditions()
			case <-m.finishSignal:
				m.FinishGame()
			}

			m.onPhaseChanged(m)
		}()
	}
}

// waitForPreparation waits a bit before the game starts.
func (m *moderator) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		m.config.PreparationDuration,
	)
	defer cancel()

	select {
	case <-m.finishSignal:
		m.FinishGame()
	case <-ctx.Done():
	}

	fmt.Println("Preparation is done")
}

// StartGame starts the game.
func (m *moderator) StartGame() int64 {
	if m.gameID.IsUnknown() || m.GameStatus() != constants.Idle {
		return -1
	}

	fmt.Println("Starting")
	m.world.Load()

	go func() {
		m.gameStatus = constants.Waiting
		m.waitForPreparation()
		m.gameStatus = constants.Starting
		// m.gameID.ObservePlayers(maps.Keys(m.world.Players()))
		go m.runScheduler()
	}()

	return 1
}

// FinishGame ends the game.
func (m *moderator) FinishGame() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.GameStatus() == constants.Finished {
		return false
	}

	// Notify winner
	fmt.Println("Winner is %w", m.winningFaction)

	m.finishSignal <- true
	close(m.finishSignal)
	close(m.nextTurnSignal)
	m.gameStatus = constants.Finished

	return true
}

// RequestPlay receives the play request from the player.
func (m *moderator) RequestPlay(
	playerID types.PlayerId,
	req *types.RoleRequest,
) *types.RoleResponse {
	if !m.mutex.TryLock() {
		return &types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "Turn is over!",
			},
		}
	}
	defer m.mutex.Unlock()

	if slices.Contains(m.playedPlayerID, playerID) {
		return &types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "You played this turn!",
			},
		}
	}

	player := m.world.Player(playerID)
	if player == nil {
		return &types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "Non-existent player!",
			},
		}
	}

	res := player.ActivateAbility(req)
	if res.Ok {
		m.playedPlayerID = append(m.playedPlayerID, playerID)

		// Move to the next turn if all players have finished their job
		// if len(m.playedPlayerID) == len(m.scheduler.PlayablePlayerID()) {
		// 	m.nextTurnSignal <- true
		// }

		// Cache player request
		// m.rdb.LPush(
		// 	context.Background(),
		// 	fmt.Sprint(m.worldID),
		// 	fmt.Sprint(res.TargetID),
		// 	fmt.Sprint(res.ActionID),
		// 	fmt.Sprint(res.RoleID),
		// 	fmt.Sprint(playerID),
		// 	fmt.Sprint(res.Round),
		// )
	}

	return res
}
