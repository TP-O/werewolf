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

func NewWerewolf(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(world, &action.VoteActionSetting{
		FactionID: constants.WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &werewolf{
		role: &role{
			id:           constants.WerewolfRoleID,
			factionID:    constants.WerewolfFactionID,
			phaseID:      constants.NightPhaseID,
			beginRoundID: constants.FirstRound,
			turnID:       constants.WerewolfTurnID,
			world:        world,
			playerID:     playerID,
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

	w.world.Poll(constants.VillagerFactionID).AddCandidates(w.playerID)
}

// OnRevoke is triggered when the role is removed from a player.
func (w *werewolf) OnRevoke() {
	w.role.OnRevoke()

	w.world.Poll(constants.VillagerFactionID).RemoveElector(w.playerID)
	w.world.Poll(constants.VillagerFactionID).RemoveCandidate(w.playerID)
	w.world.Poll(constants.WerewolfFactionID).RemoveElector(w.playerID)
}
