package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	accountpb "github.com/boseabhishek/foo-accounts/gen/go/accounts/proto/account"
	sv "github.com/boseabhishek/foo-accounts/pkg/service"
)

var (
	grpcServerEndpoint = "localhost:8080"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcServer := grpc.NewServer()
	accountServer := sv.NewAccountServer()
	accountpb.RegisterAccountServiceServer(grpcServer, accountServer)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := accountpb.RegisterAccountServiceHandlerFromEndpoint(ctx, gwmux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	fs := http.FileServer(http.Dir("../../swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", fs))
	mux.HandleFunc("/proto/account/account.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../gen/go/accounts/proto/account/account.swagger.json")
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})

	log.Println("Starting server on", grpcServerEndpoint)
	log.Println("Access Swagger UI at http://localhost:8080/swagger-ui/")
	if err := http.ListenAndServe(grpcServerEndpoint, h2c.NewHandler(handler, &http2.Server{})); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
