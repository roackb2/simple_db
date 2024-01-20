package parser

import (
	"fmt"

	logger "github.com/roackb2/simple_db/internal/log"
	stmt "github.com/roackb2/simple_db/internal/statement"
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

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lex.nextToken()
	logger.Debug("after parser.nextToken, curToken: %s, peekToken: %s\n", parser.curToken.Literal, parser.peekToken.Literal)
}

func (parser *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, parser.peekToken.Literal)
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

func (parser *Parser) parseInsertStatement() *stmt.Statement {
	insertStatement := &stmt.InsertStatement{}
	// INTO
	if !parser.expectPeek(INTO) {
		return nil
	}
	// table name
	if !parser.expectPeek(IDENTIFIER) {
		return nil
	}
	insertStatement.TableName = parser.curToken.Literal
	// column names
	if !parser.expectPeek(OPEN_PARENTHESIS) {
		return nil
	}
	if !parser.expectPeek(IDENTIFIER) {
		return nil
	}
	insertStatement.Columns = append(insertStatement.Columns, parser.curToken.Literal)
	for parser.peekToken.Type == COMMA {
		parser.nextToken()
		parser.nextToken()
		insertStatement.Columns = append(insertStatement.Columns, parser.curToken.Literal)
	}
	insertStatement.Columns = append(insertStatement.Columns, parser.curToken.Literal)
	if !parser.expectPeek(CLOSE_PARENTHESIS) {
		return nil
	}
	// VALUES
	if !parser.expectPeek(VALUES) {
		return nil
	}
	// column values
	if !parser.expectPeek(OPEN_PARENTHESIS) {
		return nil
	}
	if !parser.expectPeek(STRING) {
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
	return &stmt.Statement{stmt.PrepareSuccess, stmt.StatementInsert, "", insertStatement, nil}
}

func (parser *Parser) parseSelectStatement() *stmt.Statement {
	selectStmt := &stmt.SelectStatement{}

	// Logic to parse fields, table name, and WHERE clause

	return &stmt.Statement{
		PrepareRes:    stmt.PrepareSuccess,
		StatementType: stmt.StatementSelect,
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

func (parser *Parser) ParseStatement() *stmt.Statement {
	logger.Debug("Statement starts with: %s\n", parser.curToken.Type)
	switch parser.curToken.Type {
	case INSERT:
		return parser.parseInsertStatement()
	case SELECT:
		return parser.parseSelectStatement()
	default:
		return &stmt.Statement{stmt.PrepareFail, stmt.StatementUnknown, "", nil, nil}
	}
}
