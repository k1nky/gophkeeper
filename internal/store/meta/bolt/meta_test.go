package bolt

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
)

type metaTestSuite struct {
	suite.Suite
	bs *BoltStorage
}

func (suite *metaTestSuite) SetupTest() {
	var err error
	rootDir := suite.T().TempDir()
	if suite.bs, err = openTestDB(rootDir); err != nil {
		suite.FailNow(err.Error())
		return
	}
	suite.bs.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tb("meta"))
		suite.Assert().NotNil(b)
		return nil
	})
}

func (suite *metaTestSuite) TestNewMeta() {
	uk := vault.NewUniqueKey()
	m := vault.Meta{
		UserID: 1,
		Extra:  "some extra",
	}
	got, err := suite.bs.NewMeta(context.TODO(), uk, m)
	suite.Assert().NoError(err)
	suite.Assert().Equal(m, *got)
}

func (suite *metaTestSuite) TestGetMetaNotExists() {
	m, err := suite.bs.GetMeta(context.TODO(), vault.UniqueKey("not_exists"))
	suite.Assert().NoError(err)
	suite.Assert().Nil(m)
}

func (suite *metaTestSuite) TestGetMeta() {
	uk := vault.NewUniqueKey()
	m, err := suite.bs.NewMeta(context.TODO(), uk, vault.Meta{
		UserID: 1,
		Extra:  "some extra",
	})
	suite.Assert().NoError(err)

	got, err := suite.bs.GetMeta(context.TODO(), uk)
	suite.Assert().NoError(err)
	suite.Assert().Equal(m, got)
}

func (suite *metaTestSuite) TestListMetaByUser() {
	uks := make(vault.List, 0)
	expteced := []vault.Meta{
		{
			UserID: 1,
			Extra:  "some extra",
		},
		{
			UserID: 1,
			Extra:  "some extra #2",
		},
	}
	for _, m := range expteced {
		uk := vault.NewUniqueKey()
		nm, err := suite.bs.NewMeta(context.TODO(), uk, m)
		uks[uk] = *nm
		suite.Assert().NoError(err)

	}
	_, err := suite.bs.NewMeta(context.TODO(), vault.NewUniqueKey(), vault.Meta{
		UserID: 2,
		Extra:  "",
	})
	suite.Assert().NoError(err)
	list, err := suite.bs.ListMetaByUser(context.TODO(), 1)
	suite.Assert().NoError(err)
	suite.Assert().Equal(uks, list)
}
