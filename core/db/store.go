package db

import (
	"context"
	"database/sql"
	"strconv"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

type Store struct {
	*Queries
	db *sql.DB
}

var store *Store

func DB() *Store {
	return store
}

func NewStore(db *sql.DB) *Store {
	if store == nil {
		store = &Store{
			Queries: New(db),
			db:      db,
		}
	}
	return store
}

func (s *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	var err error

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

type StoreGameParams struct {
	GameID           types.GameID
	WinningFactionID types.FactionID
	Players          []contract.Player
	Records          []string
}

func (s *Store) StoreGame(ctx context.Context, params *StoreGameParams) error {
	return s.execTx(ctx, func(q *Queries) error {
		var err error

		if err = q.FinishGame(context.Background(), FinishGameParams{
			ID: int64(params.GameID),
			WinningFactionID: sql.NullInt16{
				Int16: int16(params.WinningFactionID),
				Valid: true,
			},
		}); err != nil {
			return err
		}

		// Store role assignments
		assignments := AssignGameRolesParams{}
		for _, player := range params.Players {
			assignments.GameID = append(assignments.GameID, int64(params.GameID))
			assignments.FactionID = append(assignments.FactionID, int16(player.FactionID()))
			assignments.PlayerID = append(assignments.PlayerID, string(player.ID()))
			assignments.RoleID = append(assignments.RoleID, int16(player.MainRoleID()))
		}
		if err = q.AssignGameRoles(context.Background(), assignments); err != nil {
			return err
		}

		logs := CreateGameLogsParams{}
		for i := len(params.Records) - 1; i >= 4; i -= 5 {
			roundID, _ := strconv.Atoi(params.Records[i])
			roleID, _ := strconv.Atoi(params.Records[i-2])
			actionID, _ := strconv.Atoi(params.Records[i-3])

			logs.GameID = append(logs.GameID, int64(params.GameID))
			logs.RoundID = append(logs.RoundID, int16(roundID))
			logs.ActorID = append(logs.ActorID, params.Records[i-1])
			logs.RoleID = append(logs.RoleID, int16(roleID))
			logs.ActionID = append(logs.ActionID, int16(actionID))
			logs.TargetID = append(logs.TargetID, params.Records[i-4])
		}
		return q.CreateGameLogs(context.Background(), logs)
	})
}
