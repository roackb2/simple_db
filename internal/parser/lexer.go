package parser

import (
	"strings"

	logger "github.com/roackb2/simple_db/internal/log"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize the first character
	return l
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	switch strings.ToUpper(ident) {
	case "INSERT":
		return INSERT
	case "INTO":
		return INTO
	case "VALUES":
		return VALUES
	case "SELECT":
		return SELECT
	case "FROM":
		return FROM
	case "WHERE":
		return WHERE
	default:
		return IDENTIFIER
	}
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.readPosition]
	}
	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) readNumber() string {
	position := lex.position
	for isDigit(lex.ch) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

func (lex *Lexer) skipWhitespace() {
	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\n' || lex.ch == '\r' {
		lex.readChar()
	}
}

func (lex *Lexer) readIdentifier() string {
	position := lex.position
	for isLetter(lex.ch) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

func (lex *Lexer) readToken(tokenType TokenType, ch byte) Token {
	tok := newToken(tokenType, ch)
	lex.readChar()
	return tok
}

func (lex *Lexer) nextToken() Token {
	var tok Token

	lex.skipWhitespace()

	logger.Debug("Reading token character %s\n", string(lex.ch))
	switch lex.ch {
	case '(':
		tok = lex.readToken(OPEN_PARENTHESIS, lex.ch)
	case ')':
		tok = lex.readToken(CLOSE_PARENTHESIS, lex.ch)
	case '\'':
		tok = lex.readToken(SINGLE_QUOTE, lex.ch)
	case ',':
		tok = lex.readToken(COMMA, lex.ch)
	case 0:
		tok = lex.readToken(EOF, byte(0))
	default:
		if isLetter(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
		} else if isDigit(lex.ch) {
			tok.Literal = lex.readNumber()
			tok.Type = NUMBER
		} else {
			tok = lex.readToken(ILLEGAL, lex.ch)
		}
	}
	logger.Debug("Parsed token '%s' is of type %s\n", tok.Literal, tok.Type)
	return tok
}
