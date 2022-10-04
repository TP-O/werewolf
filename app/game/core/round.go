package core

import (
	"golang.org/x/exp/slices"

	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type round struct {
	id               types.RoundId
	currentPhaseId   types.PhaseId
	currentTurnIndex int
	phases           map[types.PhaseId]types.Phase
}

func NewRound() contract.Round {
	round := &round{}

	// Set default values
	round.Reset()

	return round
}

func (r *round) CurrentId() types.RoundId {
	return r.id
}

func (r *round) CurrentPhaseId() types.PhaseId {
	return r.currentPhaseId
}

func (r *round) CurrentTurn() *types.Turn {
	if len(r.CurrentPhase()) == 0 || r.currentTurnIndex >= len(r.CurrentPhase()) {
		return nil
	}

	return r.CurrentPhase()[r.currentTurnIndex]
}

func (r *round) CurrentPhase() types.Phase {
	return r.phases[r.currentPhaseId]
}

func (r *round) IsAllowed(playerId types.PlayerId) bool {
	if r.CurrentTurn() == nil {
		return false
	}

	return slices.Contains(r.CurrentTurn().PlayerIds, playerId)
}

func (r *round) IsEmpty() bool {
	for _, p := range r.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

func (r *round) Reset() {
	r.id = 1
	r.currentTurnIndex = 0
	r.currentPhaseId = enum.NightPhaseId
	r.phases = make(map[types.PhaseId]types.Phase)
}

// Move to next turn and delete previous turn if it's times is out.
// Repeat from the beginning if the end is exceeded and return false
// if round is empty.
func (r *round) NextTurn(isSkipped bool) bool {
	if r.IsEmpty() {
		return false
	}

	if isSkipped {
		r.passCurrentTurn()
	}

	if r.currentTurnIndex < len(r.CurrentPhase())-1 {
		r.currentTurnIndex++

		// Skip turn if not the time
		if r.CurrentTurn().BeginRound > r.id {
			return r.NextTurn(true)
		}
	} else {
		r.currentTurnIndex = 0
		r.currentPhaseId = (r.currentPhaseId + 1) % (enum.DuskPhaseId + 1)

		// Start new round
		if r.currentPhaseId == 0 {
			r.currentPhaseId = enum.NightPhaseId
			r.id++
		}

		if r.CurrentTurn() == nil {
			return r.NextTurn(true)
		}
	}

	return true
}

func (r *round) passCurrentTurn() {
	currentTurn := r.CurrentTurn()

	if currentTurn != nil &&
		currentTurn.BeginRound <= r.id &&
		currentTurn.Expiration != enum.UnlimitedTimes {

		currentTurn.Expiration--
		r.removeTurnIfDone(currentTurn.RoleId)
	}
}

// Remove a turn if it's times is out.
func (r *round) removeTurnIfDone(roleId types.RoleId) bool {
	currentTurn := r.CurrentTurn()

	if currentTurn == nil || currentTurn.Expiration != enum.OutOfTimes {
		return false
	}

	return r.RemoveTurn(roleId)
}

// Remove turn. In the case that removed turn is the same as the current
// turn index and phase containing the turn is also  the same as current
// phase, decreases the current turn index by 1.
func (r *round) RemoveTurn(roleId types.RoleId) bool {
	for phaseId, phase := range r.phases {
		for turnIndex, turn := range phase {
			if turn.RoleId == roleId {
				r.phases[phaseId] = slices.Delete(r.phases[phaseId], turnIndex, turnIndex+1)

				if phaseId == r.currentPhaseId &&
					turnIndex <= r.currentTurnIndex {

					r.currentTurnIndex--

					// If current phase is empty, go back to the previous one
					for r.currentTurnIndex == -1 && !r.IsEmpty() {
						r.currentPhaseId--

						// Go back to previous round
						if r.currentPhaseId == 0 {
							r.currentPhaseId = enum.DuskPhaseId
							r.id--
						}

						r.currentTurnIndex = len(r.CurrentPhase()) - 1
					}

					// Reset if current turn index is still -1
					if r.currentTurnIndex == -1 {
						r.currentTurnIndex = 0
						r.currentPhaseId = 0
					}
				}

				return true
			}
		}
	}

	return false
}

func (r *round) isValidPhaseId(phaseId types.PhaseId) bool {
	return phaseId >= enum.NightPhaseId && phaseId <= enum.DuskPhaseId
}

func (r *round) AddTurn(setting *types.TurnSetting) bool {
	if !r.isValidPhaseId(setting.PhaseId) || r.existRole(setting.RoleId) {
		return false
	}

	phaseId, turnIndex := r.handleTurnSetting(setting)

	if turnIndex == -1 {
		return false
	}

	// New turn's position is less than or equal to current turn index,
	// so increase current turn index by 1
	if r.CurrentTurn() != nil &&
		turnIndex <= r.currentTurnIndex &&
		(setting.PhaseId == r.currentPhaseId ||
			setting.Position == enum.NextPosition) {

		r.currentTurnIndex++
	}

	r.phases[phaseId] = slices.Insert(
		r.phases[phaseId],
		turnIndex,
		&types.Turn{
			RoleId:     setting.RoleId,
			PlayerIds:  setting.PlayerIds,
			BeginRound: setting.BeginRound,
			Priority:   setting.Priority,
			Expiration: setting.Expiration,
		},
	)

	return true
}

func (r *round) existRole(roleId types.RoleId) bool {
	for _, phase := range r.phases {
		if slices.IndexFunc(phase, func(turn *types.Turn) bool {
			return turn.RoleId == roleId
		}) != -1 {
			return true
		}
	}

	return false
}

// Decide which phase and turn index contain new turn. Return -1 in second parameter
// if failed.
func (r *round) handleTurnSetting(setting *types.TurnSetting) (types.PhaseId, int) {
	turnIndex := -1
	phaseId := setting.PhaseId

	if setting.Position == enum.NextPosition {
		phaseId = r.currentPhaseId

		if len(r.CurrentPhase()) != 0 {
			turnIndex = r.currentTurnIndex + 1
		} else {
			turnIndex = 0
		}
	} else if setting.Position == enum.SortedPosition {
		turnIndex = slices.IndexFunc(r.phases[phaseId], func(turn *types.Turn) bool {
			return turn.Priority > setting.Priority
		})

		if turnIndex == -1 {
			turnIndex = len(r.phases[phaseId])
		}
	} else if setting.Position == enum.LastPosition {
		turnIndex = len(r.phases[phaseId])
	} else {
		if setting.Position < 0 || int(setting.Position) > len(r.phases[phaseId]) {
			return phaseId, -1
		}

		turnIndex = int(setting.Position)
	}

	return phaseId, turnIndex
}

func (r *round) AddPlayer(playerId types.PlayerId, roleId types.RoleId) bool {
	for _, phase := range r.phases {
		for _, turn := range phase {
			if turn.RoleId == roleId {
				if slices.Contains(turn.PlayerIds, playerId) {
					return false
				}

				turn.PlayerIds = append(turn.PlayerIds, playerId)

				return true
			}
		}
	}

	return false
}

func (r *round) DeletePlayer(playerId types.PlayerId, roleId types.RoleId) bool {
	for _, phase := range r.phases {
		for index, turn := range phase {
			if turn.RoleId == roleId {
				turn.PlayerIds = slices.Delete(
					turn.PlayerIds,
					index,
					index+1,
				)

				return true
			}
		}
	}

	return false
}

func (r *round) DeletePlayerFromAllTurns(playerId types.PlayerId) {
	for _, phase := range r.phases {
		for _, turn := range phase {
			deletedIndex := slices.Index(turn.PlayerIds, playerId)

			if deletedIndex != -1 {
				turn.PlayerIds = slices.Delete(
					turn.PlayerIds,
					deletedIndex,
					deletedIndex+1,
				)
			}
		}
	}
}
