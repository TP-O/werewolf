package contract

import "uwwolf/app/types"

type Round interface {
	// CurrentId returns the current round id.
	CurrentId() types.RoundId

	// CurrentPhaseId returns the current phase id.
	CurrentPhaseId() types.PhaseId

	// CurrentTurn returns the current turn instance.
	CurrentTurn() *types.Turn

	// CurrentPhase returns the turn array of the current
	// phase.
	CurrentPhase() types.Phase

	// IsAllowed returns true if it is given player turn now.
	IsAllowed(playerId types.PlayerId) bool

	// IsEmpty returns true if all phases in round are empty.
	IsEmpty() bool

	// Reset renews the round.
	Reset()

	// NextTurn moves to the next turn. Returns false if
	// the round is empty.
	NextTurn() bool

	// RemoveTurn removes the turn of given role id from
	// the round.
	RemoveTurn(roleId types.RoleId) bool

	// AddTurn adds new turn to the round. Returns false if
	// new turn is invalid.
	AddTurn(setting *types.TurnSetting) bool

	// AddPlayer adds new player to the turn of given role id.
	AddPlayer(playerId types.PlayerId, roleId types.RoleId) bool

	// DeletePlayer removes the player from the turn of given role id.
	DeletePlayer(playerId types.PlayerId, roleId types.RoleId) bool

	// DeletePlayerFromAllTurns remove the player from round totally.
	DeletePlayerFromAllTurns(playerId types.PlayerId)
}
