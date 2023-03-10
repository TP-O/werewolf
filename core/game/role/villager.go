package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

type villager struct {
	*role
}

func NewVillager(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, &action.VoteActionSetting{
		FactionID: vars.VillagerFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &villager{
		role: &role{
			id:           vars.VillagerRoleID,
			factionID:    vars.VillagerFactionID,
			phaseID:      vars.DayPhaseID,
			beginRoundID: vars.FirstRound,
			turnID:       vars.VillagerTurnID,
			game:         game,
			player:       game.Player(playerID),
			abilities: []*ability{
				{
					action:      voteAction,
					activeLimit: vars.UnlimitedTimes,
				},
			},
		},
	}, nil
}

// RegisterTurn adds role's turn to the game schedule.
func (v *villager) RegisterSlot() {
	v.role.RegisterSlot()

	v.game.Poll(vars.VillagerFactionID).AddCandidates(v.player.ID())
	v.game.Poll(vars.WerewolfFactionID).AddCandidates(v.player.ID())
}

// UnregisterSlot removes role's slot from the game schedule.
func (v *villager) UnregisterSlot() {
	v.role.UnregisterSlot()

	v.game.Poll(vars.VillagerFactionID).RemoveElector(v.player.ID())
	v.game.Poll(vars.VillagerFactionID).RemoveCandidate(v.player.ID())
	v.game.Poll(vars.WerewolfFactionID).RemoveCandidate(v.player.ID())
}
