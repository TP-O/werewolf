CREATE TABLE games (
    id uuid PRIMARY KEY,
    winning_faction_id smallint,
    record list<frozen<set<game_turn>>>,
    started_at timestamp,
    finished_at timestamp,
);
