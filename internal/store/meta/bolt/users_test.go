package bolt

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
)

type usersTestSuite struct {
	suite.Suite
	bs *BoltStorage
}

func (suite *usersTestSuite) SetupTest() {
	var err error
	rootDir := suite.T().TempDir()
	if suite.bs, err = openTestDB(rootDir); err != nil {
		suite.FailNow(err.Error())
		return
	}
	suite.bs.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(tb("users"))
		suite.Assert().NotNil(b)
		return nil
	})
}

func (suite *usersTestSuite) TestNewUser() {
	u := user.User{
		Login:    "u",
		Password: "p",
	}
	newUser, err := suite.bs.NewUser(context.TODO(), u)
	suite.NoError(err)
	suite.Equal(u.Login, newUser.Login)
	suite.Equal(u.Password, newUser.Password)
	suite.NotEqual(0, newUser.ID)
}

func (suite *usersTestSuite) TestNewUserDuplicate() {
	u := user.User{
		ID:       1,
		Login:    "u",
		Password: "p",
	}
	_, err := suite.bs.NewUser(context.TODO(), u)
	suite.Assert().NoError(err)
	got, err := suite.bs.NewUser(context.TODO(), u)
	suite.ErrorIs(err, user.ErrDuplicateLogin, "")
	suite.Nil(got, "")
}

func (suite *usersTestSuite) TestGetUserByLogin() {
	u := user.User{
		ID:       1,
		Login:    "u1",
		Password: "p1",
	}
	_, err := suite.bs.NewUser(context.TODO(), u)
	suite.Assert().NoError(err)
	got, err := suite.bs.GetUserByLogin(context.TODO(), "u1")
	suite.NoError(err)
	suite.Equal(u, *got)
}

func (suite *usersTestSuite) TestGetUserByLoginNotExists() {
	got, err := suite.bs.GetUserByLogin(context.TODO(), "u1000")
	suite.NoError(err)
	suite.Nil(got)
}
