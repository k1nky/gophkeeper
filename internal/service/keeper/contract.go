package keeper

import (
	"context"

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

type logger interface {
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}