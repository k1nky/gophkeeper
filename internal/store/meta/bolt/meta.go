package bolt

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	bolt "go.etcd.io/bbolt"
)

func (bs *BoltStorage) IsExist(ctx context.Context, meta vault.Meta) bool {
	if m, err := bs.GetMetaByID(ctx, meta.ID, meta.UserID); err == nil && m != nil {
		return true
	}
	if m, err := bs.GetMetaByAlias(ctx, meta.Alias, meta.UserID); err == nil && m != nil {
		return true
	}
	return false
}

// NewMeta добавляет новую запись мета-данных секрета. Возвращает добавленный элемент.
func (bs *BoltStorage) NewMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
	if len(meta.ID) == 0 {
		return nil, vault.ErrEmptyMetaID
	}
	if bs.IsExist(ctx, meta) {
		return nil, vault.ErrDuplicate
	}
	err := bs.DB.Update(func(tx *bolt.Tx) error {

		mb := tx.Bucket(tb("meta"))
		// группируем мета-данные по пользователям
		umb, err := mb.CreateBucketIfNotExists(tb(fmt.Sprintf("%d", meta.UserID)))
		if err != nil {
			return err
		}
		value, err := serialize(meta)
		if err != nil {
			return err
		}

		return umb.Put(tb(string(meta.ID)), value)
	})
	if err == nil {
		return &meta, nil
	}
	return nil, err
}

// GetMetaByID возвращает мета-данные секрета пользователя userID по идентификатору metaID.
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

// GetMetaByID возвращает мета-данные секрета пользователя userID по псевдониму alias.
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

// ListMetaByUser возвращает список мета-данных секретов пользователя userID.
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
