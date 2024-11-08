// Package uci handles the engine interface for communication with chess engines.
// It provides a WorkerPool to manage multiple engines concurrently and allows specifying
// the engine type for each move calculation.
package uci

import (
	"fmt"
	"sync"
)

// MoveRequest represents a request to calculate the best move for a given FEN string.
// The Type field specifies the engine type to use (e.g., "stockfish", "argo").
type MoveRequest struct {
	Fen  string
	Type string
}

// WorkerPool manages a pool of chess engines to handle move requests concurrently.
// It allows clients to submit MoveRequests specifying the engine type, and routes
// these requests to the appropriate engines.
type WorkerPool struct {
	engines     map[string][]EngineInterface
	jobChannels map[string]chan MoveRequest
	jobs        chan MoveRequest
	results     chan string
	wg          sync.WaitGroup
}

// NewWorkerPool initializes a WorkerPool with the specified number of Argo and Stockfish engines.
// It returns an error if any of the engine instances fail to initialize.
func NewWorkerPool(numArgo int, numStockfish int) (*WorkerPool, error) {
	// create a engine map
	engines := make(map[string][]EngineInterface)
	jobChannels := make(map[string]chan MoveRequest)

	engines["stockfish"] = make([]EngineInterface, 0, numStockfish)
	jobChannels["stockfish"] = make(chan MoveRequest)
	for i := 0; i < numStockfish; i++ {
		engine, err := NewStockfishEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to create argo server: %v", err)
		}
		engines["stockfish"] = append(engines["stockfish"], engine)
	}

	engines["argo"] = make([]EngineInterface, 0, numArgo)
	jobChannels["argo"] = make(chan MoveRequest)
	for i := 0; i < numArgo; i++ {
		engine, err := NewArgoEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to create argo server: %v", err)
		}
		engines["argo"] = append(engines["argo"], engine)
	}

	return &WorkerPool{
		engines:     engines,
		jobChannels: jobChannels,
		jobs:        make(chan MoveRequest),
		results:     make(chan string),
	}, nil
}

// Start launches worker goroutines for each engine and begins processing requests.
// It should be called before submitting any jobs to the WorkerPool.
func (wp *WorkerPool) Start() {
	for engineType, engines := range wp.engines {
		jobChan := wp.jobChannels[engineType]

		for _, engine := range engines {
			wp.wg.Add(1)
			go func(e EngineInterface, jobChan <-chan MoveRequest) {
				defer wp.wg.Done()
				for job := range jobChan {
					bestMove, err := e.GetBestMove(job.Fen)
					if err != nil {
						wp.results <- fmt.Sprintf("error: %v", err)
					} else {
						wp.results <- bestMove
					}
				}
			}(engine, jobChan)
		}
	}
	go wp.dispatcher()
}

func (wp *WorkerPool) dispatcher() {
	for job := range wp.jobs {
		engineType := job.Type

		if engineType != "stockfish" && engineType != "argo" {
			engineType = "stockfish"
		}

		if ch, ok := wp.jobChannels[engineType]; ok {
			ch <- job
		} else {
			wp.results <- fmt.Sprintf("error: unknown engine type %s", job.Type)
		}
	}

	// Once the jobs channel is closed, close all engine-specific job channels to signal workers to stop
	for _, ch := range wp.jobChannels {
		close(ch)
	}
}

// SubmitJob submits a MoveRequest to the WorkerPool for processing.
// The request will be routed to the appropriate engine based on its Type.
func (wp *WorkerPool) SubmitJob(move MoveRequest) {
	wp.jobs <- move
}

// GetResult retrieves a result from the WorkerPool. It blocks until a result is available.
// The result is typically the best move calculated by an engine or an error message.
func (wp *WorkerPool) GetResult() string {
	return <-wp.results
}

// Stop gracefully shuts down the WorkerPool, ensuring all pending jobs are processed.
// It closes all channels and waits for all worker goroutines to finish.
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}
