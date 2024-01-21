package filestore

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/stretchr/testify/suite"
)

type filestoreTestSuite struct {
	dir string
	suite.Suite
	fs *FileStore
}

func (suite *filestoreTestSuite) SetupTest() {
	suite.dir = suite.T().TempDir()
	suite.fs = New(suite.dir)
}

func (suite *filestoreTestSuite) TearDownTest() {}

func (suite *filestoreTestSuite) TestGetNotExisis() {
	buf := bytes.NewBuffer(nil)
	err := suite.fs.Get(context.TODO(), "not_exists", buf)
	suite.Assert().ErrorIs(err, vault.ErrObjectNotExists)
}

func (suite *filestoreTestSuite) TestPutSmallObject() {
	var expected = []byte("hello")

	obj := bytes.NewBuffer(expected)
	err := suite.fs.Put(context.TODO(), "obj1", obj)
	suite.Assert().NoError(err)

	p := suite.fs.path("obj1")
	data, err := os.ReadFile(p)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, data)
}

func (suite *filestoreTestSuite) TestGet() {
	var expected = []byte("hello")

	obj := bytes.NewBuffer(expected)
	err := suite.fs.Put(context.TODO(), "obj1", obj)
	suite.Assert().NoError(err)

	err = suite.fs.Get(context.TODO(), "obj1", obj)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, obj.Bytes())
}

func (suite *filestoreTestSuite) TestDelete() {
	var expected = []byte("hello")

	obj := bytes.NewBuffer(expected)
	err := suite.fs.Put(context.TODO(), "obj1", obj)
	suite.Assert().NoError(err)

	err = suite.fs.Delete(context.TODO(), "obj1")
	suite.Assert().NoError(err)
	p := suite.fs.path("obj1")
	_, err = os.ReadFile(p)
	suite.Assert().ErrorIs(err, os.ErrNotExist)
}

func TestFileStore(t *testing.T) {
	suite.Run(t, new(filestoreTestSuite))
}
