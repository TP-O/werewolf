package tool

// import (
// 	"testing"
// 	"uwwolf/internal/app/game/logic/types"
// 	"uwwolf/game/vars"

// 	"github.com/stretchr/testify/suite"
// )

// type SchudulerSuite struct {
// 	suite.Suite
// 	beginPhaseID types.PhaseID
// 	player1ID    types.PlayerID
// 	player2ID    types.PlayerID
// 	player3ID    types.PlayerID
// }

// func TestSchudulerSuite(t *testing.T) {
// 	suite.Run(t, new(SchudulerSuite))
// }

// func (ss *SchudulerSuite) SetupSuite() {
// 	ss.beginPhaseID = vars.NightPhaseID
// 	ss.player1ID = types.PlayerID("1")
// 	ss.player2ID = types.PlayerID("2")
// 	ss.player3ID = types.PlayerID("3")
// }

// func (ss SchudulerSuite) TestNewScheduler() {
// 	s := NewScheduler(ss.beginPhaseID)

// 	ss.NotNil(s)
// 	ss.Equal(vars.ZeroRound, s.(*scheduler).roundID)
// 	ss.Equal(vars.PreTurn, s.(*scheduler).turnID)
// 	ss.Equal(ss.beginPhaseID, s.(*scheduler).beginPhaseID)
// 	ss.Equal(ss.beginPhaseID, s.(*scheduler).phaseID)
// 	ss.NotNil(s.(*scheduler).phases)
// 	ss.Len(s.(*scheduler).phases, 3)
// 	for _, phase := range s.(*scheduler).phases {
// 		ss.Len(phase, 3)
// 	}
// }

// func (ss SchudulerSuite) TestRoundID() {
// 	s := NewScheduler(ss.beginPhaseID)

// 	expectedRoundID := vars.SecondRound
// 	s.(*scheduler).roundID = expectedRoundID

// 	ss.Equal(expectedRoundID, s.RoundID())
// }

// func (ss SchudulerSuite) TestPhaseID() {
// 	s := NewScheduler(ss.beginPhaseID)

// 	ss.Equal(ss.beginPhaseID, s.PhaseID())
// }

// func (ss SchudulerSuite) TestPhase() {
// 	s := NewScheduler(ss.beginPhaseID)

// 	expectedTurn := types.Turn(
// 		map[types.PlayerID]*types.TurnSlot{
// 			ss.player1ID: {
// 				RoleID: types.RoleID(99),
// 			},
// 		})
// 	s.(*scheduler).phases[ss.beginPhaseID][vars.PreTurn] = expectedTurn

// 	ss.Equal(expectedTurn, s.Phase()[vars.PreTurn])
// }

// func (ss SchudulerSuite) TestTurnID() {
// 	s := NewScheduler(ss.beginPhaseID)
// 	s.(*scheduler).turnID = vars.PostTurn

// 	ss.Equal(vars.PostTurn, s.TurnID())
// }

// func (ss SchudulerSuite) TestTurn() {
// 	tests := []struct {
// 		name         string
// 		expectedTurn types.Turn
// 		setup        func(*scheduler)
// 	}{
// 		{
// 			name:         "Nil (Current phase is empty)",
// 			expectedTurn: nil,
// 			setup: func(s *scheduler) {
// 				s.phases[ss.beginPhaseID] = make(map[types.TurnID]types.Turn)
// 			},
// 		},
// 		{
// 			name: "Ok",
// 			expectedTurn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						RoleID: types.RoleID(99),
// 					},
// 				}),
// 			setup: func(s *scheduler) {
// 				s.phases[ss.beginPhaseID][vars.PreTurn] = types.Turn(
// 					map[types.PlayerID]*types.TurnSlot{
// 						ss.player1ID: {
// 							RoleID: types.RoleID(99),
// 						},
// 					})
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)
// 			test.setup(s.(*scheduler))

// 			turn := s.Turn()

