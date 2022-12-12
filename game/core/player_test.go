package core

import (
	"testing"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/maps"
)

type PlayerSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	game         *gamemock.MockGame
	player       *gamemock.MockPlayer
	poll         *gamemock.MockPoll
	scheduler    *gamemock.MockScheduler
	hunterRole   *gamemock.MockRole
	villagerRole *gamemock.MockRole
	werewolfRole *gamemock.MockRole
	playerID     enum.PlayerID
	roles        map[enum.RoleID]contract.Role
}

func (ps *PlayerSuite) SetupSuite() {
	ps.playerID = "1"
}

func (ps *PlayerSuite) SetupTest() {
	ps.ctrl = gomock.NewController(ps.T())
	ps.game = gamemock.NewMockGame(ps.ctrl)
	ps.player = gamemock.NewMockPlayer(ps.ctrl)
	ps.poll = gamemock.NewMockPoll(ps.ctrl)
	ps.scheduler = gamemock.NewMockScheduler(ps.ctrl)
	ps.hunterRole = gamemock.NewMockRole(ps.ctrl)
	ps.villagerRole = gamemock.NewMockRole(ps.ctrl)
	ps.werewolfRole = gamemock.NewMockRole(ps.ctrl)

	ps.game.EXPECT().Scheduler().Return(ps.scheduler).AnyTimes()
	ps.werewolfRole.EXPECT().ID().Return(enum.WerewolfRoleID).AnyTimes()
	ps.hunterRole.EXPECT().ID().Return(enum.HunterRoleID).AnyTimes()
	ps.villagerRole.EXPECT().ID().Return(enum.VillagerRoleID).AnyTimes()
	ps.hunterRole.EXPECT().FactionID().Return(enum.VillagerFactionID).AnyTimes()
	ps.villagerRole.EXPECT().FactionID().Return(enum.VillagerFactionID).AnyTimes()
	ps.hunterRole.EXPECT().BeforeDeath().Return(false).AnyTimes()
	ps.villagerRole.EXPECT().BeforeDeath().Return(true).AnyTimes()
	ps.hunterRole.EXPECT().AfterDeath().Return().AnyTimes()
	ps.villagerRole.EXPECT().AfterDeath().Return().AnyTimes()
	ps.roles = map[enum.RoleID]contract.Role{
		enum.HunterRoleID:   ps.hunterRole,
		enum.VillagerRoleID: ps.villagerRole,
	}
}

func (ps *PlayerSuite) TearDownTest() {
	ps.ctrl.Finish()
}

func TestPlayerSuite(t *testing.T) {
	suite.Run(t, new(PlayerSuite))
}

func (ps *PlayerSuite) TestID() {
	player := NewPlayer(ps.game, ps.playerID)

	ps.Equal(ps.playerID, player.ID())
}

func (ps *PlayerSuite) TestMainRoleID() {
	p := NewPlayer(ps.game, ps.playerID)
	p.(*player).mainRoleID = enum.WerewolfRoleID

	ps.Equal(enum.WerewolfRoleID, p.MainRoleID())
}

func (ps *PlayerSuite) TestRoleIDs() {
	p := NewPlayer(ps.game, ps.playerID)
	p.(*player).roles = ps.roles

	for _, roleID := range maps.Keys(ps.roles) {
		ps.Contains(p.RoleIDs(), roleID)
	}
}

func (ps *PlayerSuite) TestRoles() {
	p := NewPlayer(ps.game, ps.playerID)
	p.(*player).roles = ps.roles

	ps.Equal(ps.roles, p.Roles())
}

func (ps *PlayerSuite) TestFactionID() {
	p := NewPlayer(ps.game, ps.playerID)
	p.(*player).factionID = enum.WerewolfFactionID

	ps.Equal(enum.WerewolfFactionID, p.FactionID())
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
			p := NewPlayer(ps.game, ps.playerID)
			test.setup(p.(*player))

			ps.Equal(test.expectedStatus, p.IsDead())
		})
	}
}

func (ps *PlayerSuite) TestDie() {
	tests := []struct {
		name           string
		isExited       bool
		expectedStatus bool
		newIsDead      bool
		setup          func(*player)
	}{
		{
			name:           "Died",
			isExited:       false,
			expectedStatus: false,
			newIsDead:      true,
			setup: func(p *player) {
				p.isDead = true
			},
		},
		{
			name:           "Still alive (Role passive saves)",
			isExited:       false,
			expectedStatus: false,
			newIsDead:      false,
			setup: func(p *player) {
				p.isDead = false
				p.roles = ps.roles
			},
		},
		{
			name:           "Dies (No role passive saves)",
			isExited:       false,
			expectedStatus: true,
			newIsDead:      true,
			setup: func(p *player) {
				p.isDead = false
				p.roles[enum.VillagerRoleID] = ps.villagerRole
			},
		},
		{
			name:           "Dies (Role passive saves but player is exited)",
			isExited:       true,
			expectedStatus: true,
			newIsDead:      true,
			setup: func(p *player) {
				p.isDead = false
				p.roles = ps.roles
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPlayer(ps.game, ps.playerID)
			test.setup(p.(*player))

			ps.Equal(test.expectedStatus, p.Die(test.isExited))
			ps.Equal(test.newIsDead, p.IsDead())
		})
	}
}

