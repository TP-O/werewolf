package state

import (
	"golang.org/x/exp/slices"

	"uwwolf/app/types"
)

type Round struct {
	id               types.RoundId
	currentPhaseId   types.PhaseId
	currentTurnIndex int
	phases           map[types.PhaseId][]*turn
}

type turn struct {
	roleId     types.RoleId
	playerIds  []types.PlayerId
	beginRound types.RoundId
	priority   int
	expiration types.NumberOfTimes
}

func (t turn) RoleId() types.RoleId {
	return t.roleId
}

func (t turn) PlayerIds() []types.PlayerId {
	return t.playerIds
}

func NewRound() *Round {
	round := &Round{}

	// Set default values
	round.Reset()

	return round
}

func (r *Round) CurrentId() types.RoundId {
	return r.id
}

func (r *Round) CurrentPhaseId() types.PhaseId {
	return r.currentPhaseId
}

func (r *Round) CurrentTurn() *turn {
	if len(r.CurrentPhase()) == 0 || r.currentTurnIndex >= len(r.CurrentPhase()) {
		return nil
	}

	return r.CurrentPhase()[r.currentTurnIndex]
}

func (r *Round) CurrentPhase() []*turn {
	return r.phases[r.currentPhaseId]
}

func (r *Round) IsAllowed(playerId types.PlayerId) bool {
	if r.CurrentTurn() == nil {
		return false
	}

	return slices.Contains(r.CurrentTurn().playerIds, playerId)
}

func (r *Round) IsEmpty() bool {
	for _, p := range r.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

func (r *Round) Reset() {
	r.id = 1
	r.currentTurnIndex = 0
	r.currentPhaseId = types.NightPhase
	r.phases = make(map[types.PhaseId][]*turn)
}

// Move to next turn and delete previous turn if it's times is out.
// Repeat from the beginning if the end is exceeded and return false
// if round is empty.
func (r *Round) NextTurn() bool {
	if r.IsEmpty() {
		return false
	}

	r.passCurrentTurn()

	if r.currentTurnIndex < len(r.CurrentPhase())-1 {
		r.currentTurnIndex++

		// Skip turn if not the time
		if r.CurrentTurn().beginRound > r.id {
			return r.NextTurn()
		}
	} else {
		r.currentTurnIndex = 0
		r.currentPhaseId = (r.currentPhaseId + 1) % (types.DuskPhase + 1)

		// Start new round
		if r.currentPhaseId == 0 {
			r.currentPhaseId = types.NightPhase
			r.id++
		}

		if r.CurrentTurn() == nil {
			return r.NextTurn()
		}
	}

	return true
}

func (r *Round) passCurrentTurn() {
	currentTurn := r.CurrentTurn()

	if currentTurn != nil &&
		currentTurn.beginRound <= r.id &&
		currentTurn.expiration != types.UnlimitedTimes {

		currentTurn.expiration--
		r.removeTurnIfDone(currentTurn.roleId)
	}
}

// Remove a turn if it's times is out.
func (r *Round) removeTurnIfDone(roleId types.RoleId) bool {
	currentTurn := r.CurrentTurn()

	if currentTurn == nil || currentTurn.expiration != types.OutOfTimes {
		return false
	}

	return r.RemoveTurn(roleId)
}

// Remove turn. In the case that removed turn is the same as the current
// turn index and phase containing the turn is also  the same as current
// phase, decreases the current turn index by 1.
func (r *Round) RemoveTurn(roleId types.RoleId) bool {
	for phaseId, phase := range r.phases {
		for turnIndex, turn := range phase {
			if turn.roleId == roleId {
				r.phases[phaseId] = slices.Delete(r.phases[phaseId], turnIndex, turnIndex+1)

				if phaseId == r.currentPhaseId &&
					turnIndex <= r.currentTurnIndex {

					r.currentTurnIndex--

					// If current phase is empty, go back to the previous one
					for r.currentTurnIndex == -1 && !r.IsEmpty() {
						r.currentPhaseId--

						// Go back to previous round
						if r.currentPhaseId == 0 {
							r.currentPhaseId = types.DuskPhase
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

func (r *Round) isValidPhaseId(phaseId types.PhaseId) bool {
	return phaseId >= types.NightPhase && phaseId <= types.DuskPhase
}

func (r *Round) AddTurn(setting *types.TurnSetting) bool {
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
			setting.Position == types.NextPosition) {

		r.currentTurnIndex++
	}

	r.phases[phaseId] = slices.Insert(
		r.phases[phaseId],
		turnIndex,
		&turn{
			roleId:     setting.RoleId,
			playerIds:  setting.PlayerIds,
			beginRound: setting.BeginRound,
			priority:   setting.Priority,
			expiration: setting.Expiration,
		},
	)

	return true
}

func (r *Round) existRole(roleId types.RoleId) bool {
	for _, phase := range r.phases {
		if slices.IndexFunc(phase, func(turn *turn) bool {
			return turn.roleId == roleId
		}) != -1 {
			return true
		}
	}

	return false
}

// Decide which phase and turn index contain new turn. Return -1 in second parameter
// if failed.
func (r *Round) handleTurnSetting(setting *types.TurnSetting) (types.PhaseId, int) {
	turnIndex := -1
	phaseId := setting.PhaseId

	if setting.Position == types.NextPosition {
		phaseId = r.currentPhaseId

		if len(r.CurrentPhase()) != 0 {
			turnIndex = r.currentTurnIndex + 1
		} else {
			turnIndex = 0
		}
	} else if setting.Position == types.SortedPosition {
		turnIndex = slices.IndexFunc(r.phases[phaseId], func(turn *turn) bool {
			return turn.priority > setting.Priority
		})

		if turnIndex == -1 {
			turnIndex = len(r.phases[phaseId])
		}
	} else if setting.Position == types.LastPosition {
		turnIndex = len(r.phases[phaseId])
	} else {
		if setting.Position < 0 || int(setting.Position) > len(r.phases[phaseId]) {
			return phaseId, -1
		}

		turnIndex = int(setting.Position)
	}

	return phaseId, turnIndex
}

func (r *Round) AddPlayer(playerId types.PlayerId, roleId types.RoleId) bool {
	for _, phase := range r.phases {
		for _, turn := range phase {
			if turn.roleId == roleId {
				if slices.Contains(turn.playerIds, playerId) {
					return false
				}

				turn.playerIds = append(turn.playerIds, playerId)

				return true
			}
		}
	}

	return false
}

func (r *Round) DeletePlayer(playerId types.PlayerId, roleId types.RoleId) bool {
	for _, phase := range r.phases {
		for index, turn := range phase {
			if turn.roleId == roleId {
				turn.playerIds = slices.Delete(
					turn.playerIds,
					index,
					index+1,
				)

				return true
			}
		}
	}

	return false
}

func (r *Round) DeletePlayerFromAllTurns(playerId types.PlayerId) {
	for _, phase := range r.phases {
		for _, turn := range phase {
			deletedIndex := slices.Index(turn.playerIds, playerId)

			if deletedIndex != -1 {
				turn.playerIds = slices.Delete(
					turn.playerIds,
					deletedIndex,
					deletedIndex+1,
				)
			}
		}
	}
}
