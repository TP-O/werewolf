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
	GameID             types.GameID
	Scheduler          contract.Scheduler
	TurnDuration       time.Duration
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

func (m *moderator) InitGame(req types.CreateGameRequest) bool {
	if m.game != nil {
		return false
	}

	// if game, err := m.db.CreateGame(context.Background()); err != nil {
	// 	return false
	// } else {
	m.game = NewGame(m.scheduler, &types.GameSetting{
		RoleIDs:          req.RoleIDs,
		RequiredRoleIDs:  req.RequiredRoleIDs,
		NumberWerewolves: req.NumberWerewolves,
		PlayerIDs:        req.PlayerIDs,
	})

	return true
}

func (m *moderator) checkWinConditions() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.game.AlivePlayerIDsWithFactionID(vars.WerewolfFactionID)) == 0 {
		m.winningFaction = vars.VillagerFactionID
	} else if len(m.game.AlivePlayerIDsWithFactionID(vars.WerewolfFactionID)) >=
		len(m.game.AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID)) {
		m.winningFaction = vars.WerewolfFactionID
	}

	if !m.winningFaction.IsUnknown() {
		m.FinishGame()
	}
}

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
				defer func() {
					m.game.Poll(vars.VillagerFactionID).Close()
					if record := m.game.Poll(vars.VillagerFactionID).Record(vars.ZeroRound); record != nil && !record.WinnerID.IsUnknown() {
						m.game.KillPlayer(record.WinnerID, false)
					}
				}()
			} else {
				duration = m.turnDuration

				if m.scheduler.PhaseID() == vars.NightPhaseID &&
					m.scheduler.TurnID() == vars.MidTurn {
					m.game.Poll(vars.WerewolfFactionID).Open()
					defer func() {
						m.game.Poll(vars.WerewolfFactionID).Close()
						if record := m.game.Poll(vars.WerewolfFactionID).Record(vars.ZeroRound); record != nil && !record.WinnerID.IsUnknown() {
							m.game.KillPlayer(record.WinnerID, false)
						}
					}()
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

func (m *moderator) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		util.Config().Game.PreparationDuration,
	)
	defer cancel()

	select {
	case <-m.finishSignal:
		m.game.Finish()
	case <-ctx.Done():
	}

	fmt.Println("Preparation is done")
}

func (m *moderator) StartGame() bool {
	if m.game.StatusID() != vars.Idle || m.game.Prepare() == -1 {
		return false
	}

	fmt.Println("Starting")

	m.nextTurnSignal = make(chan bool)
	m.finishSignal = make(chan bool)
	m.waitForPreparation()
	m.game.Start()
	go m.runScheduler()

	fmt.Println(m.game.AlivePlayerIDsWithFactionID(vars.WerewolfFactionID))
	fmt.Println(m.game.AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID))

	return true
}

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
