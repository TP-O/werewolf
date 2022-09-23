package core

import (
	"golang.org/x/exp/maps"

	"uwwolf/app/game/contract"
	"uwwolf/app/types"
	"uwwolf/app/util"
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

func (p *player) Id() types.PlayerId {
	return p.id
}

func (p *player) RoleIds() []types.RoleId {
	return maps.Keys(p.roles)
}

func (p *player) FactionId() types.FactionId {
	return p.factionId
}

func (p *player) AssignRoles(roles ...contract.Role) {
	for _, role := range roles {
		if !util.ExistKeyInMap(p.roles, role.Id()) {
			p.roles[role.Id()] = role
			p.modifyFaction(role.FactionId())
		}
	}
}

func (p *player) modifyFaction(factionId types.FactionId) {
	if factionId > p.factionId {
		p.factionId = factionId
	}
}

func (p *player) UseSkill(req *types.ActionRequest) *types.ActionResponse {
	if role := p.roles[p.game.Round().CurrentTurn().RoleId()]; role != nil {
		return role.ActivateSkill(req)
	} else {
		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag:   types.UnauthorizedErrorTag,
				Alert: "Unable to activate skill!",
			},
		}
	}
}
