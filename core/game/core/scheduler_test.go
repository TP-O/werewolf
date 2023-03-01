package core

// import (
// 	"testing"
// 	"uwwolf/game/enum"
// 	"uwwolf/game/types"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type SchudulerSuite struct {
// 	suite.Suite
// 	ctrl *gomock.Controller
// }

// func TestSchudulerSuite(t *testing.T) {
// 	suite.Run(t, new(SchudulerSuite))
// }

// func (ss *SchudulerSuite) TestNew() {
// 	scheduler := NewScheduler(enum.DuskPhaseID)

// 	ss.NotNil(scheduler)
// }

// func (ss *SchudulerSuite) TestRound() {
// 	myScheduler := NewScheduler(enum.DuskPhaseID)
// 	expectedRound := enum.Round(99)
// 	myScheduler.(*scheduler).round = expectedRound

// 	ss.Equal(expectedRound, myScheduler.Round())
// }

// func (ss *SchudulerSuite) TestPhaseID() {
// 	myScheduler := NewScheduler(enum.DuskPhaseID)
// 	expectedPhaseID := enum.PhaseID(98)
// 	myScheduler.(*scheduler).phaseID = expectedPhaseID

// 	ss.Equal(expectedPhaseID, myScheduler.PhaseID())
// }

// func (ss *SchudulerSuite) TestPhase() {
// 	myScheduler := NewScheduler(enum.DuskPhaseID)
// 	// Update current phase
// 	phaseID := enum.NightPhaseID
// 	myScheduler.(*scheduler).phaseID = phaseID
// 	expectedTurn := &types.Turn{
// 		RoleID: enum.VillagerRoleID,
// 	}
// 	myScheduler.(*scheduler).phases[phaseID] = []*types.Turn{
// 		expectedTurn,
// 	}

// 	ss.Equal(expectedTurn, myScheduler.Phase()[0])
// }

// func (ss *SchudulerSuite) TestTurn() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name         string
// 		returnNil    bool
// 		expectedTurn *types.Turn
// 		setup        func(*scheduler)
// 	}{
// 		{
// 			name:      "Nil (Current phase is empty)",
// 			returnNil: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{}
// 			},
// 		},
// 		{
// 			name:      "Nil (Turn index is out of range)",
// 			returnNil: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{},
// 					{},
// 				}
// 				myScheduler.turnIndex = 5
// 			},
// 		},
// 		{
// 			name:      "Nil (Turn index is negative)",
// 			returnNil: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{},
// 					{},
// 				}
// 				myScheduler.turnIndex = -1
// 			},
// 		},
// 		{
// 			name:      "Ok",
// 			returnNil: false,
// 			expectedTurn: &types.Turn{
// 				RoleID: enum.VillagerRoleID,
// 			},
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 				myScheduler.turnIndex = 0
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			turn := myScheduler.Turn()

