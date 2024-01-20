package storage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Record struct {
	Fields [][]byte
}

func NewRecord() *Record {
	return &Record{
		Fields: make([][]byte, 0),
	}
}

func (r *Record) AddField(field []byte) {
	r.Fields = append(r.Fields, field)
}

func (r *Record) GetField(index int) ([]byte, error) {
	if index < 0 || index >= len(r.Fields) {
		return nil, errors.New("index out of range")
	}
	return r.Fields[index], nil
}

func (r *Record) Serialize() []byte {
	var buffer bytes.Buffer
	for _, field := range r.Fields {
		fieldSize := uint32(len(field))
		binary.Write(&buffer, binary.LittleEndian, fieldSize)
		buffer.Write(field)
	}
	return buffer.Bytes()
}

func DeserializeRecord(data []byte) (*Record, error) {
	buffer := bytes.NewBuffer(data)
	var fields [][]byte

	for buffer.Len() > 0 {
		var fieldSize uint32
		if err := binary.Read(buffer, binary.LittleEndian, &fieldSize); err != nil {
			return nil, err
		}
		field := make([]byte, fieldSize)
		if _, err := buffer.Read(field); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return &Record{Fields: fields}, nil
}

func (r *Record) Size() int {
	size := 0
	for _, field := range r.Fields {
		size += 4 + len(field) // 4 bytes for the length prefix of each field
	}
	return size
}

func (r *Record) Copy() *Record {
	newRecord := NewRecord()
	for _, field := range r.Fields {
		newField := make([]byte, len(field))
		copy(newField, field)
		newRecord.AddField(newField)
	}
	return newRecord
}
