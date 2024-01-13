package parser

type TokenType string

const (
	ILLEGAL           = "ILLEGAL"
	EOF               = "EOF"
	WS                = "WS"
	INSERT            = "INSERT"
	INTO              = "INTO"
	VALUES            = "VALUES"
	IDENTIFIER        = "IDENTIFIER"
	COMMA             = "COMMA"
	OPEN_PARENTHESIS  = "OPEN_PARENTHESIS"
	CLOSE_PARENTHESIS = "CLOSE_PARENTHESIS"
	SINGLE_QUOTE      = "SINGLE_QUOTE"
	STRING            = "STRING" // string values
	NUMBER            = "NUMBER" // integer values
	SELECT            = "SELECT"
	FROM              = "FROM"
	WHERE             = "WHERE"
	EQUALS            = "EQUALS"
)

type Token struct {
	Type    TokenType
	Literal string
}
