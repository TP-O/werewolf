CREATE TYPE IF NOT EXISTS game_turn (
    actor_id varchar,
    action_id smallint,
    target_ids set<varchar>,
);
