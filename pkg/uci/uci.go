// Package uci should handle the engine interface for communication
package uci

// EngineInterface is the interface connecting to the server
type EngineInterface interface {
	GetBestMove(fen string) (string, error)
}