func (ps *PlayerSuite) TestRevive() {
	tests := []struct {
		name           string
		expectedStatus bool
		newIsDead      bool
		setup          func(*player)
	}{
		{
			name:           "Still alive",
			expectedStatus: false,
			newIsDead:      false,
			setup: func(p *player) {
				p.isDead = false
			},
		},
		{
			name:           "Ok",
			expectedStatus: true,
			newIsDead:      false,
			setup: func(p *player) {
				p.isDead = true
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPlayer(ps.game, ps.playerID)
			test.setup(p.(*player))

			ps.Equal(test.expectedStatus, p.Revive())
			ps.Equal(test.newIsDead, p.IsDead())
		})
	}
}

func (ps *PlayerSuite) TestSetFactionID() {
	player := NewPlayer(ps.game, ps.playerID)
	player.SetFactionID(enum.WerewolfFactionID)

	ps.Equal(enum.WerewolfFactionID, player.FactionID())
}

func (ps *PlayerSuite) TestAssignRole() {
	tests := []struct {
		name           string
		roleID         enum.RoleID
		expectedStatus bool
		expectedErr    string
		newMainRoleID  enum.RoleID
		newFactionID   enum.FactionID
		setup          func(*player)
	}{
		{
			name:           "Falure (Assigned role)",
			roleID:         enum.VillagerRoleID,
			expectedStatus: false,
			expectedErr:    "This role is already assigned ¯\\_(ツ)_/¯",
			newMainRoleID:  enum.VillagerRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				p.mainRoleID = enum.VillagerRoleID
				p.factionID = enum.VillagerFactionID
				p.roles[enum.VillagerRoleID] = ps.villagerRole
			},
		},
		{
			name:           "Falure (Invalid role)",
			roleID:         0,
			expectedStatus: false,
			expectedErr:    "Non-existent role ¯\\_ಠ_ಠ_/¯",
			newMainRoleID:  enum.VillagerRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				p.mainRoleID = enum.VillagerRoleID
				p.factionID = enum.VillagerFactionID
			},
		},
		{
			name:           "Ok (Doesn't update main role and faction)",
			roleID:         enum.VillagerRoleID,
			expectedStatus: true,
			newMainRoleID:  enum.HunterRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				ps.game.EXPECT().Player(ps.playerID).Return(ps.player).Times(1)
				ps.game.EXPECT().Poll(enum.VillagerFactionID).Return(ps.poll).Times(1)
				ps.poll.EXPECT().AddElectors(ps.playerID).Return(true).Times(1)
				ps.poll.EXPECT().SetWeight(ps.playerID, uint(1)).Times(1)

				p.mainRoleID = enum.HunterRoleID
				p.factionID = enum.VillagerFactionID
				p.roles[enum.HunterRoleID] = ps.hunterRole
			},
		},
		{
			name:           "Ok (Update main role)",
			roleID:         enum.HunterRoleID,
			expectedStatus: true,
			newMainRoleID:  enum.HunterRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				ps.game.EXPECT().Player(ps.playerID).Return(ps.player).Times(1)

				p.mainRoleID = enum.VillagerRoleID
				p.factionID = enum.VillagerFactionID
				p.roles[enum.VillagerRoleID] = ps.villagerRole
			},
		},
		{
			name:           "Ok (Update main role and faction)",
			roleID:         enum.WerewolfRoleID,
			expectedStatus: true,
			newMainRoleID:  enum.WerewolfRoleID,
			newFactionID:   enum.WerewolfFactionID,
			setup: func(p *player) {
				ps.game.EXPECT().Player(ps.playerID).Return(ps.player).Times(1)
				ps.game.EXPECT().Poll(enum.WerewolfFactionID).Return(ps.poll).Times(1)
				ps.poll.EXPECT().AddElectors(ps.playerID).Return(true).Times(1)
				ps.poll.EXPECT().SetWeight(ps.playerID, uint(1)).Times(1)

				p.mainRoleID = enum.VillagerRoleID
				p.factionID = enum.VillagerFactionID
				p.roles[enum.VillagerRoleID] = ps.villagerRole
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPlayer(ps.game, ps.playerID)
			test.setup(p.(*player))
			ok, err := p.AssignRole(test.roleID)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.newMainRoleID, p.MainRoleID())
			ps.Equal(test.newFactionID, p.FactionID())

			if ok {
				ps.Contains(p.RoleIDs(), test.roleID)
			} else {
				ps.Equal(test.expectedErr, err.Error())
			}
		})
	}
}

