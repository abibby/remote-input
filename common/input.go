package common

import (
	"encoding"
)

type Event interface {
	encoding.BinaryUnmarshaler
	Size() int
	InputEvent() *InputEvent
}
