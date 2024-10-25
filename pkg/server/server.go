// Package server should handle all server logic
package server

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Tecu23/engine-microservice/pkg/enginepb"
	"github.com/Tecu23/engine-microservice/pkg/uci"
)

// ChessEngineServer is a server type
type ChessEngineServer struct {
	enginepb.UnimplementedChessEngineServer
	workerPool *uci.WorkerPool
}

// CalculateBestMove is the actual implementation of the gRPC method
func (srv *ChessEngineServer) CalculateBestMove(
	ctx context.Context,
	req *enginepb.MoveRequest,
) (*enginepb.MoveResponse, error) {

	srv.workerPool.SubmitJob(req.Fen)

	bestMove := srv.workerPool.GetResult()

	return &enginepb.MoveResponse{
		BestMove: bestMove,
	}, nil
}

// RegisterServer registers the gRPC server
func RegisterServer(grpcServer *grpc.Server, pool *uci.WorkerPool) {
	enginepb.RegisterChessEngineServer(grpcServer, &ChessEngineServer{workerPool: pool})
}
