package bolt

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	bolt "go.etcd.io/bbolt"
)

// GetMetaByAlias возвращает мета-данные секрета пользователя userID по псевдониму alias.
// Поиск по псевдониму более затратный, чем по ИД, т.к. требует перебора мета-данных всех секретов указанного пользователя.
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

// IsExist возвращает true, если в хранилище уже существует запись о секрете для
// пользователя с таким ИД или псевдонимом.
func (bs *BoltStorage) IsExist(ctx context.Context, meta vault.Meta) bool {
	if m, err := bs.GetMetaByID(ctx, meta.ID, meta.UserID); err == nil && m != nil {
		return true
	}
	if m, err := bs.GetMetaByAlias(ctx, meta.Alias, meta.UserID); err == nil && m != nil {
		return true
	}
	return false
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

func (bs *BoltStorage) putMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
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

		return umb.Put([]byte(meta.ID), value)
	})
	if err == nil {
		return &meta, nil
	}
	return nil, err
}

// NewMeta добавляет новую запись мета-данных секрета. Возвращает добавленный элемент.
func (bs *BoltStorage) NewMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
	if len(meta.ID) == 0 {
		return nil, vault.ErrEmptyMetaID
	}
	if bs.IsExist(ctx, meta) {
		return nil, vault.ErrDuplicate
	}
	return bs.putMeta(ctx, meta)
}

// UpdateMeta обновляет мета-данные секрета. Возвращает обновленнный элемент.
func (bs *BoltStorage) UpdateMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
	if len(meta.ID) == 0 {
		return nil, vault.ErrEmptyMetaID
	}
	if !bs.IsExist(ctx, meta) {
		return nil, vault.ErrMetaNotExists
	}
	return bs.putMeta(ctx, meta)
}

// DeleteMeta удаляет мета-данные секрета.
func (bs *BoltStorage) DeleteMeta(ctx context.Context, meta vault.Meta) error {
	if len(meta.ID) == 0 {
		return nil
	}
	if !bs.IsExist(ctx, meta) {
		return nil
	}
	err := bs.DB.Update(func(tx *bolt.Tx) error {

		mb := tx.Bucket(tb("meta"))
		umb := mb.Bucket(tb(fmt.Sprintf("%d", meta.UserID)))
		if umb == nil {
			return nil
		}
		return umb.Delete([]byte(meta.ID))
	})
	return err
}
