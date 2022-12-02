package role

import (
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type ability struct {
	action      contract.Action
	activeLimit enum.Limit
}

type role struct {
	id         enum.RoleID
	phaseID    enum.PhaseID
	factionID  enum.FactionID
	beginRound enum.Round
	priority   enum.Priority
	game       contract.Game
	player     contract.Player
	abilities  map[enum.ActionID]*ability
}

func (r *role) ID() enum.RoleID {
	return r.id
}

func (r *role) PhaseID() enum.PhaseID {
	return r.phaseID
}

func (r *role) FactionID() enum.FactionID {
	return r.factionID
}

func (r *role) Priority() enum.Priority {
	return r.priority
}

func (r *role) BeginRound() enum.Round {
	return r.beginRound
}

func (r *role) ActiveLimit(actionID enum.ActionID) enum.Limit {
	limit := enum.ReachedLimit

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
			if ability.activeLimit == enum.ReachedLimit {
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
				ability.activeLimit != enum.Unlimited {
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
