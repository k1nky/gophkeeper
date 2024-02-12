package store

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/store.go -package=mock ObjectStore
type ObjectStore interface {
	Open(ctx context.Context) error
	Get(ctx context.Context, key string) (*vault.DataReader, error)
	Close() error
	Put(ctx context.Context, key string, obj *vault.DataReader) error
	Delete(ctx context.Context, key string) error
}

//go:generate mockgen -source=contract.go -destination=mock/store.go -package=mock MetaStore
type MetaStore interface {
	Close() error
	DeleteMeta(ctx context.Context, meta vault.Meta) error
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	NewUser(ctx context.Context, u user.User) (*user.User, error)
	NewMeta(ctx context.Context, m vault.Meta) (*vault.Meta, error)
	GetMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error)
	GetMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error)
	ListMetaByUser(ctx context.Context, userID user.ID) (vault.List, error)
	Open(ctx context.Context) (err error)
	UpdateMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error)
}

type Store interface {
	Open(ctx context.Context) error
	NewUser(ctx context.Context, u user.User) (*user.User, error)
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
	GetSecretData(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.DataReader, error)
	GetSecretMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error)
	GetSecretMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error)
	ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error)
	Close() error
	DeleteSecret(ctx context.Context, meta vault.Meta) error
	UpdateSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
	UpdateSecretMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error)
}
