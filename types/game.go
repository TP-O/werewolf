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
	NextTurnPosition TurnPosition = iota - 2
	LastTurnPosition
	FirstTurnPosition
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
	PhaseId   PhaseId
	RoleId    RoleId
	PlayerIds []PlayerId
	Times     NumberOfTimes
	Position  TurnPosition
}
