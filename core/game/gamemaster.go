package game

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
	"uwwolf/db"
	"uwwolf/game/contract"
	"uwwolf/game/core"
	"uwwolf/game/role"
	"uwwolf/game/types"
	"uwwolf/util"

	"github.com/redis/go-redis/v9"
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
	rdb                *redis.Client
	db                 db.Querier
	winningFaction     types.FactionID
}

var _ contract.Gamemaster = (*gamemaster)(nil)

func NewGamemaster(rdb *redis.Client, db db.Querier, notifier contract.Notifier) contract.Gamemaster {
	return &gamemaster{
		scheduler:          core.NewScheduler(role.NightPhaseID),
		nextTurnSignal:     make(chan bool),
		finishSignal:       make(chan bool),
		turnDuration:       time.Duration(1) * time.Second,
		discussionDuration: time.Duration(2) * time.Second,
		notifier:           notifier,
		rdb:                rdb,
		db:                 db,
	}
}

func (g *gamemaster) InitGame(req types.CreateGameRequest) bool {
	if g.game != nil {
		return false
	}

	if gameRow, err := g.db.CreateGame(context.Background()); err != nil {
		return false
	} else {
		g.game = core.NewGame(g.scheduler, types.GameSetting{
			GameID:           uint64(gameRow.ID),
			RoleIDs:          req.RoleIDs,
			RequiredRoleIDs:  req.RequiredRoleIDs,
			NumberWerewolves: req.NumberWerewolves,
			PlayerIDs:        req.PlayerIDs,
		})
	}

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

	// Store role assignments
	assignments := db.AssignGameRolesParams{}

	for _, player := range g.game.Players() {
		assignments.GameIds = append(assignments.GameIds, fmt.Sprint(g.game.ID()))
		assignments.FactionIds = append(assignments.FactionIds, string(player.FactionID()))
		assignments.PlayerIds = append(assignments.PlayerIds, string(player.ID()))
		assignments.RoleIds = append(assignments.RoleIds, string(player.MainRoleID()))
	}

	if g.db.AssignGameRoles(context.Background(), assignments) != nil {
		return false
	}

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

	g.db.FinishGame(context.Background(), db.FinishGameParams{
		ID: int64(g.game.ID()),
		WinningFactionID: sql.NullInt16{
			Int16: int16(g.winningFaction),
			Valid: true,
		},
	})

	logs := db.CreateGameLogsParams{}

	result := g.rdb.LRange(context.Background(), fmt.Sprint(g.game.ID()), 0, -1)
	for i := len(result.Val()) / 5; i >= 4; i -= 5 {
		logs.GameIds = append(logs.GameIds, fmt.Sprint(g.game.ID()))
		logs.RoundIds = append(logs.RoundIds, result.Val()[i])
		logs.ActorIds = append(logs.RoundIds, result.Val()[i-1])
		logs.RoleIds = append(logs.RoundIds, result.Val()[i-2])
		logs.ActionIds = append(logs.RoundIds, result.Val()[i-3])
		logs.TargetIds = append(logs.RoundIds, result.Val()[i-4])
	}

	g.db.CreateGameLogs(context.Background(), logs)

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

		// Record player turn
		g.rdb.LPush(
			context.Background(),
			fmt.Sprint(g.game.ID()),
			res.RoundID,
			playerID,
			res.RoleID,
			res.ActionID,
			res.TargetID,
		)
	}

	return res
}
