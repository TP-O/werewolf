package game

import (
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

// scheduler manages game's turns.
type scheduler struct {
	// roundID is the current round ID.
	roundID types.RoundID

	// beginPhaseID is the first phase of round.
	beginPhaseID types.PhaseID

	// phaseID is the current phase ID.
	phaseID types.PhaseID

	// turnID is the current turn ID.
	turnID types.TurnID

	// phases contains the phases in their turns.
	phases map[types.PhaseID]map[types.TurnID]types.Turn
}

func NewScheduler(beginPhaseID types.PhaseID) contract.Scheduler {
	return &scheduler{
		roundID:      vars.ZeroRound,
		beginPhaseID: beginPhaseID,
		phaseID:      beginPhaseID,
		turnID:       vars.PreTurn,
		phases: map[types.PhaseID]map[types.TurnID]types.Turn{
			vars.NightPhaseID: {
				vars.PreTurn:  make(types.Turn),
				vars.MidTurn:  make(types.Turn),
				vars.PostTurn: make(types.Turn),
			},
			vars.DuskPhaseID: {
				vars.PreTurn:  make(types.Turn),
				vars.MidTurn:  make(types.Turn),
				vars.PostTurn: make(types.Turn),
			},
			vars.DayPhaseID: {
				vars.PreTurn:  make(types.Turn),
				vars.MidTurn:  make(types.Turn),
				vars.PostTurn: make(types.Turn),
			},
		},
	}
}

// RoundID returns the latest round ID.
func (s scheduler) RoundID() types.RoundID {
	return s.roundID
}

// PhaseID returns the current phase ID.
func (s scheduler) PhaseID() types.PhaseID {
	return s.phaseID
}

// Phase returns the current phase.
func (s scheduler) Phase() map[types.TurnID]types.Turn {
	return s.phases[s.phaseID]
}

// TurnID returns the current turn ID.
func (s scheduler) TurnID() types.TurnID {
	return s.turnID
}

// Turn returns the current turn.
func (s scheduler) Turn() types.Turn {
	if len(s.Phase()) == 0 {
		return nil
	}
	return s.Phase()[s.turnID]
}

// CanPlay checks if the given playerID can play in the
// current turn.
func (s scheduler) CanPlay(playerID types.PlayerID) bool {
	slot := s.Turn()[playerID]

	return slot != nil &&
		((slot.BeginRoundID != vars.ZeroRound && slot.BeginRoundID <= s.roundID) ||
			(slot.PlayedRoundID != vars.ZeroRound && slot.PlayedRoundID == s.roundID)) &&
		slot.FrozenTimes == vars.OutOfTimes
}

// PlayablePlayerIDs returns playable player ID list in
// the current turn.
func (s scheduler) PlayablePlayerIDs() []types.PlayerID {
	playerIDs := []types.PlayerID{}
	for playerID := range s.Turn() {
		if s.CanPlay(playerID) {
			playerIDs = append(playerIDs, playerID)
		}
	}
	return playerIDs
}

// IsEmptyPhase check if specific phase is empty.
// Check if scheduler is empty if `phaseID` is 0.
func (s scheduler) IsEmptyPhase(phaseID types.PhaseID) bool {
	if !phaseID.IsUnknown() {
		for _, turn := range s.phases[phaseID] {
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
	if phase, ok := s.phases[newSlot.PhaseID]; !ok {
		return false
	} else {
		phase[newSlot.TurnID][newSlot.PlayerID] = &types.TurnSlot{
			BeginRoundID:  newSlot.BeginRoundID,
			FrozenTimes:   newSlot.FrozenTimes,
			RoleID:        newSlot.RoleID,
			PlayedRoundID: newSlot.PlayedRoundID,
		}

		return true
	}
}

// RemoveSlot removes a player turn from the scheduler
// by `TurnID` or `RoleID`.
//
// If `TurnID` is filled, ignore `RoleID`.
//
// If `PhaseID` is 0, removes all of turns of that player.
func (s *scheduler) RemoveSlot(removedSlot *types.RemovedTurnSlot) bool {
	if removedSlot.PhaseID.IsUnknown() {
		// Remove all player turns
		for _, phase := range s.phases {
			for _, turn := range phase {
				delete(turn, removedSlot.PlayerID)
			}
		}
	} else if removedSlot.TurnID != vars.ZeroTurn &&
		int(removedSlot.TurnID) < len(s.phases[removedSlot.PhaseID]) {
		// Remove by turn ID
		delete(s.phases[removedSlot.PhaseID][removedSlot.TurnID], removedSlot.PlayerID)
	} else if !removedSlot.RoleID.IsUnknown() {
		// Remove by role ID
		for _, turn := range s.phases[removedSlot.PhaseID] {
			if turn[removedSlot.PlayerID] != nil &&
				turn[removedSlot.PlayerID].RoleID == removedSlot.RoleID {
				delete(turn, removedSlot.PlayerID)
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
	if frozenSlot.TurnID != vars.ZeroTurn &&
		int(frozenSlot.TurnID) < len(s.phases[frozenSlot.PhaseID]) {
		// Freeze by turn ID
		s.phases[frozenSlot.PhaseID][frozenSlot.TurnID][frozenSlot.PlayerID].
			FrozenTimes = frozenSlot.FrozenTimes
	} else if !frozenSlot.RoleID.IsUnknown() {
		// Freeze by role ID
		for _, turn := range s.phases[frozenSlot.PhaseID] {
			if turn[frozenSlot.PlayerID] != nil &&
				turn[frozenSlot.PlayerID].RoleID == frozenSlot.RoleID {
				turn[frozenSlot.PlayerID].FrozenTimes = frozenSlot.FrozenTimes
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
	if s.IsEmptyPhase(types.PhaseID(0)) {
		return false
	}

	if s.roundID == vars.ZeroRound {
		s.roundID = vars.FirstRound
	} else {
		for playerID, slot := range s.Turn() {
			// Defrost player slots of the current turn
			if slot.FrozenTimes != vars.OutOfTimes &&
				slot.FrozenTimes != vars.UnlimitedTimes {
				slot.FrozenTimes--
			}
			// Remove one-round slot
			if slot.PlayedRoundID == s.roundID {
				delete(s.Turn(), playerID)
			}
		}

		if int(s.turnID) < len(s.Phase())-1 {
			s.turnID++
		} else {
			s.turnID = vars.PreTurn
			s.phaseID = s.phaseID.NextPhasePhaseID(vars.DuskPhaseID)
			if s.phaseID == s.beginPhaseID {
				s.roundID++
			}
		}
	}

	// Move to the next turn if the current is empty
	if len(s.Turn()) == 0 {
		return s.NextTurn()
	}

	return true
}
