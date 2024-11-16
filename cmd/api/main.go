// main entry of the microservice
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Tecu23/engine-microservice/pkg/pool"
	"github.com/Tecu23/engine-microservice/pkg/server"
)

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	workerPool, err := pool.NewWorkerPool(4, 4)
	if err != nil {
		log.Fatalf("Failed to create the worker pool: %v", err)
	}
	workerPool.Start()

	server.RegisterServer(grpcServer, workerPool)

	log.Println("Chess Engine gRPC server is running on port 8089...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}