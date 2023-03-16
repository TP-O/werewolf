package game

import (
	"context"
	"fmt"
	"sync"
	"time"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	"uwwolf/util"

	"golang.org/x/exp/slices"
)

type ModeratorInit struct {
	GameID    types.GameID
	Scheduler contract.Scheduler

	// TurnDuration is the duration of a turn.
	TurnDuration time.Duration

	// DiscussionDuration is the duration of the villager discussion.
	DiscussionDuration time.Duration
}

// moderator controlls a game.
type moderator struct {
	gameID             types.GameID
	game               contract.Game
	scheduler          contract.Scheduler
	mutex              *sync.Mutex
	nextTurnSignal     chan bool
	finishSignal       chan bool
	turnDuration       time.Duration
	discussionDuration time.Duration
	playedPlayerID     []types.PlayerID
	winningFaction     types.FactionID
}

func NewModerator(init *ModeratorInit) contract.Moderator {
	return &moderator{
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		mutex:              new(sync.Mutex),
		gameID:             init.GameID,
		turnDuration:       init.TurnDuration,
		discussionDuration: init.DiscussionDuration,
		scheduler:          init.Scheduler,
	}
}

// InitGame creates a new idle game instance.
func (m *moderator) InitGame(setting *types.GameSetting) bool {
	if m.game != nil {
		return false
	}

	// if game, err := m.db.CreateGame(context.Background()); err != nil {
	// 	return false
	// } else {
	m.game = NewGame(m.scheduler, setting)

	return true
}

// checkWinConditions checks if any faction satisfies its win condition,
// if any, finish the game.
func (m *moderator) checkWinConditions() {
	m.mutex.Lock()
	if len(m.game.AlivePlayerIDsWithFactionID(vars.WerewolfFactionID)) == 0 {
		// Villager wins if all werewolves are dead
		m.winningFaction = vars.VillagerFactionID
	} else if len(m.game.AlivePlayerIDsWithFactionID(vars.WerewolfFactionID)) >=
		len(m.game.AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID)) {
		// Werewolf wins if the number is overwhelming or equal to villager
		m.winningFaction = vars.WerewolfFactionID
	}
	m.mutex.Unlock()

	if !m.winningFaction.IsUnknown() {
		m.FinishGame()
	}
}

// handlePoll handles poll result of each faction.
func (m moderator) handlePoll(factionID types.FactionID) {
	if poll := m.game.Poll(factionID); poll != nil && poll.Close() {
		if record := poll.Record(vars.ZeroRound); record != nil &&
			!record.WinnerID.IsUnknown() {
			m.game.KillPlayer(record.WinnerID, false)
		}
	}
}

// runScheduler switches turns automatically.
func (m *moderator) runScheduler() {
	for m.game.StatusID() == vars.Starting {
		m.mutex.Lock()
		m.playedPlayerID = make([]types.PlayerID, 0)
		m.scheduler.NextTurn()

		func() {
			var duration time.Duration

			if m.scheduler.PhaseID() == vars.DayPhaseID &&
				m.scheduler.TurnID() == vars.MidTurn {
				duration = m.discussionDuration

				m.game.Poll(vars.VillagerFactionID).Open()
				defer m.handlePoll(vars.VillagerFactionID)
			} else {
				duration = m.turnDuration

				if m.scheduler.PhaseID() == vars.NightPhaseID &&
					m.scheduler.TurnID() == vars.MidTurn {
					m.game.Poll(vars.WerewolfFactionID).Open()
					defer m.handlePoll(vars.WerewolfFactionID)
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
		}()
	}
}

// waitForPreparation waits a bit before the game starts.
func (m *moderator) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		util.Config().Game.PreparationDuration,
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
func (m *moderator) StartGame() bool {
	if m.game.StatusID() != vars.Idle || m.game.Prepare() == -1 {
		return false
	}

	fmt.Println("Starting")

	m.waitForPreparation()
	m.game.Start()
	go m.runScheduler()

	return true
}

// FinishGame ends the game.
func (m *moderator) FinishGame() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.game.StatusID() == vars.Finished {
		return false
	}

	// Notify winner
	fmt.Println("Winner is %w", m.winningFaction)

	m.finishSignal <- true
	close(m.finishSignal)
	close(m.nextTurnSignal)
	m.game.Finish()

	return true
}

// RequestPlay receives the play request from the player.
func (m *moderator) RequestPlay(
	playerID types.PlayerID,
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

	res := m.game.Play(playerID, req)
	if res.Ok {
		m.playedPlayerID = append(m.playedPlayerID, playerID)

		// Move to the next turn if all players have finished their job
		// if len(m.playedPlayerID) == len(m.scheduler.PlayablePlayerID()) {
		// 	m.nextTurnSignal <- true
		// }

		// Cache player request
		// m.rdb.LPush(
		// 	context.Background(),
		// 	fmt.Sprint(m.gameID),
		// 	fmt.Sprint(res.TargetID),
		// 	fmt.Sprint(res.ActionID),
		// 	fmt.Sprint(res.RoleID),
		// 	fmt.Sprint(playerID),
		// 	fmt.Sprint(res.RoundID),
		// )
	}

	return res
}
