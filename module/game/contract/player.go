package contract

import "uwwolf/types"

type Player interface {
	GetId() types.PlayerId

	GetSId() types.SocketId

	GetFactionId() types.FactionId

	UseSkill(req *types.ActionRequest) *types.ActionResponse
}
