// Package uci should handle the engine interface for communication
package uci

import (
	"fmt"
	"sync"
)

// WorkerPool represents a pool of engines to handle requests concurrently
type WorkerPool struct {
	engines []*EngineInterface
	jobs    chan string
	results chan string
	wg      sync.WaitGroup
}

// NewWorkerPool initializes a pool of engines
func NewWorkerPool(numWorkers int) (*WorkerPool, error) {

	engines := make([]*EngineInterface, 0, numWorkers)

	for i := 0; i < numWorkers; i++ {
		engine, err := NewEngineInterface()
		if err != nil {
			return nil, err
		}

		engines = append(engines, engine)
	}

	return &WorkerPool{
		engines: engines,
		jobs:    make(chan string),
		results: make(chan string),
	}, nil
}

// Start initializes the worker pool to start processing FEN requests
func (wp *WorkerPool) Start() {
	for _, engine := range wp.engines {
		wp.wg.Add(1)

		go func(e *EngineInterface) {
			defer wp.wg.Done()
			for fen := range wp.jobs {
				bestMove, err := e.GetBestMove(fen)
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
func (wp *WorkerPool) SubmitJob(fen string) {
	wp.jobs <- fen
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
