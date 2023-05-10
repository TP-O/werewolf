package action

import (
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RecognizeSuite struct {
	suite.Suite
	IdentifiedFactionId types.FactionId
	IdentifiedRoleId    types.RoleId
	actorId             types.PlayerId
	IdentifiedIds       []types.PlayerId
}

func TestIdentifySuite(t *testing.T) {
	suite.Run(t, new(RecognizeSuite))
}

func (rs *RecognizeSuite) SetupSuite() {
	rs.actorId = types.PlayerId("1")
	rs.IdentifiedIds = []types.PlayerId{"1", "2"}
	rs.IdentifiedFactionId = constants.WerewolfFactionId
	rs.IdentifiedRoleId = constants.WerewolfRoleId
}

func (rs RecognizeSuite) TestNewFactionIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	fIdent := NewFactionIdentify(world, rs.IdentifiedFactionId).(*identify)

	rs.Equal(IdentifyActionId, fIdent.Id())
	rs.Equal(rs.IdentifiedFactionId, fIdent.FactionId)
	rs.Empty(fIdent.Faction)
	rs.Equal(types.RoleId(0), fIdent.RoleId)
	rs.Empty(fIdent.Role)
	rs.False(fIdent.isIdentified)
}

func (rs RecognizeSuite) TestNewRoleIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	Ident := NewRoleIdentify(world, rs.IdentifiedRoleId).(*identify)

	rs.Equal(IdentifyActionId, Ident.Id())
	rs.Equal(types.FactionId(0), Ident.FactionId)
	rs.Empty(Ident.Faction)
	rs.Equal(rs.IdentifiedRoleId, Ident.RoleId)
	rs.Empty(Ident.Role)
	rs.False(Ident.isIdentified)
}

func (rs RecognizeSuite) TestValIdateFactionIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	expectErr := "You already recognized everyone ¯\\(º_o)/¯"

	Ident := NewFactionIdentify(world, rs.IdentifiedFactionId).(*identify)
	Ident.isIdentified = true

	err := Ident.validate(&types.ActionRequest{})
	rs.Equal(expectErr, err.Error())
}

func (rs RecognizeSuite) TestPerformFactionIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	world.EXPECT().AlivePlayerIdsWithFactionId(rs.IdentifiedFactionId).
		Return(rs.IdentifiedIds)

	req := &types.ActionRequest{
		ActorId: rs.actorId,
	}
	expectRes := types.ActionResponse{
		Ok:   true,
		Data: rs.IdentifiedIds,
	}

	Ident := NewFactionIdentify(world, rs.IdentifiedFactionId).(*identify)
	res := Ident.perform(req)

	rs.Equal(expectRes, res)
	rs.Equal(rs.IdentifiedIds, Ident.Faction)
	rs.True(Ident.isIdentified)
}

func (rs RecognizeSuite) TestValIdateRoleIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	expectErr := "You already recognized everyone ¯\\(º_o)/¯"

	Ident := NewRoleIdentify(world, rs.IdentifiedRoleId).(*identify)
	Ident.isIdentified = true

	err := Ident.validate(&types.ActionRequest{})
	rs.Equal(expectErr, err.Error())
}

func (rs RecognizeSuite) TestPerformRoleIdentify() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	world := mock_game_logic.NewMockWorld(ctrl)

	world.EXPECT().AlivePlayerIdsWithRoleId(rs.IdentifiedRoleId).
		Return(rs.IdentifiedIds)

	req := &types.ActionRequest{
		ActorId: rs.actorId,
	}
	expectRes := types.ActionResponse{
		Ok:   true,
		Data: rs.IdentifiedIds,
	}

	Ident := NewRoleIdentify(world, rs.IdentifiedRoleId).(*identify)
	res := Ident.perform(req)

	rs.Equal(expectRes, res)
	rs.Equal(rs.IdentifiedIds, Ident.Role)
	rs.True(Ident.isIdentified)
}
