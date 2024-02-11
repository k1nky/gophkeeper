package http

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
)

//go:generate mockgen -source=contract.go -destination=mock/auth.go -package=mock authService
type authService interface {
	Register(ctx context.Context, u user.User) (string, error)
	Login(ctx context.Context, u user.User) (string, error)
	Authorize(token string) (user.PrivateClaims, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}
