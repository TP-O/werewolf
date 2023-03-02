package game

import (
	"context"
	"sync"
	"time"
	"uwwolf/game/contract"
	"uwwolf/game/core"
	"uwwolf/game/role"
	"uwwolf/game/types"
	"uwwolf/util"

	"golang.org/x/exp/slices"
)

// gamemaster controlls a game.
type gamemaster struct {
	game               contract.Game
	scheduler          contract.Scheduler
	mutex              *sync.Mutex
	nextTurnSignal     chan bool
	finishSignal       chan bool
	turnDuration       time.Duration
	discussionDuration time.Duration
	playedPlayerIDs    []types.PlayerID
	notifier           contract.Notifier
}

var _ contract.Gamemaster = (*gamemaster)(nil)

func NewGamemaster(notifier contract.Notifier) contract.Gamemaster {
	return &gamemaster{
		scheduler:          core.NewScheduler(role.NightPhaseID),
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		turnDuration:       time.Duration(1) * time.Second,
		discussionDuration: time.Duration(2) * time.Second,
		notifier:           notifier,
	}
}

func (g *gamemaster) InitGame(newGame types.CreateGameRequest) bool {
	if g.game != nil {
		return false
	}

	g.game = core.NewGame(g.scheduler, types.GameSetting{
		GameID:           "xxx",
		RoleIDs:          newGame.RoleIDs,
		RequiredRoleIDs:  newGame.RequiredRoleIDs,
		NumberWerewolves: newGame.NumberWerewolves,
		PlayerIDs:        newGame.PlayerIDs,
	})
	return true
}

func (g gamemaster) checkWinner() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var winningFaction types.FactionID
	if len(g.game.PlayerIDsWithFactionID(role.WerewolfFactionID, true)) == 0 {
		winningFaction = role.VillagerFactionID
	} else if len(g.game.PlayerIDsWithFactionID(role.WerewolfFactionID, true)) >=
		len(g.game.PlayerIDsWithoutFactionID(role.WerewolfFactionID, true)) {
		winningFaction = role.WerewolfFactionID
	}

	if winningFaction.IsUnknown() {
		g.scheduler.NextTurn()
		g.notifier.NotifyGame(g.game.ID(), "Next turn!")
	} else {
		g.FinishGame()
		g.notifier.NotifyGame(g.game.ID(), "Winner is found!")
	}
}

func (g *gamemaster) runScheduler() {
	// Start the scheduler
	g.scheduler.NextTurn()

	for g.game.StatusID() == core.Starting {
		// Renew played play list
		g.playedPlayerIDs = make([]types.PlayerID, 0)

		func() {
			var duration time.Duration

			if g.scheduler.PhaseID() == role.DayPhaseID &&
				g.scheduler.TurnID() == role.MidTurn {
				duration = g.discussionDuration

				g.game.Poll(role.VillagerFactionID).Open()
				defer g.game.Poll(role.VillagerFactionID).Close()
			} else {
				duration = g.turnDuration

				if g.scheduler.PhaseID() == role.NightPhaseID &&
					g.scheduler.TurnID() == role.MidTurn {
					g.game.Poll(role.WerewolfFactionID).Open()
					defer g.game.Poll(role.WerewolfFactionID).Close()
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			select {
			case <-g.nextTurnSignal:
				g.checkWinner()
			case <-ctx.Done():
				g.checkWinner()
			case <-g.finishSignal:
				g.FinishGame()
				g.notifier.NotifyGame(g.game.ID(), "Game is over!")
			}
		}()
	}
}

func (g *gamemaster) waitForPreparation() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(util.Config().Game.PreparationDuration)*time.Second,
	)
	defer cancel()

	select {
	case <-g.finishSignal:
		g.game.Finish()
		g.notifier.NotifyGame(g.game.ID(), "Game is over before preparation!")
	case <-ctx.Done():
	}
}

func (g *gamemaster) StartGame() bool {
	if g.game.StatusID() != core.Idle {
		return false
	}

	g.nextTurnSignal = make(chan bool)
	g.finishSignal = make(chan bool)
	g.game.Prepare()
	g.waitForPreparation()
	g.game.Start()
	go g.runScheduler()
	g.notifier.NotifyGame(g.game.ID(), "Game is started!")

	return true
}

func (g gamemaster) FinishGame() bool {
	if g.game.StatusID() == core.Finished {
		return false
	}

	g.finishSignal <- true
	close(g.finishSignal)
	close(g.nextTurnSignal)
	g.game.Finish()

	return true
}

func (g *gamemaster) ReceivePlayRequest(
	playerID types.PlayerID,
	req types.ActivateAbilityRequest,
) types.ActionResponse {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if slices.Contains(g.playedPlayerIDs, playerID) {
		return types.ActionResponse{
			Ok:      false,
			Message: "You played this turn!",
		}
	}

	res := g.game.Play(playerID, req)
	if res.Ok {
		g.playedPlayerIDs = append(g.playedPlayerIDs, playerID)

		// Move to the next turn if all players have finished their job
		if len(g.playedPlayerIDs) == len(g.scheduler.PlayablePlayerIDs()) {
			g.nextTurnSignal <- true
		}
	}

	return res
}
