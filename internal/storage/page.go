package storage

import (
	"errors"
)

const (
	SlotSize       = 6 // 2 bytes for offset, 4 bytes for length
	PageHeaderSize = 4 // Assume 4 bytes for free space pointer
)

// Page represents a database page.
type Page struct {
	Data              []byte
	FreeSpacePointer  uint16
	RecordDescriptors []SlotDescriptor
}

// SlotDescriptor describes the location and size of a record on the page.
type SlotDescriptor struct {
	Offset uint16
	Length uint32
}

// AddRecord adds a new record to the page.
func (p *Page) AddRecord(record *Record) error {
	// Serialize the record to a byte slice.
	recordData := record.Serialize()
	recordSize := uint32(len(recordData))
	if int(p.FreeSpacePointer)+int(recordSize) > len(p.Data) {
		return errors.New("not enough space on the page")
	}

	// Find space for the new record.
	slotDescriptor := SlotDescriptor{
		Offset: p.FreeSpacePointer,
		Length: recordSize,
	}
	p.RecordDescriptors = append(p.RecordDescriptors, slotDescriptor)

	// Write the record data to the page.
	copy(p.Data[p.FreeSpacePointer:], recordData)

	// Update the free space pointer.
	p.FreeSpacePointer += uint16(recordSize)

	return nil
}

// RetrieveRecord retrieves a record from the page by its slot index.
func (p *Page) RetrieveRecord(slotIndex int) (*Record, error) {
	if slotIndex < 0 || slotIndex >= len(p.RecordDescriptors) {
		return nil, errors.New("slot index out of range")
	}

	slotDescriptor := p.RecordDescriptors[slotIndex]
	recordData := p.Data[slotDescriptor.Offset : slotDescriptor.Offset+uint16(slotDescriptor.Length)]
	return DeserializeRecord(recordData)
}

// DeleteRecord marks a record as deleted by setting its length to zero.
func (p *Page) DeleteRecord(slotIndex int) error {
	if slotIndex < 0 || slotIndex >= len(p.RecordDescriptors) {
		return errors.New("slot index out of range")
	}

	// Mark the record as deleted.
	p.RecordDescriptors[slotIndex].Length = 0

	// p.compactPage()

	return nil
}

func (p *Page) CompactPage() error {
	// Temporary slice to hold the compacted data
	compactedData := make([]byte, len(p.Data))
	copy(compactedData, p.Data[:PageHeaderSize]) // Copy the header

	// Pointer to the start of the free space in the compacted data
	compactFreeSpacePointer := PageHeaderSize

	// Temporary slice for the new slot descriptors
	var newSlotDescriptors []SlotDescriptor

	// Iterate over the existing slot descriptors
	for _, slot := range p.RecordDescriptors {
		if slot.Length == 0 {
			// This record is deleted; skip it
			continue
		}

		// Extract the record data
		recordData := p.Data[slot.Offset : slot.Offset+uint16(slot.Length)]
		// Place it into the compacted data slice
		copy(compactedData[compactFreeSpacePointer:], recordData)

		// Create a new slot descriptor for the moved record
		newSlot := SlotDescriptor{
			Offset: uint16(compactFreeSpacePointer),
			Length: slot.Length,
		}
		newSlotDescriptors = append(newSlotDescriptors, newSlot)

		// Update the free space pointer in the compacted data
		compactFreeSpacePointer += int(slot.Length)
	}

	// Replace the page's data with the compacted data
	p.Data = compactedData
	// Update the free space pointer
	p.FreeSpacePointer = uint16(compactFreeSpacePointer)
	// Replace the old slot descriptors with the new ones
	p.RecordDescriptors = newSlotDescriptors

	return nil
}
