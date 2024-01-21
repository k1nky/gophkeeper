package store

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophkeeper/internal/adapter/store/mock"
	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/stretchr/testify/suite"
)

type adapterTestSuite struct {
	suite.Suite
	a      *Adapter
	mstore *mock.MockMetaStore
	ostore *mock.MockObjectStore
}

func TestStore(t *testing.T) {
	suite.Run(t, new(adapterTestSuite))
}

func (suite *adapterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mstore = mock.NewMockMetaStore(ctrl)
	suite.ostore = mock.NewMockObjectStore(ctrl)
	suite.a = New(suite.mstore, suite.ostore)
}

func (suite *adapterTestSuite) TestOpen() {
	ctx := context.TODO()
	suite.ostore.EXPECT().Open(ctx).Return(nil)
	suite.mstore.EXPECT().Open(ctx).Return(nil)
	err := suite.a.Open(ctx)
	suite.Assert().NoError(err)
}

func (suite *adapterTestSuite) TestOpenWithError() {
	ctx := context.TODO()
	suite.mstore.EXPECT().Close().Return(nil)
	suite.ostore.EXPECT().Open(ctx).Return(errors.New("could not open"))
	suite.mstore.EXPECT().Open(ctx).Return(nil)
	err := suite.a.Open(ctx)
	suite.Assert().Error(err)
}

func (suite *adapterTestSuite) TestNewUser() {
	ctx := context.TODO()
	u := user.User{
		Login:    "login",
		Password: "password",
		ID:       1,
	}
	suite.mstore.EXPECT().NewUser(ctx, gomock.Any()).Return(&u, nil)
	newUser, err := suite.a.NewUser(ctx, u)
	suite.Assert().NoError(err)
	suite.Assert().Equal(u, *newUser)
}

func (suite *adapterTestSuite) TestPutSecret() {
	ctx := context.TODO()
	s := vault.Secret{
		Meta: vault.Meta{
			UserID: 1,
			Extra:  "extra data",
		},
		Data: bytes.NewBufferString("some super secret"),
	}

	suite.mstore.EXPECT().NewMeta(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	suite.ostore.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	newSecret, err := suite.a.PutSecret(ctx, s)
	suite.Assert().NoError(err)
	s.Key = newSecret.Key
	suite.Assert().Equal(s, *newSecret)
}
