package role

import (
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

// role is the basis for all concreate roles. The concrete role
// must embed this struct and modify its methods as required.
type role struct {
	// id is the role ID.
	id types.RoleId

	// phaseID is the active phase ID of this role.
	phaseId types.PhaseId

	// factionID is the faction ID to which this role belongs.
	factionId types.FactionId

	// beginRoundID is the round where this role begins to playing.
	beginRound types.Round

	playerId types.PlayerId

	// turnID is play order of this role in the active phase.
	turn types.Turn

	// isBeforeDeathTriggered marks that `BeforeDeath` function
	// is called or not.
	isBeforeDeathTriggered bool

	// game is the game instance this role affects.
	moderator contract.Moderator

	// abilities is the abilities of this role.
	abilities []*ability
}

// ID returns role's ID.
func (r role) Id() types.RoleId {
	return r.id
}

// PhaseID returns role's active phase ID.
// func (r role) PhaseID() types.PhaseID {
// 	return r.phaseID
// }

func (r role) FactionId() types.FactionId {
	return r.factionId
}

// ActiveTimes returns remaining times this role can use the specific ability.
// Returns total limit if the `index` is -1.
func (r role) ActiveTimes(index int) types.Times {
	limit := constants.OutOfTimes
	if index > -1 && index < len(r.abilities) {
		limit = r.abilities[index].activeLimit
	} else if index == -1 {
		for _, ability := range r.abilities {
			limit += ability.activeLimit
		}
	}

	return limit
}

// OnAssign is triggered when the role is assigned to a player.
func (r *role) OnAfterAssign() {
	r.moderator.Scheduler().AddSlot(types.AddTurnSlot{
		PhaseId:  r.phaseId,
		Turn:     r.turn,
		PlayerId: r.playerId,
		TurnSlot: types.TurnSlot{
			BeginRound: r.beginRound,
			RoleId:     r.id,
		},
	})
}

// OnRevoke is triggered when the role is removed from a player.
func (r *role) OnAfterRevoke() {
	r.moderator.Scheduler().RemoveSlot(types.RemoveTurnSlot{
		PhaseId:  r.phaseId,
		RoleId:   r.id,
		PlayerId: r.playerId,
	})
}

// OnBeforeDeath is triggered before killing this role.
// If returns false, the player assigned it is saved.
func (r *role) OnBeforeDeath() bool {
	if r.isBeforeDeathTriggered {
		return true
	}

	// Do something...

	r.isBeforeDeathTriggered = true
	return true
}

// OnAfterDeath is triggered after killing this role.
func (r role) OnAfterDeath() {
	//
}

// ActivateAbility executes the action corresponding to the required ability.
func (r *role) Use(req types.RoleRequest) types.RoleResponse {
	if int(req.AbilityIndex) >= len(r.abilities) {
		return types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "This action is beyond your ability (╥﹏╥)",
			},
		}
	}

	ability := r.abilities[req.AbilityIndex]

	if ability.activeLimit == constants.OutOfTimes {
		return types.RoleResponse{
			ActionResponse: types.ActionResponse{
				Message: "Unable to use this action anymore ¯\\(º_o)/¯",
			},
		}
	}

	actionRes := func() types.ActionResponse {
		if ability.isImmediate ||
			req.IsSkipped {
			return ability.action.Execute(types.ActionRequest{
				ActorId:   r.playerId,
				TargetId:  req.TargetId,
				IsSkipped: req.IsSkipped,
			})
		}

		r.moderator.RegisterActionExecution(types.ExecuteActionRegistration{
			RoleId:     r.Id(),
			ActionId:   ability.action.Id(),
			CanExecute: ability.CanExecute,
			Exec: func() types.ActionResponse {
				return ability.action.Execute(types.ActionRequest{
					ActorId:   r.playerId,
					TargetId:  req.TargetId,
					IsSkipped: req.IsSkipped,
				})
			},
		})
		return types.ActionResponse{
			Ok:       true,
			ActionId: ability.action.Id(),
			ActionRequest: types.ActionRequest{
				ActorId:  r.playerId,
				TargetId: req.TargetId,
			},
			Message: "Action is registered!",
		}
	}()

	if actionRes.Ok &&
		!req.IsSkipped &&
		ability.activeLimit != constants.UnlimitedTimes {
		ability.activeLimit--

		// Remove the player turn if the limit is reached
		if r.ActiveTimes(-1) == constants.OutOfTimes {
			r.moderator.Scheduler().RemoveSlot(types.RemoveTurnSlot{
				PhaseId:  r.phaseId,
				RoleId:   r.id,
				PlayerId: r.playerId,
			})
		}
	}

	return types.RoleResponse{
		Round:          r.moderator.Scheduler().Round(),
		PhaseId:        r.moderator.Scheduler().PhaseId(),
		Turn:           r.moderator.Scheduler().Turn(),
		RoleId:         r.Id(),
		ActionResponse: actionRes,
	}
}
