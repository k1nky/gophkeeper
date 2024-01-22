package store

import (
	"bytes"
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

type Adapter struct {
	mstore MetaStore
	ostore ObjectStore
}

func New(mstore MetaStore, ostore ObjectStore) *Adapter {
	return &Adapter{
		mstore: mstore,
		ostore: ostore,
	}
}

func (a *Adapter) Open(ctx context.Context) error {
	// TODO: errgroups
	if err := a.mstore.Open(ctx); err != nil {
		return err
	}
	if err := a.ostore.Open(ctx); err != nil {
		a.mstore.Close()
		return err
	}
	return nil
}

func (a *Adapter) Close() error {
	// TODO: errgroups
	a.mstore.Close()
	a.ostore.Close()
	return nil
}

func (a *Adapter) NewUser(ctx context.Context, u user.User) (*user.User, error) {
	return a.mstore.NewUser(ctx, u)
}

func (a *Adapter) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	return a.mstore.GetUserByLogin(ctx, login)
}

func (a *Adapter) PutSecret(ctx context.Context, s vault.Secret) (*vault.Secret, error) {
	s.Key = vault.NewUniqueKey()
	if len(s.Key) == 0 {
		// TODO: wrap error
		return nil, fmt.Errorf("internal error")
	}
	if err := a.ostore.Put(ctx, string(s.Key), s.Data); err != nil {
		// TODO: wrap error
		return nil, err
	}
	if _, err := a.mstore.NewMeta(ctx, s.Key, s.Meta); err != nil {
		// TODO: wrap error
		a.ostore.Delete(ctx, string(s.Key))
		return nil, err
	}
	return &s, nil
}

func (a *Adapter) GetSecret(ctx context.Context, uk vault.UniqueKey) (*vault.Secret, error) {
	m, err := a.mstore.GetMeta(ctx, uk)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	buf := bytes.NewBuffer(nil)
	err = a.ostore.Get(ctx, string(uk), buf)
	if err != nil {
		return nil, err
	}
	s := &vault.Secret{
		Key:  uk,
		Meta: *m,
		Data: buf,
	}
	return s, nil
}

func (a *Adapter) ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	return a.mstore.ListMetaByUser(ctx, userID)
}
