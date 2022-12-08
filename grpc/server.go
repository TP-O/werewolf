package grpc

import (
	"context"
	"log"
	"net"
	"time"
	"uwwolf/config"
	"uwwolf/game/core"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	"uwwolf/grpc/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGameServer
}

func (s *server) CreateGame(ctx context.Context, req *pb.GameRequest) (*pb.GameResponse, error) {
	gameID := "aaa"
	core.Manager().AddGame(enum.GameID(gameID), &types.GameSetting{
		NumberOfWerewolves: uint8(req.NumberOfWerewolves),
		TurnDuration:       time.Duration(req.TurnDuration),
		DiscussionDuration: time.Duration(req.DiscussionDuration),
		RoleIDs:            req.RoleIds,
		RequiredRoleIDs:    req.RequiredRoleIds,
		PlayerIDs:          req.PlayerIds,
	})
	startedAt := core.Manager().Game(enum.GameID(gameID)).Start()

	return &pb.GameResponse{
		Id:             gameID,
		StartedAt:      uint32(startedAt),
		PreprationTime: uint32(config.Game().PreparationTime),
	}, nil
}

func Run() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterGameServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
