package keeper

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	log "github.com/k1nky/gophkeeper/internal/logger"
	"github.com/k1nky/gophkeeper/internal/service/keeper/mock"
	"github.com/stretchr/testify/suite"
)

type keeperServiceTestSuite struct {
	suite.Suite
	store *mock.Mockstorage
	svc   *Service
}

func TestKeeperService(t *testing.T) {
	suite.Run(t, new(keeperServiceTestSuite))
}

func (suite *keeperServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.Suite.T())
	suite.store = mock.NewMockstorage(ctrl)
	suite.svc = New(suite.store, &log.Blackhole{})
}

func (suite *keeperServiceTestSuite) TestGetSecretData() {
	metaID := vault.NewMetaID()
	expected := []byte("some secret text")
	data := vault.NewDataReader(vault.NewBytesBuffer(expected))
	suite.store.EXPECT().GetSecretMetaByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&vault.Meta{
		ID: metaID,
	}, nil)
	suite.store.EXPECT().GetSecretData(gomock.Any(), gomock.Any(), gomock.Any()).Return(data, nil)
	r, err := suite.svc.GetSecretData(context.TODO(), metaID)
	suite.Assert().NoError(err)
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(r)
	suite.NoError(err)
	suite.Equal(expected, buf.Bytes())
}

func (suite *keeperServiceTestSuite) TestGetSecretDataNotFound() {
	metaID := vault.NewMetaID()
	suite.store.EXPECT().GetSecretMetaByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	r, err := suite.svc.GetSecretData(context.TODO(), metaID)
	suite.Assert().NoError(err)
	suite.Nil(r)
}

func (suite *keeperServiceTestSuite) TestGetSecretDataWithMetaError() {
	metaID := vault.NewMetaID()
	suite.store.EXPECT().GetSecretMetaByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))
	r, err := suite.svc.GetSecretData(context.TODO(), metaID)
	suite.Error(err)
	suite.Nil(r)
}

func (suite *keeperServiceTestSuite) TestGetSecretMeta() {
	expected := &vault.Meta{
		ID:    vault.NewMetaID(),
		Alias: "alias#1",
	}
	suite.store.EXPECT().GetSecretMetaByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected, nil)
	got, err := suite.svc.GetSecretMeta(context.TODO(), expected.ID)
	suite.NoError(err)
	suite.Equal(expected, got)
}

func (suite *keeperServiceTestSuite) TestGetSecretMetaWithError() {
	suite.store.EXPECT().GetSecretMetaByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))
	got, err := suite.svc.GetSecretMeta(context.TODO(), vault.NewMetaID())
	suite.Error(err)
	suite.Nil(got)
}

func (suite *keeperServiceTestSuite) TestGetSecretMetaByAlias() {
	expected := &vault.Meta{
		ID:    vault.NewMetaID(),
		Alias: "alias#1",
	}
	suite.store.EXPECT().GetSecretMetaByAlias(gomock.Any(), gomock.Any(), gomock.Any()).Return(expected, nil)
	got, err := suite.svc.GetSecretMetaByAlias(context.TODO(), expected.Alias)
	suite.NoError(err)
	suite.Equal(expected, got)
}

func (suite *keeperServiceTestSuite) TestGetSecretMetaByAliasWithError() {
	suite.store.EXPECT().GetSecretMetaByAlias(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))
	got, err := suite.svc.GetSecretMetaByAlias(context.TODO(), "alias2")
	suite.Error(err)
	suite.Nil(got)
}

func (suite *keeperServiceTestSuite) TestPutSecret() {
	expected := vault.Meta{
		ID:    vault.NewMetaID(),
		Alias: "alias#2",
	}
	suite.store.EXPECT().PutSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(&expected, nil)
	got, err := suite.svc.PutSecret(context.TODO(), expected, nil)
	suite.NoError(err)
	suite.Equal(expected, *got)
}

func (suite *keeperServiceTestSuite) TestListSecretsByUser() {
	expected := vault.List{
		vault.Meta{
			ID:    vault.NewMetaID(),
			Alias: "alias1",
		},
		vault.Meta{
			ID:    vault.NewMetaID(),
			Alias: "alias2",
		},
	}
	suite.store.EXPECT().ListSecretsByUser(gomock.Any(), gomock.Any()).Return(expected, nil)
	got, err := suite.svc.ListSecretsByUser(context.TODO())
	suite.NoError(err)
	suite.ElementsMatch(expected, got)
}
