package logic

import (
	"errors"
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PlayerSuite struct {
	suite.Suite
	playerId types.PlayerId

	// role1_1Id has smallest weight.
	role1_1Id types.RoleId

	// role1_2Id has same faction Id and larger weight than role1_1Id
	role1_2Id types.RoleId

	// role2Id has difference factionId and larger weight than role1_2Id
	role2Id types.RoleId

	faction1_1Id types.FactionId
	faction1_2Id types.FactionId
	faction2Id   types.FactionId
}

func (ps *PlayerSuite) SetupSuite() {
	ps.playerId = types.PlayerId("1")
	ps.role1_1Id = constants.VillagerRoleId
	ps.role1_2Id = constants.HunterRoleId
	ps.role2Id = constants.WerewolfRoleId
	ps.faction1_1Id = constants.VillagerFactionId
	ps.faction1_2Id = constants.VillagerFactionId
	ps.faction2Id = constants.WerewolfFactionId
}

func TestPlayerSuite(t *testing.T) {
	suite.Run(t, new(PlayerSuite))
}

func (ps PlayerSuite) TestNewPlayer() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)

	p := NewPlayer(moderator, ps.playerId)

	ps.Equal(ps.playerId, p.Id())
	ps.Same(moderator, p.(*player).moderator)
	ps.Equal(constants.VillagerFactionId, p.FactionId())
	ps.NotNil(p.(*player).roles)
	ps.Empty(p.(*player).roles)
}

func (ps PlayerSuite) TestId() {
	player := NewPlayer(nil, ps.playerId)

	ps.Equal(ps.playerId, player.Id())
}

func (ps PlayerSuite) TestMainRoleId() {
	p := NewPlayer(nil, ps.playerId)
	p.(*player).mainRoleId = ps.role1_1Id

	ps.Equal(ps.role1_1Id, p.MainRoleId())
}

func (ps PlayerSuite) TestRoleIds() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	role1 := mock_game_logic.NewMockRole(ctrl)
	role2 := mock_game_logic.NewMockRole(ctrl)

	p := NewPlayer(nil, ps.playerId)
	p.(*player).roles = map[types.RoleId]contract.Role{
		ps.role1_1Id: role1,
		ps.role2Id:   role2,
	}

	ps.Contains(p.RoleIds(), ps.role1_1Id)
	ps.Contains(p.RoleIds(), ps.role2Id)
}

func (ps PlayerSuite) TestRoles() {
	ctrl := gomock.NewController(ps.T())
	defer ctrl.Finish()
	role1 := mock_game_logic.NewMockRole(ctrl)
	role2 := mock_game_logic.NewMockRole(ctrl)

	roles := map[types.RoleId]contract.Role{
		ps.role1_1Id: role1,
		ps.role2Id:   role2,
	}

	p := NewPlayer(nil, ps.playerId)
	p.(*player).roles = roles

	ps.Equal(roles, p.Roles())
}

func (ps *PlayerSuite) TestFactionId() {
	p := NewPlayer(nil, ps.playerId)
	p.(*player).factionId = ps.faction1_1Id

	ps.Equal(ps.faction1_1Id, p.FactionId())
}

func (ps *PlayerSuite) TestSetFactionId() {
	player := NewPlayer(nil, ps.playerId)
	player.SetFactionId(ps.faction1_1Id)

	ps.Equal(ps.faction1_1Id, player.FactionId())
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
			p := NewPlayer(nil, ps.playerId)
			test.setup(p.(*player))

			ps.Equal(test.expectedStatus, p.IsDead())
		})
	}
}

func (ps PlayerSuite) TestDie() {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedDead   bool
		setup          func(*player, *mock_game_logic.MockRole)
	}{
		{
			name:           "Already dead",
			expectedStatus: false,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game_logic.MockRole) {
				p.isDead = true
			},
		},
		{
			name:           "Alive (Role ability)",
			expectedStatus: false,
			expectedDead:   false,
			setup: func(p *player, mr *mock_game_logic.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(false)
			},
		},
		{
			name:           "Dead (No role protection)",
			expectedStatus: true,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game_logic.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(true)
				mr.EXPECT().OnAfterDeath()
				mr.EXPECT().OnAfterRevoke()
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			role1 := mock_game_logic.NewMockRole(ctrl)

			p := NewPlayer(moderator, ps.playerId)

			p.(*player).roles = map[types.RoleId]contract.Role{
				ps.role1_1Id: role1,
			}
			test.setup(p.(*player), role1)

			ps.Equal(test.expectedStatus, p.Die())
			ps.Equal(test.expectedDead, p.IsDead())
		})
	}
}

