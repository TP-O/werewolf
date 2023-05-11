package types

import "time"

// PhaseID is ID type of phase.
type PhaseId = uint8

// GameID is ID type of game.
type GameID uint64

// IsUnknown checks if faction ID is 0.
func (g GameID) IsUnknown() bool {
	return g == 0
}

// GameStatusID is ID type of game status.
type GameStatusID = uint8

// GameInitialization is the required options to start a game.
type GameInitialization struct {
	// RoleIDs is role ID list that can be played in the game.
	RoleIds []RoleId

	// RequiredRoleIDs is role ID list that must be played in the game.
	RequiredRoleIds []RoleId

	// NumberWerewolves is number of werewolves required to exist in the game.
	NumberWerewolves uint8

	// PlayerIDs is player ID list playing the game.
	PlayerIDs []PlayerId
}

// GameRegistration contains the game configuration and joined players.
type GameRegistration struct {
	// GameInitialization is the required options to start a game.
	GameInitialization

	// ID is game ID.
	ID GameID

	// TurnDuration is the duration of a turn.
	TurnDuration time.Duration

	// DiscussionDuration is the duration of the villager discussion.
	DiscussionDuration time.Duration
}
