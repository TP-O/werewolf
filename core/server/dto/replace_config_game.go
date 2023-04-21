package dto

import "uwwolf/game/types"

type ReplaceGameConfigDto struct {
	RoleIDs            []types.RoleID `json:"role_ids" binding:"required,min=1,unique,dive,gt=2"`
	RequiredRoleIDs    []types.RoleID `json:"required_role_ids" binding:"omitempty,unique"`
	NumberWerewolves   uint8          `json:"number_werewolves" binding:"required,number,gt=0,lt=8"`
	TurnDuration       uint16         `json:"turn_duration" binding:""`
	DiscussionDuration uint16         `json:"discussion_duration" binding:""`
}
