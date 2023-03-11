package game

import (
	"fmt"
	"testing"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	mock_game "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PlayerSuite struct {
	suite.Suite
	playerID types.PlayerID

	// role1_1ID has smallest weight.
	role1_1ID types.RoleID

	// role1_2ID has same faction ID and larger weight than role1_1ID
	role1_2ID types.RoleID

	// role2ID has difference factionID and larger weight than role1_2ID
	role2ID types.RoleID

	faction1_1ID types.FactionID
	faction1_2ID types.FactionID
	faction2ID   types.FactionID
}

func (ps *PlayerSuite) SetupSuite() {
	ps.playerID = types.PlayerID("1")
	ps.role1_1ID = vars.VillagerRoleID
	ps.role1_2ID = vars.HunterRoleID
	ps.role2ID = vars.WerewolfRoleID
	ps.faction1_1ID = vars.VillagerFactionID
	ps.faction1_2ID = vars.VillagerFactionID
	ps.faction2ID = vars.WerewolfFactionID
}

func TestPlayerSuite(t *testing.T) {
	suite.Run(t, new(PlayerSuite))
}

func (ps PlayerSuite) TestNewPlayer() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)

	p := NewPlayer(game, ps.playerID)

	ps.Equal(ps.playerID, p.ID())
	ps.Same(game, p.(*player).game)
	ps.Equal(vars.VillagerFactionID, p.FactionID())
	ps.NotNil(p.(*player).roles)
	ps.Empty(p.(*player).roles)
}

func (ps PlayerSuite) TestID() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)

	player := NewPlayer(game, ps.playerID)

	ps.Equal(ps.playerID, player.ID())
}

func (ps PlayerSuite) TestMainRoleID() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)

	p := NewPlayer(game, ps.playerID)
	p.(*player).mainRoleID = ps.role1_1ID

	ps.Equal(ps.role1_1ID, p.MainRoleID())
}

func (ps PlayerSuite) TestRoleIDs() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)
	role1 := mock_game.NewMockRole(ctrl)
	role2 := mock_game.NewMockRole(ctrl)

	p := NewPlayer(game, ps.playerID)
	p.(*player).roles = map[types.RoleID]contract.Role{
		ps.role1_1ID: role1,
		ps.role2ID:   role2,
	}

	ps.Contains(p.RoleIDs(), ps.role1_1ID)
	ps.Contains(p.RoleIDs(), ps.role2ID)
}

func (ps PlayerSuite) TestRoles() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)
	role1 := mock_game.NewMockRole(ctrl)
	role2 := mock_game.NewMockRole(ctrl)

	roles := map[types.RoleID]contract.Role{
		ps.role1_1ID: role1,
		ps.role2ID:   role2,
	}

	p := NewPlayer(game, ps.playerID)
	p.(*player).roles = roles

	ps.Equal(roles, p.Roles())
}

func (ps *PlayerSuite) TestFactionID() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)

	p := NewPlayer(game, ps.playerID)
	p.(*player).factionID = ps.faction1_1ID

	ps.Equal(ps.faction1_1ID, p.FactionID())
}

func (ps *PlayerSuite) TestSetFactionID() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	game := mock_game.NewMockGame(ctrl)

	player := NewPlayer(game, ps.playerID)
	player.SetFactionID(ps.faction1_1ID)

	ps.Equal(ps.faction1_1ID, player.FactionID())
}

func (ps *PlayerSuite) TestIsDead() {
	tests := []struct {
		name           string
		expectedStatus bool
		setup          func(*player)
	}{
		{
			name:           "Died",
			expectedStatus: true,
			setup: func(p *player) {
				p.isDead = true
			},
		},
		{
			name:           "Alive",
			expectedStatus: false,
			setup: func(p *player) {
				p.isDead = false
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)

			p := NewPlayer(game, ps.playerID)
			test.setup(p.(*player))

			ps.Equal(test.expectedStatus, p.IsDead())
		})
	}
}

