package parser

import (
	"fmt"
)

func PrepareStatement(input string) *Statement {
	fmt.Println("Preparing statement", input)
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	statement := parser.ParseStatement()

	if len(parser.Errors()) > 0 {
		fmt.Println("Parser errors encountered:")
		for _, err := range parser.Errors() {
			fmt.Println("\t", err)
		}
	} else {
		fmt.Printf("Parsed Statement: %#v\n", statement)
	}
	return statement
}
