package engine

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// ArgoEngine is the object that communicates with the engine
type ArgoEngine struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	reader *bufio.Reader
	mu     sync.Mutex // making the engine thread safe
}

// NewArgoEngine creates a new interface to communicate with the engine
func NewArgoEngine(path string) (*ArgoEngine, error) {
	// cmd := exec.Command("./bin/engines/argo")
	engine := &ArgoEngine{}
	engine.cmd = exec.Command(path)

	var err error
	engine.stdin, err = engine.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	engine.stdout, err = engine.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	engine.reader = bufio.NewReader(engine.stdout)

	if err := engine.cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start engine: %v", err)
	}

	if err := engine.initializeUCI(); err != nil {
		return nil, err
	}

	return engine, nil
}

// GetBestMove returns the best move from the engine
func (e *ArgoEngine) CalculateBestMove(fen string, depth int) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.sendCommand("position fen " + fen)

	e.sendCommand(fmt.Sprintf("go depth %d", depth))

	var bestMove string
	for {
		line, err := e.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading from engine: %v", err)
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "bestmove") {
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				bestMove = parts[1]
			}
			break
		}
	}

	if bestMove == "" {
		return "", fmt.Errorf("no best move found")
	}

	return bestMove, nil
}

func (e *ArgoEngine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.sendCommand("quit")
	return e.cmd.Wait()
}

func (e *ArgoEngine) Info() string {
	return "Argo Engine"
}

func (e *ArgoEngine) initializeUCI() error {
	e.sendCommand("uci")
	for {
		line, err := e.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading from engine: %v", err)
		}

		if strings.HasPrefix(line, "uciok") {
			break
		}
	}

	return nil
}

func (e *ArgoEngine) sendCommand(cmd string) error {
	_, err := io.WriteString(e.stdin, cmd+"\n")
	return err
}
