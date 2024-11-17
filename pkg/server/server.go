// Package server should handle all server logic
package server

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Tecu23/engine-microservice/pkg/api/generated"
	"github.com/Tecu23/engine-microservice/pkg/config"
	"github.com/Tecu23/engine-microservice/pkg/engine"
)

// ChessEngineServer is a server type
type ChessEngineServer struct {
	generated.UnimplementedChessEngineServer
	enginePools map[string]*engine.EnginePool
}

// CalculateBestMove is the actual implementation of the gRPC method
func (s *ChessEngineServer) CalculateBestMove(
	ctx context.Context,
	req *generated.MoveRequest,
) (*generated.MoveResponse, error) {
	pool, exists := s.enginePools[req.EngineType]
	if !exists {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"unsupported engine type: %s",
			req.EngineType,
		)
	}

	eng, err := pool.GetEngine()
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "no available engine instances: %v", err)
	}
	defer pool.ReturnEngine(eng)

	bestMove, err := eng.CalculateBestMove(req.Fen, int(req.Depth))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "engine error: %v", err)
	}

	return &generated.MoveResponse{
		BestMove:   bestMove,
		EngineInfo: eng.Info(),
	}, nil
}

// RegisterServer registers the gRPC server
func RegisterServer(grpcServer *grpc.Server, cfg *config.Config) {
	// Create engine configurations
	engineConfigs := []*engine.EngineConfig{
		{
			EngineType: "stockfish",
			Path:       cfg.EnginePathStockfish,
			PoolSize:   cfg.EnginePoolSize,
		},
		// Add configurations for other engines if needed
	}

	// Create the server with engine pools
	chessServer, err := NewChessEngineServer(engineConfigs)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	generated.RegisterChessEngineServer(grpcServer, chessServer)
}

func NewChessEngineServer(engineConfigs []*engine.EngineConfig) (*ChessEngineServer, error) {
	enginePools := make(map[string]*engine.EnginePool)

	for _, cfg := range engineConfigs {
		pool, err := engine.NewEnginePool(cfg)
		if err != nil {
			return nil, err
		}

		enginePools[cfg.EngineType] = pool
	}

	return &ChessEngineServer{enginePools: enginePools}, nil
}