// 			if test.returnNil {
// 				ss.Nil(turn)
// 			} else {
// 				ss.NotNil(turn)
// 				ss.Equal(test.expectedTurn, turn)
// 			}
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestIsEmpty() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name           string
// 		phaseID        enum.PhaseID
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 	}{

// 		{
// 			name:           "Non-empty (Check specific phase)",
// 			phaseID:        enum.DayPhaseID,
// 			expectedStatus: false,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					{},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Empty (Check specific phase)",
// 			phaseID:        enum.DayPhaseID,
// 			expectedStatus: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 		{
// 			name:           "Non-empty (Check all phases)",
// 			phaseID:        0,
// 			expectedStatus: false,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					{},
// 				}
// 				myScheduler.phases[enum.NightPhaseID] = []*types.Turn{
// 					{},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Empty (Check all phases)",
// 			phaseID:        0,
// 			expectedStatus: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					//
// 				}
// 				myScheduler.phases[enum.NightPhaseID] = []*types.Turn{
// 					//
// 				}
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.IsEmpty(test.phaseID)

// 			ss.Equal(test.expectedStatus, ok)
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestIsValidPhaseID() {
// 	myScheduler := NewScheduler(enum.NightPhaseID)

// 	for i := enum.NightPhaseID; i <= enum.DuskPhaseID; i++ {
// 		ss.True(myScheduler.(*scheduler).isValidPhaseID(i))
// 	}

// 	ss.False(myScheduler.(*scheduler).isValidPhaseID(0))
// 	ss.False(myScheduler.(*scheduler).isValidPhaseID(enum.DuskPhaseID + 1))
// }

// func (ss *SchudulerSuite) TestExistRole() {
// 	tests := []struct {
// 		name           string
// 		roleID         enum.RoleID
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 	}{
// 		{
// 			name:           "Does not exist",
// 			roleID:         enum.SeerRoleID,
// 			expectedStatus: false,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Exist",
// 			roleID:         enum.VillagerRoleID,
// 			expectedStatus: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(enum.NightPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.(*scheduler).existRole(test.roleID)

// 			ss.Equal(test.expectedStatus, ok)
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestCalculateTurnIndex() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name              string
// 		setting           *types.TurnSetting
// 		expectedPhaseID   enum.PhaseID
// 		expectedTurnIndex int
// 		setup             func(*scheduler)
// 	}{
// 		{
// 			name: "Next turn (In non-empty phase)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: enum.NextPosition,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name: "Next turn (In empty phase)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: enum.NextPosition,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = -1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 		{
// 			name: "Sorted turn (In empty phase)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: enum.SortedPosition,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 		{
// 			name: "Sorted turn (In phase containing all higher priority phases)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: enum.SortedPosition,
// 				Priority: 1,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						Priority: 6,
// 					},
// 					{
// 						Priority: 3,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name: "Sorted turn (In phase containing lower and higher priority phases)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: enum.SortedPosition,
// 				Priority: 4,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						Priority: 6,
// 					},
// 					{
// 						Priority: 3,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name: "Latest turn",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: enum.LastPosition,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						Priority: 6,
// 					},
// 					{
// 						Priority: 3,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name: "Customize-position turn (Negative position)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: -99,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: -1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 		{
// 			name: "Customize-position turn (Out of range position)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: 3,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: -1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						Priority: 6,
// 					},
// 					{
// 						Priority: 3,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name: "Customize-position turn (Ok)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				Position: 1,
// 			},
// 			expectedPhaseID:   beginPhaseID,
// 			expectedTurnIndex: 1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						Priority: 6,
// 					},
// 					{
// 						Priority: 3,
// 					},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			phaseID, turnIndex := myScheduler.(*scheduler).calculateTurnIndex(
// 				test.setting,
// 			)

// 			ss.Equal(test.expectedPhaseID, phaseID)
// 			ss.Equal(test.expectedTurnIndex, turnIndex)
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestAddTurn() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name                string
// 		setting             *types.TurnSetting
// 		expectedStatus      bool
// 		expectedTurnIndex   int
// 		newCurrentTurnIndex int
// 		setup               func(*scheduler)
// 	}{
// 		{
// 			name: "Failure (Invalid phase ID)",
// 			setting: &types.TurnSetting{
// 				PhaseID: 99,
// 			},
// 			setup: func(myScheduler *scheduler) {},
// 		},
// 		{
// 			name: "Failure (Existent role ID)",
// 			setting: &types.TurnSetting{
// 				PhaseID: enum.NightPhaseID,
// 				RoleID:  enum.HunterRoleID,
// 			},
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name: "Failure (Invalid turn position)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  enum.DayPhaseID,
// 				Position: 99,
// 			},
// 			expectedStatus: false,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phaseID = enum.DayPhaseID
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Adds before current turn)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				RoleID:   enum.VillagerRoleID,
// 				Position: 0,
// 			},
// 			expectedStatus:      true,
// 			expectedTurnIndex:   0,
// 			newCurrentTurnIndex: 2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{},
// 					{},
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Adds to current turn)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				RoleID:   enum.VillagerRoleID,
// 				Position: 1,
// 			},
// 			expectedStatus:      true,
// 			expectedTurnIndex:   1,
// 			newCurrentTurnIndex: 2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{},
// 					{},
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Adds after current turn)",
// 			setting: &types.TurnSetting{
// 				PhaseID:  beginPhaseID,
// 				RoleID:   enum.VillagerRoleID,
// 				Position: 2,
// 			},
// 			expectedStatus:      true,
// 			expectedTurnIndex:   2,
// 			newCurrentTurnIndex: 1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{},
// 					{},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.AddTurn(test.setting)

// 			ss.Equal(test.expectedStatus, ok)

