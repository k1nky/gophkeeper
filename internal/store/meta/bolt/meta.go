package bolt

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	bolt "go.etcd.io/bbolt"
)

func (bs *BoltStorage) NewMeta(ctx context.Context, m vault.Meta) (*vault.Meta, error) {
	err := bs.DB.Update(func(tx *bolt.Tx) error {

		mb := tx.Bucket(tb("meta"))
		umb, err := mb.CreateBucketIfNotExists(tb(fmt.Sprintf("%d", m.UserID)))
		if err != nil {
			return err
		}
		value, err := serialize(m)
		if err != nil {
			return err
		}

		return umb.Put(tb(string(m.ID)), value)
	})
	if err == nil {
		return &m, nil
	}
	return nil, err
}

func (bs *BoltStorage) GetMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error) {
	m := &vault.Meta{}
	err := bs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tb("meta"))
		umb := b.Bucket(tb(fmt.Sprintf("%d", userID)))
		if umb == nil {
			m = nil
			return nil
		}
		value := umb.Get([]byte(metaID))
		if value == nil {
			m = nil
			return nil
		}
		if err := deserialize(value, m); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, err
}

func (bs *BoltStorage) GetMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error) {
	if len(alias) == 0 {
		return nil, nil
	}
	m := &vault.Meta{}
	found := false
	err := bs.View(func(tx *bolt.Tx) error {
		mb := tx.Bucket(tb("meta"))
		umb := mb.Bucket(tb(fmt.Sprintf("%d", userID)))
		if umb == nil {
			m = nil
			return nil
		}

		c := umb.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := deserialize(v, m); err != nil {
				return err
			}
			if m.Alias == alias {
				found = true
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return m, err
}

func (bs *BoltStorage) ListMetaByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	list := vault.List{}
	err := bs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tb("meta"))
		umb := b.Bucket(tb(fmt.Sprintf("%d", userID)))
		if umb == nil {
			return nil
		}

		c := umb.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			m := vault.Meta{}
			if err := deserialize(v, &m); err != nil {
				return err
			}
			list = append(list, m)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
