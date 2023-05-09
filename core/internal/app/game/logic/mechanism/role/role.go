package role

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
)

// ability contains one action and its limit.
type ability struct {
	// action is a specific action.
	action contract.Action

	// activeLimit is number of times the action can be used.
	activeLimit types.Times
}

// role is the basis for all concreate roles. The concrete role
// must embed this struct and modify its methods as required.
type role struct {
	// id is the role ID.
	id types.RoleID

	// phaseID is the active phase ID of this role.
	phaseID types.PhaseID

	// factionID is the faction ID to which this role belongs.
	factionID types.FactionID

	// beginRoundID is the round where this role begins to playing.
	beginRoundID types.RoundID

	playerID types.PlayerID

	// turnID is play order of this role in the active phase.
	turnID types.TurnID

	// isBeforeDeathTriggered marks that `BeforeDeath` function
	// is called or not.
	isBeforeDeathTriggered bool

	// game is the game instance this role affects.
	world contract.World

	// abilities is the abilities of this role.
	abilities []*ability
}

// ID returns role's ID.
func (r role) ID() types.RoleID {
	return r.id
}

// PhaseID returns role's active phase ID.
// func (r role) PhaseID() types.PhaseID {
// 	return r.phaseID
// }

func (r role) FactionID() types.FactionID {
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
	limit := declare.OutOfTimes
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
	r.world.Scheduler().AddSlot(&types.NewTurnSlot{
		PhaseID:      r.phaseID,
		TurnID:       r.turnID,
		BeginRoundID: r.beginRoundID,
		PlayerID:     r.playerID,
		RoleID:       r.id,
	})
}

// OnRevoke is triggered when the role is removed from a player.
func (r *role) OnRevoke() {
	r.world.Scheduler().RemoveSlot(&types.RemovedTurnSlot{
		PhaseID:  r.phaseID,
		RoleID:   r.id,
		PlayerID: r.playerID,
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
	if ability.activeLimit == declare.OutOfTimes {
		return &types.ActionResponse{
			Ok:      false,
			Message: "Unable to use this ability anymore ¯\\(º_o)/¯",
		}
	}

	res := ability.action.Execute(&types.ActionRequest{
		ActorID:   r.playerID,
		TargetID:  req.TargetID,
		IsSkipped: req.IsSkipped,
	})
	if res.Ok &&
		!req.IsSkipped &&
		ability.activeLimit != declare.UnlimitedTimes {
		ability.activeLimit--

		// Remove the player turn if the limit is reached
		if r.ActiveTimes(-1) == declare.OutOfTimes {
			r.world.Scheduler().RemoveSlot(&types.RemovedTurnSlot{
				PhaseID:  r.phaseID,
				RoleID:   r.id,
				PlayerID: r.playerID,
			})
		}
	}

	return res
}
