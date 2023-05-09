package role

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/action"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
)

type hunter struct {
	*role
}

func NewHunter(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	return &hunter{
			role: &role{
				id:           declare.HunterRoleID,
				phaseID:      declare.DayPhaseID,
				factionID:    declare.VillagerFactionID,
				beginRoundID: declare.FirstRound,
				turnID:       declare.HunterTurnID,
				world:        world,
				playerID:     playerID,
				abilities: []*ability{
					{
						action:      action.NewKill(world),
						activeLimit: declare.OutOfTimes,
					},
				},
			},
		},
		nil
}

// OnAssign is triggered when the role is assigned to a player.
func (h *hunter) OnAssign() {
	//
}

// OnAfterDeath is triggered after killing this role.
func (h *hunter) OnAfterDeath() {
	diedAtPhaseID := h.world.Scheduler().PhaseID()

	// Ability is disabled if current round is too early
	if h.world.Scheduler().RoundID() < h.beginRoundID {
		return
	}

	// This turn can be only played in the current round
	slot := &types.NewTurnSlot{
		PhaseID:       h.phaseID,
		PlayerID:      h.playerID,
		RoleID:        h.id,
		PlayedRoundID: h.world.Scheduler().RoundID(),
	}

	if diedAtPhaseID == h.phaseID {
		// Play in next turn if he dies at his phase
		slot.TurnID = h.world.Scheduler().TurnID() + 1
	} else {
		// Play in his turn of the next day if he dies at
		// a time is not his phase
		slot.TurnID = h.turnID
	}

	h.abilities[0].activeLimit = declare.Once
	h.world.Scheduler().AddSlot(slot)
}
