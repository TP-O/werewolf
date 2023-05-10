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

	"golang.org/x/exp/slices"
)

// moderator controlls a game.
type moderator struct {
	gameID types.GameID

	// gameStatus is the current game status ID.
	gameStatus types.GameStatusID

	world              contract.World
	config             config.Game
	scheduler          contract.Scheduler
	mutex              *sync.Mutex
	nextTurnSignal     chan bool
	finishSignal       chan bool
	turnDuration       time.Duration
	discussionDuration time.Duration
	playedPlayerID     []types.PlayerId
	winningFaction     types.FactionID
	onPhaseChanged     func(mod contract.Moderator)
}

func NewModerator(config config.Game, reg *types.GameRegistration) contract.Moderator {
	m := &moderator{
		gameID:             reg.ID,
		gameStatus:         constants.Idle,
		config:             config,
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		mutex:              new(sync.Mutex),
		turnDuration:       reg.TurnDuration,
		discussionDuration: reg.DiscussionDuration,
		scheduler:          NewScheduler(constants.NightPhaseID),
	}
	m.world = NewWorld(m.scheduler, &types.GameInitialization{
		RoleIDs:          reg.RoleIDs,
		RequiredRoleIDs:  reg.RequiredRoleIDs,
		NumberWerewolves: reg.NumberWerewolves,
		PlayerIDs:        reg.PlayerIDs,
	})

	return m
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
	if len(m.world.AlivePlayerIDsWithFactionID(constants.WerewolfFactionID)) == 0 {
		// Villager wins if all werewolves are dead
		m.winningFaction = constants.VillagerFactionID
	} else if len(m.world.AlivePlayerIDsWithFactionID(constants.WerewolfFactionID)) >=
		len(m.world.AlivePlayerIDsWithoutFactionID(constants.WerewolfFactionID)) {
		// Werewolf wins if the number is overwhelming or equal to villager
		m.winningFaction = constants.WerewolfFactionID
	}
	m.mutex.Unlock()

	if !m.winningFaction.IsUnknown() {
		m.FinishGame()
	}
}

// handlePoll handles poll result of each faction.
func (m moderator) handlePoll(factionID types.FactionID) {
	if poll := m.world.Poll(factionID); poll != nil && poll.Close() {
		if record := poll.Record(constants.ZeroRound); !record.WinnerID.IsUnknown() {
			if player := m.world.Player(record.WinnerID); player != nil {
				player.Die(false)
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

		func() {
			var duration time.Duration

			if m.scheduler.PhaseID() == constants.DayPhaseID &&
				m.scheduler.TurnID() == constants.MidTurn {
				duration = m.discussionDuration

				m.world.Poll(constants.VillagerFactionID).Open() // nolint: errcheck
				defer m.handlePoll(constants.VillagerFactionID)
			} else {
				duration = m.turnDuration

				if m.scheduler.PhaseID() == constants.NightPhaseID &&
					m.scheduler.TurnID() == constants.MidTurn {
					m.world.Poll(constants.WerewolfFactionID).Open() // nolint: errcheck
					defer m.handlePoll(constants.WerewolfFactionID)
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
	req *types.ActivateAbilityRequest,
) *types.ActionResponse {
	if !m.mutex.TryLock() {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Turn is over!",
		}
	}
	defer m.mutex.Unlock()

	if slices.Contains(m.playedPlayerID, playerID) {
		return &types.ActionResponse{
			Ok:      false,
			Message: "You played this turn!",
		}
	}

	player := m.world.Player(playerID)
	if player == nil {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Non-existent player!",
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
		// 	fmt.Sprint(res.RoundID),
		// )
	}

	return res
}
