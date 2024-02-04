// Пакет bolt предоставляет хранилище мета-данных секретов в boltdb.
package bolt

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/k1nky/gophkeeper/internal/adapter/store"
	bolt "go.etcd.io/bbolt"
)

// Хранилище мета-данных секретов в boltdb. Хранимые данные серилизуются средстами пакета gob.
type BoltStorage struct {
	dsn string
	*bolt.DB
}

var _ store.MetaStore = new(BoltStorage)

// New возвращает новое хранилище мета-данных секретов в boltdb.
func New(dsn string) *BoltStorage {
	return &BoltStorage{
		dsn: dsn,
	}
}

// Open открывает хранилище.
func (bs *BoltStorage) Open(ctx context.Context) (err error) {
	if bs.DB, err = bolt.Open(bs.dsn, 0600, &bolt.Options{}); err != nil {
		return
	}
	err = bs.DB.Update(func(tx *bolt.Tx) error {
		// создаем обязательные бакеты
		for _, bucket := range []string{"users", "meta"} {
			if _, err := tx.CreateBucketIfNotExists(tb(bucket)); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

// Close закрывает хранилище.
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
