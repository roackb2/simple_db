package statement

import (
	"fmt"
)

func ExecuteStatement(stmt Statement) {
	switch stmt.StatementType {
	case StatementSelect:
		fmt.Println("This is where we would do a select.")
	case StatementInsert:
		fmt.Println("This is where we would do an insert.")
	}
}
