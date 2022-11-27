package types

import "time"

type GameID string

func (g GameID) IsValid() bool {
	return len(g) == 28
}

type PhaseID uint

func (p PhaseID) IsUnknown() bool {
	return p == 0
}

type GameStatus uint

type GameSetting struct {
	NumberOfWerewolves int           `json:"numberOfWerewolves" validate:"required,number,number_of_werewolves=PlayerIds"`
	TurnDuration       time.Duration `json:"turnDuration" validate:"required,number,gt=10"`
	DiscussionDuration time.Duration `json:"discussionDuration" validate:"required,number,gt=10"`
	RoleIDs            []RoleID      `json:"werewolfRoleIDs" validate:"required,unique,role_id=w"`
	RequiredRoleIDs    []RoleID
	PlayerIDs          []PlayerID `json:"playerIDs" validate:"required,min=3,unique,capacity,dive,len=28"`
}
