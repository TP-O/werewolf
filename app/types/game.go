package types

import "time"

type GameId uint

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
	Id                 GameId        `json:"id" validate:"required,number,gt=0"`
	NumberOfWerewolves int           `json:"numberOfWerewolves" validate:"required,number,number_of_werewolves=PlayerIds"`
	TurnDuration       time.Duration `json:"turnDuration" validate:"required,number,gt=10"`
	DiscussionDuration time.Duration `json:"discussionDuration" validate:"required,number,gt=10"`
	RolePool           []RoleId      `json:"rolePool" validate:"required,min=1,unique,role_pool"`
	PlayerIds          []PlayerId    `json:"playerIds" validate:"required,min=1,unique,capacity,dive,len=28"`
}

type TurnSetting struct {
	PhaseId    PhaseId
	RoleId     RoleId
	PlayerIds  []PlayerId
	BeginRound RoundId
	Priority   int
	Expiration NumberOfTimes
	Position   TurnPosition
}
