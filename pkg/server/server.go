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
	engineInterface *uci.Interface
}

// CalculateBestMove is the actual implementation of the gRPC method
func (srv *ChessEngineServer) CalculateBestMove(
	ctx context.Context,
	req *enginepb.MoveRequest,
) (*enginepb.MoveResponse, error) {

	fen := req.Fen

	bestMove, err := srv.engineInterface.GetBestMove(fen)
	if err != nil {
		return nil, err
	}

	return &enginepb.MoveResponse{
		BestMove: bestMove,
	}, nil
}

// RegisterServer registers the gRPC server
func RegisterServer(grpcServer *grpc.Server, i *uci.Interface) {
	enginepb.RegisterChessEngineServer(grpcServer, &ChessEngineServer{engineInterface: i})
}
