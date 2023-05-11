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
			id:           constants.VillagerRoleId,
			factionID:    constants.VillagerFactionId,
			phaseID:      constants.DayPhaseId,
			beginRoundID: constants.FirstRound,
			turnID:       constants.VillagerTurnID,
			moderator:    moderator,
			playerId:     playerId,
			abilities: []*ability{
				{
					action:      voteAction,
					activeLimit: constants.UnlimitedTimes,
				},
			},
		},
	}, nil
}

// OnAssign is triggered when the role is assigned to a player.
func (v *villager) OnAssign() {
	v.role.OnAssign()

	v.moderator.World().Poll(constants.VillagerFactionId).AddCandidates(v.playerId)
	v.moderator.World().Poll(constants.WerewolfFactionId).AddCandidates(v.playerId)
}

// OnRevoke is triggered when the role is removed from a player.
func (v *villager) OnRevoke() {
	v.role.OnRevoke()

	v.moderator.World().Poll(constants.VillagerFactionId).RemoveElector(v.playerId)
	v.moderator.World().Poll(constants.VillagerFactionId).RemoveCandidate(v.playerId)
	v.moderator.World().Poll(constants.WerewolfFactionId).RemoveCandidate(v.playerId)
}
