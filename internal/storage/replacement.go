package storage

import (
	"container/list"
)

// ReplacementPolicy is an interface for page replacement algorithms.
type ReplacementPolicy interface {
	ChoosePageToEvict(pool map[int64]*BufferPage) int64
}

// LRUPolicy implements the ReplacementPolicy interface using LRU logic.
type LRUPolicy struct {
	evictList *list.List              // A doubly linked list to implement LRU
	entries   map[int64]*list.Element // Map to keep track of page access order
}

// lruEntry is used to hold the value in the evictList.
type lruEntry struct {
	key int64
}

// NewLRUPolicy creates a new LRUPolicy.
func NewLRUPolicy() *LRUPolicy {
	return &LRUPolicy{
		evictList: list.New(),
		entries:   make(map[int64]*list.Element),
	}
}

// ChoosePageToEvict selects the least recently used page for eviction.
func (l *LRUPolicy) ChoosePageToEvict(pool map[int64]*BufferPage) int64 {
	for {
		if l.evictList.Len() == 0 {
			return -1
		}
		// Fetch the oldest accessed page from the back of the evictList
		elem := l.evictList.Back()
		if elem == nil {
			return -1
		}
		entry := elem.Value.(*lruEntry)
		if _, ok := pool[entry.key]; ok && !pool[entry.key].IsPinned {
			// If the page is not pinned, return it for eviction
			l.evictList.Remove(elem)
			delete(l.entries, entry.key)
			return entry.key
		}
		// If the page is pinned, move to the next oldest page
		l.evictList.Remove(elem)
		delete(l.entries, entry.key)
	}
}

// PageAccessed updates the LRU policy when a page is accessed.
func (l *LRUPolicy) PageAccessed(pageID int64) {
	// If the page is already in the access order map, move it to the front
	if elem, ok := l.entries[pageID]; ok {
		l.evictList.MoveToFront(elem)
		return
	}
	// If it's a new page, add it to the access order map and the front of the evictList
	elem := l.evictList.PushFront(&lruEntry{key: pageID})
	l.entries[pageID] = elem
}

// PageRemoved updates the LRU policy when a page is removed from the buffer.
func (l *LRUPolicy) PageRemoved(pageID int64) {
	if elem, ok := l.entries[pageID]; ok {
		l.evictList.Remove(elem)
		delete(l.entries, pageID)
	}
}
