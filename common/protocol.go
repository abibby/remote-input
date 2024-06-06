package common

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"time"
)

type InputEvent struct {
	Time  time.Time
	Type  uint16
	Code  uint16
	Value int32
}

var _ encoding.BinaryUnmarshaler = (*InputEvent)(nil)

func (e *InputEvent) UnmarshalBinary(b []byte) error {
	if len(b) != 24 {
		return fmt.Errorf("invalid length")
	}
	sec := binary.LittleEndian.Uint64(b[0:8])
	usec := binary.LittleEndian.Uint64(b[8:16])
	e.Time = time.Unix(int64(sec), int64(usec))

	e.Type = binary.LittleEndian.Uint16(b[16:18])
	e.Code = binary.LittleEndian.Uint16(b[18:20])
	e.Value = int32(binary.LittleEndian.Uint32(b[20:]))
	return nil
}