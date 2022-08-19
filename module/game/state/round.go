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
	roleId    types.RoleId
	playerIds []types.PlayerId
	times     types.NumberOfTimes
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

	if currentTurn == nil && currentTurn.times != types.UnlimitedTimes {
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
				index == r.currentTurnIndex {

				r.currentTurnIndex--

				// If current phase is empty, go back to the previous one
				for r.currentTurnIndex == -1 && !r.IsEmpty() {
					r.currentPhaseId = (r.currentPhaseId - 1) % (types.DuskPhase + 1)

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

// func (r *Round) AddTurn(setting *types.TurnSetting) bool {
// 	if !r.isValidPhaseId(setting.PhaseId) ||
// 		len(r.phases[setting.PhaseId]) < int(setting.Position) {

// 		return false
// 	}

// 	newTurn := &turn{
// 		roleId:    setting.RoleId,
// 		playerIds: setting.PlayerIds,
// 		times:     setting.Times,
// 	}

// 	if setting.Position == types.NextTurnPosition {
// 		r.phases[r.currentPhaseId] = slices.Insert(
// 			r.phases[r.currentPhaseId],
// 			r.currentTurnIndex+1,
// 			newTurn,
// 		)
// 	} else if setting.Position == types.LastTurnPosition {
// 		r.phases[setting.PhaseId] = append(r.phases[setting.PhaseId], newTurn)
// 	} else {
// 		r.phases[setting.PhaseId] = slices.Insert(
// 			r.phases[setting.PhaseId],
// 			int(setting.Position),
// 			newTurn,
// 		)
// 	}

// 	return true
// }

// // Add players to sepicific turn.
// func (r *Round) AddPlayers(phaseId types.PhaseId, turnId int, playerIds ...types.PlayerId) bool {
// 	if !r.isValidPhaseId(phaseId) || r.phases[phaseId][turnId] == nil {
// 		return false
// 	}

// 	r.phases[phaseId][turnId].playerIds = append(
// 		r.phases[phaseId][turnId].playerIds,
// 		playerIds...,
// 	)

// 	return true
// }

// // Remove player from specific turn.
// func (r *Round) RemovePlayer(phaseId types.PhaseId, turnId int, playerId types.PhaseId) bool {
// 	if !r.isValidPhaseId(phaseId) || r.phases[phaseId][turnId] == nil {
// 		return false
// 	}

// 	r.phases[phaseId][turnId].playerIds = slices.Delete(
// 		r.phases[phaseId][turnId].playerIds,
// 		turnId,
// 		turnId+1,
// 	)

// 	return true
// }

// func (r *Round) RemovePlayerFromAllTurns(playerId types.PlayerId) {
// 	for _, phase := range r.phases {
// 		for _, turn := range phase {
// 			deletedIndex := slices.Index(turn.playerIds, playerId)

// 			if deletedIndex != -1 {
// 				turn.playerIds = slices.Delete(
// 					turn.playerIds,
// 					deletedIndex,
// 					deletedIndex+1,
// 				)
// 			}
// 		}
// 	}
// }