// 			ss.Equal(test.expectedTurn, turn)
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestCanPlay() {
// 	tests := []struct {
// 		name           string
// 		turn           types.Turn
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 	}{
// 		{
// 			name: "False (Late begin round slot)",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.SecondRound,
// 					},
// 				}),
// 			expectedStatus: false,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.FirstRound
// 			},
// 		}, {
// 			name: "False (zero begin round slot)",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.ZeroRound,
// 					},
// 				}),
// 			expectedStatus: false,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		},
// 		{
// 			name: "True (One-round slot)",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						PlayedRoundID: vars.SecondRound,
// 					},
// 				}),
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		},
// 		{
// 			name: "False (frozen slot)",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.SecondRound,
// 						FrozenTimes:  vars.Once,
// 					},
// 				}),
// 			expectedStatus: false,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)

// 			test.setup(s.(*scheduler))
// 			s.(*scheduler).phases[ss.beginPhaseID][vars.PreTurn] = test.turn

// 			ss.Equal(test.expectedStatus, s.CanPlay(ss.player1ID))
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestPlayablePlayerIDs() {
// 	tests := []struct {
// 		name              string
// 		turn              types.Turn
// 		expectedPlayerIDs []types.PlayerID
// 		setup             func(*scheduler)
// 	}{
// 		{
// 			name: "Ignore late begin round slot",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.FirstRound,
// 					},
// 					ss.player2ID: {
// 						BeginRoundID: vars.SecondRound,
// 					},
// 					ss.player3ID: {
// 						BeginRoundID: types.RoundID(3),
// 					},
// 				}),
// 			expectedPlayerIDs: []types.PlayerID{
// 				ss.player1ID,
// 				ss.player2ID,
// 			},
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		}, {
// 			name: "Ignore zero begin round slot",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.FirstRound,
// 					},
// 					ss.player2ID: {
// 						BeginRoundID: vars.SecondRound,
// 					},
// 					ss.player3ID: {
// 						BeginRoundID: vars.ZeroRound,
// 					},
// 				}),
// 			expectedPlayerIDs: []types.PlayerID{
// 				ss.player1ID,
// 				ss.player2ID,
// 			},
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		},
// 		{
// 			name: "Include one-round slot",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.FirstRound,
// 					},
// 					ss.player2ID: {
// 						PlayedRoundID: vars.SecondRound,
// 					},
// 					ss.player3ID: {
// 						PlayedRoundID: types.RoundID(3),
// 					},
// 				}),
// 			expectedPlayerIDs: []types.PlayerID{
// 				ss.player1ID,
// 				ss.player2ID,
// 			},
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		},
// 		{
// 			name: "Ignore frozen slot",
// 			turn: types.Turn(
// 				map[types.PlayerID]*types.TurnSlot{
// 					ss.player1ID: {
// 						BeginRoundID: vars.FirstRound,
// 					},
// 					ss.player2ID: {
// 						BeginRoundID: vars.SecondRound,
// 					},
// 					ss.player3ID: {
// 						BeginRoundID: vars.SecondRound,
// 						FrozenTimes:  vars.Once,
// 					},
// 				}),
// 			expectedPlayerIDs: []types.PlayerID{
// 				ss.player1ID,
// 				ss.player2ID,
// 			},
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.SecondRound
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)

// 			test.setup(s.(*scheduler))
// 			s.(*scheduler).phases[ss.beginPhaseID][vars.PreTurn] = test.turn

