package types

type RoundId uint

type TurnPosition int

type TurnSetting struct {
	PhaseId    PhaseId
	RoleId     RoleId
	PlayerIds  []PlayerId
	BeginRound RoundId
	Priority   int
	Expiration Expiration
	Position   TurnPosition
}

type Turn struct {
	RoleId     RoleId
	PlayerIds  []PlayerId
	BeginRound RoundId
	Priority   int
	Expiration Expiration
}

type Phase []*Turn
