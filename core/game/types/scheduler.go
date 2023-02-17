package types

import "uwwolf/game/enum"

type TurnSetting struct {
	PhaseID    enum.PhaseID
	RoleID     enum.RoleID
	BeginRound enum.Round
	Priority   enum.Priority
	Position   enum.Position
	Limit      enum.Limit
}

type Turn struct {
	RoleID      enum.RoleID
	BeginRound  enum.Round
	Priority    enum.Priority
	Limit       enum.Limit
	FrozenLimit enum.Limit
}
