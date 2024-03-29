// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: game.sql

package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const assignGameRoles = `-- name: AssignGameRoles :exec
INSERT INTO role_assignments (
    game_id,
    role_id,
    faction_id,
    player_id
) VALUES (
    unnest($1::bigint[]),
    unnest($2::smallint[]),
    unnest($3::smallint[]),
    unnest($4::varchar[])
)
`

type AssignGameRolesParams struct {
	GameID    []int64  `json:"game_id"`
	RoleID    []int16  `json:"role_id"`
	FactionID []int16  `json:"faction_id"`
	PlayerID  []string `json:"player_id"`
}

func (q *Queries) AssignGameRoles(ctx context.Context, arg AssignGameRolesParams) error {
	_, err := q.db.ExecContext(ctx, assignGameRoles,
		pq.Array(arg.GameID),
		pq.Array(arg.RoleID),
		pq.Array(arg.FactionID),
		pq.Array(arg.PlayerID),
	)
	return err
}

const createGame = `-- name: CreateGame :one
INSERT INTO games DEFAULT VALUES
RETURNING id, winning_faction_id, created_at, finished_at
`

func (q *Queries) CreateGame(ctx context.Context) (Game, error) {
	row := q.db.QueryRowContext(ctx, createGame)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.WinningFactionID,
		&i.CreatedAt,
		&i.FinishedAt,
	)
	return i, err
}

const createGameLogs = `-- name: CreateGameLogs :exec
INSERT INTO game_logs (
    game_id,
    round_id,
    actor_id,
    role_id,
    action_id,
    target_id
) VALUES (
    unnest($1::bigint[]),
    unnest($2::smallint[]),
    unnest($3::varchar[]),
    unnest($4::smallint[]),
    unnest($5::smallint[]),
    unnest($6::varchar[])
)
`

type CreateGameLogsParams struct {
	GameID   []int64  `json:"game_id"`
	RoundID  []int16  `json:"round_id"`
	ActorID  []string `json:"actor_id"`
	RoleID   []int16  `json:"role_id"`
	ActionID []int16  `json:"action_id"`
	TargetID []string `json:"target_id"`
}

func (q *Queries) CreateGameLogs(ctx context.Context, arg CreateGameLogsParams) error {
	_, err := q.db.ExecContext(ctx, createGameLogs,
		pq.Array(arg.GameID),
		pq.Array(arg.RoundID),
		pq.Array(arg.ActorID),
		pq.Array(arg.RoleID),
		pq.Array(arg.ActionID),
		pq.Array(arg.TargetID),
	)
	return err
}

const finishGame = `-- name: FinishGame :exec
UPDATE games SET
    winning_faction_id = $2,
    finished_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type FinishGameParams struct {
	ID               int64         `json:"id"`
	WinningFactionID sql.NullInt16 `json:"winning_faction_id"`
}

func (q *Queries) FinishGame(ctx context.Context, arg FinishGameParams) error {
	_, err := q.db.ExecContext(ctx, finishGame, arg.ID, arg.WinningFactionID)
	return err
}

const playingGame = `-- name: PlayingGame :one
SELECT games.id, games.winning_faction_id, games.created_at, games.finished_at FROM role_assignments
INNER JOIN games ON role_assignments.game_id == games.id
WHERE role_assignments.player_id == $1 AND games.finished_at IS NULL
`

func (q *Queries) PlayingGame(ctx context.Context, playerID string) (Game, error) {
	row := q.db.QueryRowContext(ctx, playingGame, playerID)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.WinningFactionID,
		&i.CreatedAt,
		&i.FinishedAt,
	)
	return i, err
}
