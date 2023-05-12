package role

import (
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type twoSister struct {
	*role
}

func NewTwoSister(moderator contract.Moderator, playerId types.PlayerId) (contract.Role, error) {
	role := &twoSister{
		role: &role{
			id:         constants.TwoSistersRoleId,
			factionId:  constants.VillagerFactionId,
			phaseId:    constants.NightPhaseId,
			beginRound: constants.FirstRound,
			turn:       constants.TwoSistersTurn,
			moderator:  moderator,
			playerId:   playerId,
		},
	}

	// TwoSister ability will be executed automatically
	action :=
		action.NewRoleIdentify(
			moderator.World(),
			constants.TwoSistersRoleId,
		)
	moderator.RegisterActionExecution(types.ExecuteActionRegistration{
		RoleId:   role.Id(),
		ActionId: action.Id(),
		IsRoundMatched: func() bool {
			return moderator.Scheduler().Round() == role.phaseId
		},
		IsPhaseIdMatched: func() bool {
			return moderator.Scheduler().PhaseId() == role.phaseId
		},
		IsTurnMatched: func() bool {
			return moderator.Scheduler().Turn() == role.turn
		},
		Exec: func() types.ActionResponse {
			return action.Execute(types.ActionRequest{
				ActorId: playerId,
			})
		},
	})

	return role, nil
}
