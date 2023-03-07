package core

// import (
// 	"uwwolf/game/contract"
// 	"uwwolf/game/types"
// 	"uwwolf/game/vars"
// )

// // scheduler manages game's turns.
// type scheduler struct {
// 	// roundID is the current round ID.
// 	roundID types.RoundID

// 	// beginPhaseID is the first phase of round.
// 	beginPhaseID types.PhaseID

// 	// phaseID is the current phase ID.
// 	phaseID types.PhaseID

// 	// turnID is the current turn ID.
// 	turnID types.TurnID

// 	// phases contains the phases in their turns.
// 	phases map[types.PhaseID][]types.Turn
// }

// func NewScheduler(beginPhaseID types.PhaseID) contract.Scheduler {
// 	return &scheduler{
// 		roundID:      vars.FirstRound,
// 		beginPhaseID: beginPhaseID,
// 		phaseID:      beginPhaseID,
// 		turnID:       types.TurnID(-1), // Initial turn ID
// 		phases: map[types.PhaseID][]types.Turn{
// 			vars.NightPhaseID: make([]types.Turn, vars.PostTurn+1),
// 			vars.DuskPhaseID:  make([]types.Turn, vars.PostTurn+1),
// 			vars.DayPhaseID:   make([]types.Turn, vars.PostTurn+1),
// 		},
// 	}
// }

// var _ contract.Scheduler = (*scheduler)(nil)

// func (s scheduler) RoundID() types.RoundID {
// 	return s.roundID
// }

// func (s scheduler) PhaseID() types.PhaseID {
// 	return s.phaseID
// }

// func (s scheduler) Phase() []types.Turn {
// 	return s.phases[s.phaseID]
// }

// func (s scheduler) TurnID() types.TurnID {
// 	return s.turnID
// }

// func (s scheduler) Turn() types.Turn {
// 	if len(s.Phase()) == 0 {
// 		return nil
// 	}
// 	return s.Phase()[s.turnID]
// }

// func (s scheduler) PlayablePlayerIDs() []types.PlayerID {
// 	playerIDs := []types.PlayerID{}
// 	for playerID, playerTurn := range s.Turn() {
// 		if playerTurn.BeginRoundID <= s.roundID &&
// 			playerTurn.FrozenLimit == vars.ReachedLimit {
// 			playerIDs = append(playerIDs, playerID)
// 		}
// 	}
// 	return playerIDs
// }

// func (s scheduler) IsEmptyPhase(phaseID types.PhaseID) bool {
// 	if !phaseID.IsUnknown() {
// 		return len(s.phases[phaseID]) == 0
// 	}

// 	for _, p := range s.phases {
// 		if len(p) != 0 {
// 			return false
// 		}
// 	}
// 	return true
// }

// func (s *scheduler) AddPlayerTurn(newTurn types.NewPlayerTurn) bool {
// 	if phase, ok := s.phases[newTurn.PhaseID]; !ok {
// 		return false
// 	} else {
// 		phase[newTurn.TurnID][newTurn.PlayerID] = &types.PlayerTurn{
// 			BeginRoundID: newTurn.BeginRoundID,
// 			FrozenLimit:  vars.ReachedLimit,
// 			RoleID:       newTurn.RoleID,
// 		}

// 		return true
// 	}
// }

// func (s *scheduler) RemovePlayerTurn(removedTurn types.RemovedPlayerTurn) bool {
// 	if removedTurn.PhaseID.IsUnknown() {
// 		// Remove all player turns
// 		for _, phase := range s.phases {
// 			for _, turn := range phase {
// 				delete(turn, removedTurn.PlayerID)
// 			}
// 		}
// 	} else if removedTurn.TurnID > -1 &&
// 		int(removedTurn.TurnID) < len(s.phases[removedTurn.PhaseID]) {
// 		// Remove by turn ID
// 		delete(s.phases[removedTurn.PhaseID][removedTurn.TurnID], removedTurn.PlayerID)
// 	} else if !removedTurn.RoleID.IsUnknown() {
// 		// Remove by role ID
// 		for _, turn := range s.phases[removedTurn.PhaseID] {
// 			if turn[removedTurn.PlayerID].RoleID == removedTurn.RoleID {
// 				delete(turn, removedTurn.PlayerID)
// 				break
// 			}

// 		}
// 	} else {
// 		return false
// 	}

// 	return true
// }

// func (s *scheduler) NextTurn() bool {
// 	// Return false if schedule is empty
// 	if s.IsEmptyPhase(types.PhaseID(0)) {
// 		return false
// 	}

// 	if s.turnID < vars.PreTurn {
// 		s.turnID = 0
// 	} else {
// 		// Defrost player turns of the current turn
// 		for _, playerTurn := range s.Turn() {
// 			if playerTurn.FrozenLimit != vars.ReachedLimit {
// 				playerTurn.FrozenLimit--
// 			}
// 		}

// 		if int(s.turnID) < len(s.Phase())-1 {
// 			s.turnID++
// 		} else {
// 			s.turnID = 0
// 			s.phaseID = s.phaseID.NextPhasePhaseID(vars.DuskPhaseID)
// 			if s.phaseID == s.beginPhaseID {
// 				s.roundID++
// 			}
// 		}
// 	}

// 	// Move to the next turn if the current is empty
// 	isEmptyTurn := true
// 	for _, playerTurn := range s.Turn() {
// 		if playerTurn.BeginRoundID <= s.roundID {
// 			isEmptyTurn = false
// 			break
// 		}
// 	}
// 	if isEmptyTurn {
// 		return s.NextTurn()
// 	}

// 	return true
// }
