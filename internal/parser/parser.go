package parser

import (
	"fmt"

	logger "github.com/roackb2/simple_db/internal/log"
)

type Parser struct {
	lex       *Lexer
	errors    []string
	curToken  Token
	peekToken Token
}

func NewParser(lex *Lexer) *Parser {
	parser := &Parser{
		lex:    lex,
		errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (parser *Parser) PrintParser() {
	logger.Debug("curToken: %+v\n", parser.curToken)
	logger.Debug("peekToken: %+v\n", parser.peekToken)
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lex.nextToken()
	logger.Debug("after parser.nextToken, curToken: %s, peekToken: %s\n", parser.curToken.Literal, parser.peekToken.Literal)
}

func (parser *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", string(t), parser.peekToken.Literal)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) expectPeek(tokenType TokenType) bool {
	logger.Debug("expect peek token %s: %s\n", parser.peekToken.Literal, tokenType)
	if parser.peekToken.Type == tokenType {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}

func (parser *Parser) parseInsertStatement() *Statement {
	insertStatement := &InsertStatement{}
	if !parser.expectPeek(INTO) {
		return nil
	}
	if !parser.expectPeek(IDENTIFIER) {
		return nil
	}
	insertStatement.TableName = parser.curToken.Literal
	if !parser.expectPeek(OPEN_PARENTHESIS) {
		return nil
	}
	if !parser.expectPeek(IDENTIFIER) {
		return nil
	}
	insertStatement.Columns = append(insertStatement.Columns, parser.curToken.Literal)
	fmt.Println("peek token type", parser.peekToken.Type)
	for parser.peekToken.Type == COMMA {
		parser.nextToken()
		fmt.Println("next token is comma", parser.curToken.Literal)
		parser.nextToken()
		fmt.Println("next token is identifier", parser.curToken.Literal)
		insertStatement.Columns = append(insertStatement.Columns, parser.curToken.Literal)
	}
	insertStatement.Columns = append(insertStatement.Columns, parser.curToken.Literal)
	if !parser.expectPeek(CLOSE_PARENTHESIS) {
		return nil
	}
	if !parser.expectPeek(VALUES) {
		return nil
	}
	if !parser.expectPeek(OPEN_PARENTHESIS) {
		return nil
	}
	if !parser.expectPeek(STRING) || !parser.expectPeek(NUMBER) {
		return nil
	}
	insertStatement.Values = append(insertStatement.Values, parser.curToken.Literal)
	for parser.peekToken.Type == COMMA {
		parser.nextToken()
		parser.nextToken()
		insertStatement.Values = append(insertStatement.Values, parser.curToken.Literal)
	}
	insertStatement.Values = append(insertStatement.Values, parser.curToken.Literal)

	if !parser.expectPeek(CLOSE_PARENTHESIS) {
		return nil
	}
	return &Statement{PrepareSuccess, StatementInsert, "", insertStatement, nil}
}

func (parser *Parser) parseSelectStatement() *Statement {
	selectStmt := &SelectStatement{}

	// Logic to parse fields, table name, and WHERE clause

	return &Statement{
		PrepareRes:    PrepareSuccess,
		StatementType: StatementSelect,
		SelectStmt:    selectStmt,
	}
}

func PrintTokens(lexer *Lexer) {
	for {
		tok := lexer.nextToken()
		logger.Debug("%+v\n", tok)
		if tok.Type == EOF {
			break
		}
	}
}

func (parser *Parser) ParseStatement() *Statement {
	logger.Debug("Statement starts with: %s\n", parser.curToken.Type)
	switch parser.curToken.Type {
	case INSERT:
		return parser.parseInsertStatement()
	case SELECT:
		return parser.parseSelectStatement()
	default:
		return &Statement{PrepareFail, StatementUnknown, "", nil, nil}
	}
}
