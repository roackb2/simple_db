package executor

import (
	"github.com/roackb2/simple_db/internal/parser"
	"github.com/roackb2/simple_db/internal/storage"
)

// Executor is responsible for executing SQL statements.
type Executor struct {
	bufferManager *storage.BufferPool
}

// NewExecutor creates a new Executor.
func NewExecutor(bufferManager *storage.BufferPool) *Executor {
	return &Executor{
		bufferManager: bufferManager,
	}
}

// ExecuteInsertStatement takes an InsertStatement and writes it to the appropriate pages.
func (e *Executor) ExecuteInsertStatement(insertStmt *parser.InsertStatement) error {
	// TODO: Table and schema lookup should be performed here.
	// TODO: Constraint checks should be performed here.

	for _, values := range insertStmt.Values {
		record := storage.NewRecord()
		for _, value := range values {
			// Assuming value is a string, we convert it to a byte slice.
			// If it's not, you'll need to convert it accordingly.
			record.AddField([]byte(string(value)))
		}

		// Serialize the record for storage.
		recordData := record.Serialize()

		// TODO: Find the right page to insert the record.
		// This could be a new page or an existing one with enough space.
		pageID := e.findPageForRecord(recordData)
		page, err := e.bufferManager.FetchPage(pageID)
		if err != nil {
			return err
		}

		// Add the record to the page.
		if _, err := page.AddRecord(recordData); err != nil {
			return err
		}

		// Assume bufferPage is the in-memory representation of the page.
		bufferPage, err := e.bufferManager.GetBufferPage(pageID)
		if err != nil {
			return err
		}
		bufferPage.IsDirty = true
	}

	// TODO: Flushing to disk will be handled by the buffer manager.
	return nil
}

func (e *Executor) findPageForRecord(recordData []byte) int64 {
	// Placeholder: Implement logic to find the appropriate page ID for the new record.
	// This might involve checking for pages with enough free space or allocating a new page.
	return 0 // Example page ID
}
