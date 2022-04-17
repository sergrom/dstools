package bline

import (
	"bytes"
	"encoding/binary"
	"os"
	"time"
)

const (
	curVersion uint16 = 1
)

// BLine ...
type BLine struct {
	version uint16
	meta    Meta
	data    []byte
}

// NewBLine ...
func NewBLine(stepMsec int64, startTime, endTime time.Time) *BLine {
	length := (endTime.UnixMilli()-startTime.UnixMilli())/stepMsec + 1
	return &BLine{
		version: 1,
		meta: Meta{
			Start:    startTime.UnixMilli(),
			StepMsec: stepMsec,
			Length:   length,
		},
		data: make([]byte, length/8+1),
	}
}

// NewBLineFromFile ...
func NewBLineFromFile(fName string) (*BLine, error) {
	data, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}

	return NewBLineFromBytes(data)
}

// NewBLineFromBytes ...
func NewBLineFromBytes(data []byte) (*BLine, error) {
	bl := BLine{}
	buf := bytes.NewReader(data)

	err := binary.Read(buf, binary.LittleEndian, &bl.version)
	if err != nil {
		return nil, err
	}

	m, err := readMeta(buf, bl.version)
	if err != nil {
		return nil, err
	}
	bl.meta = m
	bl.data = make([]byte, bl.meta.Length/8+1)

	err = binary.Read(buf, binary.LittleEndian, &bl.data)
	if err != nil {
		return nil, err
	}

	return &bl, nil
}

// GetBytes ...
func (bl *BLine) GetBytes() ([]byte, error) {
	var buf bytes.Buffer

	// 1. Version bytes
	err := binary.Write(&buf, binary.LittleEndian, curVersion)
	if err != nil {
		return nil, err
	}

	// 2. Meta bytes
	if err := writeMeta(&buf, bl.meta, curVersion); err != nil {
		return nil, err
	}

	// 3. Data bytes
	err = binary.Write(&buf, binary.LittleEndian, bl.data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
