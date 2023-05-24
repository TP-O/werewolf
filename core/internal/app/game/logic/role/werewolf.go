package role

import (
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type werewolf struct {
	*role
}

func NewWerewolf(moderator contract.Moderator, playerId types.PlayerId) (contract.Role, error) {
	voteAction, err := action.NewVote(moderator.World(), &action.VoteActionSetting{
		FactionId: constants.WerewolfFactionId,
		PlayerId:  playerId,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &werewolf{
		role: &role{
			id:         constants.WerewolfRoleId,
			factionId:  constants.WerewolfFactionId,
			phaseId:    constants.NightPhaseId,
			beginRound: constants.FirstRound,
			turn:       constants.WerewolfTurn,
			moderator:  moderator,
			playerId:   playerId,
			abilities: []*ability{
				{
					action:      voteAction,
					activeLimit: constants.UnlimitedTimes,
					effectiveAt: effectiveAt{
						isImmediate: true,
					},
				},
			},
		},
	}, nil
}

// OnAssign is triggered when the role is assigned to a player.
func (w *werewolf) OnAfterAssign() {
	w.role.OnAfterAssign()

	w.moderator.World().Poll(constants.WerewolfFactionId).AddElectors(w.playerId)
}

// OnRevoke is triggered when the role is removed from a player.
func (w *werewolf) OnAfterRevoke() {
	w.role.OnAfterRevoke()

	w.moderator.World().Poll(constants.WerewolfFactionId).RemoveElector(w.playerId)
}
