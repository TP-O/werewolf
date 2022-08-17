package contract

import "uwwolf/types"

type Game interface {
	GetCurrentRoundId() types.RoundId

	GetCurrentRoleId() types.RoleId

	// GetCurrentPhaseId() types.PhaseId

	// GetCurrentTurnIndex() int

	GetPlayer(playerId types.PlayerId) Player

	// RemovePlayer(playerId types.PlayerId) bool

	RequestAction(playerId types.PlayerId, data *types.ActionData) *types.PerformResult
}
