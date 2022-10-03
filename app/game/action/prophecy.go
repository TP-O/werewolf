package action

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
	"uwwolf/util"
)

type knowledge struct {
	// Map factions with their players
	Factions map[types.FactionId][]types.PlayerId `json:"factions"`

	// Map players with their faction
	Players map[types.PlayerId]types.FactionId `json:"players"`
}

type prophecy struct {
	action[*knowledge]
}

func NewProphecy(game contract.Game) contract.Action {
	prophecy := prophecy{
		action: action[*knowledge]{
			id: enum.ProphecyActionId,
			state: &knowledge{
				Factions: make(map[types.FactionId][]types.PlayerId),
				Players:  make(map[types.PlayerId]types.FactionId),
			},
			game:       game,
			expiration: enum.UnlimitedTimes,
		},
	}

	return &prophecy
}

func (p *prophecy) Perform(req *types.ActionRequest) *types.ActionResponse {
	return p.action.perform(p.validate, p.execute, req)
}

func (p *prophecy) validate(req *types.ActionRequest) (alert string) {
	if req.ActorId == req.TargetIds[0] {
		alert = "WTF! You don't know who are you?"
	} else if util.ExistKeyInMap(p.state.Players, req.TargetIds[0]) {
		alert = "Already known identity!"
	}

	return
}

func (p *prophecy) execute(req *types.ActionRequest) *types.ActionResponse {
	factionId := p.game.Player(req.TargetIds[0]).FactionId()

	// Check if a player is werewolf or not
	if factionId == enum.WerewolfFactionId {
		p.state.Players[req.TargetIds[0]] = enum.WerewolfFactionId
		p.state.Factions[enum.WerewolfFactionId] = append(
			p.state.Factions[enum.WerewolfFactionId],
			req.TargetIds[0],
		)

		return &types.ActionResponse{
			Ok:   true,
			Data: true,
		}
	}

	p.state.Players[req.TargetIds[0]] = enum.VillagerFactionId
	p.state.Factions[enum.VillagerFactionId] = append(
		p.state.Factions[enum.VillagerFactionId],
		req.TargetIds[0],
	)

	return &types.ActionResponse{
		Ok:   true,
		Data: false,
	}
}
