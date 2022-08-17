package core

import (
	"golang.org/x/exp/slices"

	"uwwolf/types"
)

type turn struct {
	roleId         types.RoleId
	playerIds      []types.PlayerId
	remainingTimes types.NumberOfTimes
}

type round struct {
	id               types.RoundId
	firstPhaseId     types.PhaseId
	currentPhaseId   types.PhaseId
	currentTurnIndex int
	phases           map[types.PhaseId][]*turn
}

func NewRound() *round {
	return &round{
		id:             1,
		firstPhaseId:   types.NightPhase,
		currentPhaseId: types.NightPhase,
		phases:         make(map[types.PhaseId][]*turn),
	}
}

// Get current round id.
func (r *round) GetId() types.RoundId {
	return r.id
}

// Get curretn phase id.
func (r *round) GetPhaseId() types.PhaseId {
	return r.currentPhaseId
}

// Get current turn index.
func (r *round) GetTurnIndex() int {
	return r.currentTurnIndex
}

// Check if phases are empty
func (r *round) IsEmpty() bool {
	for _, p := range r.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

// Get the current turn.
func (r *round) GetTurn() *turn {
	return r.phases[r.currentPhaseId][r.currentTurnIndex]
}

// Move to next turn and delete previous turn if it's temporary.
// Repeat from the start if the end is exceeded and return false
// if phases are empty.
func (r *round) NextTurn() bool {
	if r.IsEmpty() {
		return false
	}

	r.phases[r.currentPhaseId][r.currentTurnIndex].remainingTimes--
	r.removeTurnIfDone(r.currentPhaseId, r.currentTurnIndex)

	if r.currentTurnIndex < len(r.phases[r.currentPhaseId])-1 {
		r.currentTurnIndex++
	} else {
		r.currentTurnIndex = 0
		r.currentPhaseId = (r.currentPhaseId + 1) % (types.NightPhase + 1)

		if r.currentPhaseId == 0 {
			r.currentPhaseId = 1
		}

		if r.currentPhaseId == r.firstPhaseId {
			r.id++
		}
	}

	if len(r.phases[r.currentPhaseId]) == 0 {
		return r.NextTurn()
	}

	return true
}

// Remove a turn if it's remaining times is done. In the case that removed
// turn is the same as current turn index and phase containing the turn is
// also the same as current phase, decreases the current turn index by 1.
func (r *round) removeTurnIfDone(phaseId types.PhaseId, turnId int) bool {
	if len(r.phases[phaseId]) <= turnId &&
		r.phases[phaseId][turnId].remainingTimes != types.OutOfTimes {

		return false
	}

	r.phases[phaseId] = slices.Delete(r.phases[phaseId], turnId, turnId+1)

	if phaseId == r.currentPhaseId &&
		turnId == r.currentTurnIndex &&
		r.currentTurnIndex != 0 {

		r.currentTurnIndex--
	}

	return false
}

// Add a new turn.
func (r *round) AddTurn(data *types.TurnData) bool {
	if !r.isValidPhaseId(data.PhaseId) ||
		len(r.phases[data.PhaseId]) < int(data.Position) {

		return false
	}

	if data.Position == types.NextTurnPosition {
		r.phases[r.currentPhaseId] = slices.Insert(
			r.phases[r.currentPhaseId],
			r.currentTurnIndex+1,
			&turn{
				roleId:         data.RoleId,
				playerIds:      data.PlayerIds,
				remainingTimes: data.Times,
			},
		)
	} else if data.Position == types.LastTurnPosition {
		r.phases[data.PhaseId] = append(r.phases[data.PhaseId], &turn{
			roleId:         data.RoleId,
			playerIds:      data.PlayerIds,
			remainingTimes: types.UnlimitedTimes,
		})
	} else {
		r.phases[data.PhaseId] = slices.Insert(
			r.phases[data.PhaseId],
			int(data.Position),
			&turn{
				roleId:         data.RoleId,
				playerIds:      data.PlayerIds,
				remainingTimes: data.Times,
			},
		)
	}

	return true
}

// Check if phase id is valid.
func (r *round) isValidPhaseId(phaseId types.PhaseId) bool {
	return phaseId >= types.NightPhase && phaseId <= types.DuskPhase
}

// Remove a turn by role id.
func (r *round) RemoveTurn(phaseId types.PhaseId, roleId types.RoleId) bool {
	if !r.isValidPhaseId(phaseId) {
		return false
	}

	for index, turn := range r.phases[phaseId] {
		if turn.roleId == roleId {
			r.phases[phaseId] = slices.Delete(r.phases[phaseId], index, index+1)

			return true
		}
	}

	return false
}

// Check if player is in current turn.
func (r *round) IsValidPlayer(playerId types.PlayerId) bool {
	return slices.Contains(r.GetTurn().playerIds, playerId)
}

// Add players to sepicific turn.
func (r *round) AddPlayers(phaseId types.PhaseId, turnId int, playerIds ...types.PlayerId) bool {
	if !r.isValidPhaseId(phaseId) || r.phases[phaseId][turnId] == nil {
		return false
	}

	r.phases[phaseId][turnId].playerIds = append(
		r.phases[phaseId][turnId].playerIds,
		playerIds...,
	)

	return true
}

// Remove player from specific turn.
func (r *round) RemovePlayer(phaseId types.PhaseId, turnId int, playerId types.PhaseId) bool {
	if !r.isValidPhaseId(phaseId) || r.phases[phaseId][turnId] == nil {
		return false
	}

	r.phases[phaseId][turnId].playerIds = slices.Delete(
		r.phases[phaseId][turnId].playerIds,
		turnId,
		turnId+1,
	)

	return true
}

// Remove player from all turns.
func (r *round) RemovePlayerFromAllTurns(playerId types.PlayerId) {
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
