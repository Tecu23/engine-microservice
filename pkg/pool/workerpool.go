package pool

import (
	"fmt"
	"sync"

	"github.com/Tecu23/engine-microservice/pkg/types/engineinterface"
	"github.com/Tecu23/engine-microservice/pkg/types/engines/argoengine"
	"github.com/Tecu23/engine-microservice/pkg/types/engines/stockfishengine"
)

type MoveRequest struct {
	ID   string
	Fen  string
	Type string
}

type WorkerPool struct {
	engines     map[string][]engineinterface.EngineInterface
	jobChannels map[string]chan MoveRequest
	jobs        chan MoveRequest
	resultsMap  sync.Map
	wg          sync.WaitGroup
}

func NewWorkerPool(numArgo int, numStockfish int) (*WorkerPool, error) {
	engines := make(map[string][]engineinterface.EngineInterface)
	jobChannels := make(map[string]chan MoveRequest)

	engines["stockfish"] = make([]engineinterface.EngineInterface, 0, numStockfish)
	jobChannels["stockfish"] = make(chan MoveRequest)
	for i := 0; i < numStockfish; i++ {
		engine, err := stockfishengine.NewStockfishEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to create argo server: %v", err)
		}
		engines["stockfish"] = append(engines["stockfish"], engine)
	}

	engines["argo"] = make([]engineinterface.EngineInterface, 0, numArgo)
	jobChannels["argo"] = make(chan MoveRequest)
	for i := 0; i < numArgo; i++ {
		engine, err := argoengine.NewArgoEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to create argo server: %v", err)
		}
		engines["argo"] = append(engines["argo"], engine)
	}

	return &WorkerPool{
		engines:     engines,
		jobChannels: jobChannels,
		jobs:        make(chan MoveRequest),
	}, nil
}

func (wp *WorkerPool) Start() {
	for engineType, engines := range wp.engines {
		jobChan := wp.jobChannels[engineType]

		for _, engine := range engines {
			wp.wg.Add(1)
			go func(e engineinterface.EngineInterface, jobChan <-chan MoveRequest) {
				defer wp.wg.Done()
				for job := range jobChan {

					value, ok := wp.resultsMap.Load(job.ID)
					if !ok {
						return
					}

					resultChan := value.(chan string)
					defer close(resultChan)

					bestMove, err := e.GetBestMove(job.Fen)
					if err != nil {
						resultChan <- fmt.Sprintf("error: %v", err)
					} else {
						resultChan <- bestMove
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

		// if engineType != "stockfish" && engineType != "argo" {
		// 	engineType = "stockfish"
		// }

		if ch, ok := wp.jobChannels[engineType]; ok {
			ch <- job
		} else {

			value, ok := wp.resultsMap.Load(job.ID)
			if !ok {
				return
			}

			resultChan := value.(chan string)
			resultChan <- fmt.Sprintf("error: unknown engine type %s", job.Type)
			close(resultChan)
		}
	}

	for _, ch := range wp.jobChannels {
		close(ch)
	}
}

func (wp *WorkerPool) SubmitJob(move MoveRequest) {
	resultChan := make(chan string, 1)
	wp.resultsMap.Store(move.ID, resultChan)
	wp.jobs <- move
}

func (wp *WorkerPool) GetResult(jobID string) string {
	value, ok := wp.resultsMap.Load(jobID)
	if !ok {
		return fmt.Sprintf("result not found for JobID: %s", jobID)
	}

	resultChan := value.(chan string)
	return <-resultChan
}

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}
