-- name: CreateGame :one
INSERT INTO games DEFAULT VALUES
RETURNING *;

-- name: FinishGame :exec
UPDATE games SET
    winning_faction_id = $2,
    finished_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: AssignGameRoles :exec
INSERT INTO role_assignments (
    game_id,
    role_id,
    faction_id,
    player_id
) VALUES (
    unnest(@game_id::bigint[]),
    unnest(@role_id::smallint[]),
    unnest(@faction_id::smallint[]),
    unnest(@player_id::varchar[])
);

-- name: CreateGameLogs :exec
INSERT INTO game_logs (
    game_id,
    round_id,
    actor_id,
    role_id,
    action_id,
    target_id
) VALUES (
    unnest(@game_id::bigint[]),
    unnest(@round_id::smallint[]),
    unnest(@actor_id::varchar[]),
    unnest(@role_id::smallint[]),
    unnest(@action_id::smallint[]),
    unnest(@target_id::varchar[])
);

-- name: PlayingGame :one
SELECT games.* FROM role_assignments
INNER JOIN games ON role_assignments.game_id == games.id
WHERE role_assignments.player_id == $1 AND games.finished_at IS NULL;