func (ps PlayerSuite) TestDie() {
	tests := []struct {
		name           string
		isExited       bool
		expectedStatus bool
		expectedDead   bool
		setup          func(*player, *mock_game.MockRole)
	}{
		{
			name:           "Already dead",
			isExited:       false,
			expectedStatus: false,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game.MockRole) {
				p.isDead = true
			},
		},
		{
			name:           "Is protected (Role ability)",
			isExited:       false,
			expectedStatus: false,
			expectedDead:   false,
			setup: func(p *player, mr *mock_game.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(false)
			},
		},
		{
			name:           "Died (No role protection)",
			isExited:       false,
			expectedStatus: true,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(true)
				mr.EXPECT().OnAfterDeath()
				mr.EXPECT().OnRevoke()
			},
		},
		{
			name:           "Died (Ignore role protection becauseof player exiting)",
			isExited:       true,
			expectedStatus: true,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(false)
				mr.EXPECT().OnAfterDeath()
				mr.EXPECT().OnRevoke()
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)
			role1 := mock_game.NewMockRole(ctrl)

			p := NewPlayer(game, ps.playerID)

			p.(*player).roles = map[types.RoleID]contract.Role{
				ps.role1_1ID: role1,
			}
			test.setup(p.(*player), role1)

			ps.Equal(test.expectedStatus, p.Die(test.isExited))
			ps.Equal(test.expectedDead, p.IsDead())
		})
	}
}

