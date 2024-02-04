package main

import (
	"context"
	"errors"
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

const (
	MaxCloseTimeout = 5 * time.Second
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

func newGRPCServer(auth *auth.Service, l *logger.Logger) *grpc.Server {
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpchandler.LoggerUnaryInterceptor(l),
		grpchandler.AuthorizationUnaryInterceptor(auth),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpchandler.LoggerStreamInterceptor(l),
		grpchandler.AuthorizationStreamInterceptor(auth),
	}
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptors...), grpc.ChainStreamInterceptor(streamInterceptors...))
	return srv
}

func main() {
	l := logger.New()
	l.SetLevel("debug")
	store := store.New(bolt.New("/tmp/server-meta.db"), filestore.New("/tmp/server-vault"))
	auth := auth.New("secret", time.Hour*175200, store, l)
	keeper := keeper.New(store, l)
	hh := httphandler.New(auth, keeper, l)
	gh := grpchandler.New(auth, keeper, l)
	grpcServer := newGRPCServer(auth, l)
	pb.RegisterKeeperServer(grpcServer, gh)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := store.Open(ctx)
	if err != nil {
		l.Errorf("could not open storage: %s", err)
		os.Exit(1)
	}
	defer store.Close()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(newMux(grpcServer, hh), &http2.Server{}),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				l.Errorf("unexpected server closing: %v", err)
			}
		}
	}()
	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), MaxCloseTimeout)
		defer cancel()
		srv.Shutdown(c)
	}()
	<-ctx.Done()
}
