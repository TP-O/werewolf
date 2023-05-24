package logic

import (
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"

	"github.com/stretchr/testify/suite"
)

type SchedulerSuite struct {
	suite.Suite
	beginPhaseId types.PhaseId
	player1Id    types.PlayerId
	player2Id    types.PlayerId
	player3Id    types.PlayerId
}

func TestSchedulerSuite(t *testing.T) {
	suite.Run(t, new(SchedulerSuite))
}

func (ss *SchedulerSuite) SetupSuite() {
	ss.beginPhaseId = constants.NightPhaseId
	ss.player1Id = types.PlayerId("1")
	ss.player2Id = types.PlayerId("2")
	ss.player3Id = types.PlayerId("3")
}

func (ss SchedulerSuite) TestNewScheduler() {
	s := NewScheduler(ss.beginPhaseId)

	ss.NotNil(s)
	ss.Equal(constants.ZeroRound, s.(*scheduler).round)
	ss.Equal(constants.PreTurn, s.(*scheduler).turn)
	ss.Equal(ss.beginPhaseId, s.(*scheduler).beginPhaseId)
	ss.Equal(ss.beginPhaseId, s.(*scheduler).phaseId)
	ss.NotNil(s.(*scheduler).phases)
	ss.Len(s.(*scheduler).phases, 3)
	for _, phase := range s.(*scheduler).phases {
		ss.Len(phase, 3)
	}
}

func (ss SchedulerSuite) TestRound() {
	s := NewScheduler(ss.beginPhaseId)

	expectedRoundId := constants.SecondRound
	s.(*scheduler).round = expectedRoundId

	ss.Equal(expectedRoundId, s.Round())
}

func (ss SchedulerSuite) TestPhaseId() {
	s := NewScheduler(ss.beginPhaseId)

	ss.Equal(ss.beginPhaseId, s.PhaseId())
}

func (ss SchedulerSuite) TestPhase() {
	s := NewScheduler(ss.beginPhaseId)

	expectedTurnSlots := types.TurnSlots(
		map[types.PlayerId]*types.TurnSlot{
			ss.player1Id: {
				RoleId: types.RoleId(99),
			},
		})
	s.(*scheduler).phases[ss.beginPhaseId][constants.PreTurn] = expectedTurnSlots

	ss.Equal(expectedTurnSlots, s.Phase()[constants.PreTurn])
}

func (ss SchedulerSuite) TestTurn() {
	s := NewScheduler(ss.beginPhaseId)
	s.(*scheduler).turn = constants.PostTurn

	ss.Equal(constants.PostTurn, s.Turn())
}

func (ss SchedulerSuite) TestTurnSlots() {
	tests := []struct {
		name              string
		expectedTurnSlots types.TurnSlots
		setup             func(*scheduler)
	}{
		{
			name:              "Nil (Current phase is empty)",
			expectedTurnSlots: nil,
			setup: func(s *scheduler) {
				s.phases[ss.beginPhaseId] = make(map[types.Turn]types.TurnSlots)
			},
		},
		{
			name: "Ok",
			expectedTurnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						RoleId: types.RoleId(99),
					},
				}),
			setup: func(s *scheduler) {
				s.phases[ss.beginPhaseId][constants.PreTurn] = types.TurnSlots(
					map[types.PlayerId]*types.TurnSlot{
						ss.player1Id: {
							RoleId: types.RoleId(99),
						},
					})
			},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)
			test.setup(s.(*scheduler))

			turnSlots := s.TurnSlots()

			ss.Equal(test.expectedTurnSlots, turnSlots)
		})
	}
}

func (ss SchedulerSuite) TestCanPlay() {
	tests := []struct {
		name           string
		turnSlots      types.TurnSlots
		expectedStatus bool
		setup          func(*scheduler)
	}{
		{
			name: "False (Late begin round slot)",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound: constants.SecondRound,
					},
				}),
			expectedStatus: false,
			setup: func(s *scheduler) {
				s.round = constants.FirstRound
			},
		}, {
			name: "False (zero begin round slot)",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound: constants.ZeroRound,
					},
				}),
			expectedStatus: false,
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		},
		{
			name: "True (One-round slot)",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						PlayedRound: constants.SecondRound,
					},
				}),
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		},
		{
			name: "False (frozen slot)",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound:  constants.SecondRound,
						FrozenTimes: constants.Once,
					},
				}),
			expectedStatus: false,
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)

			test.setup(s.(*scheduler))
			s.(*scheduler).phases[ss.beginPhaseId][constants.PreTurn] = test.turnSlots

			ss.Equal(test.expectedStatus, s.CanPlay(ss.player1Id))
		})
	}
}

