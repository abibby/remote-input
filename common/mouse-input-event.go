package common

type MouseInputEvent InputEvent

func (e *MouseInputEvent) X() int {
	return 0
}

// import (
// 	"encoding"
// 	"encoding/binary"
// 	"fmt"
// 	"time"
// )

// type MouseInputEvent struct {
// 	Time      time.Time
// 	EventType InputEventType
// 	Code      InputEventCode
// 	Value     int32
// }

// var _ encoding.BinaryUnmarshaler = (*MouseInputEvent)(nil)
// var _ encoding.BinaryMarshaler = (*MouseInputEvent)(nil)
// var _ Event = (*MouseInputEvent)(nil)

// func (e *MouseInputEvent) Size() int {
// 	return 24
// }

// // InputEvent implements Event.
// func (e *MouseInputEvent) InputEvent() *InputEvent {
// 	return &InputEvent{
// 		Mouse: e,
// 	}
// }

// var last uint64

// // UnmarshalBinary implements encoding.BinaryUnmarshaler.
// func (e *MouseInputEvent) UnmarshalBinary(b []byte) error {
// 	if len(b) != e.Size() {
// 		return fmt.Errorf("invalid length expecting 24 recieved %d", len(b))
// 	}

// 	sec := binary.LittleEndian.Uint64(b[0:8])
// 	usec := binary.LittleEndian.Uint64(b[8:16])
// 	e.Time = time.Unix(int64(sec), int64(usec))

// 	e.EventType = InputEventType(binary.LittleEndian.Uint16(b[16:18]))
// 	e.Code = InputEventCode(binary.LittleEndian.Uint16(b[18:20]))
// 	e.Value = int32(binary.LittleEndian.Uint32(b[20:]))

// 	if last != usec {
// 		fmt.Println()
// 		last = usec
// 	}

// 	fmt.Printf("%d ", e.Time.Nanosecond())
// 	fmt.Printf("%v ", e.EventType)
// 	for _, n := range b[18:] {
// 		fmt.Printf("%08b ", n) // prints 00000000 11111101
// 	}
// 	fmt.Println()
// 	return nil
// }

// // MarshalBinary implements encoding.BinaryMarshaler.
// func (e *MouseInputEvent) MarshalBinary() ([]byte, error) {
// 	data := make([]byte, 0, 24)
// 	data = binary.LittleEndian.AppendUint64(data, uint64(e.Time.Unix()))
// 	data = binary.LittleEndian.AppendUint64(data, uint64(e.Time.Nanosecond()))
// 	data = binary.LittleEndian.AppendUint16(data, uint16(e.EventType))
// 	data = binary.LittleEndian.AppendUint16(data, uint16(e.Code))
// 	data = binary.LittleEndian.AppendUint32(data, uint32(e.Value))
// 	return data, nil
// }
// func (e *MouseInputEvent) ButtonLeft() bool {
// 	// return e.Button&1 != 0
// 	return false
// }
// func (e *MouseInputEvent) ButtonRight() bool {
// 	// return e.Button&2 != 0
// 	return false
// }
// func (e *MouseInputEvent) ButtonMiddle() bool {
// 	// return e.Button&4 != 0
// 	return false
// }
