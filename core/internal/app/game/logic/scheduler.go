package logic

import (
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"
)

// scheduler manages game's turns.
type scheduler struct {
	// roundID is the current round ID.
	round types.Round

	// beginPhaseId is the first phase of round.
	beginPhaseId types.PhaseId

	// PhaseId is the current phase ID.
	phaseId types.PhaseId

	// turnID is the current turn ID.
	turn types.Turn

	// phases contains the phases in their turns.
	phases map[types.PhaseId]map[types.Turn]types.TurnRecords
}

func NewScheduler(beginPhaseId types.PhaseId) contract.Scheduler {
	return &scheduler{
		round:        constants.ZeroRound,
		beginPhaseId: beginPhaseId,
		phaseId:      beginPhaseId,
		turn:         constants.PreTurn,
		phases: map[types.PhaseId]map[types.Turn]types.TurnRecords{
			constants.NightPhaseId: {
				constants.PreTurn:  make(types.TurnRecords),
				constants.MidTurn:  make(types.TurnRecords),
				constants.PostTurn: make(types.TurnRecords),
			},
			constants.DuskPhaseId: {
				constants.PreTurn:  make(types.TurnRecords),
				constants.MidTurn:  make(types.TurnRecords),
				constants.PostTurn: make(types.TurnRecords),
			},
			constants.DayPhaseId: {
				constants.PreTurn:  make(types.TurnRecords),
				constants.MidTurn:  make(types.TurnRecords),
				constants.PostTurn: make(types.TurnRecords),
			},
		},
	}
}

// RoundID returns the latest round ID.
func (s scheduler) Round() types.Round {
	return s.round
}

// PhaseId returns the current phase ID.
func (s scheduler) PhaseId() types.PhaseId {
	return s.phaseId
}

// Phase returns the current phase.
func (s scheduler) Phase() map[types.Turn]types.TurnRecords {
	return s.phases[s.phaseId]
}

// TurnID returns the current turn ID.
func (s scheduler) Turn() types.Turn {
	return s.turn
}

// Turn returns the current turn.
func (s scheduler) TurnRecords() types.TurnRecords {
	if len(s.Phase()) == 0 {
		return nil
	}
	return s.Phase()[s.turn]
}

// CanPlay checks if the given playerID can play in the
// current turn.
func (s scheduler) CanPlay(playerID types.PlayerId) bool {
	slot := s.TurnRecords()[playerID]

	return slot != nil &&
		((!util.IsZero(slot.BeginRound) && slot.BeginRound <= s.round) ||
			(!util.IsZero(slot.PlayedRound) && slot.PlayedRound == s.round)) &&
		slot.FrozenTimes == constants.OutOfTimes
}

// PlayablePlayerIDs returns playable player ID list in
// the current turn.
func (s scheduler) PlayablePlayerIDs() []types.PlayerId {
	playerIDs := []types.PlayerId{}
	for playerID := range s.TurnRecords() {
		if s.CanPlay(playerID) {
			playerIDs = append(playerIDs, playerID)
		}
	}
	return playerIDs
}

// IsEmptyPhase check if specific phase is empty.
// Check if scheduler is empty if `PhaseId` is 0.
func (s scheduler) IsEmptyPhase(phaseId types.PhaseId) bool {
	if !util.IsZero(phaseId) {
		for _, turn := range s.phases[phaseId] {
			if len(turn) != 0 {
				return false
			}
		}
	} else {
		for _, phase := range s.phases {
			for _, turn := range phase {
				if len(turn) != 0 {
					return false
				}
			}
		}
	}

	return true
}

// AddSlot adds new player turn to the scheduler.
func (s *scheduler) AddSlot(newSlot *types.NewTurnSlot) bool {
	if phase, ok := s.phases[newSlot.PhaseId]; !ok {
		return false
	} else {
		phase[newSlot.Turn][newSlot.PlayerId] = &types.TurnSlot{
			BeginRound:  newSlot.BeginRound,
			FrozenTimes: newSlot.FrozenTimes,
			RoleId:      newSlot.RoleId,
			PlayedRound: newSlot.PlayedRound,
		}

		return true
	}
}

// RemoveSlot removes a player turn from the scheduler
// by `TurnID` or `RoleID`.
//
// If `TurnID` is provided, ignore `RoleID`.
//
// If `PhaseId` is 0, removes all of turns of that player.
func (s *scheduler) RemoveSlot(removedSlot *types.RemovedTurnSlot) bool {
	if util.IsZero(removedSlot.PhaseId) {
		// Remove all player turns
		for _, phase := range s.phases {
			for _, turn := range phase {
				delete(turn, removedSlot.PlayerId)
			}
		}
	} else if !util.IsZero(removedSlot.Turn) &&
		int(removedSlot.Turn) < len(s.phases[removedSlot.PhaseId]) {
		// Remove by turn ID
		delete(s.phases[removedSlot.PhaseId][removedSlot.Turn], removedSlot.PlayerId)
	} else if !util.IsZero(removedSlot.RoleId) {
		// Remove by role ID
		for _, turn := range s.phases[removedSlot.PhaseId] {
			if turn[removedSlot.PlayerId] != nil &&
				turn[removedSlot.PlayerId].RoleId == removedSlot.RoleId {
				delete(turn, removedSlot.PlayerId)
				break
			}
		}
	} else {
		return false
	}

	return true
}

// FreezeSlot blocks slot N times.
func (s *scheduler) FreezeSlot(frozenSlot *types.FreezeTurnSlot) bool {
	if !util.IsZero(frozenSlot.Turn) &&
		int(frozenSlot.Turn) < len(s.phases[frozenSlot.PhaseId]) {
		// Freeze by turn ID
		s.phases[frozenSlot.PhaseId][frozenSlot.Turn][frozenSlot.PlayerId].
			FrozenTimes = frozenSlot.FrozenTimes
	} else if !util.IsZero(frozenSlot.RoleId) {
		// Freeze by role ID
		for _, turn := range s.phases[frozenSlot.PhaseId] {
			if turn[frozenSlot.PlayerId] != nil &&
				turn[frozenSlot.PlayerId].RoleId == frozenSlot.RoleId {
				turn[frozenSlot.PlayerId].FrozenTimes = frozenSlot.FrozenTimes
				break
			}
		}
	} else {
		return false
	}

	return true
}

// NextTurn moves to the next turn.
// Returns false if the scheduler is empty.
func (s *scheduler) NextTurn() bool {
	// Return false if schedule is empty
	if s.IsEmptyPhase(types.PhaseId(0)) {
		return false
	}

	if util.IsZero(s.round) {
		s.round = constants.FirstRound
	} else {
		for playerID, slot := range s.TurnRecords() {
			// Defrost player slots of the current turn
			if slot.FrozenTimes != constants.OutOfTimes &&
				slot.FrozenTimes != constants.UnlimitedTimes {
				slot.FrozenTimes--
			}
			// Remove one-round slot
			if slot.PlayedRound == s.round {
				delete(s.TurnRecords(), playerID)
			}
		}

		if int(s.turn) < len(s.Phase())-1 {
			s.turn++
		} else {
			s.turn = constants.PreTurn
			s.phaseId = util.NextPhasePhaseID(s.PhaseId())
			if s.phaseId == s.beginPhaseId {
				s.round++
			}
		}
	}

	// Move to the next turn if the current is empty
	if len(s.TurnRecords()) == 0 {
		return s.NextTurn()
	}

	return true
}
