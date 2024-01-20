package executor

import (
	"fmt"

	stmt "github.com/roackb2/simple_db/internal/statement"
	"github.com/roackb2/simple_db/internal/storage"
)

// Executor is responsible for executing SQL statements.
type Executor struct {
	bufferManager *storage.BufferPool
	metadata      *storage.MetadataManager
}

// NewExecutor creates a new Executor.
func NewExecutor(bufferManager *storage.BufferPool, metadataManager *storage.MetadataManager) *Executor {
	return &Executor{
		bufferManager: bufferManager,
		metadata:      metadataManager,
	}
}

// ExecuteInsertStatement takes an InsertStatement and writes it to the appropriate pages.
func (e *Executor) ExecuteInsertStatement(insertStmt *stmt.InsertStatement) error {
	// Serialize the insert statement into a record
	record := storage.NewRecord()
	for _, value := range insertStmt.Values {
		// Convert the string value to a byte slice and add it to the record
		record.AddField([]byte(value))
	}

	// Find a page to store the record
	lastPageID := e.metadata.GetLastPageID()
	page, err := e.bufferManager.FetchPage(lastPageID)
	if err != nil {
		// Handle the error, possibly by creating a new page
	}

	// Write the record to the page
	if err := e.bufferManager.WriteRecordToPage(page, record.Serialize()); err != nil {
		return err
	}

	// Flush the page to disk to ensure it's saved
	if err := e.bufferManager.FlushPage(lastPageID); err != nil {
		return err
	}

	return nil
}

func (e *Executor) findPageForRecord(recordData []byte) int64 {
	// This is a simplified version. In a real database, you would need to
	// check the table's data pages for an available slot that fits the record size.
	// If no suitable page is found, a new page may need to be allocated.

	// For now, let's assume we are appending records to the last page or creating a new one if full.
	lastPageID := e.getLastPageID() // Method to get the last used page ID from metadata
	lastPage, err := e.bufferManager.FetchPage(lastPageID)
	if err != nil {
		// Handle error (e.g., if the page doesn't exist, allocate a new one)
		return e.allocateNewPage() // Method to allocate a new page and return its ID
	}

	// Check if the last page has enough space for the record.
	// The calculation will depend on the page's internal structure and metadata.
	if lastPage.hasEnoughSpaceFor(recordData) {
		return lastPageID
	} else {
		// If it doesn't fit, allocate a new page.
		return e.allocateNewPage()
	}
}

func (e *Executor) getLastPageID() int64 {
	// Retrieve the last page ID from the database metadata.
	// This information is typically maintained in a system table or file.
	// Placeholder logic:
	return e.metadata.GetLastPageID() // Assume this method exists and returns the last page ID.
}

func (e *Executor) allocateNewPage() int64 {
	// Logic to allocate a new page.
	// This would involve updating database metadata and possibly writing an empty page to disk.
	// Placeholder logic:
	newPageID := e.metadata.IncrementLastPageID() // Assume this increments and returns a new page ID.
	e.bufferManager.AddNewPage(newPageID)         // Assume this adds a new page to the buffer pool.
	return newPageID
}

func (e *Executor) ExecuteStatement(statement stmt.Statement) error {
	switch statement.StatementType {
	case stmt.StatementSelect:
		fmt.Println("This is where we would do a select.")
	case stmt.StatementInsert:
		return e.ExecuteInsertStatement(statement.InsertStmt)
	}
	return nil
}
