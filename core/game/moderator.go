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

func NewModerator(init *types.ModeratorInit) (contract.Moderator, error) {
	m := &moderator{
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		mutex:              new(sync.Mutex),
		turnDuration:       init.TurnDuration,
		discussionDuration: init.DiscussionDuration,
		scheduler:          NewScheduler(vars.NightPhaseID),
	}

	game, err := NewGame(m.scheduler, &types.GameSetting{
		RoleIDs:          init.RoleIDs,
		RequiredRoleIDs:  init.RequiredRoleIDs,
		NumberWerewolves: init.NumberWerewolves,
		PlayerIDs:        init.PlayerIDs,
	})
	if err != nil {
		return nil, err
	} else {
		m.game = game
	}

	return m, nil
}

func (m *moderator) SetGameID(gameID types.GameID) bool {
	if m.game.StatusID() != vars.Idle {
		return false
	}

	m.gameID = gameID
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

				m.game.Poll(vars.VillagerFactionID).Open() // nolint: errcheck
				defer m.handlePoll(vars.VillagerFactionID)
			} else {
				duration = m.turnDuration

				if m.scheduler.PhaseID() == vars.NightPhaseID &&
					m.scheduler.TurnID() == vars.MidTurn {
					m.game.Poll(vars.WerewolfFactionID).Open() // nolint: errcheck
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
func (m *moderator) StartGame() int64 {
	if m.gameID.IsUnknown() || m.game.StatusID() != vars.Idle {
		return -1
	}

	fmt.Println("Starting")

	go func() {
		m.waitForPreparation()
		m.game.Start()
		go m.runScheduler()
	}()

	return m.game.Prepare()
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
