package common

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"time"
)

type KeyboardInputEventType uint16

const (
	EV_SYN       = KeyboardInputEventType(0x00)
	EV_KEY       = KeyboardInputEventType(0x01)
	EV_REL       = KeyboardInputEventType(0x02)
	EV_ABS       = KeyboardInputEventType(0x03)
	EV_MSC       = KeyboardInputEventType(0x04)
	EV_SW        = KeyboardInputEventType(0x05)
	EV_LED       = KeyboardInputEventType(0x11)
	EV_SND       = KeyboardInputEventType(0x12)
	EV_REP       = KeyboardInputEventType(0x14)
	EV_FF        = KeyboardInputEventType(0x15)
	EV_PWR       = KeyboardInputEventType(0x16)
	EV_FF_STATUS = KeyboardInputEventType(0x17)
	EV_MAX       = KeyboardInputEventType(0x1f)
	EV_CNT       = KeyboardInputEventType((EV_MAX + 1))
)

func (t KeyboardInputEventType) String() string {
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
		return "EV_FF_"
	case EV_MAX:
		return "EV_MAX"
	case EV_CNT:
		return "EV_CNT"
	}
	return fmt.Sprintf("unknown %d", t)
}

type KeyboardInputEvent struct {
	Time      time.Time
	EventType KeyboardInputEventType
	Code      KeyCode
	Value     int32
}

var _ encoding.BinaryUnmarshaler = (*KeyboardInputEvent)(nil)
var _ encoding.BinaryMarshaler = (*KeyboardInputEvent)(nil)
var _ Event = (*KeyboardInputEvent)(nil)

func (e *KeyboardInputEvent) Size() int {
	return 24
}

// InputEvent implements Event.
func (e *KeyboardInputEvent) InputEvent() *InputEvent {
	return &InputEvent{
		Keyboard: e,
	}
}

func (e *KeyboardInputEvent) UnmarshalBinary(b []byte) error {
	if len(b) != 24 {
		return fmt.Errorf("invalid length expecting 24 recieved %d", len(b))
	}
	sec := binary.LittleEndian.Uint64(b[0:8])
	usec := binary.LittleEndian.Uint64(b[8:16])
	e.Time = time.Unix(int64(sec), int64(usec))

	e.EventType = KeyboardInputEventType(binary.LittleEndian.Uint16(b[16:18]))
	e.Code = KeyCode(binary.LittleEndian.Uint16(b[18:20]))
	e.Value = int32(binary.LittleEndian.Uint32(b[20:]))
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (e *KeyboardInputEvent) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, 24)
	data = binary.LittleEndian.AppendUint64(data, uint64(e.Time.Unix()))
	data = binary.LittleEndian.AppendUint64(data, uint64(e.Time.Nanosecond()))
	data = binary.LittleEndian.AppendUint16(data, uint16(e.EventType))
	data = binary.LittleEndian.AppendUint16(data, uint16(e.Code))
	data = binary.LittleEndian.AppendUint32(data, uint32(e.Value))
	return data, nil
}
