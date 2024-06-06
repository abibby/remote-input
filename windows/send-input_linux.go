package windows

import "log"

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
	log.Printf("SendInput(%#v)\n", i)
	return nil
}
