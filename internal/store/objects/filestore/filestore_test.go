package filestore

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type filestoreTestSuite struct {
	dir string
	suite.Suite
	fs *FileStore
}

func isReaderEqualBytes(t *testing.T, expected []byte, reader *vault.DataReader) bool {
	b1 := bytes.NewBuffer(nil)
	if _, err := b1.ReadFrom(reader); !assert.NoError(t, err) {
		return false
	}
	return assert.Equal(t, expected, b1.Bytes())
}

func (suite *filestoreTestSuite) SetupTest() {
	suite.dir = suite.T().TempDir()
	suite.fs = New(suite.dir)
}

func (suite *filestoreTestSuite) TearDownTest() {}

func (suite *filestoreTestSuite) TestGetNotExisis() {
	data, err := suite.fs.Get(context.TODO(), "not_exists")
	suite.Assert().Nil(data)
	suite.Assert().ErrorIs(err, vault.ErrObjectNotExists)
}

func (suite *filestoreTestSuite) TestPutSmallObject() {
	var expected = []byte("hello")

	data := vault.NewDataReader(vault.NewBytesBuffer(expected))
	err := suite.fs.Put(context.TODO(), "obj1", data)
	suite.Assert().NoError(err)

	p := suite.fs.path("obj1")
	fact, err := os.ReadFile(p)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, fact)
}

func (suite *filestoreTestSuite) TestMultipleGet() {
	var expected = []byte("hello")
	ctx := context.TODO()

	data := vault.NewDataReader(vault.NewBytesBuffer(expected))
	err := suite.fs.Put(context.TODO(), "obj1", data)
	suite.Assert().NoError(err)

	r1, err := suite.fs.Get(ctx, "obj1")
	suite.Assert().NoError(err)
	defer r1.Close()
	r2, err := suite.fs.Get(ctx, "obj1")
	suite.Assert().NoError(err)
	defer r2.Close()
	isReaderEqualBytes(suite.T(), expected, r1)
	isReaderEqualBytes(suite.T(), expected, r2)
}

func (suite *filestoreTestSuite) TestMultiplePut() {
	var expected1 = []byte("hello")
	var expected2 = []byte("hello2")
	ctx := context.TODO()

	err := suite.fs.Put(context.TODO(), "obj1", vault.NewDataReader(vault.NewBytesBuffer(expected1)))
	suite.Assert().NoError(err)
	r1, err := suite.fs.Get(ctx, "obj1")
	suite.Assert().NoError(err)
	defer r1.Close()

	err = suite.fs.Put(context.TODO(), "obj1", vault.NewDataReader(vault.NewBytesBuffer(expected2)))
	suite.Assert().NoError(err)
	r2, err := suite.fs.Get(ctx, "obj1")
	suite.Assert().NoError(err)
	defer r2.Close()

	isReaderEqualBytes(suite.T(), expected2, r1)
	isReaderEqualBytes(suite.T(), expected2, r2)
}

func (suite *filestoreTestSuite) TestGetClose() {
	var expected = []byte("hello")
	ctx := context.TODO()

	data := vault.NewDataReader(vault.NewBytesBuffer(expected))
	err := suite.fs.Put(context.TODO(), "obj1", data)
	suite.Assert().NoError(err)

	r1, err := suite.fs.Get(ctx, "obj1")
	suite.Assert().NoError(err)
	r1.Close()
	r2, err := suite.fs.Get(ctx, "obj1")
	suite.Assert().NoError(err)
	defer r2.Close()
	isReaderEqualBytes(suite.T(), expected, r2)
}

func (suite *filestoreTestSuite) TestGet() {
	var expected = []byte("hello")

	data := vault.NewDataReader(vault.NewBytesBuffer(expected))
	err := suite.fs.Put(context.TODO(), "obj1", data)
	suite.Assert().NoError(err)

	got, err := suite.fs.Get(context.TODO(), "obj1")
	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(got)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, buf.Bytes())
}

func (suite *filestoreTestSuite) TestDelete() {
	var expected = []byte("hello")

	data := vault.NewDataReader(vault.NewBytesBuffer(expected))
	err := suite.fs.Put(context.TODO(), "obj1", data)
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
