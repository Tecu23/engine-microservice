// Package uci should handle the engine interface for communication
package uci

import (
	"fmt"
	"sync"
)

// MoveRequest is the type of request coming from the server
type MoveRequest struct {
	Fen  string
	Type string
}

// WorkerPool represents a pool of engines to handle requests concurrently
type WorkerPool struct {
	engines []EngineInterface
	jobs    chan MoveRequest
	results chan string
	wg      sync.WaitGroup
}

// NewWorkerPool initializes a pool of engines
func NewWorkerPool(numArgo int, numStockfish int) (*WorkerPool, error) {

	engines := make([]EngineInterface, 0, numStockfish+numArgo)
	for i := 0; i < numStockfish; i++ {
		engine, err := NewStockfishEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to create argo server: %v", err)
		}
		engines = append(engines, engine)
	}
	for i := 0; i < numArgo; i++ {
		engine, err := NewArgoEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to create argo server: %v", err)
		}
		engines = append(engines, engine)
	}

	return &WorkerPool{
		engines: engines,
		jobs:    make(chan MoveRequest),
		results: make(chan string),
	}, nil
}

// Start initializes the worker pool to start processing FEN requests
func (wp *WorkerPool) Start() {
	for i := 0; i < len(wp.engines); i++ {
		engine := wp.engines[i]
		wp.wg.Add(1)
		go func(e EngineInterface) {
			defer wp.wg.Done()
			for job := range wp.jobs {
				bestMove, err := e.GetBestMove(job.Fen)
				if err != nil {
					wp.results <- fmt.Sprintf("error: %v", err)
				} else {
					wp.results <- bestMove
				}
			}
		}(engine)
	}
}

// SubmitJob submits a FEN to the worker pool for processing
func (wp *WorkerPool) SubmitJob(move MoveRequest) {
	wp.jobs <- move
}

// GetResult retrieves the result of a processed FEN
func (wp *WorkerPool) GetResult() string {
	return <-wp.results
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}
