package contract

import "uwwolf/types"

type Player interface {
	GetId() types.PlayerId

	GetSId() types.SocketId

	UseSkill(data *types.ActionData) *types.PerformResult
}
