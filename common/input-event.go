package common

import (
	"encoding"
	"fmt"
)

const version byte = 1

type InputEventType byte

const (
	TypeKeyboard = InputEventType(0)
	TypeMouse    = InputEventType(1)
)

type InputEvent struct {
	Mouse    *MouseInputEvent
	Keyboard *KeyboardInputEvent
}

var _ encoding.BinaryUnmarshaler = (*InputEvent)(nil)
var _ encoding.BinaryMarshaler = (*InputEvent)(nil)

// MarshalBinary implements encoding.BinaryMarshaler.
func (i *InputEvent) MarshalBinary() ([]byte, error) {
	data := make([]byte, 32)
	data[0] = version

	if (i.Keyboard == nil) == (i.Mouse == nil) {
		return nil, fmt.Errorf("exactly one of .Mouse or .Keyboard must be set")
	}

	var v encoding.BinaryMarshaler
	if i.Keyboard != nil {
		data[1] = byte(TypeKeyboard)
		v = i.Keyboard
	} else if i.Mouse != nil {
		data[1] = byte(TypeMouse)
		v = i.Mouse
	}

	b, err := v.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(data[2:2+len(b)], b)

	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (i *InputEvent) UnmarshalBinary(data []byte) error {
	v := data[0]
	if v != version {
		return fmt.Errorf("invalid version %d expected %d", v, version)
	}

	typ := data[1]

	switch InputEventType(typ) {
	case TypeKeyboard:
		i.Keyboard = &KeyboardInputEvent{}
		return i.Keyboard.UnmarshalBinary(data[2:26])
	case TypeMouse:
		i.Mouse = &MouseInputEvent{}
		return i.Mouse.UnmarshalBinary(data[2:5])
	default:
		return fmt.Errorf("unexpected type %d", typ)
	}
}
