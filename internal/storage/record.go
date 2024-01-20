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
		binary.Write(&buffer, binary.LittleEndian, int32(len(field)))
		buffer.Write(field)
	}
	return buffer.Bytes()
}

func DeserializeRecord(data []byte) (*Record, error) {
	buffer := bytes.NewBuffer(data)
	record := NewRecord()

	for buffer.Len() > 0 {
		var length int32
		if err := binary.Read(buffer, binary.LittleEndian, &length); err != nil {
			return nil, err
		}

		field := make([]byte, length)
		if _, err := buffer.Read(field); err != nil {
			return nil, err
		}

		record.AddField(field)
	}

	return record, nil
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
