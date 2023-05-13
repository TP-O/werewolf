package role

import (
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type villager struct {
	*role
}

func NewVillager(moderator contract.Moderator, playerId types.PlayerId) (contract.Role, error) {
	voteAction, err := action.NewVote(moderator.World(), &action.VoteActionSetting{
		FactionId: constants.VillagerFactionId,
		PlayerId:  playerId,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &villager{
		role: &role{
			id:         constants.VillagerRoleId,
			factionId:  constants.VillagerFactionId,
			phaseId:    constants.DayPhaseId,
			beginRound: constants.FirstRound,
			turn:       constants.VillagerTurn,
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
func (v *villager) OnAfterAssign() {
	v.role.OnAfterAssign()

	v.moderator.World().Poll(constants.VillagerFactionId).AddElectors(v.playerId)
	v.moderator.World().Poll(constants.VillagerFactionId).AddCandidates(v.playerId)
	v.moderator.World().Poll(constants.WerewolfFactionId).AddCandidates(v.playerId)
}

// OnRevoke is triggered when the role is removed from a player.
func (v *villager) OnAfterRevoke() {
	v.role.OnAfterRevoke()

	v.moderator.World().Poll(constants.VillagerFactionId).RemoveElector(v.playerId)
	v.moderator.World().Poll(constants.VillagerFactionId).RemoveCandidate(v.playerId)
	v.moderator.World().Poll(constants.WerewolfFactionId).RemoveCandidate(v.playerId)
}
