package windows

import (
	"syscall"
	"unsafe"
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	sendInputProc = user32.NewProc("SendInput")
)

func SendInput(key VirtualKey, flag KeyEventFlag) error {
	type keyboardInput struct {
		wVk         uint16
		wScan       uint16
		dwFlags     uint32
		time        uint32
		dwExtraInfo uint64
	}

	type input struct {
		inputType uint32
		ki        keyboardInput
		padding   uint64
	}

	i := input{
		inputType: 1,
		ki: keyboardInput{
			wVk:     uint16(key),
			dwFlags: uint32(flag),
		},
	}
	ret, _, err := sendInputProc.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Sizeof(i)),
	)
	if ret != 1 {
		return err
	}
	return nil
}
