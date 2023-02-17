package types

import "uwwolf/game/enum"

type GameSetting struct {
	TurnDuration       uint16          `json:"turn_duration" validate:"required"`
	DiscussionDuration uint16          `json:"discussion_duration" validate:"required"`
	RoleIDs            []enum.RoleID   `json:"role_idssss" validate:"required,min=2,unique,dive"`
	RequiredRoleIDs    []enum.RoleID   `json:"required_role_ids" validate:"omitempty,ltecsfield=RoleIDs,unique,dive"`
	NumberWerewolves   uint8           `json:"number_werewolves" validate:"required,gt=0"`
	PlayerIDs          []enum.PlayerID `json:"player_ids" validate:"required,unique,dive,len=20"`
}
