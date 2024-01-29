package sync

import (
	"context"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
	GetSecretData(ctx context.Context, uk vault.UniqueKey) (*vault.DataReader, error)
	GetSecretMeta(ctx context.Context, uk vault.UniqueKey) (*vault.Meta, error)
	ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error)
}

type client interface {
	ListSecrets(ctx context.Context) (vault.List, error)
	GetSecretMeta(ctx context.Context, id vault.UniqueKey) (*vault.Meta, error)
	GetSecretData(ctx context.Context, id vault.UniqueKey, w io.Writer) error
	PutSecret(ctx context.Context, meta vault.Meta, r io.Reader) (*vault.Meta, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}
