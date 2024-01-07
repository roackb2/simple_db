package main

import (
	"bufio"
	"fmt"
	"os"

	repl "github.com/roackb2/simple_db/internal/repl"
	statement "github.com/roackb2/simple_db/internal/statement"
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
				fmt.Println("Unrecognized command: ", input)
				continue
			}
		}

		stmt := *statement.PrepareStatement(input)
		switch stmt.PrepareRes {
		case statement.PrepareSuccess:
			statement.ExecuteStatement(stmt)
			continue
		case statement.PrepareFail:
			fmt.Println("Unrecognized keyword at start of ", input)
			continue
		}
	}
}
