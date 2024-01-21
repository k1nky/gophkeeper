package grpc

import (
	"context"

	pb "github.com/k1nky/gophkeeper/internal/proto"
)

type Adapter struct {
	pb.UnimplementedKeeperServer
}

func (a *Adapter) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{}, nil
}
