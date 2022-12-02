package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"
)

type ability struct {
	action      contract.Action
	activeLimit types.Limit
}

type role struct {
	id         types.RoleID
	phaseID    types.PhaseID
	factionID  types.FactionID
	beginRound types.Round
	priority   types.Priority
	game       contract.Game
	player     contract.Player
	abilities  map[types.ActionID]*ability
}

func (r *role) ID() types.RoleID {
	return r.id
}

func (r *role) PhaseID() types.PhaseID {
	return r.phaseID
}

func (r *role) FactionID() types.FactionID {
	return r.factionID
}

func (r *role) Priority() types.Priority {
	return r.priority
}

func (r *role) BeginRound() types.Round {
	return r.beginRound
}

func (r *role) ActiveLimit(actionID types.ActionID) types.Limit {
	limit := config.ReachedLimit

	if ability := r.abilities[actionID]; ability != nil {
		limit = ability.activeLimit
	}

	return limit
}

func (r *role) BeforeDeath() bool {
	return true
}

func (r *role) AfterDeath() {
	//
}

func (r *role) UseAbility(req *types.UseRoleRequest) *types.ActionResponse {
	for _, ability := range r.abilities {
		if req.ActionID == ability.action.ID() {
			if ability.activeLimit == config.ReachedLimit {
				return &types.ActionResponse{
					Ok:      false,
					Message: "Unable to use this ability anymore ¯\\(º_o)/¯",
				}
			}

			res := ability.action.Execute(&types.ActionRequest{
				ActorID:   r.player.ID(),
				TargetIDs: req.TargetIDs,
				IsSkipped: req.IsSkipped,
			})

			if res.Ok &&
				!res.IsSkipped &&
				ability.activeLimit != config.Unlimited {
				ability.activeLimit--
			}

			return res
		}
	}

	return &types.ActionResponse{
		Ok:      false,
		Message: "This is beyond your ability (╥﹏╥)",
	}
}
