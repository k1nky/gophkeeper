package http

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

//go:generate mockgen -source=contract.go -destination=mock/auth.go -package=mock authService
type authService interface {
	Register(ctx context.Context, u user.User) (string, error)
	Login(ctx context.Context, u user.User) (string, error)
	Authorize(token string) (user.PrivateClaims, error)
}

type keeperService interface {
	PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error)
	GetSecretData(ctx context.Context, uk vault.MetaID) (*vault.DataReader, error)
	GetSecretMeta(ctx context.Context, uk vault.MetaID) (*vault.Meta, error)
	ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}
