package core

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
	"uwwolf/util"
)

type player struct {
	id    types.PlayerId
	sId   types.SocketId
	game  contract.Game
	roles map[types.RoleId]contract.Role
}

func NewPlayer(sId types.SocketId, pId types.PlayerId, game contract.Game) contract.Player {
	return &player{
		id:   pId,
		sId:  sId,
		game: game,
	}
}

func (p *player) GetId() types.PlayerId {
	return p.id
}

func (p *player) GetSId() types.SocketId {
	return p.sId
}

func (p *player) UseSkill(data *types.ActionData) *types.PerformResult {
	if role, errRes := p.getRoleOfCurrentTurn(); errRes != nil {
		return errRes
	} else {
		return role.ActivateSkill(data)
	}
}

func (p *player) getRoleOfCurrentTurn() (contract.Role, *types.PerformResult) {
	roleId := p.game.GetCurrentRoleId()

	if !util.ExistKeyInMap(p.roles, roleId) {
		return nil,
			&types.PerformResult{
				ErrorTag: types.SystemErrorTag,
				Errors: map[string]string{
					types.SystemErrorProperty: "Not your turn!",
				},
			}
	}

	return p.roles[roleId], nil
}
