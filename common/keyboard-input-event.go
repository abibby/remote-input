package common

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"time"
)

type InputEventType uint16

const (
	EV_SYN       = InputEventType(0x00)
	EV_KEY       = InputEventType(0x01)
	EV_REL       = InputEventType(0x02)
	EV_ABS       = InputEventType(0x03)
	EV_MSC       = InputEventType(0x04)
	EV_SW        = InputEventType(0x05)
	EV_LED       = InputEventType(0x11)
	EV_SND       = InputEventType(0x12)
	EV_REP       = InputEventType(0x14)
	EV_FF        = InputEventType(0x15)
	EV_PWR       = InputEventType(0x16)
	EV_FF_STATUS = InputEventType(0x17)
	EV_MAX       = InputEventType(0x1f)
	EV_CNT       = InputEventType((EV_MAX + 1))
)

func (t InputEventType) String() string {
	switch t {
	case EV_SYN:
		return "EV_SYN"
	case EV_KEY:
		return "EV_KEY"
	case EV_REL:
		return "EV_REL"
	case EV_ABS:
		return "EV_ABS"
	case EV_MSC:
		return "EV_MSC"
	case EV_SW:
		return "EV_SW"
	case EV_LED:
		return "EV_LED"
	case EV_SND:
		return "EV_SND"
	case EV_REP:
		return "EV_REP"
	case EV_FF:
		return "EV_FF"
	case EV_PWR:
		return "EV_PWR"
	case EV_FF_STATUS:
		return "EV_FF_STATUS"
	case EV_MAX:
		return "EV_MAX"
	case EV_CNT:
		return "EV_CNT"
	}
	return fmt.Sprintf("unknown %d", t)
}

type InputEvent struct {
	Time      time.Time
	EventType InputEventType
	Code      uint16
	Value     int32
}

var _ encoding.BinaryUnmarshaler = (*InputEvent)(nil)
var _ encoding.BinaryMarshaler = (*InputEvent)(nil)

func (e *InputEvent) Size() int {
	return 24
}

func (e *InputEvent) UnmarshalBinary(b []byte) error {
	if len(b) != e.Size() {
		return fmt.Errorf("invalid length expecting 24 recieved %d", len(b))
	}
	sec := binary.LittleEndian.Uint64(b[0:8])
	usec := binary.LittleEndian.Uint64(b[8:16])
	e.Time = time.Unix(int64(sec), int64(usec))

	e.EventType = InputEventType(binary.LittleEndian.Uint16(b[16:18]))
	e.Code = binary.LittleEndian.Uint16(b[18:20])
	e.Value = int32(binary.LittleEndian.Uint32(b[20:]))
	return nil
}

func (e *InputEvent) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, e.Size())
	data = binary.LittleEndian.AppendUint64(data, uint64(e.Time.Unix()))
	data = binary.LittleEndian.AppendUint64(data, uint64(e.Time.Nanosecond()))
	data = binary.LittleEndian.AppendUint16(data, uint16(e.EventType))
	data = binary.LittleEndian.AppendUint16(data, uint16(e.Code))
	data = binary.LittleEndian.AppendUint32(data, uint32(e.Value))
	return data, nil
}
