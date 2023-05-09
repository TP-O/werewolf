package tool

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/types"
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

// Scheduler manages game's turns.
type Scheduler interface {
	// RoundID returns the latest round ID.
	RoundID() types.RoundID

	// PhaseID returns the current phase ID.
	PhaseID() types.PhaseID

	// Phase returns the current phase.
	Phase() map[types.TurnID]types.Turn

	// TurnID returns the current turn ID.
	TurnID() types.TurnID

	// Turn returns the current turn.
	Turn() types.Turn

	// CanPlay checks if the given playerID can play in the
	// current turn.
	CanPlay(playerID types.PlayerID) bool

	// PlayablePlayerIDs returns playable player ID list in
	// the current turn.
	PlayablePlayerIDs() []types.PlayerID

	// IsEmptyPhase check if specific phase is empty.
	// Check if scheduler is empty if `phaseID` is 0.
	IsEmptyPhase(phaseID types.PhaseID) bool

	// AddSlot adds new player turn to the scheduler.
	AddSlot(newSlot *types.NewTurnSlot) bool

	// RemoveSlot removes a player turn from the scheduler
	// by `TurnID` or `RoleID`.
	//
	// If `TurnID` is filled, ignore `RoleID`.
	//
	// If `PhaseID` is 0, removes all of turns of that player.
	RemoveSlot(removedSlot *types.RemovedTurnSlot) bool

	// FreezeSlot blocks slot N times.
	FreezeSlot(frozenSlot *types.FreezeTurnSlot) bool

	// NextTurn moves to the next turn.
	// Returns false if the scheduler is empty.
	NextTurn() bool
}

func NewScheduler(beginPhaseID types.PhaseID) Scheduler {
	return &scheduler{
		roundID:      declare.ZeroRound,
		beginPhaseID: beginPhaseID,
		phaseID:      beginPhaseID,
		turnID:       declare.PreTurn,
		phases: map[types.PhaseID]map[types.TurnID]types.Turn{
			declare.NightPhaseID: {
				declare.PreTurn:  make(types.Turn),
				declare.MidTurn:  make(types.Turn),
				declare.PostTurn: make(types.Turn),
			},
			declare.DuskPhaseID: {
				declare.PreTurn:  make(types.Turn),
				declare.MidTurn:  make(types.Turn),
				declare.PostTurn: make(types.Turn),
			},
			declare.DayPhaseID: {
				declare.PreTurn:  make(types.Turn),
				declare.MidTurn:  make(types.Turn),
				declare.PostTurn: make(types.Turn),
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
		((!slot.BeginRoundID.IsZero() && slot.BeginRoundID <= s.roundID) ||
			(!slot.PlayedRoundID.IsZero() && slot.PlayedRoundID == s.roundID)) &&
		slot.FrozenTimes == declare.OutOfTimes
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
// If `TurnID` is provided, ignore `RoleID`.
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
	} else if !removedSlot.TurnID.IsZero() &&
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
	if !frozenSlot.TurnID.IsZero() &&
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

	if s.roundID.IsZero() {
		s.roundID = declare.FirstRound
	} else {
		for playerID, slot := range s.Turn() {
			// Defrost player slots of the current turn
			if slot.FrozenTimes != declare.OutOfTimes &&
				slot.FrozenTimes != declare.UnlimitedTimes {
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
			s.turnID = declare.PreTurn
			s.phaseID = s.phaseID.NextPhasePhaseID(declare.DuskPhaseID)
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