func (ss SchedulerSuite) TestPlayablePlayerIds() {
	tests := []struct {
		name              string
		turnSlots         types.TurnSlots
		expectedPlayerIds []types.PlayerId
		setup             func(*scheduler)
	}{
		{
			name: "Ignore late begin round slot",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound: constants.FirstRound,
					},
					ss.player2Id: {
						BeginRound: constants.SecondRound,
					},
					ss.player3Id: {
						BeginRound: types.Round(3),
					},
				}),
			expectedPlayerIds: []types.PlayerId{
				ss.player1Id,
				ss.player2Id,
			},
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		}, {
			name: "Ignore zero begin round slot",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound: constants.FirstRound,
					},
					ss.player2Id: {
						BeginRound: constants.SecondRound,
					},
					ss.player3Id: {
						BeginRound: constants.ZeroRound,
					},
				}),
			expectedPlayerIds: []types.PlayerId{
				ss.player1Id,
				ss.player2Id,
			},
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		},
		{
			name: "Include one-round slot",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound: constants.FirstRound,
					},
					ss.player2Id: {
						PlayedRound: constants.SecondRound,
					},
					ss.player3Id: {
						PlayedRound: types.Round(3),
					},
				}),
			expectedPlayerIds: []types.PlayerId{
				ss.player1Id,
				ss.player2Id,
			},
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		},
		{
			name: "Ignore frozen slot",
			turnSlots: types.TurnSlots(
				map[types.PlayerId]*types.TurnSlot{
					ss.player1Id: {
						BeginRound: constants.FirstRound,
					},
					ss.player2Id: {
						BeginRound: constants.SecondRound,
					},
					ss.player3Id: {
						BeginRound:  constants.SecondRound,
						FrozenTimes: constants.Once,
					},
				}),
			expectedPlayerIds: []types.PlayerId{
				ss.player1Id,
				ss.player2Id,
			},
			setup: func(s *scheduler) {
				s.round = constants.SecondRound
			},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)

			test.setup(s.(*scheduler))
			s.(*scheduler).phases[ss.beginPhaseId][constants.PreTurn] = test.turnSlots

			ss.ElementsMatch(test.expectedPlayerIds, s.PlayablePlayerIds())
		})
	}
}

func (ss SchedulerSuite) TestIsEmpty() {
	beginPhaseId := constants.NightPhaseId
	tests := []struct {
		name           string
		phaseId        types.PhaseId
		expectedStatus bool
		setup          func(*scheduler)
	}{

		{
			name:           "Non-empty (Check specific phase)",
			phaseId:        ss.beginPhaseId,
			expectedStatus: false,
			setup: func(s *scheduler) {
				s.phases[ss.beginPhaseId][constants.PreTurn] = types.TurnSlots(
					map[types.PlayerId]*types.TurnSlot{
						ss.player1Id: {},
					},
				)
			},
		},
		{
			name:           "Empty (Check specific phase)",
			phaseId:        ss.beginPhaseId,
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[ss.beginPhaseId] = make(map[types.Turn]types.TurnSlots)
			},
		},
		{
			name:           "Non-empty (Check all phases)",
			phaseId:        types.PhaseId(0),
			expectedStatus: false,
			setup: func(s *scheduler) {
				s.phases[constants.NightPhaseId] = make(map[types.Turn]types.TurnSlots)
				s.phases[constants.DayPhaseId][constants.PreTurn] = types.TurnSlots(
					map[types.PlayerId]*types.TurnSlot{
						ss.player1Id: {},
					},
				)
				s.phases[constants.DuskPhaseId] = make(map[types.Turn]types.TurnSlots)
			},
		},
		{
			name:           "Empty (Check all phases)",
			phaseId:        types.PhaseId(0),
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[constants.NightPhaseId] = make(map[types.Turn]types.TurnSlots)
				s.phases[constants.DayPhaseId] = make(map[types.Turn]types.TurnSlots)
				s.phases[constants.DayPhaseId] = make(map[types.Turn]types.TurnSlots)
			},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(beginPhaseId)
			test.setup(s.(*scheduler))

			isEmpty := s.IsEmpty(test.phaseId)

			ss.Equal(test.expectedStatus, isEmpty)
		})
	}
}

