package grpc

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey int

const (
	keyUserClaims contextKey = iota
)

var (
	// TODO: error message
	ErrUnauthenticated = status.Error(codes.Unauthenticated, "")
)

func authorize(ctx context.Context, auth authService) (*user.PrivateClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	token := md.Get("Authorization")
	if len(token) == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	claims, err := auth.Authorize(token[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	return &claims, nil
}

type wrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrapper) Context() context.Context {
	return w.ctx
}

func (w *wrapper) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func AuthorizationUnaryInterceptor(auth authService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		claims, err := authorize(ctx, auth)
		if err != nil {
			return nil, err
		}
		newCtx := user.NewContextWithClaims(ctx, *claims)

		return handler(newCtx, req)
	}
}

func AuthorizationStreamInterceptor(auth authService) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		claims, err := authorize(ss.Context(), auth)
		if err != nil {
			return err
		}
		ctx := user.NewContextWithClaims(ss.Context(), *claims)
		w := &wrapper{
			ctx:          ctx,
			ServerStream: ss,
		}
		return handler(srv, w)
	}
}
