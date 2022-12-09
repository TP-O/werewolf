package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"uwwolf/config"
	"uwwolf/game/core"
	"uwwolf/game/types"
	"uwwolf/grpc/pb"
	"uwwolf/validator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type server struct {
	pb.UnimplementedGameServer
}

func (s *server) CreateGame(ctx context.Context, req *pb.GameRequest) (*pb.GameResponse, error) {
	setting := types.GameSetting{
		NumberWerewolves:   uint8(req.NumberOfWerewolves),
		TurnDuration:       uint16(req.TurnDuration),
		DiscussionDuration: uint16(req.DiscussionDuration),
		RoleIDs:            req.RoleIds,
		RequiredRoleIDs:    req.RequiredRoleIds,
		PlayerIDs:          req.PlayerIds,
	}

	if err := validator.ValidateStruct(setting); err != nil {
		s, _ := json.Marshal(err)

		return nil, grpc.Errorf(codes.InvalidArgument, string(s))
	}

	game, err := core.Manager().NewGame(&setting)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	if startedAt := game.Start(); startedAt == -1 {
		return nil, grpc.Errorf(codes.AlreadyExists, "Game is running ¯\\_(ツ)_/¯")
	} else {
		return &pb.GameResponse{
			Id:             game.ID(),
			StartedAt:      uint32(startedAt),
			PreprationTime: uint32(config.Game().PreparationTime),
		}, nil
	}
}

func Start() {
	port := strconv.Itoa(int(config.App().Port))
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Unable to open TCP server %v", err)
	} else {
		log.Printf("TCP server is listening on port %s", port)
	}

	s := grpc.NewServer()
	pb.RegisterGameServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Unable to serve gRPC server: %v", err)
	}
}
