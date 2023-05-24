package role

import (
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type hunter struct {
	*role
}

func NewHunter(moderator contract.Moderator, playerId types.PlayerId) (contract.Role, error) {
	return &hunter{
			role: &role{
				id:         constants.HunterRoleId,
				phaseId:    constants.DayPhaseId,
				factionId:  constants.VillagerFactionId,
				beginRound: constants.FirstRound,
				turn:       constants.HunterTurn,
				moderator:  moderator,
				playerId:   playerId,
				abilities: []*ability{
					{
						action:      action.NewKill(moderator.World()),
						activeLimit: constants.OutOfTimes,
						effectiveAt: effectiveAt{
							isImmediate: true,
						},
					},
				},
			},
		},
		nil
}

// OnAssign is triggered when the role is assigned to a player.
func (h *hunter) OnAfterAssign() {
	//
}

// OnAfterDeath is triggered after killing this role.
func (h *hunter) OnAfterDeath() {
	diedAtPhaseId := h.moderator.Scheduler().PhaseId()

	// Ability is disabled if current round is too early
	if h.moderator.Scheduler().Round() < h.beginRound {
		return
	}

	// This turn can be only played in the current round
	slot := types.AddTurnSlot{
		PhaseId:  h.phaseId,
		PlayerId: h.playerId,
		TurnSlot: types.TurnSlot{
			RoleId:      h.id,
			PlayedRound: h.moderator.Scheduler().Round(),
		},
	}

	if diedAtPhaseId == h.phaseId {
		// Play in next turn if he dies at his phase
		slot.Turn = h.moderator.Scheduler().Turn() + 1
	} else {
		// Play in his turn of the next day if he dies at
		// a time is not his phase
		slot.Turn = h.turn
	}

	h.abilities[0].activeLimit = constants.Once
	h.moderator.Scheduler().AddSlot(slot)
}
