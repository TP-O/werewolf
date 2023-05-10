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

func NewVillager(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(world, &action.VoteActionSetting{
		FactionID: constants.VillagerFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &villager{
		role: &role{
			id:           constants.VillagerRoleID,
			factionID:    constants.VillagerFactionID,
			phaseID:      constants.DayPhaseID,
			beginRoundID: constants.FirstRound,
			turnID:       constants.VillagerTurnID,
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
func (v *villager) OnAssign() {
	v.role.OnAssign()

	v.world.Poll(constants.VillagerFactionID).AddCandidates(v.playerID)
	v.world.Poll(constants.WerewolfFactionID).AddCandidates(v.playerID)
}

// OnRevoke is triggered when the role is removed from a player.
func (v *villager) OnRevoke() {
	v.role.OnRevoke()

	v.world.Poll(constants.VillagerFactionID).RemoveElector(v.playerID)
	v.world.Poll(constants.VillagerFactionID).RemoveCandidate(v.playerID)
	v.world.Poll(constants.WerewolfFactionID).RemoveCandidate(v.playerID)
}
