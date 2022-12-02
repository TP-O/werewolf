package game

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"

	"github.com/stretchr/testify/assert"
)

func TestNewScheduler(t *testing.T) {
	scheduler := NewScheduler(config.DuskPhaseID)

	assert.NotNil(t, scheduler)
}

func TestRoundScheduler(t *testing.T) {
	myScheduler := NewScheduler(config.DuskPhaseID)
	expectedRound := types.Round(99)
	myScheduler.(*scheduler).round = expectedRound

	assert.Equal(t, expectedRound, myScheduler.Round())
}

func TestPhaseIDScheduler(t *testing.T) {
	myScheduler := NewScheduler(config.DuskPhaseID)
	expectedPhaseID := types.PhaseID(99)
	myScheduler.(*scheduler).phaseID = expectedPhaseID

	assert.Equal(t, expectedPhaseID, myScheduler.PhaseID())
}

func TestPhaseScheduler(t *testing.T) {
	myScheduler := NewScheduler(config.DuskPhaseID)
	phaseID := config.NightPhaseID
	myScheduler.(*scheduler).phaseID = phaseID
	expectedTurn := &types.Turn{
		RoleID: config.VillagerRoleID,
	}
	myScheduler.(*scheduler).phases[phaseID] = []*types.Turn{
		expectedTurn,
	}

	assert.Equal(t, expectedTurn, myScheduler.Phase()[0])
}

func TestTurnScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name         string
		isNil        bool
		expectedTurn *types.Turn
		setup        func(*scheduler)
	}{
		{
			name:  "Nil (Current phase is empty)",
			isNil: true,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.phases[phaseID] = []*types.Turn{}
			},
		},
		{
			name:  "Nil (Turn index is out of range)",
			isNil: true,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.phases[phaseID] = []*types.Turn{
					{},
					{},
				}
				myScheduler.turnIndex = 5
			},
		},
		{
			name:  "Ok",
			isNil: false,
			expectedTurn: &types.Turn{
				RoleID: config.VillagerRoleID,
			},
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
				myScheduler.turnIndex = 0
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			turn := myScheduler.Turn()

			if test.isNil {
				assert.Nil(t, turn)
			} else {
				assert.NotNil(t, turn)
				assert.Equal(t, test.expectedTurn, turn)
			}
		})
	}
}

func TestIsEmptyScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name           string
		phaseID        types.PhaseID
		expectedStatus bool
		setup          func(*scheduler)
	}{
		{
			name:           "Non-empty (Check all phases)",
			phaseID:        0,
			expectedStatus: false,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{},
				}
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{},
				}
			},
		},
		{
			name:           "Empty (Check all phases)",
			phaseID:        0,
			expectedStatus: true,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					//
				}
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					//
				}
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					//
				}
			},
		},
		{
			name:           "Non-empty (Check specific phase)",
			phaseID:        config.DayPhaseID,
			expectedStatus: false,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{},
				}
			},
		},
		{
			name:           "Empty (Check specific phase)",
			phaseID:        config.DayPhaseID,
			expectedStatus: true,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					//
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			ok := myScheduler.IsEmpty(test.phaseID)

			assert.Equal(t, test.expectedStatus, ok)
		})
	}
}

func TestIsValidPhaseIDScheduler(t *testing.T) {
	myScheduler := NewScheduler(config.NightPhaseID)

	for i := config.NightPhaseID; i <= config.DuskPhaseID; i++ {
		assert.True(t, myScheduler.(*scheduler).isValidPhaseID(i))
	}

	assert.False(t, myScheduler.(*scheduler).isValidPhaseID(0))
	assert.False(t, myScheduler.(*scheduler).isValidPhaseID(config.DuskPhaseID+1))
}

func TestExistRoleScheduler(t *testing.T) {
	tests := []struct {
		name           string
		roleID         types.RoleID
		expectedStatus bool
		setup          func(*scheduler)
	}{
		{
			name:           "Does not exist",
			roleID:         config.SeerRoleID,
			expectedStatus: false,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Exist",
			roleID:         config.VillagerRoleID,
			expectedStatus: true,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(config.NightPhaseID)
			test.setup(myScheduler.(*scheduler))
			ok := myScheduler.(*scheduler).existRole(test.roleID)

			assert.Equal(t, test.expectedStatus, ok)
		})
	}
}

func TestCalculateTurnIndexScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name              string
		setting           *types.TurnSetting
		expectedPhaseID   types.PhaseID
		expectedTurnIndex int
		setup             func(*scheduler)
	}{
		{
			name: "Next turn (In non-empty phase)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: config.NextPosition,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name: "Next turn (In empty phase)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: config.NextPosition,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 0,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					//
				}
			},
		},
		{
			name: "Sorted turn (In empty phase)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: config.SortedPosition,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					//
				}
			},
		},
		{
			name: "Sorted turn (In phase containing all higher priority phases)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: config.SortedPosition,
				Priority: 1,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 2,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						Priority: 6,
					},
					{
						Priority: 3,
					},
				}
			},
		},
		{
			name: "Sorted turn (In phase containing lower and higher priority phases)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: config.SortedPosition,
				Priority: 4,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						Priority: 6,
					},
					{
						Priority: 3,
					},
				}
			},
		},
		{
			name: "Latest turn",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: config.LastPosition,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 2,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						Priority: 6,
					},
					{
						Priority: 3,
					},
				}
			},
		},
		{
			name: "Customize-position turn (Negative position)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: -99,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: -1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					//
				}
			},
		},
		{
			name: "Customize-position turn (Out of range position)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: 3,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: -1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						Priority: 6,
					},
					{
						Priority: 3,
					},
				}
			},
		},
		{
			name: "Customize-position turn (Ok)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: 1,
			},
			expectedPhaseID:   config.NightPhaseID,
			expectedTurnIndex: 1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						Priority: 6,
					},
					{
						Priority: 3,
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			phaseID, turnIndex := myScheduler.(*scheduler).calculateTurnIndex(
				test.setting,
			)

			assert.Equal(t, test.expectedPhaseID, phaseID)
			assert.Equal(t, test.expectedTurnIndex, turnIndex)
		})
	}
}

func TestAddTurnScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name                     string
		setting                  *types.TurnSetting
		expectedStatus           bool
		expectedPhaseID          types.PhaseID
		expectedTurnIndex        int
		expectedCurrentTurnIndex int
		setup                    func(*scheduler)
	}{
		{
			name: "Invalid phase ID",
			setting: &types.TurnSetting{
				PhaseID: 99,
			},
			setup: func(myScheduler *scheduler) {},
		},
		{
			name: "Existent role ID",
			setting: &types.TurnSetting{
				PhaseID: config.NightPhaseID,
				RoleID:  config.HunterRoleID,
			},
			setup: func(myScheduler *scheduler) {
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.HunterRoleID,
					},
				}
			},
		},
		{
			name: "Invalid turn position",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				Position: 99,
			},
			expectedStatus: false,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					//
				}
			},
		},
		{
			name: "Ok (Add before current turn)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				RoleID:   config.VillagerRoleID,
				Position: 0,
			},
			expectedStatus:           true,
			expectedPhaseID:          config.NightPhaseID,
			expectedTurnIndex:        0,
			expectedCurrentTurnIndex: 2,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 1
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{},
					{},
				}
			},
		},
		{
			name: "Ok (Add to current turn)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				RoleID:   config.VillagerRoleID,
				Position: 1,
			},
			expectedStatus:           true,
			expectedPhaseID:          config.NightPhaseID,
			expectedTurnIndex:        1,
			expectedCurrentTurnIndex: 2,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 1
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{},
					{},
				}
			},
		},
		{
			name: "Ok (Add after current turn)",
			setting: &types.TurnSetting{
				PhaseID:  config.NightPhaseID,
				RoleID:   config.VillagerRoleID,
				Position: 2,
			},
			expectedStatus:           true,
			expectedPhaseID:          config.NightPhaseID,
			expectedTurnIndex:        2,
			expectedCurrentTurnIndex: 1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 1
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{},
					{},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			ok := myScheduler.AddTurn(test.setting)

			assert.Equal(t, test.expectedStatus, ok)

			if test.expectedStatus == true {
				assert.Equal(
					t,
					test.setting.RoleID,
					myScheduler.(*scheduler).
						phases[test.expectedPhaseID][test.expectedTurnIndex].RoleID,
				)
				assert.Equal(t, test.expectedCurrentTurnIndex, myScheduler.(*scheduler).turnIndex)
			}
		})
	}
}

func TestRemoveTurnScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name           string
		roleID         types.RoleID
		expectedStatus bool
		newPhaseID     types.PhaseID
		newTurnIndex   int
		newRound       types.Round
		setup          func(*scheduler)
	}{
		{
			name:           "Non-existent role ID",
			roleID:         99,
			expectedStatus: false,
			newPhaseID:     phaseID,
			newTurnIndex:   1,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 1
				myScheduler.round = 1
			},
		},
		{
			name:           "Ok (Remove turn not in current phase)",
			roleID:         config.HunterRoleID,
			expectedStatus: true,
			newPhaseID:     phaseID,
			newTurnIndex:   1,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 1
				myScheduler.round = 1
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.HunterRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove turn after current turn)",
			roleID:         config.HunterRoleID,
			expectedStatus: true,
			newPhaseID:     phaseID,
			newTurnIndex:   0,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 1
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.HunterRoleID,
					},
					{
						RoleID: config.SeerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove current turn with reaching first round's head)",
			roleID:         config.HunterRoleID,
			expectedStatus: true,
			newPhaseID:     phaseID,
			newTurnIndex:   -1,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.HunterRoleID,
					},
				}
				myScheduler.phases[config.DuskPhaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove current turn without reaching first round's head)",
			roleID:         config.HunterRoleID,
			expectedStatus: true,
			newPhaseID:     config.DuskPhaseID,
			newTurnIndex:   0,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.round = 2
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.HunterRoleID,
					},
				}
				myScheduler.phases[config.DuskPhaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove current turn in one-turn phase)",
			roleID:         config.SeerRoleID,
			expectedStatus: true,
			newPhaseID:     config.DuskPhaseID,
			newTurnIndex:   1,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.round = 2
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.DuskPhaseID] = []*types.Turn{
					{
						RoleID: config.TwoSistersRoleID,
					},
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove current turn in one-turn phase without reducing round)",
			roleID:         config.SeerRoleID,
			expectedStatus: true,
			newPhaseID:     config.NightPhaseID,
			newTurnIndex:   1,
			newRound:       2,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = config.DayPhaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 2
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.NightPhaseID] = []*types.Turn{
					{
						RoleID: config.TwoSistersRoleID,
					},
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove current turn in one-turn phase with reducing round)",
			roleID:         config.SeerRoleID,
			expectedStatus: true,
			newPhaseID:     config.DayPhaseID,
			newTurnIndex:   1,
			newRound:       1,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.round = 2
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.DuskPhaseID] = []*types.Turn{
					//
				}
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.TwoSistersRoleID,
					},
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Remove current turn in one-turn phase and empty scheduler)",
			roleID:         config.SeerRoleID,
			expectedStatus: true,
			newPhaseID:     phaseID,
			newTurnIndex:   -1,
			newRound:       2,
			setup: func(myScheduler *scheduler) {
				myScheduler.turnIndex = 0
				myScheduler.round = 2
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.DuskPhaseID] = []*types.Turn{
					//
				}
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					//
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			ok := myScheduler.RemoveTurn(test.roleID)

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.newRound, myScheduler.Round())
			assert.Equal(t, test.newPhaseID, myScheduler.PhaseID())
			assert.Equal(t, test.newTurnIndex, myScheduler.(*scheduler).turnIndex)
			// Check if role ID exists
			assert.False(t, myScheduler.RemoveTurn(test.roleID))
		})
	}
}

func TestDefrostCurrentTurnScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name           string
		expectedStatus bool
		newFrozenLimit types.Limit
		setup          func(*scheduler)
	}{
		{
			name:           "Turn is not frozen",
			expectedStatus: false,
			newFrozenLimit: 0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						FrozenLimit: 0,
					},
				}
			},
		},
		{
			name:           "Turn is frozen forever",
			expectedStatus: true,
			newFrozenLimit: config.Unlimited,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						FrozenLimit: config.Unlimited,
					},
				}
			},
		},
		{
			name:           "Ok",
			expectedStatus: true,
			newFrozenLimit: 1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						FrozenLimit: 2,
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			ok := myScheduler.(*scheduler).defrostCurrentTurn()

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.newFrozenLimit, myScheduler.Turn().FrozenLimit)
		})
	}
}

func TestNextTurnScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name                  string
		isRemoved             bool
		expectedStatus        bool
		expectedRemovedRoleID types.RoleID
		newRound              types.Round
		newPhaseID            types.PhaseID
		newTurnIndex          int
		setup                 func(*scheduler)
	}{
		{
			name:           "Empty scheduler",
			isRemoved:      false,
			expectedStatus: false,
			newRound:       1,
			newPhaseID:     phaseID,
			newTurnIndex:   0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
			},
		},
		{
			name:           "Ok (Don't remove current turn)",
			isRemoved:      false,
			expectedStatus: true,
			newRound:       1,
			newPhaseID:     phaseID,
			newTurnIndex:   1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
					{
						RoleID: config.TwoSistersRoleID,
					},
				}
			},
		},
		{
			name:                  "Ok (Remove current turn and scheduler is non-empty)",
			isRemoved:             true,
			expectedStatus:        true,
			expectedRemovedRoleID: config.SeerRoleID,
			newRound:              1,
			newPhaseID:            phaseID,
			newTurnIndex:          0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
					{
						RoleID: config.TwoSistersRoleID,
					},
				}
			},
		},
		{
			name:                  "Ok (Remove current turn and scheduler is empty)",
			isRemoved:             true,
			expectedStatus:        false,
			expectedRemovedRoleID: config.SeerRoleID,
			newRound:              1,
			newPhaseID:            phaseID,
			newTurnIndex:          -1,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Skip round have begin round greater)",
			isRemoved:      false,
			expectedStatus: true,
			newRound:       1,
			newPhaseID:     phaseID,
			newTurnIndex:   2,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
					{
						RoleID:     config.TwoSistersRoleID,
						BeginRound: 5,
					},
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Next phase)",
			isRemoved:      false,
			expectedStatus: true,
			newRound:       1,
			newPhaseID:     config.DayPhaseID,
			newTurnIndex:   0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID: config.HunterRoleID,
					},
				}
			},
		},
		{
			name:           "Ok (Next round)",
			isRemoved:      false,
			expectedStatus: true,
			newRound:       2,
			newPhaseID:     phaseID,
			newTurnIndex:   0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
					{
						RoleID:     config.TwoSistersRoleID,
						BeginRound: 2,
					},
				}
			},
		},
		{
			name:           "Ok (Skip frozen turn)",
			isRemoved:      false,
			expectedStatus: true,
			newRound:       1,
			newPhaseID:     config.DuskPhaseID,
			newTurnIndex:   0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID:      config.HunterRoleID,
						FrozenLimit: 1,
					},
				}
				myScheduler.phases[config.DuskPhaseID] = []*types.Turn{
					{
						RoleID: config.VillagerRoleID,
					},
				}
			},
		},
		{
			name:                  "Ok (Defrost turn)",
			isRemoved:             true,
			expectedStatus:        true,
			expectedRemovedRoleID: config.SeerRoleID,
			newRound:              2,
			newPhaseID:            config.DayPhaseID,
			newTurnIndex:          0,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.round = 1
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID: config.SeerRoleID,
					},
				}
				myScheduler.phases[config.DayPhaseID] = []*types.Turn{
					{
						RoleID:      config.HunterRoleID,
						FrozenLimit: 1,
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myScheduler := NewScheduler(phaseID)
			test.setup(myScheduler.(*scheduler))
			ok := myScheduler.NextTurn(test.isRemoved)

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.newPhaseID, myScheduler.PhaseID())
			assert.Equal(t, test.newRound, myScheduler.Round())
			assert.Equal(t, test.newTurnIndex, myScheduler.(*scheduler).turnIndex)

			// Check if role ID is removed
			if !test.expectedRemovedRoleID.IsUnknown() {
				assert.False(t, myScheduler.RemoveTurn(test.expectedRemovedRoleID))
			}
		})
	}
}

func TestFreezeTurnScheduler(t *testing.T) {
	phaseID := config.NightPhaseID
	tests := []struct {
		name           string
		roleID         types.RoleID
		frozenLimit    types.Limit
		expectedStatus bool
		setup          func(*scheduler)
	}{
		{
			name:           "Role ID does not exist",
			roleID:         config.HunterRoleID,
			expectedStatus: false,
			setup:          func(myScheduler *scheduler) {},
		},
		{
			name:           "Ok",
			roleID:         config.SeerRoleID,
			frozenLimit:    5,
			expectedStatus: true,
			setup: func(myScheduler *scheduler) {
				myScheduler.phaseID = phaseID
				myScheduler.turnIndex = 0
				myScheduler.phases[phaseID] = []*types.Turn{
					{
						RoleID:      config.SeerRoleID,
						FrozenLimit: 0,
					},
				}
			},
		},
	}

	for _, test := range tests {
		myScheduler := NewScheduler(phaseID)
		test.setup(myScheduler.(*scheduler))
		ok := myScheduler.FreezeTurn(test.roleID, test.frozenLimit)

		assert.Equal(t, test.expectedStatus, ok)

		if test.expectedStatus == true {
			assert.Equal(t, test.frozenLimit, myScheduler.Turn().FrozenLimit)
		}
	}
}
