package state

import (
	"golang.org/x/exp/slices"

	"uwwolf/types"
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
	times      types.NumberOfTimes
}

func NewRound() *Round {
	return &Round{
		id:             1,
		currentPhaseId: types.NightPhase,
		phases:         make(map[types.PhaseId][]*turn),
	}
}

func (r *Round) GetId() types.RoundId {
	return r.id
}

func (r *Round) GetPhaseId() types.PhaseId {
	return r.currentPhaseId
}

func (r *Round) GetCurrentRoleId() types.RoleId {
	return r.phases[r.currentPhaseId][r.currentTurnIndex].roleId
}

func (r *Round) GetCurrentTurn() *turn {
	if len(r.phases[r.currentPhaseId]) == 0 {
		return nil
	}

	return r.phases[r.currentPhaseId][r.currentTurnIndex]
}

func (r *Round) GetCurrentPhase() []*turn {
	return r.phases[r.currentPhaseId]
}

func (r *Round) IsAllowed(playerId types.PlayerId) bool {
	return slices.Contains(r.GetCurrentTurn().playerIds, playerId)
}

func (r *Round) IsEmpty() bool {
	for _, p := range r.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

// Move to next turn and delete previous turn if it's times is out.
// Repeat from the beginning if the end is exceeded and return false
// if round is empty.
func (r *Round) NextTurn() bool {
	if r.IsEmpty() {
		return false
	}

	r.passCurrentTurn()

	if r.currentTurnIndex < len(r.GetCurrentPhase())-1 {
		r.currentTurnIndex++

		// Skip turn if not the time
		if r.GetCurrentTurn().beginRound < r.id {
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

		if r.GetCurrentTurn() == nil {
			return r.NextTurn()
		}
	}

	return true
}

func (r *Round) passCurrentTurn() {
	currentTurn := r.GetCurrentTurn()

	if currentTurn == nil &&
		currentTurn.beginRound <= r.id &&
		currentTurn.times != types.UnlimitedTimes {

		currentTurn.times--
		r.removeTurnIfDone(r.currentPhaseId, currentTurn.roleId)
	}
}

// Remove a turn if it's times is out.
func (r *Round) removeTurnIfDone(phaseId types.PhaseId, roleId types.RoleId) bool {
	currentTurn := r.GetCurrentTurn()

	if currentTurn == nil || currentTurn.times != types.OutOfTimes {
		return false
	}

	return r.RemoveTurn(phaseId, roleId)
}

// Remove turn. In the case that removed turn is the same as the current
// turn index and phase containing the turn is also  the same as current
// phase, decreases the current turn index by 1.
func (r *Round) RemoveTurn(phaseId types.PhaseId, roleId types.RoleId) bool {
	if !r.isValidPhaseId(phaseId) {
		return false
	}

	for index, turn := range r.phases[phaseId] {
		if turn.roleId == roleId {
			r.phases[phaseId] = slices.Delete(r.phases[phaseId], index, index+1)

			if phaseId == r.currentPhaseId &&
				index <= r.currentTurnIndex {

				r.currentTurnIndex--

				// If current phase is empty, go back to the previous one
				for r.currentTurnIndex == -1 && !r.IsEmpty() {
					r.currentPhaseId--

					// Go back to previous round
					if r.currentPhaseId == 0 {
						r.currentPhaseId = types.DuskPhase
						r.id--
					}

					r.currentTurnIndex = len(r.GetCurrentPhase()) - 1
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

	return false
}

func (r *Round) isValidPhaseId(phaseId types.PhaseId) bool {
	return phaseId >= types.NightPhase && phaseId <= types.DuskPhase
}

func (r *Round) AddTurn(setting *types.TurnSetting) bool {
	if !r.isValidPhaseId(setting.PhaseId) ||
		int(setting.Position) > len(r.phases[setting.PhaseId])-1 {

		return false
	}

	var turnIndex int

	if setting.Position == types.NextPosition {
		turnIndex = r.currentTurnIndex + 1
	} else if setting.Position == types.SortedPosition {
		turnIndex = slices.IndexFunc(r.phases[setting.PhaseId], func(turn *turn) bool {
			return turn.priority > setting.Priority
		})

		if turnIndex == -1 {
			turnIndex = len(r.phases[setting.PhaseId])
		}
	} else if setting.Position == types.LastPosition {
		turnIndex = len(r.phases[setting.PhaseId])
	} else {
		turnIndex = int(setting.Position)
	}

	// Check valid turn index
	if turnIndex < 0 {
		return false
	} else if turnIndex <= r.currentTurnIndex {
		// New turn's position is less than or equal to current turn index,
		// so increase current turn index by 1
		r.currentTurnIndex++
	}

	r.phases[r.currentPhaseId] = slices.Insert(
		r.phases[r.currentPhaseId],
		turnIndex,
		&turn{
			roleId:     setting.RoleId,
			playerIds:  setting.PlayerIds,
			beginRound: setting.BeginRound,
			priority:   setting.Priority,
			times:      setting.Times,
		},
	)

	return true
}

func (r *Round) AddPlayer(playerId types.PlayerId, roleId types.RoleId) bool {
	for _, phase := range r.phases {
		for _, turn := range phase {
			if turn.roleId == roleId {
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
