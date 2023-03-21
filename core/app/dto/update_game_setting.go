package dto

import (
	"time"
	"uwwolf/game/types"
)

type UpdateGameSettingDto struct {
	RoleIDs            []types.RoleID `json:"role_ids"`
	RequiredRoleIDs    []types.RoleID `json:"required_role_ids"`
	NumberWerewolves   uint8          `json:"number_werewolves"`
	TurnDuration       time.Duration  `json:"turn_duration"`
	DiscussionDuration time.Duration  `json:"discussion_duration"`
}
