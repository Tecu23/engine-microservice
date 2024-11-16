package engineinterface

// EngineInterface is the interface connecting to the server
type EngineInterface interface {
	GetBestMove(fen string) (string, error)
}
