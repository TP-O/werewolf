package data

import (
	"time"
	"uwwolf/game/types"
)

type WaitingRoom struct {
	ID                 string           `json:"id"`
	OwnerID            types.PlayerID   `json:"owner_id"`
	PlayerIDs          []types.PlayerID `json:"player_ids"`
	RoleIDs            []types.RoleID   `json:"role_ids"`
	RequiredRoleIDs    []types.RoleID   `json:"required_role_ids"`
	NumberWerewolves   uint8            `json:"number_werewolves"`
	TurnDuration       time.Duration    `json:"turn_duration"`
	DiscussionDuration time.Duration    `json:"discussion_duration"`
}
