package dto

import "uwwolf/internal/app/game/logic/types"

type UpdateGameSetting struct {
	RoleIDs            []types.RoleId `json:"role_ids" binding:"required,min=1,unique,dive,gt=2"`
	RequiredRoleIDs    []types.RoleId `json:"required_role_ids" binding:"omitempty,unique"`
	NumberWerewolves   uint8          `json:"number_werewolves" binding:"required,number,gt=0,lt=8"`
	TurnDuration       uint16         `json:"turn_duration" binding:""`
	DiscussionDuration uint16         `json:"discussion_duration" binding:""`
}
