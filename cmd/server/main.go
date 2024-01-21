package main

import (
	"net/http"
	"strings"

	grpchandler "github.com/k1nky/gophkeeper/internal/adapter/grpc"
	httphandler "github.com/k1nky/gophkeeper/internal/adapter/http"
	pb "github.com/k1nky/gophkeeper/internal/proto"
	"google.golang.org/grpc"
)

func main() {
	hh := &httphandler.Adapter{}
	gh := &grpchandler.Adapter{}
	grpcServer := grpc.NewServer()
	pb.RegisterKeeperServer(grpcServer, gh)
	srv := &http.Server{
		Addr: ":80",
		Handler: func(grpcServer *grpc.Server, httpServer http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(w, r)
				} else {
					httpServer.ServeHTTP(w, r)
				}
			})
		}(grpcServer, hh),
	}
	srv.ListenAndServe()
}
