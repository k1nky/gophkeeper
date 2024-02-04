package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	pb "github.com/k1nky/gophkeeper/internal/protocol/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	// UnimplementedKeeperServer must be embedded to have forward compatible implementations.
	pb.UnimplementedKeeperServer
	auth   authService
	keeper keeperService
	log    logger
}

var (
	ErrUnexpected = errors.New("unexpected error")
)

func New(auth authService, keeper keeperService, log logger) *Adapter {
	return &Adapter{
		auth:   auth,
		keeper: keeper,
		log:    log,
	}
}

func (a *Adapter) GetSecretMeta(ctx context.Context, in *pb.GetSecretMetaRequest) (*pb.Meta, error) {
	var (
		m   *vault.Meta
		err error
	)
	switch key := in.Key.(type) {
	case *pb.GetSecretMetaRequest_Alias:
		m, err = a.keeper.GetSecretMetaByAlias(ctx, key.Alias)
	case *pb.GetSecretMetaRequest_Id:
		m, err = a.keeper.GetSecretMeta(ctx, vault.MetaID(key.Id))
	}
	if err != nil {
		a.log.Errorf("grpc: GetSecretMeta: %v", err)
		return nil, status.Error(codes.Internal, ErrUnexpected.Error())
	}
	if m == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("[%s] not found", in.Key))
	}
	return &pb.Meta{
		Id:    string(m.ID),
		Extra: m.Extra,
	}, nil
}

func (a *Adapter) GetSecretData(in *pb.GetSecretDataRequest, stream pb.Keeper_GetSecretDataServer) error {
	reader, err := a.keeper.GetSecretData(stream.Context(), vault.MetaID(in.Id))
	if err != nil {
		a.log.Errorf("grpc: GetSecretData: %v", err)
		return status.Error(codes.Internal, ErrUnexpected.Error())
	}
	defer reader.Close()

	// отправлять данные секрета будем по частям в потоке
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			a.log.Errorf("grpc: GetSecretData: reading data %v", err)
			return status.Error(codes.Unknown, "reading data")
		}
		chunk := &pb.Data{
			ChunkData: buffer[:n],
		}
		if err := stream.Send(chunk); err != nil {
			a.log.Errorf("grpc: GetSecretData: sending chunk %v", err)
			return status.Error(codes.Unknown, "sending chunk")
		}
	}

	return nil
}

func (a *Adapter) PutSecret(stream pb.Keeper_PutSecretServer) error {
	// первым запросом получаем мета-данные секрета
	req, err := stream.Recv()
	if err != nil {
		a.log.Errorf("grpc: PutSecret: %v", err)
		return status.Error(codes.Unknown, ErrUnexpected.Error())
	}
	claims, _ := user.GetEffectiveUser(stream.Context())
	meta := vault.Meta{
		ID:     vault.MetaID(req.GetMeta().Id),
		Alias:  req.GetMeta().Alias,
		UserID: claims.ID,
		Extra:  req.GetMeta().Extra,
	}
	// Данные секрета будут приходить частями в потоке stream.
	// С помощью Pipe будем передавать данные также по частям в хранилище.
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				a.log.Errorf("grpc: PutSecret: receiving chunk %v", err)
				return
			}
			chunk := req.GetChunkData().ChunkData
			w.Write(chunk)
		}
	}()
	data := vault.NewDataReader(r)
	m, err := a.keeper.PutSecret(stream.Context(), meta, data)
	if err != nil {
		a.log.Errorf("grpc: PutSecret: saving data %v", err)
		return status.Error(codes.Unknown, "saving data")
	}
	// отправляем в ответ мета-данные добавленного секрета
	return stream.SendAndClose(&pb.Meta{
		Id:    string(m.ID),
		Alias: m.Alias,
		Extra: m.Extra,
	})
}

func (a *Adapter) ListSecrets(ctx context.Context, in *pb.ListSecretRequest) (*pb.ListSecretResponse, error) {
	secrets, err := a.keeper.ListSecretsByUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	list := &pb.ListSecretResponse{}
	for _, v := range secrets {
		list.Meta = append(list.Meta, &pb.Meta{
			Id:    string(v.ID),
			Extra: v.Extra,
		})
	}
	return list, nil
}
