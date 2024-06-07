package common

import (
	"encoding"
	"fmt"
)

type MouseInputEvent struct {
	Button byte
	X      int8
	Y      int8
}

var _ encoding.BinaryUnmarshaler = (*MouseInputEvent)(nil)
var _ encoding.BinaryMarshaler = (*MouseInputEvent)(nil)
var _ Event = (*MouseInputEvent)(nil)

func (e *MouseInputEvent) Size() int {
	return 3
}

// InputEvent implements Event.
func (e *MouseInputEvent) InputEvent() *InputEvent {
	return &InputEvent{
		Mouse: e,
	}
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (e *MouseInputEvent) UnmarshalBinary(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("invalid length expecting 24 recieved %d", len(b))
	}
	e.Button = b[0]
	e.X = int8(b[1])
	e.Y = int8(b[2])
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (e *MouseInputEvent) MarshalBinary() (data []byte, err error) {
	return []byte{
		e.Button,
		byte(e.X),
		byte(e.Y),
	}, nil
}

func (e *MouseInputEvent) ButtonLeft() bool {
	return e.Button&0x1 != 0
}
func (e *MouseInputEvent) ButtonRight() bool {
	return e.Button&0x2 != 0
}
func (e *MouseInputEvent) ButtonMiddle() bool {
	return e.Button&0x4 != 0
}
