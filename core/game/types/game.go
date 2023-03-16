package types

// PhaseID is ID type of phase.
type PhaseID uint8

// IsUnknown checks if phase ID is 0.
func (p PhaseID) IsUnknown() bool {
	return p == 0
}

// NextPhasePhaseID returns the next phase ID based on the given lastPhaseID.
func (p PhaseID) NextPhasePhaseID(lastPhaseID PhaseID) PhaseID {
	if p == lastPhaseID {
		return 1
	}
	return p + 1
}

// PreviousPhasePhaseID returns the previous phase ID based on the given lastPhaseID.
func (p PhaseID) PreviousPhaseID(lastPhaseID PhaseID) PhaseID {
	if p == 1 {
		return lastPhaseID
	}
	return p - 1
}

// GameID is ID type of game.
type GameID = uint64

// GameStatusID is ID type of game status.
type GameStatusID = uint8

// GameSetting is setting of game.
type GameSetting struct {
	// RoleIDs is role ID list that can be played in the game.
	RoleIDs []RoleID

	// RequiredRoleIDs is role ID list that must be played in the game.
	RequiredRoleIDs []RoleID

	// NumberWerewolves is number of werewolves required to exist in the game.
	NumberWerewolves uint8

	// PlayerIDs is player ID list playing the game.
	PlayerIDs []PlayerID
}
