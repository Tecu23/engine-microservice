package engine

// EngineInterface is the interface connecting to the server
type Engine interface {
	Initialize() error
	CalculateBestMove(fen string, depth int) (string, error)
	Info() string
	Close() error
}
