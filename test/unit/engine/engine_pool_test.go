package engine_test

import (
	"sync"
	"testing"

	"github.com/Tecu23/engine-microservice/pkg/engine"
)

func TestEnginePool(t *testing.T) {
	cfg := &engine.EngineConfig{
		EngineType: "stockfish",
		Path:       "/path/to/stockfish",
		PoolSize:   2,
	}

	pool, err := engine.NewEnginePool(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine pool: %v", err)
	}
	defer pool.Close()

	var wg sync.WaitGroup
	numRequests := 5
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()
			eng, err := pool.GetEngine()
			if err != nil {
				t.Errorf("Request %d: failed to get engine: %v", i, err)
				return
			}
			defer pool.ReturnEngine(eng)

			fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			depth := 5

			bestMove, err := eng.CalculateBestMove(fen, depth)
			if err != nil {
				t.Errorf("Request %d: CalculateBestMove failed: %v", i, err)
				return
			}

			if bestMove == "" {
				t.Errorf("Request %d: expected a best move, got empty string", i)
			} else {
				t.Logf("Request %d: Best Move: %s", i, bestMove)
			}
		}(i)
	}

	wg.Wait()
}