func (ps PlayerSuite) TestAssignRole() {
	tests := []struct {
		name               string
		roleID             types.RoleID
		expectedStatus     bool
		expectedErr        error
		expectedMainRoleID types.RoleID
		expectedFactionID  types.FactionID
		setup              func(*player)
	}{
		{
			name:               "Falure (Assigned roleID)",
			roleID:             ps.role1_1ID,
			expectedStatus:     false,
			expectedErr:        fmt.Errorf("This role is already assigned ¯\\_(ツ)_/¯"),
			expectedMainRoleID: ps.role1_1ID,
			expectedFactionID:  ps.faction1_1ID,
			setup: func(p *player) {
				p.mainRoleID = ps.role1_1ID
				p.factionID = ps.faction1_1ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: nil,
				}
			},
		},
		{
			name:               "Falure (Invalid roleID)",
			roleID:             types.RoleID(0),
			expectedStatus:     false,
			expectedErr:        fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯"),
			expectedMainRoleID: ps.role1_1ID,
			expectedFactionID:  ps.faction1_1ID,
			setup: func(p *player) {
				p.mainRoleID = ps.role1_1ID
				p.factionID = ps.faction1_1ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: nil,
				}
			},
		},
		{
			name:               "Ok (Don't update mainRoleID and factionID)",
			roleID:             ps.role1_1ID,
			expectedStatus:     true,
			expectedMainRoleID: ps.role1_2ID,
			expectedFactionID:  ps.faction1_2ID,
			setup: func(p *player) {
				p.mainRoleID = ps.role1_2ID
				p.factionID = ps.faction1_2ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_2ID: nil,
				}
			},
		},
		{
			name:               "Ok (Update mainRoleID and factionID)",
			roleID:             ps.role2ID,
			expectedStatus:     true,
			expectedMainRoleID: ps.role2ID,
			expectedFactionID:  ps.faction2ID,
			setup: func(p *player) {
				p.mainRoleID = ps.role1_1ID
				p.factionID = ps.faction1_1ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: nil,
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)
			poll := mock_game.NewMockPoll(ctrl)
			mplayer := mock_game.NewMockPlayer(ctrl)
			schduler := mock_game.NewMockScheduler(ctrl)

			// Mock for role creation
			game.EXPECT().Player(ps.playerID).Return(mplayer).AnyTimes()
			mplayer.EXPECT().ID().Return(ps.playerID).AnyTimes()
			game.EXPECT().Poll(gomock.Any()).Return(poll).AnyTimes()
			poll.EXPECT().AddElectors(gomock.Any()).AnyTimes()
			poll.EXPECT().SetWeight(gomock.Any(), gomock.Any()).AnyTimes()

			// Make sure that assinged role register it turn in scheduler
			if test.expectedStatus == true {
				game.EXPECT().Scheduler().Return(schduler)
				schduler.EXPECT().AddSlot(gomock.Any()).Return(true)

				if test.roleID == vars.VillagerRoleID {
					// Mock for villager role registration
					poll.EXPECT().AddCandidates(ps.playerID).Times(2)
				} else if test.roleID == vars.WerewolfRoleID {
					// Mock for werewolf role registration
					poll.EXPECT().AddCandidates(ps.playerID)
				}
			}

			p := NewPlayer(game, ps.playerID)
			test.setup(p.(*player))

			ok, err := p.AssignRole(test.roleID)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedMainRoleID, p.MainRoleID())
			ps.Equal(test.expectedFactionID, p.FactionID())

			if test.expectedStatus == true {
				ps.Contains(p.RoleIDs(), test.roleID)
				ps.Nil(err)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ps PlayerSuite) TestRevokeRole() {
	tests := []struct {
		name               string
		roleID             types.RoleID
		expectedStatus     bool
		expectedErr        error
		expectedMainRoleID types.RoleID
		expectedFactionID  types.FactionID
		setup              func(*player, *mock_game.MockRole, *mock_game.MockRole)
	}{
		{
			name:               "Falure (Must have at least one role)",
			roleID:             ps.role1_1ID,
			expectedStatus:     false,
			expectedErr:        fmt.Errorf("Player must player at least one role ヾ(⌐■_■)ノ♪"),
			expectedMainRoleID: ps.role1_1ID,
			expectedFactionID:  ps.faction1_1ID,
			setup: func(p *player, mr1 *mock_game.MockRole, mr2 *mock_game.MockRole) {
				p.mainRoleID = ps.role1_1ID
				p.factionID = ps.faction1_1ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: mr1,
				}
			},
		},
		{
			name:               "Falure (Non-existent role)",
			roleID:             types.RoleID(0),
			expectedStatus:     false,
			expectedErr:        fmt.Errorf("Non-existent role ID  ¯\\_(ツ)_/¯"),
			expectedMainRoleID: ps.role1_1ID,
			expectedFactionID:  ps.faction1_1ID,
			setup: func(p *player, mr1 *mock_game.MockRole, mr2 *mock_game.MockRole) {
				p.mainRoleID = ps.role1_1ID
				p.factionID = ps.faction1_1ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: mr1,
					ps.role1_2ID: mr2,
				}
			},
		},
		{
			name:               "Ok (Don't update mainRoleID and factionID)",
			roleID:             ps.role1_1ID,
			expectedStatus:     true,
			expectedMainRoleID: ps.role2ID,
			expectedFactionID:  ps.faction2ID,
			setup: func(p *player, mr1 *mock_game.MockRole, mr2 *mock_game.MockRole) {
				p.mainRoleID = ps.role2ID
				p.factionID = ps.faction2ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: mr1, // Prevent nil checking
					ps.role2ID:   mr1,
				}

				mr1.EXPECT().OnRevoke()
			},
		},
		{
			name:               "Ok (Update mainRoleID and factionID)",
			roleID:             ps.role2ID,
			expectedStatus:     true,
			expectedMainRoleID: ps.role1_1ID,
			expectedFactionID:  ps.faction1_1ID,
			setup: func(p *player, mr1 *mock_game.MockRole, mr2 *mock_game.MockRole) {
				p.mainRoleID = ps.role2ID
				p.factionID = ps.faction2ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: mr1,
					ps.role2ID:   mr1, // Prevent nil checking
				}

				mr1.EXPECT().OnRevoke()
				mr1.EXPECT().ID().Return(ps.role1_1ID)
				mr1.EXPECT().FactionID().Return(ps.faction1_1ID)
			},
		},
		{
			name:               "Ok (Update mainRoleID and factionID)",
			roleID:             ps.role2ID,
			expectedStatus:     true,
			expectedMainRoleID: ps.role1_2ID,
			expectedFactionID:  ps.faction1_2ID,
			setup: func(p *player, mr1 *mock_game.MockRole, mr2 *mock_game.MockRole) {
				p.mainRoleID = ps.role2ID
				p.factionID = ps.faction2ID
				p.roles = map[types.RoleID]contract.Role{
					ps.role1_1ID: mr1,
					ps.role1_2ID: mr2,
					ps.role2ID:   mr1, // Prevent nil checking
				}

				mr1.EXPECT().OnRevoke()
				mr1.EXPECT().ID().Return(ps.role1_1ID).MaxTimes(2)
				mr1.EXPECT().FactionID().Return(ps.faction1_1ID).MaxTimes(1)
				mr2.EXPECT().ID().Return(ps.role1_2ID).MaxTimes(2)
				mr2.EXPECT().FactionID().Return(ps.faction1_2ID).MaxTimes(1)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)
			role1 := mock_game.NewMockRole(ctrl)
			role2 := mock_game.NewMockRole(ctrl)

			p := NewPlayer(game, ps.playerID)
			test.setup(p.(*player), role1, role2)
			ok, err := p.RevokeRole(test.roleID)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedMainRoleID, p.MainRoleID())
			ps.Equal(test.expectedFactionID, p.FactionID())

			if test.expectedStatus == true {
				ps.NotContains(p.RoleIDs(), test.roleID)
				ps.Nil(err)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ps PlayerSuite) TestActivateAbility() {
	tests := []struct {
		name        string
		req         *types.ActivateAbilityRequest
		expectedRes *types.ActionResponse
		setup       func(*player, *mock_game.MockScheduler, *mock_game.MockRole)
	}{
		{
			name: "Falure (Dead)",
			req:  &types.ActivateAbilityRequest{},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "You're died (╥﹏╥)",
			},
			setup: func(p *player, ms *mock_game.MockScheduler, mr *mock_game.MockRole) {
				p.isDead = true
			},
		},
		{
			name: "Falure (Not in turn)",
			req:  &types.ActivateAbilityRequest{},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Wait for your turn, OK??",
			},
			setup: func(p *player, ms *mock_game.MockScheduler, mr *mock_game.MockRole) {
				ms.EXPECT().CanPlay(ps.playerID).Return(false)
			},
		},
		{
			name: "Ok",
			req:  &types.ActivateAbilityRequest{},
			expectedRes: &types.ActionResponse{
				Ok: true,
			},
			setup: func(p *player, ms *mock_game.MockScheduler, mr *mock_game.MockRole) {
				turn := types.Turn(map[types.PlayerID]*types.TurnSlot{
					ps.playerID: {
						RoleID: ps.role1_1ID,
					},
				})
				p.roles[ps.role1_1ID] = mr

				ms.EXPECT().Turn().Return(turn)
				ms.EXPECT().CanPlay(ps.playerID).Return(true)
				mr.EXPECT().
					ActivateAbility(&types.ActivateAbilityRequest{}).
					Return(&types.ActionResponse{
						Ok: true,
					}).
					Times(1)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			game := mock_game.NewMockGame(ctrl)
			scheduler := mock_game.NewMockScheduler(ctrl)
			role1_1 := mock_game.NewMockRole(ctrl)

			game.EXPECT().Scheduler().Return(scheduler).MaxTimes(2)

			p := NewPlayer(game, ps.playerID)
			test.setup(p.(*player), scheduler, role1_1)

			res := p.ActivateAbility(test.req)

			ps.Equal(test.expectedRes, res)
		})
	}
}
