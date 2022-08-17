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

type GameData struct {
	Id                 GameId
	Capacity           uint
	NumberOfWerewolves uint
	TimeForTurn        time.Duration
	TimeForDiscussion  time.Duration
	RolePool           []RoleId
	SocketId2PlayerId  map[SocketId]PlayerId
}

type TurnData struct {
	PhaseId   PhaseId
	RoleId    RoleId
	PlayerIds []PlayerId
	Times     NumberOfTimes
	Position  TurnPosition
}
