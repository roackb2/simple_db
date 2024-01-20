package parser

type PrepareResultCode int64
type StatementTypeCode int64

const (
	PrepareFail    PrepareResultCode = 0
	PrepareSuccess PrepareResultCode = 1
)

const (
	StatementUnknown StatementTypeCode = 0
	StatementSelect  StatementTypeCode = 1
	StatementInsert  StatementTypeCode = 2
)

type WhereClause struct {
	Column   string
	Operator string
	Value    string
}

type SelectStatement struct {
	Fields    []string
	TableName string
	Where     *WhereClause
}

type InsertStatement struct {
	TableName string
	Columns   []string
	Values    []string
}

type Statement struct {
	PrepareRes    PrepareResultCode
	StatementType StatementTypeCode
	Raw           string
	InsertStmt    *InsertStatement
	SelectStmt    *SelectStatement
}
