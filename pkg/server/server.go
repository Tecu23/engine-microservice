// Package server should handle all server logic
package server

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Tecu23/engine-microservice/pkg/enginepb"
	"github.com/Tecu23/engine-microservice/pkg/pool"
)

// ChessEngineServer is a server type
type ChessEngineServer struct {
	enginepb.UnimplementedChessEngineServer
	workerPool *pool.WorkerPool
}

// CalculateBestMove is the actual implementation of the gRPC method
func (srv *ChessEngineServer) CalculateBestMove(
	ctx context.Context,
	req *enginepb.MoveRequest,
) (*enginepb.MoveResponse, error) {
	srv.workerPool.SubmitJob(pool.MoveRequest{Fen: req.Fen, Type: req.Type, ID: req.Id})

	bestMove := srv.workerPool.GetResult(req.Id)

	return &enginepb.MoveResponse{
		BestMove: bestMove,
	}, nil
}

// RegisterServer registers the gRPC server
func RegisterServer(grpcServer *grpc.Server, pool *pool.WorkerPool) {
	enginepb.RegisterChessEngineServer(grpcServer, &ChessEngineServer{workerPool: pool})
}
