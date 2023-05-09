package role

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/action"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
)

type villager struct {
	*role
}

func NewVillager(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(world, &action.VoteActionSetting{
		FactionID: declare.VillagerFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &villager{
		role: &role{
			id:           declare.VillagerRoleID,
			factionID:    declare.VillagerFactionID,
			phaseID:      declare.DayPhaseID,
			beginRoundID: declare.FirstRound,
			turnID:       declare.VillagerTurnID,
			world:        world,
			playerID:     playerID,
			abilities: []*ability{
				{
					action:      voteAction,
					activeLimit: declare.UnlimitedTimes,
				},
			},
		},
	}, nil
}

// OnAssign is triggered when the role is assigned to a player.
func (v *villager) OnAssign() {
	v.role.OnAssign()

	v.world.Poll(declare.VillagerFactionID).AddCandidates(v.playerID)
	v.world.Poll(declare.WerewolfFactionID).AddCandidates(v.playerID)
}

// OnRevoke is triggered when the role is removed from a player.
func (v *villager) OnRevoke() {
	v.role.OnRevoke()

	v.world.Poll(declare.VillagerFactionID).RemoveElector(v.playerID)
	v.world.Poll(declare.VillagerFactionID).RemoveCandidate(v.playerID)
	v.world.Poll(declare.WerewolfFactionID).RemoveCandidate(v.playerID)
}
