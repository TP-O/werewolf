package types

type PhaseID uint8

func (p PhaseID) IsUnknown() bool {
	return p == 0
}

func (p PhaseID) NextPhasePhaseID(lastPhaseID PhaseID) PhaseID {
	if p == lastPhaseID {
		return 1
	}
	return p + 1
}

func (p PhaseID) PreviousPhaseID(lastPhaseID PhaseID) PhaseID {
	if p == 1 {
		return lastPhaseID
	}
	return p - 1
}

type GameID = uint64

type GameStatusID = uint8

type GameSetting struct {
	GameID           GameID     `json:"game_id"`
	RoleIDs          []RoleID   `json:"role_ids"`
	RequiredRoleIDs  []RoleID   `json:"required_role_ids"`
	NumberWerewolves uint8      `json:"number_werewolves"`
	PlayerIDs        []PlayerID `json:"player_ids"`
}

type CreateGameRequest struct {
	TurnDuration       uint16     `json:"turn_duration" validate:"required"`
	DiscussionDuration uint16     `json:"discussion_duration" validate:"required"`
	RoleIDs            []RoleID   `json:"role_ids" validate:"required,min=2,unique,dive"`
	RequiredRoleIDs    []RoleID   `json:"required_role_ids" validate:"omitempty,ltecsfield=RoleIDs,unique,dive"`
	NumberWerewolves   uint8      `json:"number_werewolves" validate:"required,gt=0"`
	PlayerIDs          []PlayerID `json:"player_ids" validate:"required,unique,dive,len=20"`
}
