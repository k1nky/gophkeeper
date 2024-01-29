package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/k1nky/gophkeeper/internal/adapter/grpc/mock"
	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	pb "github.com/k1nky/gophkeeper/internal/protocol/proto"
)

type adapterTestSuite struct {
	suite.Suite
	a      *Adapter
	s      *grpc.Server
	l      *bufconn.Listener
	auth   *mock.MockauthService
	keeper *mock.MockkeeperService
}

func TestStore(t *testing.T) {
	suite.Run(t, new(adapterTestSuite))
}

func (suite *adapterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.auth = mock.NewMockauthService(ctrl)
	suite.keeper = mock.NewMockkeeperService(ctrl)
	suite.a = New(suite.auth, suite.keeper)
	suite.l = bufconn.Listen(1024 * 1024)
	suite.s = grpc.NewServer()
	pb.RegisterKeeperServer(suite.s, suite.a)
	go func() {
		if err := suite.s.Serve(suite.l); err != nil {
			suite.FailNow(err.Error())
			return
		}
	}()
}

func (suite *adapterTestSuite) TearDownTest() {
	suite.s.Stop()
	suite.l.Close()
}

func (suite *adapterTestSuite) dial(ctx context.Context, conn string) (net.Conn, error) {
	return suite.l.Dial()
}

func (suite *adapterTestSuite) TestListSecrets() {
	ctx := user.NewContextWithClaims(context.Background(), user.PrivateClaims{
		ID:    1,
		Login: "u",
	})
	conn, err := grpc.DialContext(ctx, "buffer", grpc.WithContextDialer(suite.dial), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.FailNow(err.Error())
		return
	}
	defer conn.Close()
	client := pb.NewKeeperClient(conn)
	expected := vault.List{
		{
			UserID: 1,
			Key:    vault.NewUniqueKey(),
			Extra:  "extra data",
		},
	}
	suite.keeper.EXPECT().ListSecretsByUser(gomock.Any(), gomock.Any()).Return(expected, nil)
	resp, err := client.ListSecrets(ctx, &pb.ListSecretRequest{
		UserId: 1,
	})
	suite.NoError(err)
	got := vault.List{}
	for _, v := range resp.Meta {
		got = append(got, vault.Meta{
			UserID: 1, Key: vault.UniqueKey(v.Id), Extra: v.Extra,
		})
	}
	suite.NoError(err)
	suite.Equal(expected, got)
}
