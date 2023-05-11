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
				id:           constants.HunterRoleId,
				phaseID:      constants.DayPhaseId,
				factionID:    constants.VillagerFactionId,
				beginRoundID: constants.FirstRound,
				turnID:       constants.HunterTurnID,
				moderator:    moderator,
				playerId:     playerId,
				abilities: []*ability{
					{
						action:      action.NewKill(moderator.World()),
						activeLimit: constants.OutOfTimes,
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
	diedAtPhaseID := h.moderator.World().Scheduler().PhaseID()

	// Ability is disabled if current round is too early
	if h.moderator.World().Scheduler().RoundID() < h.beginRoundID {
		return
	}

	// This turn can be only played in the current round
	slot := &types.NewTurnSlot{
		PhaseID:       h.phaseID,
		PlayerId:      h.playerId,
		RoleId:        h.id,
		PlayedRoundID: h.moderator.World().Scheduler().RoundID(),
	}

	if diedAtPhaseID == h.phaseID {
		// Play in next turn if he dies at his phase
		slot.TurnId = h.moderator.World().Scheduler().TurnID() + 1
	} else {
		// Play in his turn of the next day if he dies at
		// a time is not his phase
		slot.TurnId = h.turnID
	}

	h.abilities[0].activeLimit = constants.Once
	h.moderator.World().Scheduler().AddSlot(slot)
}
