package bolt

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	bolt "go.etcd.io/bbolt"
)

// GetUserByLogin возвращает пользователя с именем login.
func (bs *BoltStorage) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	u := &user.User{}
	err := bs.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tb("users"))
		v := b.Get(tb(login))
		if v == nil {
			u = nil
			return nil
		}
		return deserialize(v, u)
	})
	return u, err
}

// NewUser создает нового пользователя u и возвращает созданный элемент.
func (bs *BoltStorage) NewUser(ctx context.Context, u user.User) (*user.User, error) {
	err := bs.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tb("users"))
		if v := b.Get(tb(u.Login)); v != nil {
			return fmt.Errorf("%s %w", u.Login, user.ErrDuplicateLogin)
		}
		id, _ := b.NextSequence()
		u.ID = user.ID(id)
		d, err := serialize(u)
		if err != nil {
			return err
		}
		return b.Put(tb(u.Login), d)
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}