func (ss SchedulerSuite) TestAddSlot() {
	tests := []struct {
		name           string
		newSlot        types.AddTurnSlot
		expectedStatus bool
	}{
		{
			name: "Failure (InvalId phase Id)",
			newSlot: types.AddTurnSlot{
				PhaseId: types.PhaseId(99),
			},
			expectedStatus: false,
		},
		{
			name: "Ok",
			newSlot: types.AddTurnSlot{
				PhaseId:  constants.NightPhaseId,
				Turn:     constants.MidTurn,
				PlayerId: ss.player1Id,
				TurnSlot: types.TurnSlot{
					BeginRound:  constants.SecondRound,
					FrozenTimes: constants.Once,
					RoleId:      constants.SeerRoleId,
				},
			},
			expectedStatus: true,
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)

			ok := s.AddSlot(test.newSlot)

			ss.Equal(test.expectedStatus, ok)
			if test.expectedStatus == true {
				ss.Equal(
					test.newSlot.BeginRound,
					s.(*scheduler).
						phases[test.newSlot.PhaseId][test.newSlot.Turn][test.newSlot.PlayerId].BeginRound,
				)
				ss.Equal(
					test.newSlot.PlayedRound,
					s.(*scheduler).
						phases[test.newSlot.PhaseId][test.newSlot.Turn][test.newSlot.PlayerId].PlayedRound,
				)
				ss.Equal(
					test.newSlot.FrozenTimes,
					s.(*scheduler).
						phases[test.newSlot.PhaseId][test.newSlot.Turn][test.newSlot.PlayerId].FrozenTimes,
				)
				ss.Equal(
					test.newSlot.RoleId,
					s.(*scheduler).
						phases[test.newSlot.PhaseId][test.newSlot.Turn][test.newSlot.PlayerId].RoleId,
				)
			}
		})
	}
}

func (ss SchedulerSuite) TestRemoveSlot() {
	tests := []struct {
		name           string
		removeSlot     types.RemoveTurnSlot
		expectedStatus bool
		setup          func(*scheduler)
		check          func(s *scheduler)
	}{
		{
			name: "Failure (Ignore TurnId and RoleId)",
			removeSlot: types.RemoveTurnSlot{
				PhaseId: constants.DayPhaseId,
				Turn:    types.Turn(0),
				RoleId:  types.RoleId(0),
			},
			expectedStatus: false,
			setup:          func(s *scheduler) {},
			check:          func(s *scheduler) {},
		},
		{
			name: "Ok (Remove all slots)",
			removeSlot: types.RemoveTurnSlot{
				PhaseId:  types.PhaseId(0),
				PlayerId: ss.player1Id,
			},
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[constants.NightPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					//
				}
				s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					//
				}
			},
			check: func(s *scheduler) {
				for _, phase := range s.phases {
					for _, turn := range phase {
						ss.Nil(turn[ss.player1Id])
					}
				}
			},
		},
		{
			name: "Failure (InvalId turn Id)",
			removeSlot: types.RemoveTurnSlot{
				PhaseId:  constants.DayPhaseId,
				PlayerId: ss.player1Id,
				Turn:     types.Turn(99),
			},
			expectedStatus: false,
			setup:          func(s *scheduler) {},
			check:          func(s *scheduler) {},
		},
		{
			name: "Ok (Remove by turn Id)",
			removeSlot: types.RemoveTurnSlot{
				PhaseId:  constants.DayPhaseId,
				PlayerId: ss.player1Id,
				Turn:     constants.PreTurn,
			},
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[constants.NightPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					//
				}
				s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					//
				}
			},
			check: func(s *scheduler) {
				ss.Nil(s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id])
				ss.NotNil(s.phases[constants.NightPhaseId][constants.PreTurn][ss.player1Id])
			},
		},
		{
			name: "Ok (Remove by role Id)",
			removeSlot: types.RemoveTurnSlot{
				PhaseId:  constants.DayPhaseId,
				PlayerId: ss.player1Id,
				RoleId:   constants.HunterRoleId,
			},
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[constants.NightPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					RoleId: constants.HunterRoleId,
				}
				s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					RoleId: constants.HunterRoleId,
				}
			},
			check: func(s *scheduler) {
				ss.Nil(s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id])
				ss.NotNil(s.phases[constants.NightPhaseId][constants.PreTurn][ss.player1Id])
			},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)

			test.setup(s.(*scheduler))
			ok := s.RemoveSlot(test.removeSlot)

			ss.Equal(test.expectedStatus, ok)
			test.check(s.(*scheduler))
		})
	}
}

