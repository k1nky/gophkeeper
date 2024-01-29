package gophkeeper

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (a *Adapter) AuthorizationStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if len(a.token) != 0 {
			newCtx := metadata.AppendToOutgoingContext(ctx, "Authorization", a.token)
			return streamer(newCtx, desc, cc, method, opts...)
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (a *Adapter) AuthorizationUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if len(a.token) != 0 {
			newCtx := metadata.AppendToOutgoingContext(ctx, "Authorization", a.token)
			return invoker(newCtx, method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
