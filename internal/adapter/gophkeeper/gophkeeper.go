package gophkeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	pb "github.com/k1nky/gophkeeper/internal/protocol/proto"
	"github.com/k1nky/gophkeeper/internal/protocol/rest"
	"github.com/k1nky/gophkeeper/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// MaxRequestTimeout максимальный таймаут для унарного запроса
	MaxRequestTimeout = 10 * time.Second
	// MaxAttempts максимальное количество попыток переподключения
	MaxAttempts = 2
)

type Adapter struct {
	Endpoint string
	Path     string
	cc       *grpc.ClientConn
	token    string
}

func New(endpoint string, path string) *Adapter {
	return &Adapter{
		Endpoint: endpoint,
		Path:     path,
	}
}

func NewPBMeta(m vault.Meta) *pb.Meta {
	return &pb.Meta{
		Id:        string(m.ID),
		Extra:     m.Extra,
		Alias:     m.Alias,
		Type:      int32(m.Type),
		Revision:  m.Revision,
		IsDeleted: m.IsDeleted,
	}
}

func NewMeta(pbm *pb.Meta) *vault.Meta {
	return &vault.Meta{
		Alias:     pbm.Alias,
		Extra:     pbm.Extra,
		ID:        vault.MetaID(pbm.Id),
		Type:      vault.SecretType(pbm.Type),
		Revision:  pbm.Revision,
		IsDeleted: pbm.IsDeleted,
	}
}

func (a *Adapter) Open(ctx context.Context) error {
	var retryPolicy = fmt.Sprintf(`{
		"methodConfig": [{
			"name": [{"service": ""}],
			"waitForReady": true,
			"retryPolicy": {
				"MaxAttempts": %d,
				"InitialBackoff": ".01s",
				"MaxBackoff": ".01s",
				"BackoffMultiplier": 1.0,
				"RetryableStatusCodes": [ "UNAVAILABLE" ]
			}
		}]
	}`, MaxAttempts)
	u, _ := url.Parse(a.Endpoint)
	unaryInterceptors := []grpc.UnaryClientInterceptor{a.AuthorizationUnaryInterceptor()}
	streamInterceptors := []grpc.StreamClientInterceptor{a.AuthorizationStreamInterceptor()}
	conn, err := grpc.Dial(fmt.Sprintf(u.Host),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
	)
	if err != nil {
		return err
	}
	a.cc = conn
	return nil
}

func (a *Adapter) SetToken(token string) {
	a.token = token
}

func (a *Adapter) Login(ctx context.Context, username string, password string) (*user.PrivateClaims, error) {
	u := rest.LoginUserRequest{
		Login:    username,
		Password: password,
	}
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(&u); err != nil {
		return nil, err
	}
	endpoint, err := url.JoinPath(a.Endpoint, "/api/user/login")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	token := resp.Header.Get("Authorization")
	claims := &auth.Claims{}
	if _, err := jwt.ParseWithClaims(token, claims, nil); !errors.Is(err, jwt.ErrTokenUnverifiable) {
		return nil, err
	}
	a.token = token
	return &claims.PrivateClaims, nil
}

func (a *Adapter) ListSecrets(ctx context.Context) (vault.List, error) {
	cli := pb.NewKeeperClient(a.cc)
	resp, err := cli.ListSecrets(ctx, &pb.ListSecretRequest{})
	if err != nil {
		return nil, err
	}
	list := make(vault.List, 0)
	for _, v := range resp.Meta {
		list = append(list, *NewMeta(v))
	}
	return list, nil
}

func (a *Adapter) GetSecretMeta(ctx context.Context, id vault.MetaID) (*vault.Meta, error) {
	cli := pb.NewKeeperClient(a.cc)
	meta, err := cli.GetSecretMeta(ctx, &pb.GetSecretMetaRequest{
		Key: &pb.GetSecretMetaRequest_Id{
			Id: string(id),
		},
	})
	if err != nil {
		return nil, err
	}
	return NewMeta(meta), nil
}

func (a *Adapter) GetSecretMetaByAlias(ctx context.Context, alias string) (*vault.Meta, error) {
	cli := pb.NewKeeperClient(a.cc)
	meta, err := cli.GetSecretMeta(ctx, &pb.GetSecretMetaRequest{
		Key: &pb.GetSecretMetaRequest_Alias{
			Alias: alias,
		},
	})
	if err != nil {
		return nil, err
	}
	return NewMeta(meta), nil
}

func (a *Adapter) GetSecretData(ctx context.Context, id vault.MetaID, w io.Writer) error {
	cli := pb.NewKeeperClient(a.cc)
	stream, err := cli.GetSecretData(ctx, &pb.GetSecretDataRequest{
		Id: string(id),
	})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		chunk := resp.GetChunkData()
		if _, err := w.Write(chunk); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) PutSecret(ctx context.Context, meta vault.Meta, r io.Reader) (*vault.Meta, error) {
	cli := pb.NewKeeperClient(a.cc)
	stream, err := cli.PutSecret(ctx)
	if err != nil {
		return nil, err
	}
	defer stream.CloseAndRecv()
	req := &pb.PutSecretRequest{
		Data: &pb.PutSecretRequest_Meta{
			Meta: NewPBMeta(meta),
		},
	}
	if err = stream.Send(req); err != nil {
		return nil, err
	}
	buffer := make([]byte, 1024)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		req := &pb.PutSecretRequest{
			Data: &pb.PutSecretRequest_ChunkData{
				ChunkData: &pb.Data{
					ChunkData: buffer[:n],
				},
			},
		}
		if err = stream.Send(req); err != nil {
			return nil, err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return NewMeta(resp), nil
}
