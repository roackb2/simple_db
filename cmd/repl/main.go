package main

import (
	"bufio"
	"os"

	logger "github.com/roackb2/simple_db/internal/log"
	"github.com/roackb2/simple_db/internal/parser"
	"github.com/roackb2/simple_db/internal/repl"
)

func main() {
	repl.PrintUsage()
	for {
		repl.PrintPrompt()
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if repl.IsMetaCommand(input) {
			switch repl.HandleMetaCommand(input) {
			case repl.CmdSuccess:
				continue
			case repl.CmdUnrecognized:
				logger.Debug("Unrecognized command: ", input)
				continue
			}
		}

		stmt := *parser.PrepareStatement(input)
		switch stmt.PrepareRes {
		case parser.PrepareSuccess:
			parser.ExecuteStatement(stmt)
			continue
		case parser.PrepareFail:
			logger.Debug("Unrecognized keyword at start of ", input)
			continue
		}
	}
}
