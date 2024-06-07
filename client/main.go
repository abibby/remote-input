package main

import (
	"io"
	"log"
	"net"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/windows"
	"github.com/stephen-fox/user32util"
)

func main() {
	log.Printf("started")

	serverIP := "192.168.2.38:38808"
	// serverIP := "localhost:38808"

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("connected to %s", serverIP)

	b := make([]byte, 32)
	for {
		e := &common.InputEvent{}
		_, err := conn.Read(b)
		if err == io.EOF {
			log.Print("disconnected")
			return
		} else if err != nil {
			log.Print(err)
			continue
		}

		err = e.UnmarshalBinary(b)
		if err != nil {
			log.Print(err)
			continue
		}

		if e.Keyboard != nil {
			handleKeyboard(e.Keyboard)
		} else if e.Mouse != nil {
			handleMouse(e.Mouse)
		}

	}

}

func handleMouse(e *common.MouseInputEvent) {
	flags := user32util.MouseEventFMove

	if e.ButtonLeft() {
		flags &= user32util.MouseEventFLeftDown
	}
	if e.ButtonRight() {
		flags &= user32util.MouseEventFRightDown
	}

	windows.SendMouseInput(int32(e.X), int32(e.Y)*-1, flags)
}

func handleKeyboard(e *common.KeyboardInputEvent) {
	if e.EventType != common.EV_KEY {
		return
	}

	vKey, ok := keyMap[e.Code]
	if !ok {
		log.Printf("no map for key code %d", e.Code)
		return
	}
	var flag windows.KeyEventFlag
	if e.Value == 1 {
		flag = windows.KEYEVENTF_KEYPRESS
	} else if e.Value == 0 {
		flag = windows.KEYEVENTF_KEYUP
	}
	err := windows.SendInput(vKey, flag)
	if err != nil {
		log.Printf("failed to send input %v", err)
		return
	}

}
