package repl

import (
	"fmt"
	"os"
	"strings"
)

func PrintUsage() {
	fmt.Println("Simple DB 0.0.1")
	fmt.Println("Type .exit to exit")
}

func PrintPrompt() {
	fmt.Print("db > ")
}

func IsMetaCommand(input string) bool {
	return input[0] == '.'
}

func HandleMetaCommand(input string) CmdRes {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return CmdUnrecognized
	}
	cmd := lines[0][1:]
	switch cmd {
	case CmdExit:
		fmt.Println("Exiting")
		os.Exit(0)
	case CmdListTable:
		fmt.Println("Listing tables (not yet implemented)")
		return CmdSuccess
	}
	return CmdUnrecognized
}
