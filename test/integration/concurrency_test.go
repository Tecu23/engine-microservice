package integration_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/Tecu23/engine-microservice/pkg/api/generated"
)

func TestConcurrentRequests(t *testing.T) {
	// Setup the client connection
	conn, err := grpc.NewClient(
		"localhost:8089",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := generated.NewChessEngineClient(conn)

	// Define the number of concurrent requests
	numRequests := 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()
			maxRetries := 3
			attempt := 0

			for {
				attempt++
				ctx := context.Background()
				// Include authentication if required

				req := &generated.MoveRequest{
					Id:         "1l",
					EngineType: "stockfish",
					Fen:        "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2o",
					Depth:      10,
				}

				res, err := client.CalculateBestMove(ctx, req)
				if err != nil {
					statusCode, ok := status.FromError(err)
					if ok && statusCode.Code() == codes.ResourceExhausted && attempt <= maxRetries {
						// Wait before retrying
						time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
						continue
					}
					t.Errorf("Request %d failed after %d attempts: %v", i, attempt, err)
					return
				}

				if res.BestMove == "" {
					t.Errorf("Request %d returned empty best move", i)
				}
				break
			}
		}(i)
	}

	wg.Wait()
}
