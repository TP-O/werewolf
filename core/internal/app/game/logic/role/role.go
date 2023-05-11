package role

import (
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"
)

// ability contains one action and its limit.
type ability struct {
	// action is a specific action.
	action contract.Action

	// activeLimit is number of times the action can be used.
	activeLimit types.Times

	round types.RoundID

	phaseID types.PhaseID

	turnId types.TurnId
}

// role is the basis for all concreate roles. The concrete role
// must embed this struct and modify its methods as required.
type role struct {
	// id is the role ID.
	id types.RoleId

	// phaseID is the active phase ID of this role.
	phaseID types.PhaseID

	// factionID is the faction ID to which this role belongs.
	factionID types.FactionId

	// beginRoundID is the round where this role begins to playing.
	beginRoundID types.RoundID

	playerId types.PlayerId

	// turnID is play order of this role in the active phase.
	turnID types.TurnId

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
	return r.factionID
}

// TurnID returns role's turn order in active phase.
// func (r role) TurnID() types.TurnID {
// 	return r.turnID
// }

// BeginRoundID returns round in which this role be able to
// use its abilities.
// func (r role) BeginRoundID() types.RoundID {
// 	return r.beginRoundID
// }

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
func (r *role) OnAssign() {
	r.moderator.World().Scheduler().AddSlot(&types.NewTurnSlot{
		PhaseID:      r.phaseID,
		TurnId:       r.turnID,
		BeginRoundID: r.beginRoundID,
		PlayerId:     r.playerId,
		RoleId:       r.id,
	})
}

// OnRevoke is triggered when the role is removed from a player.
func (r *role) OnRevoke() {
	r.moderator.World().Scheduler().RemoveSlot(&types.RemovedTurnSlot{
		PhaseID:  r.phaseID,
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
func (r *role) ActivateAbility(req *types.ActivateAbilityRequest) *types.ActionResponse {
	if int(req.AbilityIndex) >= len(r.abilities) {
		return &types.ActionResponse{
			Ok:      false,
			Message: "This is beyond your ability (╥﹏╥)",
		}
	}

	ability := r.abilities[req.AbilityIndex]
	if ability.activeLimit == constants.OutOfTimes {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Unable to use this ability anymore ¯\\(º_o)/¯",
		}
	}

	res := func() types.ActionResponse {
		if util.IsZero(ability.round) || util.IsZero(ability.phaseID) || util.IsZero(ability.turnId) || req.IsSkipped {
			return ability.action.Execute(types.ActionRequest{
				ActorId:   r.playerId,
				TargetId:  req.TargetID,
				IsSkipped: req.IsSkipped,
			})
		}

		r.moderator.RegisterActionExecution(ability.round, ability.phaseID, ability.turnId, func() {
			ability.action.Execute(types.ActionRequest{
				ActorId:   r.playerId,
				TargetId:  req.TargetID,
				IsSkipped: req.IsSkipped,
			})
		})
		return types.ActionResponse{
			Ok:       true,
			ActionId: ability.action.Id(),
			ActionRequest: types.ActionRequest{
				ActorId:  r.playerId,
				TargetId: req.TargetID,
			},
		}
	}()

	if res.Ok &&
		!req.IsSkipped &&
		ability.activeLimit != constants.UnlimitedTimes {
		ability.activeLimit--

		// Remove the player turn if the limit is reached
		if r.ActiveTimes(-1) == constants.OutOfTimes {
			r.moderator.World().Scheduler().RemoveSlot(&types.RemovedTurnSlot{
				PhaseID:  r.phaseID,
				RoleId:   r.id,
				PlayerId: r.playerId,
			})
		}
	}

	return &res
}
