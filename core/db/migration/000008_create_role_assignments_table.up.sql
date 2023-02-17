CREATE TABLE role_assignments (
    game_id varchar,
    player_id varchar,
    role_id smallint,
    is_leader boolean,
    is_dead boolean,
    PRIMARY KEY (game_id, player_id),
);
