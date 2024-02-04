package store

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"golang.org/x/sync/errgroup"
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
	g := new(errgroup.Group)
	g.Go(func() error {
		return a.mstore.Open(ctx)
	})
	g.Go(func() error {
		return a.ostore.Open(ctx)
	})
	err := g.Wait()
	if err != nil {
		a.Close()
	}
	return err
}

func (a *Adapter) Close() error {
	g := new(errgroup.Group)
	if a.mstore != nil {
		g.Go(a.mstore.Close)
	}
	if a.ostore != nil {
		g.Go(a.ostore.Close)
	}
	return g.Wait()
}

func (a *Adapter) NewUser(ctx context.Context, u user.User) (*user.User, error) {
	return a.mstore.NewUser(ctx, u)
}

func (a *Adapter) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	return a.mstore.GetUserByLogin(ctx, login)
}

func (a *Adapter) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	if len(meta.ID) == 0 {
		meta.ID = vault.NewMetaID()
	}
	key := fmt.Sprintf("%d-%s", meta.UserID, meta.ID)
	if err := a.ostore.Put(ctx, key, data); err != nil {
		return nil, err
	}
	if _, err := a.mstore.NewMeta(ctx, meta); err != nil {
		a.ostore.Delete(ctx, key)
		return nil, err
	}
	return &meta, nil
}

func (a *Adapter) GetSecretData(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.DataReader, error) {
	key := fmt.Sprintf("%d-%s", userID, metaID)
	return a.ostore.Get(ctx, key)
}

func (a *Adapter) GetSecretMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error) {
	m, err := a.mstore.GetMetaByID(ctx, metaID, userID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return m, err
}

func (a *Adapter) GetSecretMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error) {
	m, err := a.mstore.GetMetaByAlias(ctx, alias, userID)
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
