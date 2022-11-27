package types

type Round int

func (r Round) IsStarted() bool {
	return r != 0
}

type Position int

type TurnSetting struct {
	PhaseID    PhaseID
	RoleID     RoleID
	BeginRound Round
	Priority   Priority
	Position   Position
}

type Turn struct {
	RoleID      RoleID
	BeginRound  Round
	Priority    Priority
	FrozenLimit Limit
}
