package sync

import (
	"context"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	GetSecretData(ctx context.Context, id vault.MetaID) (*vault.DataReader, error)
	GetSecretMeta(ctx context.Context, id vault.MetaID) (*vault.Meta, error)
	GetSecretMetaByAlias(ctx context.Context, alias string) (*vault.Meta, error)
	ListSecretsByUser(ctx context.Context) (vault.List, error)
	PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
}

type client interface {
	ListSecrets(ctx context.Context) (vault.List, error)
	GetSecretMeta(ctx context.Context, id vault.MetaID) (*vault.Meta, error)
	GetSecretData(ctx context.Context, id vault.MetaID, w io.Writer) error
	PutSecret(ctx context.Context, meta vault.Meta, r io.Reader) (*vault.Meta, error)
}

// TODO: unused
// type logger interface {
// 	Errorf(template string, args ...interface{})
// 	Debugf(template string, args ...interface{})
// }
