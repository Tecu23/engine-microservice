package stockfishengine

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

// StockfishEngine creates a new interface to communicate with the engine
type StockfishEngine struct {
	cmd    *exec.Cmd
	writer *bufio.Writer
	reader *bufio.Scanner
	mu     sync.Mutex // making the engine thread safe
}

// NewStockfishEngine creates a new interface to communicate with the engine
func NewStockfishEngine() (*StockfishEngine, error) {
	cmd := exec.Command("./bin/engines/stockfish")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(stdin)
	reader := bufio.NewScanner(stdout)

	writer.WriteString("uci\n")
	writer.Flush()

	for reader.Scan() {
		if reader.Text() == "uciok" {
			break
		}
	}

	return &StockfishEngine{
		cmd:    cmd,
		writer: writer,
		reader: reader,
	}, nil
}

// GetBestMove returns the best move from the engine
func (i *StockfishEngine) GetBestMove(fen string) (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.writer.WriteString(fmt.Sprintf("position fen %s\n", fen))
	i.writer.WriteString("go depth 10\n")
	i.writer.Flush()

	for i.reader.Scan() {
		line := i.reader.Text()
		if len(line) > 8 && line[:8] == "bestmove" {
			return line[9:13], nil
		}
	}

	return "", fmt.Errorf("failed to get the best move from engine")
}