// 			if test.expectedStatus == true {
// 				ss.Equal(
// 					test.setting.RoleID,
// 					myScheduler.(*scheduler).
// 						phases[test.setting.PhaseID][test.expectedTurnIndex].RoleID,
// 				)
// 				ss.Equal(test.newCurrentTurnIndex, myScheduler.(*scheduler).turnIndex)
// 			}
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestRemoveTurn() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name           string
// 		roleID         enum.RoleID
// 		expectedStatus bool
// 		newPhaseID     enum.PhaseID
// 		newTurnIndex   int
// 		newRound       enum.Round
// 		setup          func(*scheduler)
// 	}{
// 		{
// 			name:           "Failure (Non-existent role ID)",
// 			roleID:         99,
// 			expectedStatus: false,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   1,
// 			newRound:       1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 1
// 				myScheduler.round = 1
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes turn not in current phase)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   0,
// 			newRound:       1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 				myScheduler.phases[enum.DayPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes turn after current turn)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   0,
// 			newRound:       1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 					{
// 						RoleID: enum.SeerRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes current turn which is not the first turn)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   0,
// 			newRound:       1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.round = 1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes current turn which is the first turn of the starting round)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   -1,
// 			newRound:       1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 				myScheduler.phases[enum.NextPhasePhaseID(beginPhaseID)] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes current turn which is the first turn of the second round)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     enum.NextPhasePhaseID(beginPhaseID),
// 			newTurnIndex:   0,
// 			newRound:       1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.round = 2
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 				myScheduler.phases[enum.NextPhasePhaseID(beginPhaseID)] = []*types.Turn{
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes current turn which is the first turn of the second phase)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   1,
// 			newRound:       2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phaseID = enum.NextPhasePhaseID(beginPhaseID)
// 				myScheduler.turnIndex = 0
// 				myScheduler.round = 2
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.TwoSistersRoleID,
// 					},
// 					{
// 						RoleID: enum.VillagerRoleID,
// 					},
// 				}
// 				myScheduler.phases[enum.NextPhasePhaseID(beginPhaseID)] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Removes current turn in one-turn phase and empty scheduler)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: true,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   -1,
// 			newRound:       2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.round = 2
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 					},
// 				}
// 				myScheduler.phases[enum.NextPhasePhaseID(beginPhaseID)] = []*types.Turn{
// 					//
// 				}
// 				myScheduler.phases[enum.NextPhasePhaseID(enum.NextPhasePhaseID(beginPhaseID))] = []*types.Turn{
// 					//
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.RemoveTurn(test.roleID)

// 			ss.Equal(test.expectedStatus, ok)
// 			ss.Equal(test.newRound, myScheduler.Round())
// 			ss.Equal(test.newPhaseID, myScheduler.PhaseID())
// 			ss.Equal(test.newTurnIndex, myScheduler.(*scheduler).turnIndex)
// 			// Check if role ID exists
// 			ss.False(myScheduler.RemoveTurn(test.roleID))
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestDefrostCurrentTurn() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name           string
// 		expectedStatus bool
// 		newFrozenLimit enum.Limit
// 		setup          func(*scheduler)
// 	}{
// 		{
// 			name:           "Turn is not frozen",
// 			expectedStatus: false,
// 			newFrozenLimit: 0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						FrozenLimit: 0,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Turn is frozen forever",
// 			expectedStatus: true,
// 			newFrozenLimit: enum.Unlimited,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						FrozenLimit: enum.Unlimited,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Reduces limit by 1)",
// 			expectedStatus: true,
// 			newFrozenLimit: 1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						FrozenLimit: 2,
// 					},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {

// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.(*scheduler).defrostCurrentTurn()

