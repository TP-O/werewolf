package role

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/action"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
)

type werewolf struct {
	*role
}

func NewWerewolf(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(world, &action.VoteActionSetting{
		FactionID: declare.WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &werewolf{
		role: &role{
			id:           declare.WerewolfRoleID,
			factionID:    declare.WerewolfFactionID,
			phaseID:      declare.NightPhaseID,
			beginRoundID: declare.FirstRound,
			turnID:       declare.WerewolfTurnID,
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
func (w *werewolf) OnAssign() {
	w.role.OnAssign()

	w.world.Poll(declare.VillagerFactionID).AddCandidates(w.playerID)
}

// OnRevoke is triggered when the role is removed from a player.
func (w *werewolf) OnRevoke() {
	w.role.OnRevoke()

	w.world.Poll(declare.VillagerFactionID).RemoveElector(w.playerID)
	w.world.Poll(declare.VillagerFactionID).RemoveCandidate(w.playerID)
	w.world.Poll(declare.WerewolfFactionID).RemoveElector(w.playerID)
}
