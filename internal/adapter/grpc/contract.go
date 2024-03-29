package grpc

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/service.go -package=mock authService
type authService interface {
	Authorize(token string) (user.PrivateClaims, error)
}

//go:generate mockgen -source=contract.go -destination=mock/service.go -package=mock keeperService
type keeperService interface {
	GetSecretData(ctx context.Context, id vault.MetaID) (*vault.DataReader, error)
	GetSecretMeta(ctx context.Context, id vault.MetaID) (*vault.Meta, error)
	GetSecretMetaByAlias(ctx context.Context, alias string) (*vault.Meta, error)
	ListSecretsByUser(ctx context.Context) (vault.List, error)
	PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}
