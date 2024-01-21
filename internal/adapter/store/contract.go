package store

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/store.go -package=mock ObjectStore
type ObjectStore interface {
	Open(ctx context.Context) error
	Get(ctx context.Context, key string, obj vault.Object) error
	Close() error
	Put(ctx context.Context, key string, obj vault.Object) error
	Delete(ctx context.Context, key string) error
}

//go:generate mockgen -source=contract.go -destination=mock/store.go -package=mock MetaStore
type MetaStore interface {
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	Close() error
	NewUser(ctx context.Context, u user.User) (*user.User, error)
	NewMeta(ctx context.Context, uk vault.UniqueKey, m vault.Meta) (*vault.Meta, error)
	GetMeta(ctx context.Context, uk vault.UniqueKey) (*vault.Meta, error)
	ListMetaByUser(ctx context.Context, id user.ID) (vault.List, error)
	Open(ctx context.Context) (err error)
}

type Store interface {
	Open(ctx context.Context) error
	NewUser(ctx context.Context, u user.User) (*user.User, error)
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	PutSecret(ctx context.Context, s vault.Secret) error
	GetSecret(ctx context.Context, uk vault.UniqueKey) (*vault.Secret, error)
	ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error)
	Close() error
}
