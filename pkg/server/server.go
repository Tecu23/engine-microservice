// Package server should handle all server logic
package server

import (
	"context"
	"time"

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

	engineCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	eng, err := pool.GetEngine(engineCtx)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, status.Errorf(
				codes.ResourceExhausted,
				"no available engine instances: %v",
				err,
			)
		}
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
func RegisterServer(grpcServer *grpc.Server, cfg *config.Config) (*ChessEngineServer, error) {
	// Create the server with engine pools
	chessServer, err := NewChessEngineServer(&cfg.Engine)
	if err != nil {
		return nil, err
	}

	generated.RegisterChessEngineServer(grpcServer, chessServer)

	return chessServer, nil
}

func NewChessEngineServer(cfg *config.EngineConfig) (*ChessEngineServer, error) {
	enginePools := make(map[string]*engine.EnginePool)
	for engineType, path := range cfg.Paths {
		poolConfig := &engine.EngineConfig{
			EngineType: engineType,
			Path:       path,
			PoolSize:   cfg.PoolSize,
		}

		pool, err := engine.NewEnginePool(poolConfig)
		if err != nil {
			return nil, err
		}

		enginePools[engineType] = pool
	}

	return &ChessEngineServer{enginePools: enginePools}, nil
}

func Shutdown(ctx context.Context, grpcServer *grpc.Server, server *ChessEngineServer) error {
	grpcServer.GracefulStop()

	for _, pool := range server.enginePools {
		pool.Close()
	}

	return nil
}
