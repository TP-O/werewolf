package types

import (
	"time"
)

type GameId uint

type GameSetting struct {
	NumberOfWerewolves int           `json:"numberOfWerewolves" validate:"required,number,number_of_werewolves=PlayerIds"`
	TurnDuration       time.Duration `json:"turnDuration" validate:"required,number,gt=10"`
	DiscussionDuration time.Duration `json:"discussionDuration" validate:"required,number,gt=10"`
	WerewolfRoleIds    []RoleId      `json:"werewolfRoleIds" validate:"required,unique,role_id=w"`
	NonWerewolfRoleIds []RoleId      `json:"nonWerewolfRoleIds" validate:"required,min=2,unique,role_id"`
	PlayerIds          []PlayerId    `json:"playerIds" validate:"required,min=3,unique,capacity,dive,len=28"`
}
