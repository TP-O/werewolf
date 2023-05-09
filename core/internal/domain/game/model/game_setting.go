package model

import "uwwolf/internal/app/game/logic/types"

type GameSetting struct {
	RoleIDs            []types.RoleID `json:"role_ids"`
	RequiredRoleIDs    []types.RoleID `json:"required_role_ids"`
	NumberWerewolves   uint8          `json:"number_werewolves"`
	TurnDuration       uint16         `json:"turn_duration"`
	DiscussionDuration uint16         `json:"discussion_duration"`
}