// 			ss.ElementsMatch(test.expectedPlayerIDs, s.PlayablePlayerIDs())
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestIsEmptyPhase() {
// 	beginPhaseID := vars.NightPhaseID
// 	tests := []struct {
// 		name           string
// 		phaseID        types.PhaseID
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 	}{

// 		{
// 			name:           "Non-empty (Check specific phase)",
// 			phaseID:        ss.beginPhaseID,
// 			expectedStatus: false,
// 			setup: func(s *scheduler) {
// 				s.phases[ss.beginPhaseID][vars.PreTurn] = types.Turn(
// 					map[types.PlayerID]*types.TurnSlot{
// 						ss.player1ID: {},
// 					},
// 				)
// 			},
// 		},
// 		{
// 			name:           "Empty (Check specific phase)",
// 			phaseID:        ss.beginPhaseID,
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[ss.beginPhaseID] = make(map[types.TurnID]types.Turn)
// 			},
// 		},
// 		{
// 			name:           "Non-empty (Check all phases)",
// 			phaseID:        types.PhaseID(0),
// 			expectedStatus: false,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.NightPhaseID] = make(map[types.TurnID]types.Turn)
// 				s.phases[vars.DayPhaseID][vars.PreTurn] = types.Turn(
// 					map[types.PlayerID]*types.TurnSlot{
// 						ss.player1ID: {},
// 					},
// 				)
// 				s.phases[vars.DuskPhaseID] = make(map[types.TurnID]types.Turn)
// 			},
// 		},
// 		{
// 			name:           "Empty (Check all phases)",
// 			phaseID:        types.PhaseID(0),
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.NightPhaseID] = make(map[types.TurnID]types.Turn)
// 				s.phases[vars.DayPhaseID] = make(map[types.TurnID]types.Turn)
// 				s.phases[vars.DayPhaseID] = make(map[types.TurnID]types.Turn)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(beginPhaseID)
// 			test.setup(s.(*scheduler))

// 			isEmpty := s.IsEmptyPhase(test.phaseID)

// 			ss.Equal(test.expectedStatus, isEmpty)
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestAddSlot() {
// 	tests := []struct {
// 		name           string
// 		newSlot        *types.NewTurnSlot
// 		expectedStatus bool
// 	}{
// 		{
// 			name: "Failure (Invalid phase ID)",
// 			newSlot: &types.NewTurnSlot{
// 				PhaseID: types.PhaseID(99),
// 			},
// 			expectedStatus: false,
// 		},
// 		{
// 			name: "Ok",
// 			newSlot: &types.NewTurnSlot{
// 				PhaseID:      vars.NightPhaseID,
// 				TurnID:       vars.MidTurn,
// 				BeginRoundID: vars.SecondRound,
// 				FrozenTimes:  vars.Once,
// 				PlayerID:     ss.player1ID,
// 				RoleID:       vars.SeerRoleID,
// 			},
// 			expectedStatus: true,
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)

// 			ok := s.AddSlot(test.newSlot)

