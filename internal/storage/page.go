package storage

import (
	"encoding/binary"
	"errors"
)

const (
	PageSize       = 4096
	SlotSize       = 6 // 2 bytes for offset, 4 bytes for length
	PageHeaderSize = 4 // Assume 4 bytes for free space pointer
)

// SlotDescriptor describes the location and size of a record on the page.
type SlotDescriptor struct {
	Offset int16 // using int16 to allow for -1 sentinel value
	Length uint32
}

// Page represents a database page that manages variable-length records.
type Page struct {
	Data              []byte
	FreeSpacePointer  int16
	RecordDescriptors []SlotDescriptor
}

// NewPage initializes a new Page with a given size.
func NewPage() *Page {
	pageData := make([]byte, PageSize)
	// Initialize free space pointer to the start of the data area
	binary.LittleEndian.PutUint16(pageData, uint16(PageHeaderSize))
	return &Page{
		Data:              pageData,
		FreeSpacePointer:  PageHeaderSize,
		RecordDescriptors: []SlotDescriptor{},
	}
}

// AddRecord adds a new record to the page.
func (p *Page) AddRecord(recordData []byte) (int, error) {
	recordSize := uint32(len(recordData))
	neededSpace := int(recordSize) + SlotSize

	// Check if there's enough space for the record and the slot descriptor
	if int(p.FreeSpacePointer)+neededSpace > PageSize {
		return -1, errors.New("not enough space on the page")
	}

	// Find an empty slot or create a new one
	slotIndex := -1
	for i, slot := range p.RecordDescriptors {
		if slot.Offset == -1 { // Deleted record slot
			slotIndex = i
			break
		}
	}
	if slotIndex == -1 { // No empty slot found, create a new one
		slotIndex = len(p.RecordDescriptors)
		p.RecordDescriptors = append(p.RecordDescriptors, SlotDescriptor{})
	}

	// Write the record to the data area
	copy(p.Data[p.FreeSpacePointer:], recordData)

	// Set the slot descriptor
	p.RecordDescriptors[slotIndex] = SlotDescriptor{
		Offset: int16(p.FreeSpacePointer),
		Length: recordSize,
	}

	// Update the free space pointer
	p.FreeSpacePointer += int16(neededSpace)

	return slotIndex, nil
}

// RetrieveRecord retrieves a record from the page by its slot index.
func (p *Page) RetrieveRecord(slotIndex int) ([]byte, error) {
	if slotIndex < 0 || slotIndex >= len(p.RecordDescriptors) {
		return nil, errors.New("slot index out of range")
	}

	slot := p.RecordDescriptors[slotIndex]
	if slot.Offset == -1 {
		return nil, errors.New("record has been deleted")
	}

	// Extract the record data
	recordData := p.Data[slot.Offset : slot.Offset+int16(slot.Length)]
	return recordData, nil
}

// DeleteRecord marks a record as deleted by setting its offset to -1.
func (p *Page) DeleteRecord(slotIndex int) error {
	if slotIndex < 0 || slotIndex >= len(p.RecordDescriptors) {
		return errors.New("slot index out of range")
	}

	p.RecordDescriptors[slotIndex].Offset = -1
	p.RecordDescriptors[slotIndex].Length = 0

	return nil
}

// CompactPage compacts the page by removing gaps left by deleted records.
func (p *Page) CompactPage() {
	compactedData := make([]byte, PageSize)
	copy(compactedData, p.Data[:PageHeaderSize]) // Copy the header
	var newDescriptors []SlotDescriptor

	compactPointer := PageHeaderSize
	for _, descriptor := range p.RecordDescriptors {
		if descriptor.Offset == -1 {
			// Skip deleted records
			continue
		}
		recordData := p.Data[descriptor.Offset : descriptor.Offset+int16(descriptor.Length)]
		copy(compactedData[compactPointer:], recordData)

		newDescriptor := SlotDescriptor{
			Offset: int16(compactPointer),
			Length: descriptor.Length,
		}
		newDescriptors = append(newDescriptors, newDescriptor)
		compactPointer += int(descriptor.Length) + SlotSize
	}

	// Update the page's data and descriptors
	p.Data = compactedData
	p.RecordDescriptors = newDescriptors
	p.FreeSpacePointer = int16(compactPointer)
}

// Serialize converts the Page into a byte slice for storage on disk.
func (p *Page) Serialize() []byte {
	buf := make([]byte, PageSize)

	// Write the free space pointer at the beginning
	binary.LittleEndian.PutUint16(buf, uint16(p.FreeSpacePointer))

	// Write the records based on the descriptor information
	for _, descriptor := range p.RecordDescriptors {
		if descriptor.Offset != -1 { // Record is not marked deleted
			start := int(descriptor.Offset)
			end := start + int(descriptor.Length)
			copy(buf[start:end], p.Data[start:end])
		}
	}

	// Write the slot directory at the end of the page
	for i, descriptor := range p.RecordDescriptors {
		offset := PageSize - (i+1)*SlotSize
		binary.LittleEndian.PutUint16(buf[offset:], uint16(descriptor.Offset))
		binary.LittleEndian.PutUint32(buf[offset+2:], descriptor.Length)
	}

	return buf
}

// DeserializePage converts a byte slice from disk into a Page structure.
func DeserializePage(buf []byte) (*Page, error) {
	if len(buf) != PageSize {
		return nil, errors.New("incorrect buffer size for page")
	}

	// Read the free space pointer from the beginning
	freeSpacePointer := binary.LittleEndian.Uint16(buf)

	// Initialize an empty Page structure
	p := &Page{
		Data:              make([]byte, PageSize),
		FreeSpacePointer:  int16(freeSpacePointer),
		RecordDescriptors: make([]SlotDescriptor, 0),
	}

	// Read the slot directory from the end of the page
	for i := 0; i < (PageSize-PageHeaderSize)/SlotSize; i++ {
		offset := PageSize - (i+1)*SlotSize
		recordOffset := int16(binary.LittleEndian.Uint16(buf[offset:]))
		recordLength := binary.LittleEndian.Uint32(buf[offset+2:])

		// Add the descriptor if the record isn't deleted
		if recordOffset != -1 {
			p.RecordDescriptors = append(p.RecordDescriptors, SlotDescriptor{
				Offset: recordOffset,
				Length: recordLength,
			})
		}
	}

	// Copy the record data into the Page structure
	copy(p.Data, buf[:PageSize])

	return p, nil
}

func (page *Page) hasEnoughSpaceFor(recordData []byte) bool {
	// Logic to check if the page has enough space for the record.
	// This will depend on the page format and how free space is tracked.
	// Placeholder logic:
	return len(page.Data)+len(recordData) <= PageSize
}
