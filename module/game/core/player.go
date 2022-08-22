package core

import (
	"fmt"
	"uwwolf/module/game/contract"
	"uwwolf/types"
	"uwwolf/util"
)

type player struct {
	id        types.PlayerId
	factionId types.FactionId
	game      contract.Game
	roles     map[types.RoleId]contract.Role
}

func NewPlayer(game contract.Game, id types.PlayerId) contract.Player {
	return &player{
		id:        id,
		factionId: types.UnknownFaction,
		game:      game,
		roles:     make(map[types.RoleId]contract.Role),
	}
}

func (p *player) GetId() types.PlayerId {
	return p.id
}

func (p *player) GetFactionId() types.FactionId {
	return p.factionId
}

func (p *player) AssignRoles(roles ...contract.Role) {
	for _, role := range roles {
		if !util.ExistKeyInMap(p.roles, role.GetId()) {
			p.roles[role.GetId()] = role
			p.ModifyFaction(role.GetFactionId())

			fmt.Println("Player: ", p.id, " - Role: ", role.GetName())
			// fmt.Println(role.GetName())

			// Chage faction id...
		}
	}
}

func (p *player) ModifyFaction(factionId types.FactionId) {
	//
}

func (p *player) UseSkill(req *types.ActionRequest) *types.ActionResponse {
	if role := p.roles[p.game.GetCurrentRoleId()]; role != nil {
		return role.ActivateSkill(req)
	} else {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.UnauthorizedErrorTag,
				Msg: map[string]string{
					types.AlertErrorField: "Unable to activate skill!",
				},
			},
		}
	}
}
