package contract

import "uwwolf/types"

type Game interface {
	GetCurrentRoundId() types.RoundId

	GetCurrentRoleId() types.RoleId

	GetCurrentPhaseId() types.PhaseId

	GetPlayer(playerId types.PlayerId) Player

	KillPlayer(playerId types.PlayerId) Player

	RequestAction(playerId types.PlayerId, req *types.ActionRequest) *types.ActionResponse
}
