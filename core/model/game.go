package model

import "uwwolf/game/types"

type GameTurn struct {
	ActorID   types.PlayerID   `cql:"actor_id" json:"actorID"`
	ActionID  types.ActionID   `cql:"action_id" json:"actionID"`
	TargetIDs []types.PlayerID `cql:"target_ids" json:"targetIDs"`
}

type GameRound map[types.PhaseID][][]GameTurn

type Game struct {
	ID               types.GameID    `cql:"id" json:"id"`
	WinningFactionID types.FactionID `cql:"winning_faction_id" json:"winningFactionID"`
	Record           []GameRound     `cql:"record" json:"record"`
	StartedAt        int64           `cql:"started_at" json:"startedAt"`
	FinishedAt       int64           `cql:"finished_at" json:"finishedAt"`
}