func (ss SchedulerSuite) TestFreezeSlot() {
	tests := []struct {
		name           string
		frozenSlot     types.FreezeTurnSlot
		expectedStatus bool
		setup          func(*scheduler)
	}{
		{
			name: "Failure (Ignore TurnId and RoleId)",
			frozenSlot: types.FreezeTurnSlot{
				PhaseId:     constants.DayPhaseId,
				Turn:        types.Turn(0),
				RoleId:      types.RoleId(0),
				FrozenTimes: constants.Twice,
			},
			expectedStatus: false,
			setup:          func(s *scheduler) {},
		},
		{
			name: "Failure (InvalId turn Id)",
			frozenSlot: types.FreezeTurnSlot{
				PhaseId:     constants.DayPhaseId,
				PlayerId:    ss.player1Id,
				Turn:        types.Turn(99),
				FrozenTimes: constants.Twice,
			},
			expectedStatus: false,
			setup:          func(s *scheduler) {},
		},
		{
			name: "Ok (Remove by turn Id)",
			frozenSlot: types.FreezeTurnSlot{
				PhaseId:     constants.DayPhaseId,
				PlayerId:    ss.player1Id,
				Turn:        constants.PreTurn,
				FrozenTimes: constants.Twice,
			},
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					//
				}
			},
		},
		{
			name: "Ok (Remove by role Id)",
			frozenSlot: types.FreezeTurnSlot{
				PhaseId:     constants.DayPhaseId,
				PlayerId:    ss.player1Id,
				RoleId:      constants.HunterRoleId,
				FrozenTimes: constants.Twice,
			},
			expectedStatus: true,
			setup: func(s *scheduler) {
				s.phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					RoleId: constants.HunterRoleId,
				}
			},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)

			test.setup(s.(*scheduler))
			ok := s.FreezeSlot(test.frozenSlot)

			ss.Equal(test.expectedStatus, ok)
			if test.expectedStatus == true {
				ss.Equal(
					test.frozenSlot.FrozenTimes,
					s.(*scheduler).phases[constants.DayPhaseId][constants.PreTurn][ss.player1Id].FrozenTimes,
				)
			}
		})
	}
}

