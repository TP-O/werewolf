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
    unnest(@game_ids::text[]),
    unnest(@role_ids::text[]),
    unnest(@faction_ids::text[]),
    unnest(@player_ids::text[])
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
    unnest(@game_ids::text[]),
    unnest(@round_ids::text[]),
    unnest(@actor_ids::text[]),
    unnest(@role_ids::text[]),
    unnest(@action_ids::text[]),
    unnest(@target_ids::text[])
);
