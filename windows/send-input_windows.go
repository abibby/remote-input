package windows

import (
	"github.com/stephen-fox/user32util"
)

// var (
// 	user32        = syscall.NewLazyDLL("user32.dll")
// 	sendInputProc = user32.NewProc("SendInput")
// )

var user32 *user32util.User32DLL

func init() {
	var err error
	user32, err = user32util.LoadUser32DLL()
	if err != nil {
		panic(err)
	}
}

func SendMouseInput(dx, dy int32, data int32, flags uint32) error {
	return user32util.SendMouseInput(user32util.MouseInput{
		Dx:        dx,
		Dy:        dy,
		MouseData: uint32(data),
		DwFlags:   flags,
	}, user32)
}

func SendInput(key VirtualKey, flag KeyEventFlag) error {
	return user32util.SendKeydbInput(user32util.KeybdInput{
		WVK:     uint16(key),
		DwFlags: uint32(flag),
	}, user32)
	// type keyboardInput struct {
	// 	wVk         uint16
	// 	wScan       uint16
	// 	dwFlags     uint32
	// 	time        uint32
	// 	dwExtraInfo uint64
	// }

	// type input struct {
	// 	inputType uint32
	// 	ki        keyboardInput
	// 	padding   uint64
	// }

	// i := input{
	// 	inputType: 1,
	// 	ki: keyboardInput{
	// 		wVk:     uint16(key),
	// 		dwFlags: uint32(flag),
	// 	},
	// }
	// ret, _, err := sendInputProc.Call(
	// 	uintptr(1),
	// 	uintptr(unsafe.Pointer(&i)),
	// 	uintptr(unsafe.Sizeof(i)),
	// )
	// if ret != 1 {
	// 	return err
	// }
	// return nil
}
