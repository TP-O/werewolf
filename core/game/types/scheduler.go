package types

type TurnID uint8

type RoundID uint8

type TurnSlot struct {
	BeginRoundID  RoundID
	PlayedRoundID RoundID
	FrozenTimes   Times
	RoleID
}

type Turn map[PlayerID]*TurnSlot

type NewTurnSlot struct {
	PhaseID
	TurnID
	BeginRoundID  RoundID
	PlayedRoundID RoundID
	FrozenTimes   Times
	PlayerID
	RoleID
}

type RemovedTurnSlot struct {
	PhaseID
	TurnID
	PlayerID
	RoleID
}

type FreezeTurnSlot struct {
	PhaseID
	TurnID
	PlayerID
	RoleID
	FrozenTimes Times
}
