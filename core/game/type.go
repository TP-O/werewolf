package game

type GameID = string

func IsValidGameID(g GameID) bool {
	return len(g) == 28
}

type GameStatus = uint8

type GameSetting struct {
	TurnDuration       uint16     `json:"turn_duration" validate:"required"`
	DiscussionDuration uint16     `json:"discussion_duration" validate:"required"`
	RoleIDs            []RoleID   `json:"role_ids" validate:"required,min=2,unique,dive"`
	RequiredRoleIDs    []RoleID   `json:"required_role_ids" validate:"omitempty,ltecsfield=RoleIDs,unique,dive"`
	NumberWerewolves   uint8      `json:"number_werewolves" validate:"required,gt=0"`
	PlayerIDs          []PlayerID `json:"player_ids" validate:"required,unique,dive,len=20"`
}