func (ps PlayerSuite) TestExit() {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedDead   bool
		setup          func(*player, *mock_game_logic.MockRole)
	}{
		{
			name:           "Already dead",
			expectedStatus: false,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game_logic.MockRole) {
				p.isDead = true
			},
		},
		{
			name:           "Dead (No role protection)",
			expectedStatus: true,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game_logic.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(true)
				mr.EXPECT().OnAfterDeath()
				mr.EXPECT().OnAfterRevoke()
			},
		},
		{
			name:           "Dead (Ignore role protection)",
			expectedStatus: true,
			expectedDead:   true,
			setup: func(p *player, mr *mock_game_logic.MockRole) {
				mr.EXPECT().OnBeforeDeath().Return(false)
				mr.EXPECT().OnAfterDeath()
				mr.EXPECT().OnAfterRevoke()
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			role1 := mock_game_logic.NewMockRole(ctrl)

			p := NewPlayer(moderator, ps.playerId)

			p.(*player).roles = map[types.RoleId]contract.Role{
				ps.role1_1Id: role1,
			}
			test.setup(p.(*player), role1)

			ps.Equal(test.expectedStatus, p.Exit())
			ps.Equal(test.expectedDead, p.IsDead())
		})
	}
}

func (ps PlayerSuite) TestAssignRole() {
	tests := []struct {
		name               string
		roleId             types.RoleId
		expectedStatus     bool
		expectedErr        error
		expectedMainRoleId types.RoleId
		expectedFactionId  types.FactionId
		setup              func(
			*player,
			*mock_game_logic.MockModerator,
			*mock_game_logic.MockRoleFactory,
			*mock_game_logic.MockRole,
		)
	}{
		{
			name:               "Failure (Assigned roleId)",
			roleId:             ps.role1_1Id,
			expectedStatus:     false,
			expectedErr:        errors.New("This role is already assigned ¯\\_(ツ)_/¯"),
			expectedMainRoleId: ps.role1_1Id,
			expectedFactionId:  ps.faction1_1Id,
			setup: func(
				p *player,
				mm *mock_game_logic.MockModerator,
				mf *mock_game_logic.MockRoleFactory,
				mr *mock_game_logic.MockRole,
			) {
				p.mainRoleId = ps.role1_1Id
				p.factionId = ps.faction1_1Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: nil,
				}
			},
		},
		{
			name:               "Failure (Invalid roleId)",
			roleId:             types.RoleId(0),
			expectedStatus:     false,
			expectedErr:        errors.New("Non-existent role ¯\\_ಠ_ಠ_/¯"),
			expectedMainRoleId: ps.role1_1Id,
			expectedFactionId:  ps.faction1_1Id,
			setup: func(
				p *player,
				mm *mock_game_logic.MockModerator,
				mf *mock_game_logic.MockRoleFactory,
				mr *mock_game_logic.MockRole,
			) {
				p.mainRoleId = ps.role1_1Id
				p.factionId = ps.faction1_1Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: nil,
				}

				mm.EXPECT().RoleFactory().Return(mf)
				mf.EXPECT().CreateById(types.RoleId(0), mm, ps.playerId).
					Return(nil, errors.New("Non-existent role ¯\\_ಠ_ಠ_/¯"))
			},
		},
		{
			name:               "Ok (Don't update mainRoleId and factionId)",
			roleId:             ps.role1_1Id,
			expectedStatus:     true,
			expectedMainRoleId: ps.role1_2Id,
			expectedFactionId:  ps.faction1_2Id,
			setup: func(
				p *player,
				mm *mock_game_logic.MockModerator,
				mf *mock_game_logic.MockRoleFactory,
				mr *mock_game_logic.MockRole,
			) {
				p.mainRoleId = ps.role1_2Id
				p.factionId = ps.faction1_2Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_2Id: nil,
				}

				mm.EXPECT().RoleFactory().Return(mf)
				mf.EXPECT().CreateById(ps.role1_1Id, mm, ps.playerId).
					Return(mr, nil)
				mr.EXPECT().Id().Return(ps.role1_1Id)
				mr.EXPECT().OnAfterAssign()
			},
		},
		{
			name:               "Ok (Update mainRoleId and factionId)",
			roleId:             ps.role2Id,
			expectedStatus:     true,
			expectedMainRoleId: ps.role2Id,
			expectedFactionId:  ps.faction2Id,
			setup: func(
				p *player,
				mm *mock_game_logic.MockModerator,
				mf *mock_game_logic.MockRoleFactory,
				mr *mock_game_logic.MockRole,
			) {
				p.mainRoleId = ps.role1_1Id
				p.factionId = ps.faction1_1Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: nil,
				}

				mm.EXPECT().RoleFactory().Return(mf)
				mf.EXPECT().CreateById(ps.role2Id, mm, ps.playerId).
					Return(mr, nil)
				mr.EXPECT().Id().Return(ps.role2Id).Times(2)
				mr.EXPECT().FactionId().Return(ps.faction2Id)
				mr.EXPECT().OnAfterAssign()
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			factory := mock_game_logic.NewMockRoleFactory(ctrl)
			role := mock_game_logic.NewMockRole(ctrl)

			p := NewPlayer(moderator, ps.playerId)
			test.setup(p.(*player), moderator, factory, role)

			ok, err := p.AssignRole(test.roleId)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedMainRoleId, p.MainRoleId())
			ps.Equal(test.expectedFactionId, p.FactionId())

			if test.expectedStatus == true {
				ps.Contains(p.RoleIds(), test.roleId)
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
		roleId             types.RoleId
		expectedStatus     bool
		expectedErr        error
		expectedMainRoleId types.RoleId
		expectedFactionId  types.FactionId
		setup              func(*player, *mock_game_logic.MockRole, *mock_game_logic.MockRole)
	}{
		{
			name:               "Failure (Must have at least one role)",
			roleId:             ps.role1_1Id,
			expectedStatus:     false,
			expectedErr:        errors.New("Player must player at least one role ヾ(⌐■_■)ノ♪"),
			expectedMainRoleId: ps.role1_1Id,
			expectedFactionId:  ps.faction1_1Id,
			setup: func(p *player, mr1 *mock_game_logic.MockRole, mr2 *mock_game_logic.MockRole) {
				p.mainRoleId = ps.role1_1Id
				p.factionId = ps.faction1_1Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: mr1,
				}
			},
		},
		{
			name:               "Failure (Non-existent role)",
			roleId:             types.RoleId(0),
			expectedStatus:     false,
			expectedErr:        errors.New("Non-existent role ID  ¯\\_(ツ)_/¯"),
			expectedMainRoleId: ps.role1_1Id,
			expectedFactionId:  ps.faction1_1Id,
			setup: func(p *player, mr1 *mock_game_logic.MockRole, mr2 *mock_game_logic.MockRole) {
				p.mainRoleId = ps.role1_1Id
				p.factionId = ps.faction1_1Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: mr1,
					ps.role1_2Id: mr2,
				}
			},
		},
		{
			name:               "Ok (Don't update mainRoleId and factionId)",
			roleId:             ps.role1_1Id,
			expectedStatus:     true,
			expectedMainRoleId: ps.role2Id,
			expectedFactionId:  ps.faction2Id,
			setup: func(p *player, mr1 *mock_game_logic.MockRole, mr2 *mock_game_logic.MockRole) {
				p.mainRoleId = ps.role2Id
				p.factionId = ps.faction2Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: mr1, // Prevent nil checking
					ps.role2Id:   mr1,
				}

				mr1.EXPECT().OnAfterRevoke()
			},
		},
		{
			name:               "Ok (Update mainRoleId and factionId)",
			roleId:             ps.role2Id,
			expectedStatus:     true,
			expectedMainRoleId: ps.role1_1Id,
			expectedFactionId:  ps.faction1_1Id,
			setup: func(p *player, mr1 *mock_game_logic.MockRole, mr2 *mock_game_logic.MockRole) {
				p.mainRoleId = ps.role2Id
				p.factionId = ps.faction2Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: mr1,
					ps.role2Id:   mr1, // Prevent nil checking
				}

				mr1.EXPECT().OnAfterRevoke()
				mr1.EXPECT().Id().Return(ps.role1_1Id)
				mr1.EXPECT().FactionId().Return(ps.faction1_1Id)
			},
		},
		{
			name:               "Ok (Update mainRoleId and factionId)",
			roleId:             ps.role2Id,
			expectedStatus:     true,
			expectedMainRoleId: ps.role1_2Id,
			expectedFactionId:  ps.faction1_2Id,
			setup: func(p *player, mr1 *mock_game_logic.MockRole, mr2 *mock_game_logic.MockRole) {
				p.mainRoleId = ps.role2Id
				p.factionId = ps.faction2Id
				p.roles = map[types.RoleId]contract.Role{
					ps.role1_1Id: mr1,
					ps.role1_2Id: mr2,
					ps.role2Id:   mr1, // Prevent nil checking
				}

				mr1.EXPECT().OnAfterRevoke()
				mr1.EXPECT().Id().Return(ps.role1_1Id).MaxTimes(2)
				mr1.EXPECT().FactionId().Return(ps.faction1_1Id).MaxTimes(1)
				mr2.EXPECT().Id().Return(ps.role1_2Id).MaxTimes(2)
				mr2.EXPECT().FactionId().Return(ps.faction1_2Id).MaxTimes(1)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			role1 := mock_game_logic.NewMockRole(ctrl)
			role2 := mock_game_logic.NewMockRole(ctrl)

			p := NewPlayer(nil, ps.playerId)
			test.setup(p.(*player), role1, role2)
			ok, err := p.RevokeRole(test.roleId)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedMainRoleId, p.MainRoleId())
			ps.Equal(test.expectedFactionId, p.FactionId())

			if test.expectedStatus == true {
				ps.NotContains(p.RoleIds(), test.roleId)
				ps.Nil(err)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ps PlayerSuite) TestUseRole() {
	tests := []struct {
		name        string
		req         types.RoleRequest
		expectedRes types.RoleResponse
		setup       func(*player, *mock_game_logic.MockScheduler, *mock_game_logic.MockRole)
	}{
		{
			name: "Failure (Dead)",
			req:  types.RoleRequest{},
			expectedRes: types.RoleResponse{
				ActionResponse: types.ActionResponse{
					Message: "You're died (╥﹏╥)",
				},
			},
			setup: func(p *player, ms *mock_game_logic.MockScheduler, mr *mock_game_logic.MockRole) {
				p.isDead = true
			},
		},
		{
			name: "Failure (Not in turn)",
			req:  types.RoleRequest{},
			expectedRes: types.RoleResponse{
				ActionResponse: types.ActionResponse{
					Message: "Wait for your turn, OK??",
				},
			},
			setup: func(p *player, ms *mock_game_logic.MockScheduler, mr *mock_game_logic.MockRole) {
				ms.EXPECT().CanPlay(ps.playerId).Return(false)
			},
		},
		{
			name: "Ok",
			req:  types.RoleRequest{},
			expectedRes: types.RoleResponse{
				ActionResponse: types.ActionResponse{
					Ok: true,
				},
			},
			setup: func(p *player, ms *mock_game_logic.MockScheduler, mr *mock_game_logic.MockRole) {
				turnSlots := types.TurnSlots(map[types.PlayerId]*types.TurnSlot{
					ps.playerId: {
						RoleId: ps.role1_1Id,
					},
				})
				p.roles[ps.role1_1Id] = mr

				ms.EXPECT().TurnSlots().Return(turnSlots)
				ms.EXPECT().CanPlay(ps.playerId).Return(true)
				mr.EXPECT().
					Use(types.RoleRequest{}).
					Return(types.RoleResponse{
						ActionResponse: types.ActionResponse{
							Ok: true,
						},
					}).
					Times(1)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			ctrl := gomock.NewController(ps.T())
			defer ctrl.Finish()
			moderator := mock_game_logic.NewMockModerator(ctrl)
			scheduler := mock_game_logic.NewMockScheduler(ctrl)
			role1_1 := mock_game_logic.NewMockRole(ctrl)

			moderator.EXPECT().Scheduler().Return(scheduler).MaxTimes(2)

			p := NewPlayer(moderator, ps.playerId)
			test.setup(p.(*player), scheduler, role1_1)

			res := p.UseRole(test.req)

			ps.Equal(test.expectedRes, res)
		})
	}
}
