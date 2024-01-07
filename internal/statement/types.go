package statement

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

type Statement struct {
	PrepareRes    PrepareResultCode
	StatementType StatementTypeCode
	Raw           string
}
