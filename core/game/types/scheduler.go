package types

type TurnSetting struct {
	PhaseID
	TurnID
	BeginRoundID RoundID
	PlayerIDs    []PlayerID
}

type Turn struct {
	ID           TurnID
	RoleID       RoleID
	BeginRoundID RoundID
	FrozenLimit  Limit
}
