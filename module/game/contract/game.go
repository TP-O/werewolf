package contract

import (
	"uwwolf/module/game/state"
	"uwwolf/types"
)

type Game interface {
	IsStarted() bool

	GetCurrentRoundId() types.RoundId

	GetCurrentRoleId() types.RoleId

	GetCurrentPhaseId() types.PhaseId

	GetPoll(factionId types.FactionId) *state.Poll

	Start() bool

	GetPlayer(playerId types.PlayerId) Player

	KillPlayer(playerId types.PlayerId) Player

	RequestAction(playerId types.PlayerId, req *types.ActionRequest) *types.ActionResponse
}
