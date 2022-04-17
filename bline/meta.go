package bline

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Meta ...
type Meta struct {
	Start    int64
	StepMsec int64
	Length   int64
	Comment  [100]byte
}

// MetaV1 ...
type MetaV1 struct {
	Start    int64
	StepMsec int64
	Length   int64
	Comment  [100]byte
}

func readMeta(r io.Reader, version uint16) (Meta, error) {
	meta := Meta{}
	switch version {
	case 1:
		metaV1 := MetaV1{}
		err := binary.Read(r, binary.LittleEndian, &metaV1)
		if err != nil {
			return meta, err
		}
		meta.Start = metaV1.Start
		meta.StepMsec = metaV1.StepMsec
		meta.Length = metaV1.Length
		meta.Comment = metaV1.Comment
	default:
		return Meta{}, fmt.Errorf("unknown Meta version: %d", version)
	}

	return meta, nil
}

func writeMeta(w io.Writer, meta Meta, version uint16) error {
	switch version {
	case 1:
		metaV1 := MetaV1{
			Start:    meta.Start,
			StepMsec: meta.StepMsec,
			Length:   meta.Length,
			Comment:  meta.Comment,
		}
		return binary.Write(w, binary.LittleEndian, metaV1)
	}

	return fmt.Errorf("unknown Meta version: %d", version)
}