func (ss SchedulerSuite) TestNextTurn() {
	tests := []struct {
		name            string
		expectedStatus  bool
		expectedRound   types.Round
		expectedPhaseId types.PhaseId
		expectedTurn    types.Turn
		setup           func(*scheduler)
		check           func(*scheduler)
	}{
		{
			name:           "Failure (Empty scheduler)",
			expectedStatus: false,
			setup:          func(s *scheduler) {},
			check:          func(s *scheduler) {},
		},
		{
			name:            "Ok (First call)",
			expectedStatus:  true,
			expectedRound:   constants.FirstRound,
			expectedPhaseId: ss.beginPhaseId,
			expectedTurn:    constants.PreTurn,
			setup: func(s *scheduler) {
				s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}
			},
			check: func(s *scheduler) {},
		},
		{
			name:            "Ok (Increase TurnId)",
			expectedStatus:  true,
			expectedRound:   constants.FirstRound,
			expectedPhaseId: ss.beginPhaseId,
			expectedTurn:    constants.MidTurn,
			setup: func(s *scheduler) {
				s.round = constants.FirstRound
				s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}
				s.phases[ss.beginPhaseId][constants.MidTurn][ss.player2Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}
			},
			check: func(s *scheduler) {},
		},
		{
			name:            "Ok (Increase PhaseId)",
			expectedStatus:  true,
			expectedRound:   constants.FirstRound,
			expectedPhaseId: util.NextPhasePhaseID(ss.beginPhaseId),
			expectedTurn:    constants.PreTurn,
			setup: func(s *scheduler) {
				s.round = constants.FirstRound
				s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}
				s.phases[util.NextPhasePhaseID(ss.beginPhaseId)][constants.PreTurn][ss.player2Id] =
					&types.TurnSlot{
						BeginRound: constants.FirstRound,
					}
			},
			check: func(s *scheduler) {},
		},
		{
			name:            "Ok (Increase RoundId)",
			expectedStatus:  true,
			expectedRound:   constants.SecondRound,
			expectedPhaseId: ss.beginPhaseId,
			expectedTurn:    constants.PreTurn,
			setup: func(s *scheduler) {
				s.round = constants.FirstRound
				s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}
			},
			check: func(s *scheduler) {},
		},
		// 3 test cases below is for skipping turn if no play in that turn can play,
		// but it isn't implemented yet.
		// {
		// 	name:            "Ok (Skip empty turn by BeginRound)",
		// 	expectedStatus:  true,
		// 	expectedRound:   constants.FirstRound,
		// 	expectedPhaseId: ss.beginPhaseId,
		// 	expectedTurn:    constants.MidTurn,
		// 	setup: func(s *scheduler) {
		// 		s.round = constants.FirstRound
		// 		s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
		// 			BeginRound: constants.SecondRound,
		// 		}
		// 		s.phases[ss.beginPhaseId][constants.MidTurn][ss.player2Id] = &types.TurnSlot{
		// 			BeginRound: constants.FirstRound,
		// 		}
		// 	},
		//     check: func(s *scheduler) {},
		// },
		// {
		// 	name:            "Ok (Skip empty turn by PlayedRoundId)",
		// 	expectedStatus:  true,
		// 	expectedRound:   constants.FirstRound,
		// 	expectedPhaseId: ss.beginPhaseId,
		// 	expectedTurn:    constants.MidTurn,
		// 	setup: func(s *scheduler) {
		// 		s.round = constants.FirstRound
		// 		s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
		// 			PlayedRound: constants.SecondRound,
		// 		}
		// 		s.phases[ss.beginPhaseId][constants.MidTurn][ss.player2Id] = &types.TurnSlot{
		// 			BeginRound: constants.FirstRound,
		// 		}
		// 	},
		//     check: func(s *scheduler) {},
		// },
		// {
		// 	name:            "Ok (Skip empty turn by FrozenTimes)",
		// 	expectedStatus:  true,
		// 	expectedRound:   constants.FirstRound,
		// 	expectedPhaseId: ss.beginPhaseId,
		// 	expectedTurn:    constants.PreTurn,
		// 	setup: func(s *scheduler) {
		// 		s.round = constants.FirstRound
		// 		s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
		// 			PlayedRound: constants.SecondRound,
		// 		}
		// 		s.phases[ss.beginPhaseId][constants.MidTurn][ss.player2Id] = &types.TurnSlot{
		// 			BeginRound:  constants.FirstRound,
		// 			FrozenTimes: constants.Once,
		// 		}
		// 	},
		// 	check: func(s *scheduler) {},
		// },
		{
			name:            "Ok (Defrost slots)",
			expectedStatus:  true,
			expectedRound:   constants.FirstRound,
			expectedPhaseId: ss.beginPhaseId,
			expectedTurn:    constants.MidTurn,
			setup: func(s *scheduler) {
				s.round = constants.FirstRound
				s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					BeginRound:  constants.FirstRound,
					FrozenTimes: constants.Once,
				}
				s.phases[ss.beginPhaseId][constants.MidTurn][ss.player2Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}
			},
			check: func(s *scheduler) {
				ss.Empty(s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id].FrozenTimes)
			},
		},
		{
			name:            "Ok (Remove one-round slot)",
			expectedStatus:  true,
			expectedRound:   constants.SecondRound,
			expectedPhaseId: ss.beginPhaseId,
			expectedTurn:    constants.MidTurn,
			setup: func(s *scheduler) {
				s.round = constants.FirstRound
				s.phases[ss.beginPhaseId][constants.PreTurn][ss.player1Id] = &types.TurnSlot{
					PlayedRound: constants.FirstRound,
				}
				s.phases[ss.beginPhaseId][constants.MidTurn][ss.player2Id] = &types.TurnSlot{
					BeginRound: constants.FirstRound,
				}

				s.NextTurn()
			},
			check: func(s *scheduler) {},
		},
	}

	for _, test := range tests {
		ss.Run(test.name, func() {
			s := NewScheduler(ss.beginPhaseId)
			test.setup(s.(*scheduler))

			ok := s.NextTurn()

			ss.Equal(test.expectedStatus, ok)
			test.check(s.(*scheduler))
			if test.expectedStatus == true {
				ss.Equal(test.expectedRound, s.Round())
				ss.Equal(test.expectedPhaseId, s.PhaseId())
				ss.Equal(test.expectedTurn, s.Turn())
			}
		})
	}
}