// 			ss.Equal(test.expectedStatus, ok)
// 			if test.expectedStatus == true {
// 				ss.Equal(
// 					test.newSlot.BeginRoundID,
// 					s.(*scheduler).
// 						phases[test.newSlot.PhaseID][test.newSlot.TurnID][test.newSlot.PlayerID].BeginRoundID,
// 				)
// 				ss.Equal(
// 					test.newSlot.PlayedRoundID,
// 					s.(*scheduler).
// 						phases[test.newSlot.PhaseID][test.newSlot.TurnID][test.newSlot.PlayerID].PlayedRoundID,
// 				)
// 				ss.Equal(
// 					test.newSlot.FrozenTimes,
// 					s.(*scheduler).
// 						phases[test.newSlot.PhaseID][test.newSlot.TurnID][test.newSlot.PlayerID].FrozenTimes,
// 				)
// 				ss.Equal(
// 					test.newSlot.RoleID,
// 					s.(*scheduler).
// 						phases[test.newSlot.PhaseID][test.newSlot.TurnID][test.newSlot.PlayerID].RoleID,
// 				)
// 			}
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestRemoveSlot() {
// 	tests := []struct {
// 		name           string
// 		removedSlot    *types.RemovedTurnSlot
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 		check          func(s *scheduler)
// 	}{
// 		{
// 			name: "Failure (Ignore TurnID and RoleID)",
// 			removedSlot: &types.RemovedTurnSlot{
// 				PhaseID: vars.DayPhaseID,
// 				TurnID:  types.TurnID(0),
// 				RoleID:  types.RoleID(0),
// 			},
// 			expectedStatus: false,
// 			setup:          func(s *scheduler) {},
// 			check:          func(s *scheduler) {},
// 		},
// 		{
// 			name: "Ok (Remove all slots)",
// 			removedSlot: &types.RemovedTurnSlot{
// 				PhaseID:  types.PhaseID(0),
// 				PlayerID: ss.player1ID,
// 			},
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.NightPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					//
// 				}
// 				s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					//
// 				}
// 			},
// 			check: func(s *scheduler) {
// 				for _, phase := range s.phases {
// 					for _, turn := range phase {
// 						ss.Nil(turn[ss.player1ID])
// 					}
// 				}
// 			},
// 		},
// 		{
// 			name: "Failure (Invalid turn ID)",
// 			removedSlot: &types.RemovedTurnSlot{
// 				PhaseID:  vars.DayPhaseID,
// 				PlayerID: ss.player1ID,
// 				TurnID:   types.TurnID(99),
// 			},
// 			expectedStatus: false,
// 			setup:          func(s *scheduler) {},
// 			check:          func(s *scheduler) {},
// 		},
// 		{
// 			name: "Ok (Remove by turn ID)",
// 			removedSlot: &types.RemovedTurnSlot{
// 				PhaseID:  vars.DayPhaseID,
// 				PlayerID: ss.player1ID,
// 				TurnID:   vars.PreTurn,
// 			},
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.NightPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					//
// 				}
// 				s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					//
// 				}
// 			},
// 			check: func(s *scheduler) {
// 				ss.Nil(s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID])
// 				ss.NotNil(s.phases[vars.NightPhaseID][vars.PreTurn][ss.player1ID])
// 			},
// 		},
// 		{
// 			name: "Ok (Remove by role ID)",
// 			removedSlot: &types.RemovedTurnSlot{
// 				PhaseID:  vars.DayPhaseID,
// 				PlayerID: ss.player1ID,
// 				RoleID:   vars.HunterRoleID,
// 			},
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.NightPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					RoleID: vars.HunterRoleID,
// 				}
// 				s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					RoleID: vars.HunterRoleID,
// 				}
// 			},
// 			check: func(s *scheduler) {
// 				ss.Nil(s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID])
// 				ss.NotNil(s.phases[vars.NightPhaseID][vars.PreTurn][ss.player1ID])
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)

// 			test.setup(s.(*scheduler))
// 			ok := s.RemoveSlot(test.removedSlot)

// 			ss.Equal(test.expectedStatus, ok)
// 			test.check(s.(*scheduler))
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestFreezeSlot() {
// 	tests := []struct {
// 		name           string
// 		frozenSlot     *types.FreezeTurnSlot
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 	}{
// 		{
// 			name: "Failure (Ignore TurnID and RoleID)",
// 			frozenSlot: &types.FreezeTurnSlot{
// 				PhaseID:     vars.DayPhaseID,
// 				TurnID:      types.TurnID(0),
// 				RoleID:      types.RoleID(0),
// 				FrozenTimes: vars.Twice,
// 			},
// 			expectedStatus: false,
// 			setup:          func(s *scheduler) {},
// 		},
// 		{
// 			name: "Failure (Invalid turn ID)",
// 			frozenSlot: &types.FreezeTurnSlot{
// 				PhaseID:     vars.DayPhaseID,
// 				PlayerID:    ss.player1ID,
// 				TurnID:      types.TurnID(99),
// 				FrozenTimes: vars.Twice,
// 			},
// 			expectedStatus: false,
// 			setup:          func(s *scheduler) {},
// 		},
// 		{
// 			name: "Ok (Remove by turn ID)",
// 			frozenSlot: &types.FreezeTurnSlot{
// 				PhaseID:     vars.DayPhaseID,
// 				PlayerID:    ss.player1ID,
// 				TurnID:      vars.PreTurn,
// 				FrozenTimes: vars.Twice,
// 			},
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					//
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Remove by role ID)",
// 			frozenSlot: &types.FreezeTurnSlot{
// 				PhaseID:     vars.DayPhaseID,
// 				PlayerID:    ss.player1ID,
// 				RoleID:      vars.HunterRoleID,
// 				FrozenTimes: vars.Twice,
// 			},
// 			expectedStatus: true,
// 			setup: func(s *scheduler) {
// 				s.phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					RoleID: vars.HunterRoleID,
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)

