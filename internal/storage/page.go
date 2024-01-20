package storage

import (
	"encoding/binary"
	"errors"
	"math"
)

const (
	SlotSize       = 6 // 2 bytes for offset, 4 bytes for length
	PageHeaderSize = 4 // Assume 4 bytes for free space pointer
)

type Page struct {
	Data             []byte
	FreeSpacePointer uint16
}

type Slot struct {
	Offset uint16
	Length uint32
}

func NewPage(pageSize int) *Page {
	return &Page{
		Data:             make([]byte, pageSize),
		FreeSpacePointer: PageHeaderSize,
	}
}

func (p *Page) AddRecord(record []byte) error {
	if p.FreeSpacePointer+uint16(len(record))+SlotSize > uint16(len(p.Data)) {
		return errors.New("not enough space")
	}
	slotOffset := p.FreeSpacePointer
	// Write record
	copy(p.Data[slotOffset:], record)
	// Update free space pointer
	p.FreeSpacePointer += uint16(len(record))
	// Write slot
	slot := Slot{
		Offset: slotOffset,
		Length: uint32(len(record)),
	}
	p.writeSlot(slot)
	return nil
}

func (p *Page) GetRecord(slotIndex int) ([]byte, error) {
	slot, err := p.readSlot(slotIndex)
	if err != nil {
		return nil, err
	}
	if slot.Length == 0 {
		return nil, errors.New("deleted or empty record")
	}
	record := p.Data[slot.Offset : slot.Offset+uint16(slot.Length)]
	return record, nil
}

func (p *Page) writeSlot(slot Slot) {
	// Convert the free space pointer to the slot index
	slotIndex := int(p.FreeSpacePointer / SlotSize)
	// Write the slot at the end of the page
	offsetBytes := make([]byte, 2)
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint16(offsetBytes, slot.Offset)
	binary.LittleEndian.PutUint32(lengthBytes, slot.Length)
	slotPosition := len(p.Data) - (slotIndex+1)*SlotSize
	copy(p.Data[slotPosition:], offsetBytes)
	copy(p.Data[slotPosition+2:], lengthBytes)
}

func (p *Page) readSlot(slotIndex int) (Slot, error) {
	if slotIndex < 0 || slotIndex >= int(math.Ceil(float64(len(p.Data)-PageHeaderSize)/float64(SlotSize))) {
		return Slot{}, errors.New("slot index out of range")
	}
	slotPosition := len(p.Data) - (slotIndex+1)*SlotSize
	offset := binary.LittleEndian.Uint16(p.Data[slotPosition:])
	length := binary.LittleEndian.Uint32(p.Data[slotPosition+2:])
	return Slot{Offset: offset, Length: length}, nil
}

func (p *Page) DeleteRecord(slotIndex int) error {
	// Read the slot to get the record's offset and length
	slot, err := p.readSlot(slotIndex)
	if err != nil {
		return err
	}

	// Check if the record is already deleted
	if slot.Length == 0 {
		return errors.New("record is already deleted")
	}

	// Mark the slot as deleted by setting the length to 0
	slot.Length = 0
	p.writeSlot(slot)

	// p.compactPage()

	return nil
}

func (p *Page) compactPage() {
	newData := make([]byte, len(p.Data))
	copy(newData, p.Data[:PageHeaderSize]) // Copy the header

	newFreeSpacePointer := uint16(PageHeaderSize)
	var newSlots []Slot

	for _, slot := range p.readAllSlots() {
		if slot.Length == 0 {
			// Skip deleted records
			continue
		}

		// Copy active record to new data slice
		copy(newData[newFreeSpacePointer:], p.Data[slot.Offset:slot.Offset+uint16(slot.Length)])

		// Create new slot with updated offset
		newSlot := Slot{
			Offset: newFreeSpacePointer,
			Length: slot.Length,
		}
		newSlots = append(newSlots, newSlot)

		// Update free space pointer
		newFreeSpacePointer += uint16(slot.Length)
	}

	// Write new slots to the end of the newData slice
	for i, slot := range newSlots {
		slotPosition := len(newData) - (i+1)*SlotSize
		binary.LittleEndian.PutUint16(newData[slotPosition:], slot.Offset)
		binary.LittleEndian.PutUint32(newData[slotPosition+2:], slot.Length)
	}

	// Update the page's data slice and free space pointer
	p.Data = newData
	p.FreeSpacePointer = newFreeSpacePointer
}

func (p *Page) readAllSlots() []Slot {
	// Helper function to read all slots from the slot directory
	var slots []Slot
	for slotIndex := 0; slotIndex < int(math.Ceil(float64(len(p.Data)-PageHeaderSize)/float64(SlotSize))); slotIndex++ {
		slot, _ := p.readSlot(slotIndex) // Error handling can be improved here
		slots = append(slots, slot)
	}
	return slots
}
