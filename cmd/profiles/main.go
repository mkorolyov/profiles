package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	profile "github.com/mkorolyov/profiles"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
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
	grpcL := m.Match(cmux.HTTP2())
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

//func main() {
//	// Create the main listener.
//	l, err := net.Listen("tcp", "0.0.0.0:9090")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	grpcS := grpc.NewServer()
//	profile.RegisterProfileServer(grpcS, profile.NewService())
//
//	if err := grpcS.Serve(l); err != nil {
//		log.Fatalf("grpc serve error %v", err)
//	}
//}
