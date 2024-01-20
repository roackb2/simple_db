package storage

type MetadataManager struct {
	lastPageID int64
}

func NewMetadataManager() *MetadataManager {
	// In a real scenario, you would load this from disk
	return &MetadataManager{lastPageID: -1} // Start with -1 to indicate no pages
}

func (m *MetadataManager) GetLastPageID() int64 {
	// In a real scenario, you would load this from disk
	if m.lastPageID == -1 {
		m.lastPageID = 0 // Initialize to the first page
	}
	return m.lastPageID
}

func (m *MetadataManager) IncrementLastPageID() int64 {
	m.lastPageID++
	return m.lastPageID
}
