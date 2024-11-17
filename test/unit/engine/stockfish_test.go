package engine_test

import (
	"testing"

	"github.com/Tecu23/engine-microservice/pkg/engine"
)

func TestStockfishEngine(t *testing.T) {
	eng, err := engine.NewStockfishEngine("/path/to/stockfish")
	if err != nil {
		t.Fatalf("Failed to create Stockfish engine: %v", err)
	}
	defer eng.Close()

	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	depth := 10

	bestMove, err := eng.CalculateBestMove(fen, depth)
	if err != nil {
		t.Fatalf("CalculateBestMove failed: %v", err)
	}

	if bestMove == "" {
		t.Fatalf("Expected a best move, got empty string")
	}

	t.Logf("Best Move: %s", bestMove)
}
