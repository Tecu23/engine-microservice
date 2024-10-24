// main entry of the microservice
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Tecu23/engine-microservice/package/server"
	"github.com/Tecu23/engine-microservice/package/uci"
)

func main() {
	lis, err := net.Listen("tcp", ":8089")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	engineInterface, err := uci.NewInterface()
	if err != nil {
		log.Fatalf("failed to initialize the chess engine interface: %v", err)
	}

	// register the chess engine server with the initialized engine
	server.RegisterServer(grpcServer, engineInterface)

	log.Println("Chess Engine gRPC server is running on port 8089...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
