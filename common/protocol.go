package common

import (
	"encoding"
	"encoding/binary"
	"fmt"

	"github.com/abibby/remote-input/windows"
)

const protoVersion = 0

type KeyEvent struct {
	Key   windows.VirtualKey
	Flags windows.KeyEventFlag
}

var _ encoding.BinaryMarshaler = (*KeyEvent)(nil)
var _ encoding.BinaryUnmarshaler = (*KeyEvent)(nil)

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (k *KeyEvent) UnmarshalBinary(data []byte) error {
	if len(data) != 7 {
		return fmt.Errorf("invalid length")
	}

	version := data[0]
	if version != protoVersion {
		return fmt.Errorf("invalid version: expected %d got %d", protoVersion, version)
	}
	key := binary.BigEndian.Uint16(data[1:])
	flags := binary.BigEndian.Uint32(data[3:])
	k.Key = windows.VirtualKey(key)
	k.Flags = windows.KeyEventFlag(flags)
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (k *KeyEvent) MarshalBinary() (data []byte, err error) {
	version := byte(1)
	data = make([]byte, 7)
	data[0] = version
	binary.BigEndian.PutUint16(data[1:], uint16(k.Key))
	binary.BigEndian.PutUint32(data[3:], uint32(k.Flags))
	return data, nil
}
