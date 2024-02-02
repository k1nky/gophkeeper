package store

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

type Adapter struct {
	mstore MetaStore
	ostore ObjectStore
}

var _ Store = new(Adapter)

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

func (a *Adapter) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	if len(meta.ID) == 0 {
		// TODO: wrap error
		return nil, fmt.Errorf("internal error")
	}
	key := fmt.Sprintf("%d-%s", meta.UserID, meta.ID)
	if err := a.ostore.Put(ctx, key, data); err != nil {
		// TODO: wrap error
		return nil, err
	}
	if _, err := a.mstore.NewMeta(ctx, meta); err != nil {
		// TODO: wrap error
		a.ostore.Delete(ctx, key)
		return nil, err
	}
	return &meta, nil
}

func (a *Adapter) GetSecretData(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.DataReader, error) {
	key := fmt.Sprintf("%d-%s", userID, metaID)
	return a.ostore.Get(ctx, key)
}

func (a *Adapter) GetSecretMeta(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error) {
	m, err := a.mstore.GetMeta(ctx, metaID, userID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return m, err
}

func (a *Adapter) ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	return a.mstore.ListMetaByUser(ctx, userID)
}
