package bolt

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

func openTestDB(root string) (*BoltStorage, error) {
	bs := New(path.Join(root, "bolt.db"))
	return bs, bs.Open(context.TODO())
}

func TestBoltStorage(t *testing.T) {
	suite.Run(t, new(usersTestSuite))
	suite.Run(t, new(metaTestSuite))
}
