package pool

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Tecu23/engine-microservice/pkg/types/engineinterface"
)

// MockEngineInterface is a mock implementation of EngineInterface for testing.
type MockEngineInterface struct {
	bestMove string
	err      error
}

// GetBestMove mocks the behavior of getting the best move.
func (m *MockEngineInterface) GetBestMove(_ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.bestMove, nil
}

func TestWorkerPool_SubmitAndRetrieveResult(t *testing.T) {
	// Setup mock engines
	mockArgo := &MockEngineInterface{bestMove: "e2e4", err: nil}
	mockStockfish := &MockEngineInterface{bestMove: "d7d5", err: nil}

	// Create a worker pool with mock engines
	wp := &WorkerPool{
		engines: map[string][]engineinterface.EngineInterface{
			"argo":      {mockArgo},
			"stockfish": {mockStockfish},
		},
		jobChannels: map[string]chan MoveRequest{
			"argo":      make(chan MoveRequest, 10),
			"stockfish": make(chan MoveRequest, 10),
		},
		jobs: make(chan MoveRequest, 10),
	}

	// Start the worker pool
	wp.Start()

	// Submit jobs
	job1 := MoveRequest{ID: "1", Fen: "some_fen_1", Type: "argo"}
	job2 := MoveRequest{ID: "2", Fen: "some_fen_2", Type: "stockfish"}
	wp.SubmitJob(job1)
	wp.SubmitJob(job2)

	// Retrieve results
	result1 := wp.GetResult("1")
	result2 := wp.GetResult("2")

	// Verify results
	if result1 != "e2e4" {
		t.Errorf("Expected result for job 1: e2e4, got: %s", result1)
	}
	if result2 != "d7d5" {
		t.Errorf("Expected result for job 2: d7d5, got: %s", result2)
	}

	// Stop the worker pool
	wp.Stop()
}

func TestWorkerPool_UnknownEngineType(t *testing.T) {
	// Setup a worker pool with no engines
	wp := &WorkerPool{
		engines: map[string][]engineinterface.EngineInterface{
			"stockfish": {},
		},
		jobChannels: map[string]chan MoveRequest{
			"stockfish": make(chan MoveRequest, 10),
		},
		jobs: make(chan MoveRequest, 10),
	}

	// Start the worker pool
	wp.Start()

	// Submit a job with an unknown engine type
	job := MoveRequest{ID: "1", Fen: "some_fen", Type: "unknown_engine"}
	wp.SubmitJob(job)

	// Retrieve the result
	result := wp.GetResult("1")

	// Verify result
	expected := "error: unknown engine type unknown_engine"
	if result != expected {
		t.Errorf("Expected result: %s, got: %s", expected, result)
	}

	// Stop the worker pool
	wp.Stop()
}

func TestWorkerPool_ErrorHandling(t *testing.T) {
	// Setup a mock engine that returns an error
	mockEngine := &MockEngineInterface{bestMove: "", err: fmt.Errorf("engine error")}

	// Create a worker pool with the mock engine
	wp := &WorkerPool{
		engines: map[string][]engineinterface.EngineInterface{
			"argo": {mockEngine},
		},
		jobChannels: map[string]chan MoveRequest{
			"argo": make(chan MoveRequest, 10),
		},
		jobs: make(chan MoveRequest, 10),
	}

	// Start the worker pool
	wp.Start()

	// Submit a job
	job := MoveRequest{ID: "1", Fen: "some_fen", Type: "argo"}
	wp.SubmitJob(job)

	// Retrieve the result
	result := wp.GetResult("1")

	// Verify the error is returned
	expected := "error: engine error"
	if result != expected {
		t.Errorf("Expected result: %s, got: %s", expected, result)
	}

	// Stop the worker pool
	wp.Stop()
}

func TestWorkerPool_Concurrency(t *testing.T) {
	// Setup mock engines
	mockArgo := &MockEngineInterface{bestMove: "e2e4", err: nil}
	mockStockfish := &MockEngineInterface{bestMove: "d7d5", err: nil}

	// Create a worker pool with mock engines
	wp := &WorkerPool{
		engines: map[string][]engineinterface.EngineInterface{
			"argo":      {mockArgo},
			"stockfish": {mockStockfish},
		},
		jobChannels: map[string]chan MoveRequest{
			"argo":      make(chan MoveRequest, 10),
			"stockfish": make(chan MoveRequest, 10),
		},
		jobs: make(chan MoveRequest, 100),
	}

	// Start the worker pool
	wp.Start()

	// Submit jobs concurrently
	numJobs := 100
	wg := sync.WaitGroup{}
	wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		go func(i int) {
			defer wg.Done()
			job := MoveRequest{
				ID:   fmt.Sprintf("%d", i),
				Fen:  fmt.Sprintf("fen_%d", i),
				Type: "argo",
			}
			wp.SubmitJob(job)
		}(i)
	}
	wg.Wait()

	// Verify results
	for i := 0; i < numJobs; i++ {
		jobID := fmt.Sprintf("%d", i)
		result := wp.GetResult(jobID)
		if result != "e2e4" {
			t.Errorf("Expected result for job %d: e2e4, got: %s", i, result)
		}
	}

	// Stop the worker pool
	wp.Stop()
}
