package types

import "time"

type GameId string

type RoundId uint

const (
	FirstRound RoundId = 1
)

type PhaseId uint

const (
	UnknownPhase PhaseId = iota
	NightPhase
	DayPhase
	DuskPhase
)

type TurnPosition int

const (
	NextPosition TurnPosition = iota - 3
	SortedPosition
	LastPosition
	FirstPosition
)

type GameSetting struct {
	Id                 GameId
	NumberOfWerewolves uint
	TimeForTurn        time.Duration
	TimeForDiscussion  time.Duration
	RolePool           []RoleId
	PlayerIds          []PlayerId
}

type TurnSetting struct {
	PhaseId    PhaseId
	RoleId     RoleId
	PlayerIds  []PlayerId
	BeginRound RoundId
	Priority   int
	Times      NumberOfTimes
	Position   TurnPosition
}
