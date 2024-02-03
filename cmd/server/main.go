package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	grpchandler "github.com/k1nky/gophkeeper/internal/adapter/grpc"
	httphandler "github.com/k1nky/gophkeeper/internal/adapter/http"
	"github.com/k1nky/gophkeeper/internal/adapter/store"
	"github.com/k1nky/gophkeeper/internal/logger"
	pb "github.com/k1nky/gophkeeper/internal/protocol/proto"
	"github.com/k1nky/gophkeeper/internal/service/auth"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
	"github.com/k1nky/gophkeeper/internal/store/meta/bolt"
	"github.com/k1nky/gophkeeper/internal/store/objects/filestore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func newMux(grpcServer *grpc.Server, httpServer http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpServer.ServeHTTP(w, r)
		}
	})
}

func newGRPCServer(auth *auth.Service) *grpc.Server {
	// unaryInterceptors := []grpc.UnaryServerInterceptor{grpcmw.LoggerUnaryInterceptor(l)}
	// streamInterceptors := []grpc.StreamServerInterceptor{grpcmw.LoggerStreamInterceptor(l)}
	unaryInterceptors := []grpc.UnaryServerInterceptor{grpchandler.AuthorizationUnaryInterceptor(auth)}
	streamInterceptors := []grpc.StreamServerInterceptor{grpchandler.AuthorizationStreamInterceptor(auth)}
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptors...), grpc.ChainStreamInterceptor(streamInterceptors...))
	return srv
}

func main() {
	store := store.New(bolt.New("/tmp/server-meta.db"), filestore.New("/tmp/server-vault"))
	auth := auth.New("secret", time.Hour*175200, store, &logger.Blackhole{})
	keeper := keeper.New(store, &logger.Blackhole{})
	hh := httphandler.New(auth, keeper, &logger.Blackhole{})
	gh := grpchandler.New(auth, keeper)
	grpcServer := newGRPCServer(auth)
	pb.RegisterKeeperServer(grpcServer, gh)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	store.Open(ctx)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(newMux(grpcServer, hh), &http2.Server{}),
	}
	srv.ListenAndServe()
	<-ctx.Done()
}