// 			ss.Equal(test.expectedStatus, ok)
// 			ss.Equal(test.newFrozenLimit, myScheduler.Turn().FrozenLimit)
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestNextTurn() {
// 	beginPhaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name                  string
// 		isRemoved             bool
// 		expectedStatus        bool
// 		expectedRemovedRoleID enum.RoleID
// 		newRound              enum.Round
// 		newPhaseID            enum.PhaseID
// 		newTurnIndex          int
// 		setup                 func(*scheduler)
// 	}{
// 		{
// 			name:           "Failure (Empty scheduler)",
// 			isRemoved:      false,
// 			expectedStatus: false,
// 			newRound:       1,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   -1,
// 			setup:          func(myScheduler *scheduler) {},
// 		},
// 		{
// 			name:           "Ok (Doesn't remove current turn)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       1,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 						Limit:  1,
// 					},
// 					{
// 						RoleID: enum.TwoSistersRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:                  "Ok (Removes the current turn and the scheduler is non-empty)",
// 			isRemoved:             true,
// 			expectedStatus:        true,
// 			expectedRemovedRoleID: enum.HunterRoleID,
// 			newRound:              1,
// 			newPhaseID:            beginPhaseID,
// 			newTurnIndex:          0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 						Limit:  1,
// 					},
// 					{
// 						RoleID: enum.TwoSistersRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:                  "Ok (Removes the current turn and the scheduler is empty)",
// 			isRemoved:             true,
// 			expectedStatus:        false,
// 			expectedRemovedRoleID: enum.HunterRoleID,
// 			newRound:              1,
// 			newPhaseID:            beginPhaseID,
// 			newTurnIndex:          -1,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.round = 1
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Next turn because the current has BeginRound which is greater than the current round)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       1,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  1,
// 					},
// 					{
// 						RoleID:     enum.TwoSistersRoleID,
// 						BeginRound: 5,
// 						Limit:      1,
// 					},
// 					{
// 						RoleID: enum.VillagerRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Next turn because the current reaches limit)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       1,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  1,
// 					},
// 					{
// 						RoleID: enum.TwoSistersRoleID,
// 						Limit:  0,
// 					},
// 					{
// 						RoleID: enum.VillagerRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Reduces turn limit)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       3,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  5,
// 					},
// 					{
// 						RoleID: enum.TwoSistersRoleID,
// 						Limit:  1,
// 					},
// 				}
// 				myScheduler.NextTurn(false)
// 				myScheduler.NextTurn(false)
// 			},
// 		},
// 		{
// 			name:           "Ok (Moves to the next phase)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       1,
// 			newPhaseID:     enum.NextPhasePhaseID(beginPhaseID),
// 			newTurnIndex:   0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  1,
// 					},
// 				}
// 				myScheduler.phases[enum.NextPhasePhaseID(beginPhaseID)] = []*types.Turn{
// 					{
// 						RoleID: enum.HunterRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Moves to the next round)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       2,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.phaseID = enum.NextPhasePhaseID(beginPhaseID)
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[enum.NextPhasePhaseID(beginPhaseID)] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  1,
// 					},
// 				}
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Skips frozen turn)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       1,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   2,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID: enum.SeerRoleID,
// 						Limit:  1,
// 					},
// 					{
// 						RoleID:      enum.HunterRoleID,
// 						FrozenLimit: 1,
// 					},
// 					{
// 						RoleID: enum.VillagerRoleID,
// 						Limit:  1,
// 					},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (Defrosts turn)",
// 			isRemoved:      false,
// 			expectedStatus: true,
// 			newRound:       3,
// 			newPhaseID:     beginPhaseID,
// 			newTurnIndex:   0,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPhaseID] = []*types.Turn{
// 					{
// 						RoleID:      enum.SeerRoleID,
// 						FrozenLimit: 1,
// 						Limit:       1,
// 					},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPhaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.NextTurn(test.isRemoved)

// 			ss.Equal(test.expectedStatus, ok)
// 			ss.Equal(test.newPhaseID, myScheduler.PhaseID())
// 			ss.Equal(test.newRound, myScheduler.Round())
// 			ss.Equal(test.newTurnIndex, myScheduler.(*scheduler).turnIndex)

// 			// Check if role ID is removed
// 			if !enum.IsUnknownRoleID(test.expectedRemovedRoleID) {
// 				ss.False(myScheduler.RemoveTurn(test.expectedRemovedRoleID))
// 			}
// 		})
// 	}
// }

// func (ss *SchudulerSuite) TestFreezeTurn() {
// 	beginPaseID := enum.NightPhaseID
// 	tests := []struct {
// 		name           string
// 		roleID         enum.RoleID
// 		frozenLimit    enum.Limit
// 		expectedStatus bool
// 		setup          func(*scheduler)
// 	}{
// 		{
// 			name:           "Failure (Role ID does not exist)",
// 			roleID:         enum.HunterRoleID,
// 			expectedStatus: false,
// 			setup:          func(myScheduler *scheduler) {},
// 		},
// 		{
// 			name:           "Ok",
// 			roleID:         enum.SeerRoleID,
// 			frozenLimit:    5,
// 			expectedStatus: true,
// 			setup: func(myScheduler *scheduler) {
// 				myScheduler.turnIndex = 0
// 				myScheduler.phases[beginPaseID] = []*types.Turn{
// 					{
// 						RoleID:      enum.SeerRoleID,
// 						FrozenLimit: 0,
// 					},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ss.Run(test.name, func() {
// 			myScheduler := NewScheduler(beginPaseID)
// 			test.setup(myScheduler.(*scheduler))
// 			ok := myScheduler.FreezeTurn(test.roleID, test.frozenLimit)

// 			ss.Equal(test.expectedStatus, ok)

// 			if test.expectedStatus == true {
// 				ss.Equal(test.frozenLimit, myScheduler.Turn().FrozenLimit)
// 			}
// 		})
// 	}
// }
