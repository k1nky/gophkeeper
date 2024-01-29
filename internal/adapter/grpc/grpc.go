package grpc

import (
	"context"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	pb "github.com/k1nky/gophkeeper/internal/protocol/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	pb.UnimplementedKeeperServer
	auth   authService
	keeper keeperService
}

func New(auth authService, keeper keeperService) *Adapter {
	return &Adapter{
		auth:   auth,
		keeper: keeper,
	}
}

func (a *Adapter) GetSecretMeta(ctx context.Context, in *pb.GetSecretRequest) (*pb.Meta, error) {
	m, err := a.keeper.GetSecretMeta(ctx, vault.UniqueKey(in.Key))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Meta{
		Id:    string(m.Key),
		Extra: m.Extra,
	}, nil
}

func (a *Adapter) GetSecretData(in *pb.GetSecretRequest, stream pb.Keeper_GetSecretDataServer) error {
	reader, err := a.keeper.GetSecretData(stream.Context(), vault.UniqueKey(in.Key))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer reader.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return status.Error(codes.Unknown, "")
		}
		chunk := &pb.Data{
			ChunkData: buffer[:n],
		}
		if err := stream.Send(chunk); err != nil {
			return status.Error(codes.Unknown, "")
		}
	}

	return nil
}

func (a *Adapter) PutSecret(stream pb.Keeper_PutSecretServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Unknown, "")
	}
	claims, _ := user.GetEffectiveUser(stream.Context())
	meta := vault.Meta{
		UserID: claims.ID,
		Extra:  req.GetMeta().Extra,
	}
	buf := vault.NewBytesBuffer(nil)
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return status.Error(codes.Unknown, "")
		}
		chunk := req.GetChunkData().ChunkData
		buf.Write(chunk)
	}
	data := vault.NewDataReader(buf)
	m, err := a.keeper.PutSecret(stream.Context(), meta, data)
	if err != nil {
		return status.Error(codes.Unknown, "")
	}
	return stream.SendAndClose(&pb.Meta{
		Id:    string(m.Key),
		Extra: m.Extra,
	})
}

func (a *Adapter) ListSecrets(ctx context.Context, in *pb.ListSecretRequest) (*pb.ListSecretResponse, error) {
	secrets, err := a.keeper.ListSecretsByUser(ctx, user.ID(in.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	list := &pb.ListSecretResponse{}
	for _, v := range secrets {
		list.Meta = append(list.Meta, &pb.Meta{
			Id:    string(v.Key),
			Extra: v.Extra,
		})
	}
	return list, nil
}
