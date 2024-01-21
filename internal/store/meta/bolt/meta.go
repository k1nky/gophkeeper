package bolt

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	bolt "go.etcd.io/bbolt"
)

func (bs *BoltStorage) NewMeta(ctx context.Context, uk vault.UniqueKey, m vault.Meta) (*vault.Meta, error) {
	err := bs.DB.Update(func(tx *bolt.Tx) error {
		// TODO: userid has existed
		// TODO: uk must be unique
		b := tx.Bucket(tb("user_meta_list"))
		uml := make([]vault.UniqueKey, 0)
		value := b.Get(tb(fmt.Sprintf("%d", m.UserID)))

		if value != nil {
			if err := deserialize(value, &uml); err != nil {
				return err
			}
		}
		uml = append(uml, uk)
		if value, err := serialize(uml); err != nil {
			return err
		} else {
			b.Put(tb(fmt.Sprintf("%d", m.UserID)), value)
		}
		b, err := tx.CreateBucketIfNotExists(tb("meta"))
		if err != nil {
			return err
		}
		if value, err = serialize(m); err != nil {
			return err
		}

		return b.Put(tb(string(uk)), value)
	})
	if err == nil {
		return &m, nil
	}
	return nil, err
}

func (bs *BoltStorage) GetMeta(ctx context.Context, uk vault.UniqueKey) (*vault.Meta, error) {
	m := &vault.Meta{}
	err := bs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tb("meta"))
		value := b.Get([]byte(uk))
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

func (bs *BoltStorage) ListMetaByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	list := vault.List{}
	err := bs.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tb("user_meta_list"))
		uml := make([]vault.UniqueKey, 0)
		value := b.Get(tb(fmt.Sprintf("%d", userID)))

		if value != nil {
			if err := deserialize(value, &uml); err != nil {
				return err
			}
		}
		b = tx.Bucket(tb("meta"))
		for _, uk := range uml {
			m := vault.Meta{}
			if err := deserialize(b.Get([]byte(uk)), &m); err != nil {
				return err
			}
			list[uk] = m
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
