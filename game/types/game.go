package types

import (
	"time"
	"uwwolf/game/enum"
)

type GameSetting struct {
	NumberOfWerewolves uint8           `json:"number_of_werewolves" binding:"required"`
	TurnDuration       time.Duration   `json:"turn_duration" binding:"required"`
	DiscussionDuration time.Duration   `json:"discussion_duration" binding:"required"`
	RoleIDs            []enum.RoleID   `json:"role_ids" binding:"required"`
	RequiredRoleIDs    []enum.RoleID   `json:"required_role_ids" binding:"required"`
	PlayerIDs          []enum.PlayerID `json:"player_ids" binding:"required"`
}
