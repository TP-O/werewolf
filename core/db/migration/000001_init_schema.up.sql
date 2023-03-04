CREATE TABLE players (
    id VARCHAR(128) PRIMARY KEY
);

CREATE TABLE actions (
    id SMALLSERIAL PRIMARY KEY
);

CREATE TABLE roles (
    id SMALLSERIAL PRIMARY KEY
);

CREATE TABLE factions (
    id SMALLSERIAL PRIMARY KEY
);

CREATE TABLE games (
    id BIGSERIAL PRIMARY KEY,
    winning_faction_id SMALLINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,

    FOREIGN KEY (winning_faction_id) REFERENCES factions(id)
);

CREATE TABLE role_assignments (
    game_id BIGINT,
    role_id SMALLINT,
    faction_id SMALLINT,
    player_id VARCHAR(128),

    PRIMARY KEY (game_id, player_id),
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (faction_id) REFERENCES factions(id),
    FOREIGN KEY (player_id) REFERENCES players(id)
);

CREATE TABLE game_logs (
    game_id BIGINT,
    round_id SMALLINT NOT NULL,
    actor_id VARCHAR(128) NOT NULL,
    role_id SMALLINT NOT NULL,
    action_id SMALLINT NOT NULL,
    target_id VARCHAR(128) NOT NULL,

    PRIMARY KEY (game_id, round_id, actor_id, role_id),
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (actor_id) REFERENCES players(id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (action_id) REFERENCES actions(id),
    FOREIGN KEY (target_id) REFERENCES players(id)
)
