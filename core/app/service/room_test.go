package service_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"uwwolf/app/data"
	"uwwolf/app/service"
	"uwwolf/game/types"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
)

type RoomServiceSuite struct {
	suite.Suite
	roomID   string
	playerID types.PlayerID
}

func (rss *RoomServiceSuite) SetupSuite() {
	rss.roomID = "room_id"
	rss.playerID = types.PlayerID("1")
}

func TestRoomServiceSuite(t *testing.T) {
	suite.Run(t, new(RoomServiceSuite))
}

func (rss RoomServiceSuite) TestNewRoomService() {
	rdb, _ := redismock.NewClusterMock()

	svc := service.NewRoomService(rdb)

	rss.Require().NotNil(svc)
	rss.Require().False(reflect.ValueOf(svc).Elem().FieldByName("rdb").IsNil())
}

func (rss RoomServiceSuite) TestPlayerWaitingRoom() {
	tests := []struct {
		name           string
		expectedResult data.WaitingRoom
		expectedStatus bool
		setup          func(rdb redismock.ClusterClientMock)
	}{
		{
			name:           "Failure (Encoded failed)",
			expectedResult: data.WaitingRoom{},
			expectedStatus: false,
			setup: func(rdb redismock.ClusterClientMock) {
				rdb.ExpectEval(service.GetWaitingRoomScript, []string{}, rss.playerID).
					SetVal("Invalid json")
			},
		},
		{
			name: "Ok",
			expectedResult: data.WaitingRoom{
				ID:        rss.roomID,
				OwnerID:   rss.playerID,
				PlayerIDs: []types.PlayerID{rss.playerID},
			},
			expectedStatus: true,
			setup: func(rdb redismock.ClusterClientMock) {
				encodedExpectedResult, _ := json.Marshal(data.WaitingRoom{
					ID:        rss.roomID,
					OwnerID:   rss.playerID,
					PlayerIDs: []types.PlayerID{rss.playerID},
				})

				rdb.ExpectEval(service.GetWaitingRoomScript, []string{}, rss.playerID).
					SetVal(string(encodedExpectedResult))
			},
		},
	}

	for _, test := range tests {
		rss.Run(test.name, func() {
			rdb, rmock := redismock.NewClusterMock()

			test.setup(rmock)

			svc := service.NewRoomService(rdb)
			room, ok := svc.PlayerWaitingRoom(rss.playerID)

			rss.Require().Equal(test.expectedResult, room)
			rss.Require().Equal(test.expectedStatus, ok)
		})
	}
}
