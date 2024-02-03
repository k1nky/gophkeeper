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
	uk := vault.NewMetaID()
	m := vault.Meta{
		UserID: 1,
		Extra:  "some extra",
		ID:     uk,
	}
	got, err := suite.bs.NewMeta(context.TODO(), m)
	suite.Assert().NoError(err)
	suite.Assert().Equal(m, *got)
}

func (suite *metaTestSuite) TestGetMetaByIDNotExists() {
	m, err := suite.bs.GetMetaByID(context.TODO(), vault.MetaID("not_exist"), 0)
	suite.Assert().NoError(err)
	suite.Assert().Nil(m)
}

func (suite *metaTestSuite) TestGetMetaByAliasNotExists() {
	m, err := suite.bs.GetMetaByAlias(context.TODO(), "", 0)
	suite.Assert().NoError(err)
	suite.Assert().Nil(m)

	m, err = suite.bs.GetMetaByAlias(context.TODO(), "not_exist", 0)
	suite.Assert().NoError(err)
	suite.Assert().Nil(m)

}

func (suite *metaTestSuite) TestGetMetaByID() {
	id := vault.NewMetaID()
	m, err := suite.bs.NewMeta(context.TODO(), vault.Meta{
		UserID: 1,
		Extra:  "some extra",
		ID:     id,
	})
	suite.Assert().NoError(err)
	suite.bs.NewMeta(context.TODO(), vault.Meta{
		UserID: 1,
		Extra:  "some extra#2",
		ID:     vault.NewMetaID(),
	})

	got, err := suite.bs.GetMetaByID(context.TODO(), id, 1)
	suite.Assert().NoError(err)
	suite.Assert().Equal(m, got)
}

func (suite *metaTestSuite) TestGetMetaByAlias() {
	id := vault.NewMetaID()
	m, err := suite.bs.NewMeta(context.TODO(), vault.Meta{
		UserID: 1,
		Alias:  "alias",
		Extra:  "some extra",
		ID:     id,
	})
	suite.Assert().NoError(err)
	suite.bs.NewMeta(context.TODO(), vault.Meta{
		UserID: 1,
		Alias:  "alias#2",
		Extra:  "some extra#2",
		ID:     vault.NewMetaID(),
	})

	got, err := suite.bs.GetMetaByAlias(context.TODO(), "alias", 1)
	suite.Assert().NoError(err)
	suite.Assert().Equal(m, got)
}

func (suite *metaTestSuite) TestListMetaByUser() {
	expteced := vault.List{
		{
			UserID: 1,
			Extra:  "some extra",
			ID:     vault.NewMetaID(),
		},
		{
			UserID: 1,
			Extra:  "some extra #2",
			ID:     vault.NewMetaID(),
		},
	}
	for _, m := range expteced {
		_, err := suite.bs.NewMeta(context.TODO(), m)
		suite.Assert().NoError(err)

	}
	_, err := suite.bs.NewMeta(context.TODO(), vault.Meta{
		ID:     vault.NewMetaID(),
		UserID: 2,
		Extra:  "",
	})
	suite.Assert().NoError(err)
	list, err := suite.bs.ListMetaByUser(context.TODO(), 1)
	suite.Assert().NoError(err)
	suite.Assert().ElementsMatch(expteced, list)
}
