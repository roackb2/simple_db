package storage

import (
	"errors"
	"os"
	"sync"
)

// BufferPage wraps around the logical Page to include buffer-specific metadata.
type BufferPage struct {
	PageID   int64 // Unique identifier for the page
	PageData *Page // The logical Page structure, defined in page.go
	IsDirty  bool  // Indicates if the page has been modified
	IsPinned bool  // Indicates if the page is currently being used
}

// BufferPool holds the buffered pages in memory.
type BufferPool struct {
	mu                sync.RWMutex
	pool              map[int64]*BufferPage
	capacity          int
	diskFile          *os.File          // The file descriptor for the database file on disk
	replacementPolicy ReplacementPolicy // Interface for the page replacement policy
}

// NewBufferPool initializes a new BufferPool.
func NewBufferPool(diskFilePath string, capacity int) (*BufferPool, error) {
	file, err := os.OpenFile(diskFilePath, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &BufferPool{
		pool:              make(map[int64]*BufferPage),
		capacity:          capacity,
		diskFile:          file,
		replacementPolicy: NewLRUPolicy(), // Initialize LRU or any other policy
	}, nil
}

// FetchPage retrieves a page from the buffer pool or disk.
func (bp *BufferPool) FetchPage(pageID int64) (*Page, error) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	// If page is in pool, return it
	if page, exists := bp.pool[pageID]; exists {
		page.IsPinned = true // Pin the page, indicating it is in use
		return page.PageData, nil
	}

	// If not, and if the pool is full, evict a page
	if len(bp.pool) >= bp.capacity {
		err := bp.evictPage()
		if err != nil {
			return nil, err
		}
	}

	// Read page from disk
	pageData, err := bp.readPageFromDisk(pageID)
	if err != nil {
		return nil, err
	}

	// Add the page to the pool
	bufferPage := &BufferPage{
		PageID:   pageID,
		PageData: pageData,
		IsPinned: true,
	}
	bp.pool[pageID] = bufferPage
	return pageData, nil
}

// GetBufferPage retrieves a buffered page by its page ID.
func (bp *BufferPool) GetBufferPage(pageID int64) (*BufferPage, error) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	bufferPage, exists := bp.pool[pageID]
	if !exists {
		return nil, errors.New("page not found in buffer pool")
	}
	return bufferPage, nil
}

// FlushPage writes a page back to disk if it's dirty.
func (bp *BufferPool) FlushPage(pageID int64) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	page, exists := bp.pool[pageID]
	if !exists {
		return errors.New("page not found in buffer pool")
	}

	if page.IsDirty {
		err := bp.writePageToDisk(page.PageID, page.PageData)
		if err != nil {
			return err
		}
		page.IsDirty = false
	}

	return nil
}

func (bp *BufferPool) WriteRecordToPage(page *Page, recordData []byte) error {
	// Append the record to the page data.
	page.Data = append(page.Data, recordData...)
	return nil
}

// writePageToDisk writes a given page to disk.
func (bp *BufferPool) writePageToDisk(pageID int64, page *Page) error {
	// Serialize the page data
	pageData := page.Serialize()

	// Write to disk at the correct offset
	_, err := bp.diskFile.WriteAt(pageData, pageID*int64(PageSize))
	return err
}

// readPageFromDisk reads a page from disk.
func (bp *BufferPool) readPageFromDisk(pageID int64) (*Page, error) {
	pageData := make([]byte, PageSize)
	_, err := bp.diskFile.ReadAt(pageData, pageID*int64(PageSize))
	if err != nil {
		return nil, err
	}
	// Assuming DeserializePage is a function defined in page.go that converts bytes to a Page structure
	return DeserializePage(pageData)
}

// evictPage selects and evicts a page from the buffer pool based on the replacement policy.
func (bp *BufferPool) evictPage() error {
	evictPageID := bp.replacementPolicy.ChoosePageToEvict(bp.pool)
	if evictPageID == -1 {
		return errors.New("no page to evict")
	}

	err := bp.FlushPage(evictPageID)
	if err != nil {
		return err
	}

	delete(bp.pool, evictPageID)
	return nil
}
