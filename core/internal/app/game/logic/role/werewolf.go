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

func NewWerewolf(world contract.World, playerId types.PlayerId) (contract.Role, error) {
	voteAction, err := action.NewVote(world, &action.VoteActionSetting{
		FactionId: constants.WerewolfFactionId,
		PlayerId:  playerId,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &werewolf{
		role: &role{
			id:           constants.WerewolfRoleId,
			factionID:    constants.WerewolfFactionId,
			phaseID:      constants.NightPhaseId,
			beginRoundID: constants.FirstRound,
			turnID:       constants.WerewolfTurnID,
			world:        world,
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
func (w *werewolf) OnAssign() {
	w.role.OnAssign()

	w.world.Poll(constants.VillagerFactionId).AddCandidates(w.playerId)
}

// OnRevoke is triggered when the role is removed from a player.
func (w *werewolf) OnRevoke() {
	w.role.OnRevoke()

	w.world.Poll(constants.VillagerFactionId).RemoveElector(w.playerId)
	w.world.Poll(constants.VillagerFactionId).RemoveCandidate(w.playerId)
	w.world.Poll(constants.WerewolfFactionId).RemoveElector(w.playerId)
}
