package types

import (
	"time"
)

type GameId uint

type GameSetting struct {
	Id                 GameId        `json:"id" validate:"required,number,gt=0"`
	NumberOfWerewolves int           `json:"numberOfWerewolves" validate:"required,number,number_of_werewolves=PlayerIds"`
	TurnDuration       time.Duration `json:"turnDuration" validate:"required,number,gt=10"`
	DiscussionDuration time.Duration `json:"discussionDuration" validate:"required,number,gt=10"`
	WerewolfRoleIds    []RoleId      `json:"werewolfRoleIds" validate:"required,min=1,unique,role_pool"`
	NonWerewolfRoleIds []RoleId      `json:"nonWerewolfRoleIds" validate:"required,min=1,unique,role_pool"`
	PlayerIds          []PlayerId    `json:"playerIds" validate:"required,min=1,unique,capacity,dive,len=28"`
}