func (ps *PlayerSuite) TestRevokeRole() {
	tests := []struct {
		name           string
		roleID         enum.RoleID
		expectedStatus bool
		expectedErr    string
		newMainRoleID  enum.RoleID
		newFactionID   enum.FactionID
		setup          func(*player)
	}{
		{
			name:           "Falure (Non-existent role)",
			roleID:         99,
			expectedStatus: false,
			expectedErr:    "Non-existent role ID  ¯\\_(ツ)_/¯",
			newMainRoleID:  enum.VillagerRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				p.mainRoleID = enum.VillagerRoleID
				p.factionID = enum.VillagerFactionID
				p.roles = ps.roles
			},
		},
		{
			name:           "Falure (Must have at least one role)",
			roleID:         enum.VillagerRoleID,
			expectedStatus: false,
			expectedErr:    "Player must player at least one role ヾ(⌐■_■)ノ♪",
			newMainRoleID:  enum.VillagerRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				p.mainRoleID = enum.VillagerRoleID
				p.factionID = enum.VillagerFactionID
				p.roles[enum.VillagerRoleID] = ps.villagerRole
			},
		},
		{
			name:           "Ok (Doesn't update main role and faction)",
			roleID:         enum.VillagerRoleID,
			expectedStatus: true,
			newMainRoleID:  enum.HunterRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				p.mainRoleID = enum.HunterRoleID
				p.factionID = enum.VillagerFactionID
				p.roles = maps.Clone(ps.roles)
			},
		},
		{
			name:           "Ok (Update main role)",
			roleID:         enum.HunterRoleID,
			expectedStatus: true,
			newMainRoleID:  enum.VillagerRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {

				p.mainRoleID = enum.HunterRoleID
				p.factionID = enum.VillagerFactionID
				p.roles = maps.Clone(ps.roles)
			},
		},
		{
			name:           "Ok (Update main role and faction)",
			roleID:         enum.WerewolfRoleID,
			expectedStatus: true,
			newMainRoleID:  enum.VillagerRoleID,
			newFactionID:   enum.VillagerFactionID,
			setup: func(p *player) {
				p.mainRoleID = enum.WerewolfRoleID
				p.factionID = enum.WerewolfFactionID
				p.roles[enum.VillagerRoleID] = ps.villagerRole
				p.roles[enum.WerewolfRoleID] = ps.werewolfRole
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPlayer(ps.game, ps.playerID)
			test.setup(p.(*player))
			ok, err := p.RevokeRole(test.roleID)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.newMainRoleID, p.MainRoleID())
			ps.Equal(test.newFactionID, p.FactionID())

			if ok {
				ps.NotContains(p.RoleIDs(), test.roleID)
			} else {
				ps.Equal(test.expectedErr, err.Error())
			}
		})
	}
}

func (ps *PlayerSuite) TestUseAbility() {
	tests := []struct {
		name        string
		req         *types.UseRoleRequest
		expectedRes *types.ActionResponse
		setup       func(*player)
	}{
		{
			name: "Falure (Scheduler hasn't started yet)",
			req:  new(types.UseRoleRequest),
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Wait until game starts ノ(ジ)ー'",
			},
			setup: func(p *player) {
				ps.scheduler.EXPECT().Turn().Return(nil).Times(1)
			},
		},
		{
			name: "Falure (Not in turn)",
			req:  new(types.UseRoleRequest),
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Wait for your turn, OK??",
			},
			setup: func(p *player) {
				turn := &types.Turn{
					RoleID: 0,
				}
				ps.scheduler.EXPECT().Turn().Return(turn).Times(1)
			},
		},
		{
			name: "Ok",
			req:  new(types.UseRoleRequest),
			expectedRes: &types.ActionResponse{
				Ok: true,
			},
			setup: func(p *player) {
				turn := &types.Turn{
					RoleID: enum.VillagerRoleID,
				}
				ps.scheduler.EXPECT().Turn().Return(turn).Times(1)
				ps.villagerRole.
					EXPECT().
					UseAbility(new(types.UseRoleRequest)).
					Return(&types.ActionResponse{
						Ok: true,
					}).Times(1)
				p.roles = ps.roles
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPlayer(ps.game, ps.playerID)
			test.setup(p.(*player))
			res := p.UseAbility(test.req)

			ps.Equal(test.expectedRes, res)
		})
	}
}
