package keeper

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
	GetSecretData(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.DataReader, error)
	GetSecretMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error)
	GetSecretMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error)
	ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}
