package model

import (
	"uwwolf/game/enum"
)

type GameTurn struct {
	ActorID   enum.PlayerID   `cql:"actor_id" json:"actorID"`
	ActionID  enum.ActionID   `cql:"action_id" json:"actionID"`
	TargetIDs []enum.PlayerID `cql:"target_ids" json:"targetIDs"`
}

type GameRound map[enum.PhaseID][][]GameTurn

type Game struct {
	ID               enum.GameID    `cql:"id" json:"id"`
	WinningFactionID enum.FactionID `cql:"winning_faction_id" json:"winningFactionID"`
	Record           []GameRound    `cql:"record" json:"record"`
	StartedAt        int64          `cql:"started_at" json:"startedAt"`
	FinishedAt       int64          `cql:"finished_at" json:"finishedAt"`
}
