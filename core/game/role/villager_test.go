package role

// import (
// 	"testing"
// 	"uwwolf/game/enum"
// 	gamemock "uwwolf/mock/game"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// )

// type VillagerSuite struct {
// 	suite.Suite
// 	ctrl     *gomock.Controller
// 	game     *gamemock.MockGame
// 	player   *gamemock.MockPlayer
// 	poll     *gamemock.MockPoll
// 	playerID enum.PlayerID
// }

// func TestVillagerSuite(t *testing.T) {
// 	suite.Run(t, new(VillagerSuite))
// }

// func (vs *VillagerSuite) SetupSuite() {
// 	vs.playerID = "1"
// }

// func (vs *VillagerSuite) SetupTest() {
// 	vs.ctrl = gomock.NewController(vs.T())
// 	vs.game = gamemock.NewMockGame(vs.ctrl)
// 	vs.player = gamemock.NewMockPlayer(vs.ctrl)
// 	vs.poll = gamemock.NewMockPoll(vs.ctrl)
// 	vs.game.EXPECT().Player(vs.playerID).Return(vs.player).AnyTimes()
// }

// func (vs *VillagerSuite) TearDownTest() {
// 	vs.ctrl.Finish()
// }

// func (vs *VillagerSuite) TestNewVillager() {
// 	tests := []struct {
// 		name        string
// 		expectedErr string
// 		returnNil   bool
// 		setup       func()
// 	}{
// 		{
// 			name:        "Failure (Poll does not exist)",
// 			returnNil:   true,
// 			expectedErr: "Poll does not exist ¯\\_(ツ)_/¯",
// 			setup: func() {
// 				vs.game.EXPECT().Poll(enum.VillagerFactionID).Return(nil).Times(1)
// 			},
// 		},
// 		{
// 			name:        "Failure (Unable to join to the poll)",
// 			returnNil:   true,
// 			expectedErr: "Unable to join to the poll ಠ_ಠ",
// 			setup: func() {
// 				vs.game.EXPECT().Poll(enum.VillagerFactionID).Return(vs.poll).Times(1)
// 				vs.poll.EXPECT().AddElectors(vs.playerID).Return(false).Times(1)
// 			},
// 		},
// 		{
// 			name:      "Ok",
// 			returnNil: false,
// 			setup: func() {
// 				vs.game.EXPECT().Poll(enum.VillagerFactionID).Return(vs.poll).Times(1)
// 				vs.poll.EXPECT().AddElectors(vs.playerID).Return(true).Times(1)
// 				vs.poll.EXPECT().SetWeight(vs.playerID, uint(1)).Times(1)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		vs.Run(test.name, func() {
// 			test.setup()
// 			villager, err := NewVillager(vs.game, vs.playerID)

// 			if test.returnNil {
// 				vs.Nil(villager)
// 				vs.NotNil(err)
// 				vs.Equal(test.expectedErr, err.Error())
// 			} else {
// 				vs.Nil(err)
// 				vs.Equal(enum.VillagerRoleID, villager.ID())
// 				vs.Equal(enum.DayPhaseID, villager.PhaseID())
// 				vs.Equal(enum.VillagerFactionID, villager.FactionID())
// 				vs.Equal(enum.VillagerTurnPriority, villager.Priority())
// 				vs.Equal(enum.FirstRound, villager.BeginRound())
// 				vs.Equal(enum.Unlimited, villager.ActiveLimit(enum.VoteActionID))
// 			}
// 		})
// 	}
// }
