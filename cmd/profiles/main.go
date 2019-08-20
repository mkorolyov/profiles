package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	profile "humans.net/ms/profiles"
	"log"
	"net"
	"net/http"
)

func main() {
	// Create the main listener.
	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	m := cmux.New(l)

	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	// Create your protocol servers.
	grpcS := grpc.NewServer()
	profile.RegisterProfileServer(grpcS, profile.NewService())

	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()
	if err := profile.RegisterProfileHandlerFromEndpoint(context.Background(), mux, "0.0.0.0:9090",
		opts); err != nil {
		panic(fmt.Sprint("failed to register http handler for service :%v", err))
	}
	httpS := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := grpcS.Serve(grpcL); err != nil {
			log.Fatalf("grpc server error %v", err)
		}
	}()

	go func() {
		if err := httpS.Serve(httpL); err != nil {
			log.Fatalf("https server error %v", err)
		}
	}()

	// Start serving!
	if err := m.Serve(); err != nil {
		log.Fatalf("cmux server error: %v", err)
	}
}