// 			test.setup(s.(*scheduler))
// 			ok := s.FreezeSlot(test.frozenSlot)

// 			ss.Equal(test.expectedStatus, ok)
// 			if test.expectedStatus == true {
// 				ss.Equal(
// 					test.frozenSlot.FrozenTimes,
// 					s.(*scheduler).phases[vars.DayPhaseID][vars.PreTurn][ss.player1ID].FrozenTimes,
// 				)
// 			}
// 		})
// 	}
// }

// func (ss SchudulerSuite) TestNextTurn() {
// 	tests := []struct {
// 		name            string
// 		expectedStatus  bool
// 		expectedRoundID types.RoundID
// 		expectedPhaseID types.PhaseID
// 		expectedTurnID  types.TurnID
// 		setup           func(*scheduler)
// 		check           func(*scheduler)
// 	}{
// 		{
// 			name:           "Failure (Empty scheduler)",
// 			expectedStatus: false,
// 			setup:          func(s *scheduler) {},
// 			check:          func(s *scheduler) {},
// 		},
// 		{
// 			name:            "Ok (First call)",
// 			expectedStatus:  true,
// 			expectedRoundID: vars.FirstRound,
// 			expectedPhaseID: ss.beginPhaseID,
// 			expectedTurnID:  vars.PreTurn,
// 			setup: func(s *scheduler) {
// 				s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}
// 			},
// 			check: func(s *scheduler) {},
// 		},
// 		{
// 			name:            "Ok (Increase TurnID)",
// 			expectedStatus:  true,
// 			expectedRoundID: vars.FirstRound,
// 			expectedPhaseID: ss.beginPhaseID,
// 			expectedTurnID:  vars.MidTurn,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.FirstRound
// 				s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}
// 				s.phases[ss.beginPhaseID][vars.MidTurn][ss.player2ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}
// 			},
// 			check: func(s *scheduler) {},
// 		},
// 		{
// 			name:            "Ok (Increase PhaseID)",
// 			expectedStatus:  true,
// 			expectedRoundID: vars.FirstRound,
// 			expectedPhaseID: ss.beginPhaseID.NextPhasePhaseID(vars.DuskPhaseID),
// 			expectedTurnID:  vars.PreTurn,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.FirstRound
// 				s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}
// 				s.phases[ss.beginPhaseID.NextPhasePhaseID(vars.DuskPhaseID)][vars.PreTurn][ss.player2ID] =
// 					&types.TurnSlot{
// 						BeginRoundID: vars.FirstRound,
// 					}
// 			},
// 			check: func(s *scheduler) {},
// 		},
// 		{
// 			name:            "Ok (Increase RoundID)",
// 			expectedStatus:  true,
// 			expectedRoundID: vars.SecondRound,
// 			expectedPhaseID: ss.beginPhaseID,
// 			expectedTurnID:  vars.PreTurn,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.FirstRound
// 				s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}
// 			},
// 			check: func(s *scheduler) {},
// 		},
// 		// {
// 		// 	name:            "Ok (Skip empty turn by BeginRoundID)",
// 		// 	expectedStatus:  true,
// 		// 	expectedRoundID: vars.FirstRound,
// 		// 	expectedPhaseID: ss.beginPhaseID,
// 		// 	expectedTurnID:  vars.MidTurn,
// 		// 	setup: func(s *scheduler) {
// 		// 		s.roundID = vars.FirstRound
// 		// 		s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 		// 			BeginRoundID: vars.SecondRound,
// 		// 		}
// 		// 		s.phases[ss.beginPhaseID][vars.MidTurn][ss.player2ID] = &types.TurnSlot{
// 		// 			BeginRoundID: vars.FirstRound,
// 		// 		}
// 		// 	},
// 		// },
// 		// {
// 		// 	name:            "Ok (Skip empty turn by PlayedRoundID)",
// 		// 	expectedStatus:  true,
// 		// 	expectedRoundID: vars.FirstRound,
// 		// 	expectedPhaseID: ss.beginPhaseID,
// 		// 	expectedTurnID:  vars.MidTurn,
// 		// 	setup: func(s *scheduler) {
// 		// 		s.roundID = vars.FirstRound
// 		// 		s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 		// 			PlayedRoundID: vars.SecondRound,
// 		// 		}
// 		// 		s.phases[ss.beginPhaseID][vars.MidTurn][ss.player2ID] = &types.TurnSlot{
// 		// 			BeginRoundID: vars.FirstRound,
// 		// 		}
// 		// 	},
// 		// },
// 		// {
// 		// 	name:            "Ok (Skip empty turn by FrozenTimes)",
// 		// 	expectedStatus:  true,
// 		// 	expectedRoundID: vars.FirstRound,
// 		// 	expectedPhaseID: ss.beginPhaseID,
// 		// 	expectedTurnID:  vars.PreTurn,
// 		// 	setup: func(s *scheduler) {
// 		// 		s.roundID = vars.FirstRound
// 		// 		s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 		// 			PlayedRoundID: vars.SecondRound,
// 		// 		}
// 		// 		s.phases[ss.beginPhaseID][vars.MidTurn][ss.player2ID] = &types.TurnSlot{
// 		// 			BeginRoundID: vars.FirstRound,
// 		// 			FrozenTimes:  vars.Once,
// 		// 		}
// 		// 	},
// 		// },
// 		{
// 			name:            "Ok (Defrost slots)",
// 			expectedStatus:  true,
// 			expectedRoundID: vars.FirstRound,
// 			expectedPhaseID: ss.beginPhaseID,
// 			expectedTurnID:  vars.MidTurn,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.FirstRound
// 				s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 					FrozenTimes:  vars.Once,
// 				}
// 				s.phases[ss.beginPhaseID][vars.MidTurn][ss.player2ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}
// 			},
// 			check: func(s *scheduler) {
// 				ss.Empty(s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID].FrozenTimes)
// 			},
// 		},
// 		{
// 			name:            "Ok (Remove one-round slot)",
// 			expectedStatus:  true,
// 			expectedRoundID: vars.SecondRound,
// 			expectedPhaseID: ss.beginPhaseID,
// 			expectedTurnID:  vars.MidTurn,
// 			setup: func(s *scheduler) {
// 				s.roundID = vars.FirstRound
// 				s.phases[ss.beginPhaseID][vars.PreTurn][ss.player1ID] = &types.TurnSlot{
// 					PlayedRoundID: vars.FirstRound,
// 				}
// 				s.phases[ss.beginPhaseID][vars.MidTurn][ss.player2ID] = &types.TurnSlot{
// 					BeginRoundID: vars.FirstRound,
// 				}

// 				s.NextTurn()
// 			},
// 			check: func(s *scheduler) {},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			s := NewScheduler(ss.beginPhaseID)
// 			test.setup(s.(*scheduler))

// 			ok := s.NextTurn()

// 			ss.Equal(test.expectedStatus, ok)
// 			test.check(s.(*scheduler))
// 			if test.expectedStatus == true {
// 				ss.Equal(test.expectedRoundID, s.RoundID())
// 				ss.Equal(test.expectedPhaseID, s.PhaseID())
// 				ss.Equal(test.expectedTurnID, s.TurnID())
// 			}
// 		})
// 	}
// }
