package grpc

import (
	"context"
	"time"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	// TODO: error message
	ErrUnauthenticated = status.Error(codes.Unauthenticated, "")
)

type requestLogger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
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

func LoggerStreamInterceptor(l requestLogger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		md, _ := metadata.FromIncomingContext(ss.Context())
		err := handler(srv, ss)
		if err != nil {
			l.Infof("%s %v %v", info.FullMethod, md, time.Since(start))
		} else {
			l.Errorf("%s %v %v", info.FullMethod, md, err)
		}
		return err
	}
}

func LoggerUnaryInterceptor(l requestLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		md, _ := metadata.FromIncomingContext(ctx)
		resp, err = handler(ctx, req)
		if err == nil {
			l.Infof("%s %v %v", info.FullMethod, md, time.Since(start))
		} else {
			l.Errorf("%s %v %v", info.FullMethod, md, err)
		}
		return resp, err
	}
}
