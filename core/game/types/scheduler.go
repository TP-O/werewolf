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
	ExpiredAfter Limit
}

type RemovedPlayerTurn struct {
	PhaseID
	TurnID
	PlayerID
	RoleID
}
