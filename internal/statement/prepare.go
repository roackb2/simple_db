package statement

import (
	"strings"
)

func PrepareStatement(input string) *Statement {
	if strings.HasPrefix(input, "select") {
		return &Statement{PrepareSuccess, StatementSelect, input}
	}
	if strings.HasPrefix(input, "insert") {
		return &Statement{PrepareSuccess, StatementInsert, input}
	}
	return &Statement{PrepareFail, StatementUnknown, input}
}
