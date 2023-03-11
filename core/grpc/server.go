package grpc

import (
	"context"
	"log"
	"net"
	"strconv"
	"uwwolf/grpc/pb"
	"uwwolf/util"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGameServer
}

func (s *server) StartGame(ctx context.Context, req *pb.StartGameRequest) (*pb.StartGameResponse, error) {
	// setting := &types.GameSetting{
	// 	NumberWerewolves:   uint8(req.NumberOfWerewolves),
	// 	TurnDuration:       uint16(req.TurnDuration),
	// 	DiscussionDuration: uint16(req.DiscussionDuration),
	// 	RoleIDs:            req.RoleIds,
	// 	RequiredRoleIDs:    req.RequiredRoleIds,
	// 	PlayerIDs:          req.PlayerIds,
	// }

	// if err := validator.ValidateStruct(setting); len(err.FieldViolations) != 0 {
	// 	status := status.New(codes.InvalidArgument, "Invalid")
	// 	ss, _ := status.WithDetails(err)

	// 	return nil, ss.Err()
	// }

	// game, err := core.Manager().NewGame(setting)

	// if err != nil {
	// 	return nil, grpc.Errorf(codes.Internal, err.Error())
	// }

	// if startedAt := game.Start(); startedAt == -1 {
	// 	return nil, grpc.Errorf(codes.AlreadyExists, "Game is running ¯\\_(ツ)_/¯")
	// } else {
	// 	return &pb.StartGameResponse{
	// 		GameId:              game.ID(),
	// 		StartedAt:           uint32(startedAt),
	// 		PreparationDuration: uint32(util.Config().Game.PreparationDuration),
	// 	}, nil
	// }

	return nil, nil
}

func (s *server) PerformAction(ctx context.Context, req *pb.PerformActionRequest) (*pb.PerformActionResponse, error) {
	// use := &types.UseRoleRequest{
	// 	ActionID:  req.ActionId,
	// 	TargetIDs: req.TargetIds,
	// 	IsSkipped: req.IsSkipped,
	// }

	// if err := validator.ValidateStruct(use); len(err.FieldViolations) != 0 {
	// 	status := status.New(codes.InvalidArgument, "Invalid")
	// 	ss, _ := status.WithDetails(err)

	// 	return nil, ss.Err()
	// }

	// if game := core.Manager().Game(req.GameId); game == nil {
	// 	return nil, grpc.Errorf(codes.InvalidArgument, "Game does not exist")
	// } else {
	// 	res := game.UsePlayerRole(req.PlayerId, use)
	// 	if _, err := proto.Marshal(proto.MessageV1("cc")); err != nil {
	// 		return &pb.PerformActionResponse{
	// 				Ok:        res.Ok,
	// 				IsSkipped: res.IsSkipped,
	// 				Message:   res.Message,
	// 				Data:      "",
	// 			},
	// 			nil
	// 	} else {
	// 		return &pb.PerformActionResponse{
	// 				Ok:        res.Ok,
	// 				IsSkipped: res.IsSkipped,
	// 				Message:   res.Message,
	// 				Data:      "",
	// 			},
	// 			nil
	// 	}
	// }

	return nil, nil
}

func Start() {
	port := strconv.Itoa(int(util.Config().App.Port))
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
