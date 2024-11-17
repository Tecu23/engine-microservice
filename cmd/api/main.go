// main entry of the microservice
package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Tecu23/engine-microservice/pkg/auth"
	"github.com/Tecu23/engine-microservice/pkg/config"
	"github.com/Tecu23/engine-microservice/pkg/server"
)

func main() {
	cfg := config.LoadConfig()

	auth.Initialize(cfg.AuthTokens)

	creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
	if err != nil {
		log.Fatalf("failed to setup TLS: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)

	server.RegisterServer(grpcServer, cfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Chess Engine gRPC server is running on port %d...", cfg.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-quit
	// 	log.Println("Shutting down server...")
	// 	grpcServer.GracefulStop()
	// 	for _, pool := range engine.EnginePool() {
	// 		pool.Close()
	// 	}
	// 	log.Println("Server stopped.")
	// }()
	//
	// // Start serving
	// if err := grpcServer.Serve(lis); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }
}
