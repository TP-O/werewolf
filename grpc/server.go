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
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedGameServer
}

func (s *server) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	setting := &types.GameSetting{
		NumberWerewolves:   uint8(req.NumberWerewolves),
		TurnDuration:       uint16(req.TurnDuration),
		DiscussionDuration: uint16(req.DiscussionDuration),
		RoleIDs:            req.RoleIds,
		RequiredRoleIDs:    req.RequiredRoleIds,
		PlayerIDs:          req.PlayerIds,
	}

	if err := validator.ValidateStruct(setting); len(err.FieldViolations) != 0 {
		status := status.New(codes.InvalidArgument, "Invalid")
		ss, _ := status.WithDetails(err)

		return nil, ss.Err()
	}

	game, err := core.Manager().NewGame(setting)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	if startedAt := game.Start(); startedAt == -1 {
		return nil, grpc.Errorf(codes.AlreadyExists, "Game is running ¯\\_(ツ)_/¯")
	} else {
		return &pb.CreateGameResponse{
			Id:             game.ID(),
			StartedAt:      uint32(startedAt),
			PreprationTime: uint32(config.Game().PreparationTime),
		}, nil
	}
}

func (s *server) UseRole(ctx context.Context, req *pb.UseRoleRequest) (*pb.UseRoleResponse, error) {
	use := &types.UseRoleRequest{
		ActionID:  req.ActionId,
		TargetIDs: req.TargetIds,
		IsSkipped: req.IsSkipped,
	}

	if err := validator.ValidateStruct(use); len(err.FieldViolations) != 0 {
		status := status.New(codes.InvalidArgument, "Invalid")
		ss, _ := status.WithDetails(err)

		return nil, ss.Err()
	}

	if game := core.Manager().Game(req.GameId); game == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Game does not exist")
	} else {
		res := game.UsePlayerRole(req.PlayerId, use)

		if data, err := json.Marshal(res.Data); err != nil {
			return &pb.UseRoleResponse{
					Ok:        res.Ok,
					IsSkipped: res.IsSkipped,
					Message:   res.Message,
					Data:      "",
				},
				nil
		} else {
			return &pb.UseRoleResponse{
					Ok:        res.Ok,
					IsSkipped: res.IsSkipped,
					Message:   res.Message,
					Data:      string(data),
				},
				nil
		}
	}
}

func (s *server) ConnectPlayer(ctx context.Context, req *pb.ConnectPlayerRequest) (*pb.ConnectPlayerResponse, error) {
	if game := core.Manager().Game(req.GameId); game == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Game does not exist")
	} else if req.IsExited {
		return &pb.ConnectPlayerResponse{
			Ok: game.ExitPlayer(req.PlayerId),
		}, nil
	} else {
		return &pb.ConnectPlayerResponse{
			Ok: game.ConnectPlayer(req.PlayerId, req.IsConnected),
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
