package types

import (
	"time"
	"uwwolf/game/enum"
)

type GameSetting struct {
	NumberOfWerewolves int
	TurnDuration       time.Duration
	DiscussionDuration time.Duration
	RoleIDs            []enum.RoleID
	RequiredRoleIDs    []enum.RoleID
	PlayerIDs          []enum.PlayerID
}
