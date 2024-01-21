package bolt

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/k1nky/gophkeeper/internal/adapter/store"
	bolt "go.etcd.io/bbolt"
)

type BoltStorage struct {
	dsn string
	*bolt.DB
}

var _ store.MetaStore = new(BoltStorage)

func New(dsn string) *BoltStorage {
	return &BoltStorage{
		dsn: dsn,
	}
}

func (bs *BoltStorage) Open(ctx context.Context) (err error) {
	if bs.DB, err = bolt.Open(bs.dsn, 0600, &bolt.Options{}); err != nil {
		return
	}
	err = bs.DB.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{"users", "user_meta_list", "meta"} {
			if _, err := tx.CreateBucketIfNotExists(tb(bucket)); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (bs *BoltStorage) Close() error {
	return bs.DB.Close()
}

func serialize(a any) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(a)
	return b.Bytes(), err
}

func deserialize(b []byte, a any) error {
	buf := bytes.NewBuffer(b)
	return gob.NewDecoder(buf).Decode(a)
}

func tb(s string) []byte {
	return []byte(s)
}
