package role

import (
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

// ability contains one action and its limit.
type ability struct {
	// action is a specific action.
	action contract.Action

	// activeLimit is number of times the action can be used.
	activeLimit types.Limit
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

	// turnID is play order of this role in the active phase.
	turnID types.TurnID

	// game is the game instance this role affects.
	game contract.Game

	// player is the player assigned this role.
	player contract.Player

	// abilities is the abilities of this role.
	abilities []ability
}

var _ contract.Role = (*role)(nil)

func (r role) ID() types.RoleID {
	return r.id
}

func (r role) PhaseID() types.PhaseID {
	return r.phaseID
}

func (r role) FactionID() types.FactionID {
	return r.factionID
}

func (r role) TurnID() types.TurnID {
	return r.turnID
}

func (r role) BeginRoundID() types.RoundID {
	return r.beginRoundID
}

func (r role) ActiveLimit(index int) types.Limit {
	limit := ReachedLimit
	if index > -1 && index < len(r.abilities) {
		limit = r.abilities[index].activeLimit
	} else if index == -1 {
		for _, ability := range r.abilities {
			limit += ability.activeLimit
		}
	}

	return limit
}

func (r role) BeforeDeath() bool {
	return true
}

func (r role) AfterDeath() {
	//
}

func (r role) AfterSaved() {
	//
}

func (r *role) ActivateAbility(req types.ActivateAbilityRequest) types.ActionResponse {
	if int(req.AbilityIndex) >= len(r.abilities) {
		return types.ActionResponse{
			Ok:      false,
			Message: "This is beyond your ability (╥﹏╥)",
		}
	}

	ability := r.abilities[req.AbilityIndex]
	if ability.activeLimit == ReachedLimit {
		return types.ActionResponse{
			Ok:      false,
			Message: "Unable to use this ability anymore ¯\\(º_o)/¯",
		}
	}

	res := ability.action.Execute(types.ActionRequest{
		ActorID:   r.player.ID(),
		TargetID:  req.TargetID,
		IsSkipped: req.IsSkipped,
	})
	if res.Ok &&
		!req.IsSkipped &&
		ability.activeLimit != Unlimited {
		ability.activeLimit--

		// Remove the player turn if the limit is reached
		if r.ActiveLimit(-1) == ReachedLimit {
			r.game.Scheduler().RemovePlayerTurn(types.RemovedPlayerTurn{
				PhaseID:  r.phaseID,
				RoleID:   r.id,
				PlayerID: r.player.ID(),
			})
		}
	}

	return res
}
