package parser

import (
	"fmt"

	stmt "github.com/roackb2/simple_db/internal/statement"
)

func PrepareStatement(input string) *stmt.Statement {
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
