package types

type PlayerTurn struct {
	BeginRoundID RoundID
	FrozenLimit  Limit
	RoleID
}

type Turn map[PlayerID]*PlayerTurn

type NewPlayerTurn struct {
	PhaseID
	TurnID
	BeginRoundID RoundID
	PlayerID
	RoleID
}

type RemovedPlayerTurn struct {
	PhaseID
	TurnID
	PlayerID
	RoleID
}
