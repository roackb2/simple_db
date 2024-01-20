package main

import (
	"bufio"
	"errors"
	"os"

	"github.com/roackb2/simple_db/internal/executor"
	logger "github.com/roackb2/simple_db/internal/log"
	"github.com/roackb2/simple_db/internal/parser"
	"github.com/roackb2/simple_db/internal/repl"
	"github.com/roackb2/simple_db/internal/storage"
)

func main() {
	repl.PrintUsage()

	filePath := "./db"

	// Initialize the BufferPool and Executor
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		os.Create(filePath)
	}
	bufferPool, err := storage.NewBufferPool(filePath, 1024) // Example capacity
	if err != nil {
		logger.Error("Failed to create buffer pool: %v", err)
		return
	}
	metadataManager := storage.NewMetadataManager()
	exec := executor.NewExecutor(bufferPool, metadataManager)

	for {
		repl.PrintPrompt()
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if repl.IsMetaCommand(input) {
			switch repl.HandleMetaCommand(input) {
			case repl.CmdSuccess:
				continue
			case repl.CmdUnrecognized:
				logger.Debug("Unrecognized command: %s", input)
				continue
			}
		}

		stmt := parser.PrepareStatement(input)

		// Handle execution of the statement
		if err := exec.ExecuteStatement(*stmt); err != nil {
			logger.Error("Failed to execute statement: %v\n", err)
			continue
		}
	}
}
