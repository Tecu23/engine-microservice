package engine

import (
	"fmt"
	"sync"
)

type EnginePool struct {
	engines chan Engine
	mu      sync.Mutex
	config  *EngineConfig
}

type EngineConfig struct {
	EngineType string
	Path       string
	PoolSize   int
}

func NewEnginePool(config *EngineConfig) (*EnginePool, error) {
	pool := &EnginePool{
		engines: make(chan Engine, config.PoolSize),
		config:  config,
	}

	for i := 0; i < config.PoolSize; i++ {
		eng, err := pool.createEngineInstance()
		if err != nil {
			return nil, fmt.Errorf("failed to create engine instace: %v", err)
		}

		pool.engines <- eng
	}

	return pool, nil
}

func (p *EnginePool) createEngineInstance() (Engine, error) {
	switch p.config.EngineType {
	case "stockfish":
		eng, err := NewStockfishEngine(p.config.Path)
		if err != nil {
			return nil, err
		}
		return eng, nil
	// Add cases for other engine types
	default:
		return nil, fmt.Errorf("unsupported engine type: %s", p.config.EngineType)
	}
}

func (p *EnginePool) GetEngine() (Engine, error) {
	select {
	case eng := <-p.engines:
		return eng, nil
	default:
		// Optionally, you can create a new engine instance if the pool is empty
		return nil, fmt.Errorf("no available engine instances")
	}
}

func (p *EnginePool) ReturnEngine(eng Engine) {
	p.engines <- eng
}

func (p *EnginePool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.engines)
	for eng := range p.engines {
		eng.Close()
	}
}
