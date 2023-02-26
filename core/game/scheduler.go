package game

import (
	"uwwolf/game/contract"
	"uwwolf/game/role"
	"uwwolf/game/types"

	"golang.org/x/exp/slices"
)

// scheduler manages the game's turn.
type scheduler struct {
	roundID      types.RoundID
	beginPhaseID types.PhaseID
	phaseID      types.PhaseID
	turnID       types.TurnID
	phases       map[types.PhaseID][]types.Turn
}

func NewScheduler(beginPhaseID types.PhaseID) contract.Scheduler {
	return &scheduler{
		roundID:      types.RoundID(0),
		beginPhaseID: beginPhaseID,
		phaseID:      beginPhaseID,
		turnID:       role.PreTurn,
		phases: map[types.PhaseID][]types.Turn{
			role.DayPhaseID:   {},
			role.NightPhaseID: {},
			role.DuskPhaseID:  {},
		},
	}
}

func (s *scheduler) RoundID() types.RoundID {
	return s.round
}

func (s *scheduler) PhaseID() types.PhaseID {
	return s.phaseID
}

func (s *scheduler) Phase() []*types.Turn {
	return s.phases[s.phaseID]
}

func (s *scheduler) Turn() types.Turn {
	if len(s.Phase()) == 0 ||
		s.turnIndex >= len(s.Phase()) ||
		s.turnIndex < 0 {
		return nil
	}

	return s.Phase()[s.turnIndex]
}

func (s *scheduler) IsEmpty(phaseID types.PhaseID) bool {
	if !types.IsUnknownPhaseID(phaseID) {
		return len(s.phases[phaseID]) == 0
	}

	for _, p := range s.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

func (s *scheduler) isValidPhaseID(phaseID types.PhaseID) bool {
	return phaseID >= NightPhaseID && phaseID <= DayPhaseID
}

func (s *scheduler) existRole(roleID types.RoleID) bool {
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
func (s *scheduler) calculateTurnIndex(setting *types.TurnSetting) (types.PhaseID, int) {
	turnIndex := -1
	phaseID := setting.PhaseID

	if setting.Position == NextPosition {
		phaseID = s.phaseID

		if len(s.Phase()) != 0 {
			turnIndex = s.turnIndex + 1
		} else {
			// Become the first turn if the previous one doesn't exist
			turnIndex = 0
		}
	} else if setting.Position == SortedPosition {
		turnIndex = slices.IndexFunc(s.phases[phaseID], func(turn *types.Turn) bool {
			// return turn.Priority < setting.Priority
			return true
		})

		// Become the first turn if phase is empty or become the last turn
		// if all existed turns have higher priority, respectively
		if turnIndex == -1 {
			turnIndex = len(s.phases[phaseID])
		}
	} else if setting.Position == LastPosition {
		turnIndex = len(s.phases[phaseID])
	} else {
		if setting.Position >= 0 && int(setting.Position) <= len(s.phases[phaseID]) {
			turnIndex = int(setting.Position)
		}
	}

	return phaseID, turnIndex
}

func (s *scheduler) AddTurn(setting types.TurnSetting) bool {
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
			Limit:      setting.Limit,
		},
	)

	return true
}

func (s *scheduler) RemoveTurn(roleID types.RoleID) bool {
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
					for s.turnIndex == -1 && !s.IsEmpty(types.PhaseID(0)) {
						// Increase round if current phase is begin phase
						if s.phaseID == s.beginPhaseID {
							// No previous round to go back
							if s.round == 1 {
								break
							} else {
								s.round--
							}
						}

						s.phaseID = types.PreviousPhaseID(s.phaseID, 1)
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

	if turn.FrozenLimit != ReachedLimit {
		if turn.FrozenLimit != Unlimited {
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
	if s.IsEmpty(types.PhaseID(0)) {
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
	} else {
		s.turnIndex = 0
		s.phaseID = types.NextPhasePhaseID(s.phaseID, 1)

		// Start new round
		if s.phaseID == s.beginPhaseID {
			s.round++
		}

		if s.Turn() == nil {
			return s.NextTurn(false)
		}
	}

	// Skip turn if not the time
	if s.Turn().BeginRound > s.round || s.Turn().Limit == ReachedLimit {
		return s.NextTurn(false)
	}

	// Skip turn if it's frozen
	if s.defrostCurrentTurn() {
		return s.NextTurn(false)
	}

	if s.Turn().Limit != Unlimited {
		s.Turn().Limit--
	}

	return true
}

func (s *scheduler) FreezeTurn(roleID types.RoleID, limit types.Limit) bool {
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
