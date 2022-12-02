package core

import (
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"

	"golang.org/x/exp/slices"
)

type scheduler struct {
	round        enum.Round
	beginPhaseID enum.PhaseID
	phaseID      enum.PhaseID
	turnIndex    int
	phases       map[enum.PhaseID][]*types.Turn
}

func NewScheduler(beginPhaseID enum.PhaseID) contract.Scheduler {
	return &scheduler{
		round:        enum.FirstRound,
		beginPhaseID: beginPhaseID,
		phaseID:      beginPhaseID,
		turnIndex:    -1,
		phases: map[enum.PhaseID][]*types.Turn{
			enum.DayPhaseID:   {},
			enum.NightPhaseID: {},
			enum.DuskPhaseID:  {},
		},
	}
}

func (s *scheduler) Round() enum.Round {
	return s.round
}

func (s *scheduler) PhaseID() enum.PhaseID {
	return s.phaseID
}

func (s *scheduler) Phase() []*types.Turn {
	return s.phases[s.phaseID]
}

func (s *scheduler) Turn() *types.Turn {
	if len(s.Phase()) == 0 ||
		s.turnIndex >= len(s.Phase()) ||
		s.turnIndex < 0 {
		return nil
	}

	return s.Phase()[s.turnIndex]
}

func (s *scheduler) IsEmpty(phaseID enum.PhaseID) bool {
	if !phaseID.IsUnknown() {
		return len(s.phases[phaseID]) == 0
	}

	for _, p := range s.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

func (s *scheduler) isValidPhaseID(phaseID enum.PhaseID) bool {
	return phaseID > enum.LowerPhaseID && phaseID < enum.UpperPhaseID
}

func (s *scheduler) existRole(roleID enum.RoleID) bool {
	for _, phase := range s.phases {
		if slices.IndexFunc(phase, func(turn *types.Turn) bool {
			return turn.RoleID == roleID
		}) != -1 {
			return true
		}
	}

	return false
}

// calculateTurnIndex decides which phase and turn index contain new turn.
// Return -1 in second parameter if failed.
func (s *scheduler) calculateTurnIndex(setting *types.TurnSetting) (enum.PhaseID, int) {
	turnIndex := -1
	phaseID := setting.PhaseID

	if setting.Position == enum.NextPosition {
		phaseID = s.phaseID

		if len(s.Phase()) != 0 {
			turnIndex = s.turnIndex + 1
		} else {
			// Become the first turn if the previous one doesn't exist
			turnIndex = 0
		}
	} else if setting.Position == enum.SortedPosition {
		turnIndex = slices.IndexFunc(s.phases[phaseID], func(turn *types.Turn) bool {
			return turn.Priority < setting.Priority
		})

		// Become the first turn if phase is empty or become the last turn
		// if all existed turns have higher priority, respectively
		if turnIndex == -1 {
			turnIndex = len(s.phases[phaseID])
		}
	} else if setting.Position == enum.LastPosition {
		turnIndex = len(s.phases[phaseID])
	} else {
		if setting.Position >= 0 && int(setting.Position) <= len(s.phases[phaseID]) {
			turnIndex = int(setting.Position)
		}
	}

	return phaseID, turnIndex
}

func (s *scheduler) AddTurn(setting *types.TurnSetting) bool {
	if !s.isValidPhaseID(setting.PhaseID) || s.existRole(setting.RoleID) {
		return false
	}

	phaseID, turnIndex := s.calculateTurnIndex(setting)

	if turnIndex == -1 {
		return false
	}

	// Increase current turn index by 1 if new turn's position
	// is less than or equal to current turn index
	if s.Turn() != nil && turnIndex <= s.turnIndex && phaseID == s.phaseID {
		s.turnIndex++
	}

	s.phases[phaseID] = slices.Insert(
		s.phases[phaseID],
		turnIndex,
		&types.Turn{
			RoleID:     setting.RoleID,
			BeginRound: setting.BeginRound,
			Priority:   setting.Priority,
		},
	)

	return true
}

func (s *scheduler) RemoveTurn(roleID enum.RoleID) bool {
	for phaseID, phase := range s.phases {
		for removedTurnIndex, turn := range phase {
			if turn.RoleID == roleID {
				s.phases[phaseID] = slices.Delete(
					s.phases[phaseID],
					removedTurnIndex,
					removedTurnIndex+1,
				)

				// Decrease current turn index by 1 if removed turn's position
				// is less than or equal to current turn index
				if phaseID == s.phaseID && removedTurnIndex <= s.turnIndex {
					s.turnIndex--

					// Move the current turn index to the previous phase's last turn
					// if the current phase is empty
					for s.turnIndex == -1 && !s.IsEmpty(enum.PhaseID(0)) {
						// Increase round if current phase is begin phase
						if s.phaseID == s.beginPhaseID {
							// No previous round to go back
							if s.round == 1 {
								break
							} else {
								s.round--
							}
						}

						s.phaseID = s.phaseID.PreviousPhase()
						s.turnIndex = len(s.Phase()) - 1
					}

					// Back to begin phase if current turn index is still -1
					if s.turnIndex == -1 {
						s.phaseID = s.beginPhaseID
					}
				}

				return true
			}
		}
	}

	return false
}

func (s *scheduler) defrostCurrentTurn() bool {
	turn := s.Turn()

	if turn.FrozenLimit != enum.ReachedLimit {
		if turn.FrozenLimit != enum.Unlimited {
			turn.FrozenLimit--
		}

		return true
	}

	return false
}

// Move to next turn and delete previous turn if it's times is out.
// Repeat from the beginning if the end is exceeded and return false
// if round is empty.
func (s *scheduler) NextTurn(isRemoved bool) bool {
	if s.IsEmpty(enum.PhaseID(0)) {
		return false
	}

	if isRemoved {
		s.RemoveTurn(s.Turn().RoleID)

		// Avoid next turn if the scheduler is empty after removing
		return s.NextTurn(false)
	}

	// Increase turn index by 1 if the new one does not
	// reach the length of the current phase
	if s.turnIndex < len(s.Phase())-1 {
		s.turnIndex++

		// Skip turn if not the time
		if s.Turn().BeginRound > s.round {
			return s.NextTurn(false)
		}
	} else {
		s.turnIndex = 0
		s.phaseID = s.phaseID.NextPhase()

		// Start new round
		if s.phaseID == s.beginPhaseID {
			s.round++
		}

		if s.Turn() == nil {
			return s.NextTurn(false)
		}
	}

	// Skip turn if it's frozen
	if s.defrostCurrentTurn() {
		return s.NextTurn(false)
	}

	return true
}

func (s *scheduler) FreezeTurn(roleID enum.RoleID, limit enum.Limit) bool {
	for _, phase := range s.phases {
		for _, turn := range phase {
			if turn.RoleID == roleID {
				turn.FrozenLimit = limit

				return true
			}
		}
	}

	return false
}
