package dto

import (
	"time"
	"uwwolf/game/types"
)

type StartGameDto struct {
	// TurnDuration is the duration of a turn.
	TurnDuration time.Duration `json:"turn_duration"`

	// DiscussionDuration is the duration of the villager discussion.
	DiscussionDuration time.Duration `json:"discussion_duration"`

	// RoleIDs is role ID list that can be played in the game.
	RoleIDs []types.RoleID `json:"role_ids"`

	// RequiredRoleIDs is role ID list that must be played in the game.
	RequiredRoleIDs []types.RoleID `json:"required_role_ids"`

	// NumberWerewolves is number of werewolves required to exist in the game.
	NumberWerewolves uint8 `json:"number_werewolves"`
}

type Room struct {
	ID        string           `json:"id"`
	OwnerID   types.PlayerID   `json:"owner_id"`
	PlayerIDs []types.PlayerID `json:"player_ids"`
	types.ModeratorInit
}
